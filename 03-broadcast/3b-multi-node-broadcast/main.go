package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	node *maelstrom.Node
	mu   sync.RWMutex

	seen      map[int]bool
	neighbors []string
}

func (s *server) HandleBroadcast(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	toAdd := int(body["message"].(float64))

	s.mu.Lock()
	if _, ok := s.seen[toAdd]; ok {
		s.mu.Unlock()
		return nil
	}
	s.seen[toAdd] = true
	s.mu.Unlock()

	// broadcast this message to our neighbors
	go s.Broadcast(body)

	return s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

func (s *server) Broadcast(body any) {
	for _, dest := range s.neighbors {
		s.node.Send(dest, body)
	}
}

func (s *server) HandleRead(msg maelstrom.Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]int, 0, len(s.seen))
	for k, _ := range s.seen {
		keys = append(keys, k)
	}

	return s.node.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": keys,
	})
}

type TopologyReq struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func (s *server) HandleTopology(msg maelstrom.Message) error {
	var body TopologyReq
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.mu.Lock()
	s.neighbors = body.Topology[s.node.ID()]
	s.mu.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}

func main() {
	n := maelstrom.NewNode()
	s := &server{
		node: n,
		seen: make(map[int]bool),
	}

	n.Handle("broadcast_ok", nil)
	n.Handle("broadcast", s.HandleBroadcast)
	n.Handle("read", s.HandleRead)
	n.Handle("topology", s.HandleTopology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
