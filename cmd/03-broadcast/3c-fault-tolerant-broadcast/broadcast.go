package main

import (
	"time"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (s *Server) HandleBroadcast(msg maelstrom.Message) error {
	body, err := app.Unmarshal(msg.Body)
	if err != nil {
		return err
	}

	key := body.Message

	s.mu.Lock()
	if _, ok := s.store[key]; ok {
		s.mu.Unlock()
		return nil
	}
	s.store[key] = true
	s.mu.Unlock()

	// broadcast this message to our neighbors
	s.Broadcast(msg.Src, body)

	return s.node.Reply(msg, map[string]any{
		"type": app.MessageBroadcastOk,
	})
}

func (s *Server) Broadcast(src string, body any) {
	for _, dest := range s.neighbors {
		if dest == src {
			continue
		}

		go func(dst string) {
			// retry with backoff until dest acks
			ack := false

			for i := 0; !ack && i < Retries; i++ {
				s.node.RPC(dst, body, func(_ maelstrom.Message) error {
					ack = true
					return nil
				})

				time.Sleep(time.Duration(i*BackoffMultiplier) * time.Millisecond)
			}
		}(dest)
	}
}
