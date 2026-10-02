package main

import (
	"fmt"
	"sync"

	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type uniqueID struct {
	mu   sync.Mutex
	time time.Time
	ID   int
}

func NewUniqueID() *uniqueID {
	return &uniqueID{
		ID: 0,
	}
}

func (uid *uniqueID) generate(MaelstromID string) string {
	uid.mu.Lock()
	defer uid.mu.Unlock()
	uid.ID++
	timestamp := time.Now().Unix()
	idString := fmt.Sprintf("%s-%d-%d", MaelstromID, uid.ID, timestamp)
	return idString
}

func main() {
	n := maelstrom.NewNode()
	uid := NewUniqueID()

	n.Handle("generate", func(msg maelstrom.Message) error {
		mID := n.ID()
		id := uid.generate(mID)
		res := map[string]any{
			"type": "generate_ok",
			"id":   id,
		}
		return n.Reply(msg, res)
	})
	n.Run()
}
