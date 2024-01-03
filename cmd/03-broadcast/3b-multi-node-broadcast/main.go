package main

import (
	"sync"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	node *maelstrom.Node
	mu   sync.RWMutex

	store     map[int]bool
	neighbors []string
}

func (s *server) HandleBroadcast(msg maelstrom.Message) error {
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
	go s.Broadcast(msg.Src, body)

	return s.node.Reply(msg, map[string]any{
		"type": app.MessageBroadcastOk,
	})
}

func (s *server) Broadcast(src string, body any) {
	for _, dest := range s.neighbors {
		if dest == src {
			continue
		}
		s.node.Send(dest, body)
	}
}

func (s *server) HandleRead(msg maelstrom.Message) error {
	s.mu.RLock()

	keys := make([]int, 0, len(s.store))
	for k, _ := range s.store {
		keys = append(keys, k)
	}

	s.mu.RUnlock()

	return s.node.Reply(msg, map[string]any{
		"type":     app.MessageReadOk,
		"messages": keys,
	})
}

func (s *server) HandleTopology(msg maelstrom.Message) error {
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

func main() {
	n := maelstrom.NewNode()
	s := &server{
		node:  n,
		store: make(map[int]bool),
	}

	n.Handle(app.MessageBroadcast, s.HandleBroadcast)
	n.Handle(app.MessageRead, s.HandleRead)
	n.Handle(app.MessageTopology, s.HandleTopology)

	if err := n.Run(); err != nil {
		panic(err)
	}
}
