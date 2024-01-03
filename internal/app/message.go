package app

import "encoding/json"

const (
	MessageEcho         = "echo"
	MessageGenerate     = "generate"
	MessageBroadcast    = "broadcast"
	MessageRead         = "read"
	MessageTopology     = "topology"
	MessageBroadcastAll = "broadcast_all"

	MessageEchoOk         = MessageEcho + "_ok"
	MessageGenerateOk     = MessageGenerate + "_ok"
	MessageBroadcastOk    = MessageBroadcast + "_ok"
	MessageReadOk         = MessageRead + "_ok"
	MessageTopologyOk     = MessageTopology + "_ok"
	MessageBroadcastAllOk = MessageBroadcastAll + "_ok"
)

type Message struct {
	ID string `json:"id,omitempty"`

	// echo fields
	MsgID int    `json:"msg_id,omitempty"`
	Echo  string `json:"echo,omitempty"`

	Message  int                 `json:"message,omitempty"`
	Topology map[string][]string `json:"topology,omitempty"`
	Type     string              `json:"type,omitempty"`

	// gossip fields
	Messages []int `json:"messages,omitempty"`
}

func (m *Message) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &m)
}

func Unmarshal(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
