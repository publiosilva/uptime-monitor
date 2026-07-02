package notification

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type TopicMessage struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}

type TopicSubscription struct {
	Topic string
	Conn  *websocket.Conn
}

type TopicHub struct {
	mu           sync.RWMutex
	subscritions map[string]map[*websocket.Conn]bool

	Broadcast   chan TopicMessage
	Subscribe   chan TopicSubscription
	Unsubscribe chan TopicSubscription
}

func NewTopicHub() *TopicHub {
	return &TopicHub{
		subscritions: make(map[string]map[*websocket.Conn]bool),
		Broadcast:    make(chan TopicMessage, 250),
		Subscribe:    make(chan TopicSubscription, 50),
		Unsubscribe:  make(chan TopicSubscription, 50),
	}
}

func (h *TopicHub) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sub := <-h.Subscribe:
			h.mu.Lock()
			if _, ok := h.subscritions[sub.Topic]; !ok {
				h.subscritions[sub.Topic] = make(map[*websocket.Conn]bool)
			}
			h.subscritions[sub.Topic][sub.Conn] = true
			h.mu.Unlock()
			log.Printf("Client subscribed to topic: %s", sub.Topic)
		case sub := <-h.Unsubscribe:
			h.mu.Lock()
			if _, ok := h.subscritions[sub.Topic]; ok {
				delete(h.subscritions[sub.Topic], sub.Conn)
				if len(h.subscritions[sub.Topic]) == 0 {
					delete(h.subscritions, sub.Topic)
				}
			}
			h.mu.Unlock()
			log.Printf("Client unsubscribed from topic: %s", sub.Topic)
		case msg := <-h.Broadcast:
			payload, err := json.Marshal(msg.Payload)
			if err != nil {
				log.Printf("Serialization error: %v", err)
				continue
			}

			h.mu.RLock()
			clients, clientsExist := h.subscritions[msg.Topic]
			if !clientsExist || len(clients) == 0 {
				h.mu.RUnlock()
				continue
			}

			for client := range clients {
				go func(c *websocket.Conn, t string) {
					if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
						log.Printf("Write error on topic %s, cleaning up connection", t)
						client.Close()
						h.Unsubscribe <- TopicSubscription{Topic: t, Conn: c}
					}
				}(client, msg.Topic)
			}
			h.mu.RUnlock()
		}
	}
}
