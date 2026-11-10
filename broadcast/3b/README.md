# Multi-Node Broadcast

A solution to [Challenge #3b: Multi-Node Broadcast](https://fly.io/dist-sys/3b/) from the Fly.io Distributed Systems series. Implementing a decentralized gossip protocol across a 5-node cluster to propagate messages and achieve eventual consistency.

---

## What It Does

In Challenge #3b, Maelstrom runs **5 independent node processes** simultaneously (`n0`, `n1`, `n2`, `n3`, `n4`).

Clients send `broadcast` and `read` requests to random nodes across the cluster. When a node receives a write from a client, it must propagate (gossip) that value to other nodes so that any node can answer a subsequent `read` request with the complete set of all cluster messages:

```
Client A ----( broadcast: 1 )----> [ Node n0 ]
                                         |
                                (gossip broadcast: 1)
                                         v
Client B <---( read: [1, 2] )----- [ Node n3 ] <---( gossip: 2 )--- [ Node n1 ]
```

---

## The Distributed Systems Challenges

Moving from single-node (3a) to multi-node (3b) introduces three classic distributed systems hurdles:

1. **Topology Awareness:**
   Each node receives a `topology` RPC mapping the network. A node extracts its designated neighbors (`body.Topology[n.ID()]`) and stores them to know where to forward messages.

2. **Preventing the Infinite Gossip Storm:**
   If Node A sends a message to Node B, and Node B forwards it back to Node A, the cluster enters an infinite ping-pong loop that overwhelms the network.
   - **Set Deduplication (`map[int64]bool`):** An atomic check-and-insert method (`s.save(msg) bool`) checks if the message is brand new.
   - **Gossip Guard:** Nodes **only** forward messages if `isNew == true`.
   - **Sender Suppression:** Nodes skip forwarding back to the node that originally sent them the message (`if neighbor == msg.Src { continue }`).

3. **Inter-Node Peer Protocol (`broadcast_ok`):**
   In 3b, nodes communicate with each other using `n.Send(neighbor, ...)`. When a neighbor receives a broadcast, it replies with `{"type": "broadcast_ok"}`. The initiating node must register a handler for `"broadcast_ok"` to silently acknowledge receipt without crashing.

---

## How It Works

- **Thread-Safe Set Storage (`store` struct):**
  Uses `map[int64]bool` protected by a `sync.Mutex` for $O(1)$ lookups and deduplication.
- **Atomic State Transitions:**
  `save(msg int64) bool` locks the mutex, checks if the key already exists, marks it seen, and unlocks—returning whether the value was newly inserted.
- **Gossip Loop (`n.Send`):**
  When a write is confirmed new, the node loops through its configured `neighbors` and dispatches `broadcast` messages asynchronously without blocking the client reply.
- **Safe Read Clones:**
  `read()` locks the store and iterates through the map keys to return a clean `[]int64` slice, guaranteeing that concurrent broadcasts cannot trigger read/write data races during JSON serialization.

---

## Build & Run

Using the included `Makefile`:

### On Linux / Ubuntu WSL:

```bash
make run-linux
```

Or manually:

```bash
go build -o broadcast .
java -Djava.awt.headless=true -jar ../../../maelstrom/lib/maelstrom.jar test -w broadcast --bin ./broadcast --node-count 5 --time-limit 20 --rate 10
```

### On Windows:

```powershell
make run-win
```

Or manually:

```powershell
go build -o broadcast.exe .
java -jar ..\..\..\maelstrom\lib\maelstrom.jar test -w broadcast --bin .\broadcast.exe --node-count 5 --time-limit 20 --rate 10
```

---

## Results

```clojure
:perf {:latency-graph {:valid? true},
       :rate-graph {:valid? true},
       :valid? true},
:timeline {:valid? true},
:exceptions {:valid? true},
:stats {:valid? true,
        :count 183,
        :ok-count 183,
        :fail-count 0,
        :info-count 0,
        :by-f {:broadcast {:valid? true,
                           :count 85,
                           :ok-count 85,
                           :fail-count 0,
                           :info-count 0},
               :read {:valid? true,
                      :count 98,
                      :ok-count 98,
                      :fail-count 0,
                      :info-count 0}}},
:availability {:valid? true, :ok-fraction 1.0},
:net {:all {:send-count 1406,
            :recv-count 1406,
            :msg-count 1406,
            :msgs-per-op 7.68306},
      :clients {:send-count 386, :recv-count 386, :msg-count 386},
      :servers {:send-count 1020,
                :recv-count 1020,
                :msg-count 1020,
                :msgs-per-op 5.5737705},
      :valid? true},
:workload {:worst-stale (),
           :duplicated-count 0,
           :valid? true,
           :lost-count 0,
           :lost (),
           :stable-count 85,
           :stale-count 0,
           :stale (),
           :never-read-count 0,
           :stable-latencies {0 0, 0.5 0, 0.95 0, 0.99 0, 1 0},
           :attempt-count 85,
           :never-read (),
           :duplicated {}},
:valid? true
```

```text
Everything looks good! ヽ(‘ー`)ノ
```

Passed across **5 nodes** with **100% availability**, **0 lost messages**, and full cross-cluster convergence.
