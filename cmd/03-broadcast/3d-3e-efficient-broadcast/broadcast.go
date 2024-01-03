package main

import (
	"time"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (s *Server) HandleBroadcastAll(msg maelstrom.Message) error {
	body, err := app.Unmarshal(msg.Body)
	if err != nil {
		return err
	}

	s.mu.Lock()
	for _, key := range body.Messages {
		s.store[key] = true
	}
	s.mu.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type": app.MessageBroadcastAllOk,
	})
}

func (s *Server) StartBroadcastLoop() {
	go func() {
		for {
			select {
			case <-s.ticker.C:
				// buffer items and batch send once every interval
				s.SendBroadcastAll()
				break
			case <-s.done:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Server) StopBroadcastLoop() {
	s.done <- true
}

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

	return s.node.Reply(msg, map[string]any{
		"type": app.MessageBroadcastOk,
	})
}

func (s *Server) SendBroadcastAll() {
	s.mu.RLock()
	body := &app.Message{
		Type:     app.MessageBroadcastAll,
		Messages: s.GetKeys(),
	}
	s.mu.RUnlock()

	for _, dest := range s.node.NodeIDs() {
		if dest == s.node.ID() {
			continue
		}

		go func(dst string) {
			// retry with backoff until dest acks
			ack := false

			for i := 1; !ack && i <= Retries; i++ {
				s.node.RPC(dst, body, func(_ maelstrom.Message) error {
					ack = true
					return nil
				})

				time.Sleep(time.Duration(i*BackoffMultiplier) * time.Millisecond)
			}
		}(dest)
	}
}
