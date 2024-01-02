package main

import (
	"gossip-glomers/internal/app"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle(app.MessageEcho, func(msg maelstrom.Message) error {
		body, err := app.Unmarshal(msg.Body)
		if err != nil {
			return err
		}

		body.Type = app.MessageEchoOk

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		panic(err)
	}
}
