package main

import (
	"sync"
	"time"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	Retries           = 10
	BackoffMultiplier = 5000 // ms
)

type Server struct {
	node *maelstrom.Node
	mu   sync.RWMutex

	store     map[int]bool
	neighbors []string

	ticker *time.Ticker
	done   chan bool
}

func main() {
	s := &Server{
		node:   maelstrom.NewNode(),
		store:  make(map[int]bool),
		ticker: time.NewTicker(1800 * time.Millisecond),
		done:   make(chan bool),
	}

	s.StartBroadcastLoop()
	defer s.StopBroadcastLoop()

	s.node.Handle(app.MessageRead, s.HandleRead)
	s.node.Handle(app.MessageBroadcast, s.HandleBroadcast)
	s.node.Handle(app.MessageTopology, s.HandleTopology)
	s.node.Handle(app.MessageBroadcastAll, s.HandleBroadcastAll)

	if err := s.node.Run(); err != nil {
		panic(err)
	}
}
