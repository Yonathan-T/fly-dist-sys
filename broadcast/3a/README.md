# Single-Node Broadcast

A solution to [Challenge #3a: Single-Node Broadcast](https://fly.io/dist-sys/3a/) from the Fly.io Distributed Systems series. Implementing an in-memory, thread-safe broadcast and message retrieval service on a single node.

## What it does

Maelstrom sends your node three types of RPC messages over stdin:

1. **`broadcast`**: Store an incoming message (an integer) and return `broadcast_ok`:
   ```json
   in:  {"type": "broadcast", "message": 1000, "msg_id": 1}
   out: {"type": "broadcast_ok", "msg_id": 1, "in_reply_to": 1}
   ```

2. **`read`**: Return a list of all messages seen and stored so far:
   ```json
   in:  {"type": "read", "msg_id": 2}
   out: {"type": "read_ok", "msg_id": 2, "in_reply_to": 2, "messages": [1000, 1001]}
   ```

3. **`topology`**: Provide the network topology (which neighbors can communicate):
   ```json
   in:  {"type": "topology", "topology": {"n0": []}, "msg_id": 3}
   out: {"type": "topology_ok", "msg_id": 3, "in_reply_to": 3}
   ```

---

## How It Works

Because Maelstrom dispatches concurrent requests across multiple parallel goroutines, all reads and writes to in-memory state must be thread-safe:

- **Encapsulated Store (`store` struct):** Bundles a `sync.Mutex` with a slice of numbers (`keys []int64`).
- **Thread-Safe Writes (`save`):** Locks the mutex before appending new numbers to the slice, preventing concurrent slice append corruption.
- **Thread-Safe Reads (`read`):** Locks the mutex and makes a defensive copy of the slice before returning it. This guarantees that JSON serialization never races against concurrent incoming broadcasts.
- **Single-Node Topology:** In Challenge 3a, the cluster runs on a single node (`n0`). The node acknowledges the `topology` RPC immediately with `topology_ok`, serving as the warmup before multi-node gossip in 3b.

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
java -Djava.awt.headless=true -jar ../../../maelstrom/lib/maelstrom.jar test -w broadcast --bin ./broadcast --node-count 1 --time-limit 20 --rate 10
```

### On Windows:

```powershell
make run-win
```

Or manually:

```powershell
go build -o broadcast.exe .
java -jar ..\..\..\maelstrom\lib\maelstrom.jar test -w broadcast --bin .\broadcast.exe --node-count 1 --time-limit 20 --rate 10
```

---

## Running on Windows vs. Linux

- **Windows:** Maelstrom creates a `store/current` symlink. On Windows, you need either **Developer Mode** enabled in Windows Settings or run your terminal as Administrator.
- **Linux / WSL:** Works smoothly out of the box. Running headless Java requires `-Djava.awt.headless=true` if no graphical display is attached.

---

## Results

```clojure
:perf {:latency-graph {:valid? true},
       :rate-graph {:valid? true},
       :valid? true},
:timeline {:valid? true},
:exceptions {:valid? true},
:stats {:valid? true,
        :count 200,
        :ok-count 200,
        :fail-count 0,
        :info-count 0,
        :by-f {:broadcast {:valid? true,
                           :count 106,
                           :ok-count 106,
                           :fail-count 0,
                           :info-count 0},
               :read {:valid? true,
                      :count 94,
                      :ok-count 94,
                      :fail-count 0,
                      :info-count 0}}},
:availability {:valid? true, :ok-fraction 1.0},
:net {:all {:send-count 404,
            :recv-count 404,
            :msg-count 404,
            :msgs-per-op 2.02},
      :clients {:send-count 404, :recv-count 404, :msg-count 404},
      :servers {:send-count 0,
                :recv-count 0,
                :msg-count 0,
                :msgs-per-op 0.0},
      :valid? true},
:workload {:worst-stale (),
           :duplicated-count 0,
           :valid? true,
           :lost-count 0,
           :lost (),
           :stable-count 106,
           :stale-count 0,
           :stale (),
           :never-read-count 0,
           :stable-latencies {0 0, 0.5 0, 0.95 0, 0.99 0, 1 0},
           :attempt-count 106,
           :never-read (),
           :duplicated {}},
:valid? true
```

```text
Everything looks good! ヽ(‘ー`)ノ
```

Passed with **200 operations**, **100% availability**, and **0 lost messages**.
