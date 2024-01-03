package main

import (
	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (s *Server) HandleRead(msg maelstrom.Message) error {
	return s.node.Reply(msg, map[string]any{
		"type":     app.MessageReadOk,
		"messages": s.GetKeys(),
	})
}

func (s *Server) GetKeys() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]int, 0, len(s.store))
	for k, _ := range s.store {
		keys = append(keys, k)
	}

	return keys
}
