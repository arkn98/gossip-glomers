package main

import (
	"sync"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	node *maelstrom.Node
	mu   sync.RWMutex

	store []int
}

func (s *server) HandleBroadcast(msg maelstrom.Message) error {
	body, err := app.Unmarshal(msg.Body)
	if err != nil {
		return err
	}

	toAdd := body.Message

	s.mu.Lock()
	s.store = append(s.store, toAdd)
	s.mu.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type": app.MessageBroadcastOk,
	})
}

func (s *server) HandleRead(msg maelstrom.Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.node.Reply(msg, map[string]any{
		"type":     app.MessageReadOk,
		"messages": s.store,
	})
}

func (s *server) HandleTopology(msg maelstrom.Message) error {
	// do nothing for now
	return s.node.Reply(msg, map[string]any{
		"type": app.MessageTopologyOk,
	})
}

func main() {
	n := maelstrom.NewNode()
	s := &server{node: n}

	n.Handle(app.MessageBroadcast, s.HandleBroadcast)
	n.Handle(app.MessageRead, s.HandleRead)
	n.Handle(app.MessageTopology, s.HandleTopology)

	if err := n.Run(); err != nil {
		panic(err)
	}
}
