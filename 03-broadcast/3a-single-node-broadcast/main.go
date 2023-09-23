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
	seen []int
}

func (s *server) HandleBroadcast(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.mu.Lock()
	s.seen = append(s.seen, int(body["message"].(float64)))
	s.mu.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

func (s *server) HandleRead(msg maelstrom.Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.node.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": s.seen,
	})
}

func (s *server) HandleTopology(msg maelstrom.Message) error {
	// do nothing for now
	return s.node.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}

func main() {
	n := maelstrom.NewNode()
	s := &server{node: n}

	n.Handle("broadcast", s.HandleBroadcast)
	n.Handle("read", s.HandleRead)
	n.Handle("topology", s.HandleTopology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
