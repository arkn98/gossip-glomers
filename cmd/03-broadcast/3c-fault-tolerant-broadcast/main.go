package main

import (
	"sync"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	Retries           = 100
	BackoffMultiplier = 100
)

type Server struct {
	node *maelstrom.Node
	mu   sync.RWMutex

	store     map[int]bool
	neighbors []string
}

func main() {
	s := &Server{
		node:  maelstrom.NewNode(),
		store: make(map[int]bool),
	}

	s.node.Handle(app.MessageRead, s.HandleRead)
	s.node.Handle(app.MessageBroadcast, s.HandleBroadcast)
	s.node.Handle(app.MessageTopology, s.HandleTopology)

	if err := s.node.Run(); err != nil {
		panic(err)
	}
}
