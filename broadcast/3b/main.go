package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type store struct {
	mu        sync.Mutex
	keys      map[int64]bool
	neighbors []string
}

func NewStore() *store {
	return &store{
		keys: make(map[int64]bool),
	}
}

func (s *store) save(msg int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.keys[msg] {
		return false
	}
	s.keys[msg] = true
	return true
}

func (s *store) read() []int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]int64, 0, len(s.keys))
	for k := range s.keys {
		res = append(res, k)
	}
	return res
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
		isNew := s.save(body.Message)

		if isNew {
			for _, neighbor := range s.neighbors {
				if neighbor == msg.Src {
					continue
				}
				n.Send(neighbor, map[string]any{
					"type":    "broadcast",
					"message": body.Message,
				})
			}
		}
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
		s.neighbors = body.Topology[n.ID()]

		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	n.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		return nil
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
