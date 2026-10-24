package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type store struct {
	mu   sync.Mutex
	keys []int64
}

func NewStore() *store {
	return &store{
		keys: make([]int64, 0),
	}
}

func (s *store) save(msg int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.keys = append(s.keys, msg)
}

func (s *store) read() []int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	msgs := make([]int64, len(s.keys))
	copy(msgs, s.keys)
	return msgs
}

func main() {
	n := maelstrom.NewNode()
	s := NewStore()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body struct {
			Type    string `json:"type"`
			Message int64  `json:"message"`
		}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		s.save(body.Message)

		return n.Reply(msg, map[string]any{
			"type": "broadcast_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": s.read(),
		})
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body struct {
			Type     string              `json:"type"`
			Topology map[string][]string `json:"topology"`
		}
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
