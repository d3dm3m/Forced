package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"biohorror/internal/db"
	"biohorror/internal/game"
	"biohorror/internal/game/space"
	"biohorror/pkg/protocol"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for now
		},
	}
	projMgr *game.ProjectileManager
	mu      sync.RWMutex
)

// Packet structure (Local definition or use protocol package)
// We will use protocol.Packet for consistency, but main needs to marshal/unmarshal
type Packet struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type LoginWithTokenPayload struct {
	Token string `json:"token"`
}

type LoginResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Player  *game.Player `json:"player,omitempty"`
}

func main() {
	// Initialize Game Data
	game.LoadGameData()
	log.Println("Space Core: Game Data Loaded")

	// Initialize Database
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Pool.Close()
	log.Println("Space Core: Database Connected")

	// Initialize Redis
	if err := db.ConnectRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Space Core: Redis Connected")

	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Space Core Server listening on WebSocket :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("New Client Connected. Waiting for Token...")

	// 1. Wait for LOGIN_WITH_TOKEN
	player, err := waitForLogin(conn)
	if err != nil {
		log.Printf("Login failed: %v", err)
		conn.WriteJSON(Packet{Type: "LOGIN_FAILURE", Payload: json.RawMessage(`{"message":"Login failed"}`)})
		return
	}

	log.Printf("Player %s authenticated for Space Core.", player.Username)

	// Send LOGIN_SUCCESS
	payloadBytes, _ := json.Marshal(LoginResponse{
		Success: true,
		Message: "Login successful",
		Player:  player,
	})
	conn.WriteJSON(Packet{Type: "LOGIN_SUCCESS", Payload: payloadBytes})

	// 2. Enter Game Loop
	gameLoop(conn, player)
}

func waitForLogin(conn *websocket.Conn) (*game.Player, error) {
	// Set read deadline to prevent hanging connections
	// conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var packet Packet
	if err := json.Unmarshal(message, &packet); err != nil {
		return nil, err
	}

	if packet.Type != "LOGIN_WITH_TOKEN" {
		return nil, log.Output(2, "Expected LOGIN_WITH_TOKEN packet")
	}

	var loginPayload LoginWithTokenPayload
	if err := json.Unmarshal(packet.Payload, &loginPayload); err != nil {
		return nil, err
	}

	// Validate Token
	username, err := game.ValidateTransferToken(loginPayload.Token)
	if err != nil {
		return nil, err
	}

	// Load Player
	repo := game.NewPostgresPlayerRepository()
	player, err := repo.LoadPlayer(username)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, log.Output(2, "Player not found after token validation")
	}

	return player, nil
}

func gameLoop(conn *websocket.Conn, player *game.Player) {
	// Initialize Services
	repo := game.NewPostgresPlayerRepository()
	marketService := game.NewMarketService(repo)
	surgeryService := game.NewSurgeryService(repo)

	// Init Global Projectile Manager if nil (Lazy init for single-instance in this MVP structure)
	// In production, this would be injected.
	if projMgr == nil {
		projMgr = game.NewProjectileManager()
	}

	// Setup Ticker for Game Logic (10Hz)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Channel to signal connection close
	done := make(chan struct{})

	// Action Channel to handle requests from the Reader Goroutine in the Main Loop
	// This ensures thread safety for the Player struct.
	actionChan := make(chan Packet, 10)

	// Start Reader Goroutine
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read error for %s: %v", player.Username, err)
				return
			}

			var packet Packet
			if err := json.Unmarshal(message, &packet); err != nil {
				log.Printf("Invalid packet from %s: %v", player.Username, err)
				continue
			}

			// Forward actionable packets to Main Loop
			actionChan <- packet
		}
	}()

	tickCount := 0

	// Main Loop
	for {
		select {
		case <-done:
			return // Connection closed

		case packet := <-actionChan:
			// Handle Input Packet Safely
			handleActionPacket(packet, conn, player, marketService, surgeryService)

		case <-ticker.C:
			tickCount++

			// 1. Tick Projectiles (Every Tick - 10Hz)
			// Note: 10Hz is slow for projectiles, client interpolation is key.
			mu.Lock()
			impacts := projMgr.UpdateSimulation(0.1) // 100ms
			for _, p := range impacts {
				// Handle Impact (Self is target?)
				// In Space Core, we mostly care if WE got hit or if we hit someone else.
				// Since this is a single-client loop MVP, we check if target is self.
				// But packets come from client saying "Fire at X".
				// Space Logic: If I am target, take damage.
				if p.TargetID == player.ID && p.Outcome == "Hit" {
					player.CurrentHealth -= p.Damage
					sendCombatHit(conn, p.ID, p.TargetID, p.Damage)
				}
			}
			mu.Unlock()

			// 2. Every 10 ticks (1 second), perform Stat Calculation & Bio-Load Logic
			if tickCount%10 == 0 {
				stats := game.CalculateShipStats(player)

				// Apply Rejection / Decay
				if stats.RejectionRate > 0 {
					player.CurrentHealth -= stats.RejectionRate
					log.Printf("Player %s taking rejection damage: %.2f (Health: %.2f)", player.Username, stats.RejectionRate, player.CurrentHealth)
				}

				// Cap Health
				if player.CurrentHealth > stats.MaxHealth {
					player.CurrentHealth = stats.MaxHealth
				}

				// Death Check
				if player.CurrentHealth <= 0 {
					log.Printf("Player %s DIED due to Bio-Load Failure.", player.Username)
					// Respawn / Reset logic would go here
					player.CurrentHealth = 0 // Clamp
				}

				// Broadcast Stats Update
				sendShipStats(conn, player, stats)
			}

			// 3. Save State Periodically (Every 5 seconds / 50 ticks)
			if tickCount%50 == 0 {
				// Note: Services save state on action, but passive decay needs saving too.
				repo.SavePlayerState(player)
			}
		}
	}
}

func handleActionPacket(packet Packet, conn *websocket.Conn, player *game.Player, market *game.MarketService, surgery *game.SurgeryService) {
	switch packet.Type {
	case protocol.PACKET_TYPE_BUY_ITEM:
		var payload protocol.BuyItemPayload
		if err := json.Unmarshal(packet.Payload, &payload); err == nil {
			if err := market.BuyItem(player, payload.ItemID); err != nil {
				log.Printf("BuyItem failed: %v", err)
			} else {
				log.Printf("Player %s bought %s", player.Username, payload.ItemID)
			}
		}
	case protocol.PACKET_TYPE_GRAFT_ORGAN:
		var payload protocol.GraftOrganPayload
		if err := json.Unmarshal(packet.Payload, &payload); err == nil {
			if err := surgery.GraftOrgan(player, payload.InventoryIndex, payload.SlotType, payload.SlotIndex); err != nil {
				log.Printf("GraftOrgan failed: %v", err)
			} else {
				log.Printf("Player %s grafted organ into %s[%d]", player.Username, payload.SlotType, payload.SlotIndex)
			}
		}
	case protocol.PACKET_TYPE_FIRE_WEAPON:
		// Space Combat Fire
		// For MVP, target is dummy or self-test.
		// We'll spawn a Missile aimed at a fixed point for visuals.
		pkt := handleSpaceAttack(player)
		if pkt != nil {
			// Convert to local Packet type for main.go
			// Ideally we use protocol.Packet throughout, but main defines its own struct with identical json tags
			// We can just marshal and send.
			// conn.WriteJSON expects interface{}.
			// We can pass pkt directly if `protocol.Packet` matches `Packet`.
			// `main.Packet` is identical structure.
			conn.WriteJSON(pkt)
		}

	case "PACKET_TYPE_SPACE_STATE":
		// Validation placeholder
		_ = space.Vector3{}
	}
}

func handleSpaceAttack(player *game.Player) *protocol.Packet {
	// Dummy Target (100 units ahead)

	// Assuming Player struct uses X/Y as X/Z plane for space, or X/Y screen plane.
	// Physics uses Vector3. Let's default Z=0.
	startX, startY, startZ := player.PositionX, player.PositionY, 0.0
	targetX, targetY, targetZ := player.PositionX + 100, player.PositionY, 0.0

	mu.Lock()
	proj := projMgr.SpawnProjectile(
		player.ID, "dummy_target",
		startX, startY, startZ,
		targetX, targetY, targetZ,
		50.0, 20.0, 1.0,
		protocol.BEHAVIOR_MISSILE, "Hit",
	)
	mu.Unlock()

	// Build Packet
	payload := protocol.ProjectileSpawnPayload{
		ProjectileID: proj.ID,
		ShooterID:    proj.ShooterID,
		TargetID:     proj.TargetID,
		Behavior:     proj.Behavior,
		Outcome:      proj.Outcome,
		StartX:       proj.X,
		StartY:       proj.Y,
		StartZ:       proj.Z,
		EndX:         proj.TargetX,
		EndY:         proj.TargetY,
		EndZ:         proj.TargetZ,
		Speed:        proj.Speed,
	}

	bytes, _ := json.Marshal(payload)
	log.Printf("Spawned Space Missile: %s", proj.ID)

	return &protocol.Packet{
		Type:    protocol.PACKET_TYPE_PROJECTILE_SPAWN,
		Payload: bytes,
	}
}

func sendCombatHit(conn *websocket.Conn, projID, targetID string, damage float64) {
	payload := protocol.CombatHitPayload{
		ProjectileID: projID,
		TargetID:     targetID,
		Damage:       damage,
	}
	bytes, _ := json.Marshal(payload)
	conn.WriteJSON(Packet{Type: protocol.PACKET_TYPE_COMBAT_HIT, Payload: bytes})
}

func sendShipStats(conn *websocket.Conn, player *game.Player, stats game.DerivedStats) {
	payload := protocol.ShipStatsPayload{
		CurrentHealth: player.CurrentHealth,
		MaxHealth:     stats.MaxHealth,
		BioLoad:       stats.CurrentBioLoad,
		BioCapacity:   stats.BioCapacity,
		Speed:         stats.Speed,
	}

	payloadBytes, _ := json.Marshal(payload)
	packet := Packet{
		Type:    protocol.PACKET_TYPE_SHIP_STATS,
		Payload: payloadBytes,
	}

	conn.WriteJSON(packet)
}
