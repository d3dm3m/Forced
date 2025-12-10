package game

import (
	"encoding/json"
	"math"
	"net"
	"sync"
)

type ChatMessage struct {
	SenderID string `json:"sender_id"`
	Channel  string `json:"channel"` // "local", "system", "global"
	Text     string `json:"text"`
}

type ChatSession struct {
	PlayerID string
	Addr     *net.UDPAddr
	X, Y     float64 // Last known position
}

type ChatManager struct {
	sessions map[string]*ChatSession
	mu       sync.RWMutex
	conn     *net.UDPConn
}

func NewChatManager(conn *net.UDPConn) *ChatManager {
	return &ChatManager{
		sessions: make(map[string]*ChatSession),
		conn:     conn,
	}
}

func (cm *ChatManager) RegisterSession(playerID string, addr *net.UDPAddr, x, y float64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.sessions[playerID] = &ChatSession{
		PlayerID: playerID,
		Addr:     addr,
		X:        x,
		Y:        y,
	}
}

func (cm *ChatManager) UpdatePosition(playerID string, x, y float64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if session, ok := cm.sessions[playerID]; ok {
		session.X = x
		session.Y = y
	}
}

func (cm *ChatManager) BroadcastLocal(senderID string, text string, rangeTiles float64) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	sender, ok := cm.sessions[senderID]
	if !ok {
		return
	}

	msg := ChatMessage{
		SenderID: senderID,
		Channel:  "local",
		Text:     text,
	}
	packetBytes, _ := json.Marshal(struct {
		Type    string      `json:"type"`
		Payload ChatMessage `json:"payload"`
	}{
		Type:    "PACKET_TYPE_CHAT_MESSAGE",
		Payload: msg,
	})

	for _, session := range cm.sessions {
		// Calculate distance
		dx := session.X - sender.X
		dy := session.Y - sender.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		if dist <= rangeTiles {
			cm.conn.WriteToUDP(packetBytes, session.Addr)
		}
	}
}

func (cm *ChatManager) BroadcastSystem(text string) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	msg := ChatMessage{
		SenderID: "SYSTEM",
		Channel:  "system",
		Text:     text,
	}
	packetBytes, _ := json.Marshal(struct {
		Type    string      `json:"type"`
		Payload ChatMessage `json:"payload"`
	}{
		Type:    "PACKET_TYPE_CHAT_MESSAGE",
		Payload: msg,
	})

	for _, session := range cm.sessions {
		cm.conn.WriteToUDP(packetBytes, session.Addr)
	}
}
