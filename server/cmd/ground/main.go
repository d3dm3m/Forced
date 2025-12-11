package main

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
	"biohorror/internal/db"
	"biohorror/internal/game"
	"biohorror/internal/game/ground"
	"biohorror/pkg/protocol"
)

// Session Management
type GroundSession struct {
	ID        string
	Addr      *net.UDPAddr
	Player    *game.Player
	LastSeen  time.Time
}

var (
	sessions   = make(map[string]*GroundSession)
	mu         sync.RWMutex
	projMgr    *game.ProjectileManager
	gatekeeper game.IGatekeeper
)

// Main
func main() {
	game.LoadGameData()
	log.Println("Ground Core: Game Data Loaded")

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Pool.Close()
	log.Println("Ground Core: Database Connected")

	if err := db.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Ground Core: Redis Connected")

	addr, err := net.ResolveUDPAddr("udp", ":5000")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Ground Core Server listening on UDP :5000")

	repo := game.NewPostgresPlayerRepository()
	projMgr = game.NewProjectileManager()
	gatekeeper = game.NewRedisGatekeeper()

	go gameLoop(conn, repo)

	for {
		buffer := make([]byte, 4096)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		go handlePacket(conn, remoteAddr, buffer[:n], repo)
	}
}

func handlePacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte, repo game.PlayerRepository) {
	var packet protocol.Packet
	if err := json.Unmarshal(data, &packet); err != nil {
		return
	}

	switch packet.Type {
	case "REQUEST_LOGIN":
		var loginPayload struct {
			Username string `json:"username"`
		}
		if err := json.Unmarshal(packet.Payload, &loginPayload); err == nil {
			player, err := repo.LoadPlayer(loginPayload.Username)
			if err != nil || player == nil {
				if player == nil { player, _ = repo.CreatePlayer(loginPayload.Username, "marine") }
			}
			if player == nil { return }

			mu.Lock()
			sessions[loginPayload.Username] = &GroundSession{
				ID:       loginPayload.Username,
				Addr:     addr,
				Player:   player,
				LastSeen: time.Now(),
			}
			mu.Unlock()

			loginSuccess := protocol.LoginSuccessPayload{
				Message:    "Logged in",
				GroundGear: player.GroundGear,
			}
			payloadBytes, _ := json.Marshal(loginSuccess)
			response := protocol.Packet{ Type: "LOGIN_SUCCESS", Payload: payloadBytes }
			sendPacket(conn, addr, response)
		}

	case protocol.PACKET_TYPE_GROUND_MOVEMENT:
		handleMovement(packet.Payload, addr)

	case protocol.PACKET_TYPE_GROUND_ATTACK:
		handleAttack(conn, packet.Payload, addr)

	case "REQUEST_LAUNCH":
		handleLaunchRequest(conn, addr)
	}
}

func handleLaunchRequest(conn *net.UDPConn, addr *net.UDPAddr) {
	mu.Lock()
	defer mu.Unlock()

	var session *GroundSession
	for _, s := range sessions {
		if s.Addr.String() == addr.String() {
			session = s
			break
		}
	}

	if session == nil {
		return
	}

	// Generate Token
	token, err := game.GenerateTransferToken(session.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return
	}

	// Determine Destination Server
	systemID := session.Player.SystemID
	if systemID == "" {
		systemID = "Sol-0" // Default
	}

	destURL, err := gatekeeper.GetServerForSystem(systemID)
	if err != nil {
		log.Printf("Failed to resolve gatekeeper route for system %s: %v", systemID, err)
		// Fallback or error handling
		destURL = "ws://localhost:8080/ws"
	}

	// Send Back
	payload := map[string]string{
		"token": token,
		"url":   destURL,
	}
	bytes, _ := json.Marshal(payload)

	packet := protocol.Packet{
		Type:    "PACKET_TYPE_LAUNCH_GRANTED",
		Payload: bytes,
	}
	sendPacket(conn, addr, packet)

	log.Printf("Launch Granted for %s. Token: %s", session.ID, token)
}

func handleMovement(payload json.RawMessage, addr *net.UDPAddr) {
	var move protocol.GroundMovementPayload
	if err := json.Unmarshal(payload, &move); err != nil { return }

	mu.Lock()
	defer mu.Unlock()

	var session *GroundSession
	for _, s := range sessions {
		if s.Addr.String() == addr.String() {
			session = s
			break
		}
	}

	if session == nil { return }

	ctx := ground.ValidationContext{
		LastPosition: ground.Vector2{X: session.Player.PositionX, Y: session.Player.PositionY},
		LastTimestamp: session.LastSeen,
		MaxSpeed: 20.0,
	}
	newState := ground.ClientState{ PositionX: move.X, PositionY: move.Y }

	if valid, _ := ground.ValidateMovement(newState, ctx); valid {
		session.Player.PositionX = move.X
		session.Player.PositionY = move.Y
		session.LastSeen = time.Now()
	}
}

func handleAttack(conn *net.UDPConn, payload json.RawMessage, addr *net.UDPAddr) {
	var attack protocol.GroundAttackPayload
	if err := json.Unmarshal(payload, &attack); err != nil { return }

	mu.RLock()
	var attacker *GroundSession
	var target *GroundSession
	for _, s := range sessions {
		if s.Addr.String() == addr.String() {
			attacker = s
			break
		}
	}
	if attack.TargetID != "" { target = sessions[attack.TargetID] }
	mu.RUnlock() // Unlock early, careful with pointer validity

	if attacker == nil || target == nil { return }

	// Calculate Hit/Miss
	hit := projMgr.CalculateHitChance(attacker.Player, target.Player)
	outcome := "Miss"
	if hit { outcome = "Hit" }

	// Spawn Projectile
	speed := 20.0 // Units per sec
	damage := 10.0

	// Create Projectile in Manager (Thread Safe?)
	// projMgr is global, but map is not thread safe.
	// Need mutex for projMgr or use global mu.
	// Re-using global mu for simplicity in MVP.
	mu.Lock()
	// Using Z=0 for Ground (2D Plane)
	// Behavior = Linear (Bullet), TurnRate = 0
	proj := projMgr.SpawnProjectile(
		attacker.ID, target.ID,
		attacker.Player.PositionX, attacker.Player.PositionY, 0,
		target.Player.PositionX, target.Player.PositionY, 0,
		speed, damage, 0.0,
		protocol.BEHAVIOR_LINEAR, outcome,
	)
	mu.Unlock()

	// Broadcast Spawn
	payloadData := protocol.ProjectileSpawnPayload{
		ProjectileID: proj.ID,
		ShooterID:    attacker.ID,
		TargetID:     target.ID,
		Behavior:     proj.Behavior,
		Outcome:      outcome,
		StartX:       proj.X,
		StartY:       proj.Y,
		StartZ:       proj.Z,
		EndX:         proj.TargetX,
		EndY:         proj.TargetY,
		EndZ:         proj.TargetZ,
		Speed:        speed,
	}

	bytes, _ := json.Marshal(payloadData)
	packet := protocol.Packet{ Type: protocol.PACKET_TYPE_PROJECTILE_SPAWN, Payload: bytes }
	packetBytes, _ := json.Marshal(packet)

	mu.RLock()
	for _, s := range sessions {
		conn.WriteToUDP(packetBytes, s.Addr)
	}
	mu.RUnlock()
}

func gameLoop(conn *net.UDPConn, repo game.PlayerRepository) {
	ticker := time.NewTicker(50 * time.Millisecond) // 20Hz
	defer ticker.Stop()

	tickCount := 0

	for range ticker.C {
		tickCount++

		// 1. Tick Projectiles
		mu.Lock()
		impacts := projMgr.UpdateSimulation(0.05) // 50ms

		// Handle Impacts
		for _, p := range impacts {
			if p.Outcome == "Hit" {
				if target, ok := sessions[p.TargetID]; ok {
					target.Player.CurrentHealth -= p.Damage

					// Send Hit Event
					hitPayload := protocol.CombatHitPayload{
						ProjectileID: p.ID,
						TargetID:     p.TargetID,
						Damage:       p.Damage,
					}
					bytes, _ := json.Marshal(hitPayload)
					packet := protocol.Packet{ Type: protocol.PACKET_TYPE_COMBAT_HIT, Payload: bytes }
					pktBytes, _ := json.Marshal(packet)

					// Inner loop broadcast, inefficient but works for MVP
					for _, s := range sessions {
						conn.WriteToUDP(pktBytes, s.Addr)
					}
				}
			}
		}
		mu.Unlock()

		// 2. Broadcast State
		broadcastState(conn)

		// 3. Auto Save
		if tickCount % 100 == 0 {
			saveAllSessions(repo)
		}
	}
}

func saveAllSessions(repo game.PlayerRepository) {
	mu.RLock()
	var playersToSave []*game.Player
	for _, s := range sessions {
		playersToSave = append(playersToSave, s.Player)
	}
	mu.RUnlock()

	for _, p := range playersToSave {
		if err := repo.SavePlayerState(p); err != nil {
			log.Printf("Failed to auto-save %s: %v", p.Username, err)
		}
	}
}

func broadcastState(conn *net.UDPConn) {
	mu.RLock()
	defer mu.RUnlock()

	var entities []protocol.GroundEntityState
	now := time.Now()
	for _, s := range sessions {
		if now.Sub(s.LastSeen) < 5*time.Second {
			entities = append(entities, protocol.GroundEntityState{
				ID: s.ID,
				X:  s.Player.PositionX,
				Y:  s.Player.PositionY,
			})
		}
	}

	if len(entities) == 0 { return }

	payload := protocol.GroundStatePayload{ Entities: entities }
	payloadBytes, _ := json.Marshal(payload)

	packet := protocol.Packet{ Type: protocol.PACKET_TYPE_GROUND_STATE, Payload: payloadBytes }
	packetBytes, _ := json.Marshal(packet)

	for _, s := range sessions {
		if now.Sub(s.LastSeen) < 5*time.Second {
			conn.WriteToUDP(packetBytes, s.Addr)
		}
	}
}

func sendPacket(conn *net.UDPConn, addr *net.UDPAddr, packet protocol.Packet) {
	bytes, _ := json.Marshal(packet)
	conn.WriteToUDP(bytes, addr)
}
