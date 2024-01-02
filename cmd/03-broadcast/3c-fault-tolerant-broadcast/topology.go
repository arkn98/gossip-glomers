package main

import (
	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (s *Server) HandleTopology(msg maelstrom.Message) error {
	body, err := app.Unmarshal(msg.Body)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.neighbors = body.Topology[s.node.ID()]
	s.mu.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type": app.MessageTopologyOk,
	})
}
