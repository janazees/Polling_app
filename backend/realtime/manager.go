package realtime

import (
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	Hubs  map[string]*Hub
	Mutex sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Hubs: make(map[string]*Hub),
	}
}

func (m *Manager) GetHub(pollCode string) *Hub {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	hub, exists := m.Hubs[pollCode]

	if !exists {
		hub = NewHub()
		m.Hubs[pollCode] = hub
	}

	return hub
}

func (m *Manager) StartRedisSubscriber(
	redisClient *redis.Client,
	pollCode string,
) {
	channel := "poll:" + pollCode + ":updates"

	pubsub := redisClient.Subscribe(context.Background(), channel)

	log.Println("Subscribed to Redis channel:", channel)

	go func() {
		defer pubsub.Close()

		for {
			message, err := pubsub.ReceiveMessage(context.Background())

			if err != nil {
				log.Println("Redis subscriber error:", err)
				continue
			}

			log.Println("Received live update:", message.Payload)

			hub := m.GetHub(pollCode)

			hub.Broadcast([]byte(message.Payload))
		}
	}()
}
