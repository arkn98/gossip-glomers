package main

import (
	"fmt"
	"time"

	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle(app.MessageGenerate, func(msg maelstrom.Message) error {
		body, err := app.Unmarshal(msg.Body)
		if err != nil {
			return err
		}

		body.Type = app.MessageGenerateOk
		body.ID = fmt.Sprintf("%s-%d", n.ID(), time.Now().UnixNano())

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		panic(err)
	}
}
