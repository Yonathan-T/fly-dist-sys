# Unique-ID's

a solution to [Challenge #2: Unique ID Generation](https://fly.io/dist-sys/2/) from the Fly.io Distributed Systems series. Building a globally-unique ID generation service that works under network partitions.

## What it does

Maelstrom sends your node `generate` requests over stdin, and your node responds with a globally unique ID on stdout:

```
in:  {"type": "generate", "msg_id": 1}
out: {"type": "generate_ok", "msg_id": 1, "in_reply_to": 1, "id": "n0-1-1790892627"}
```

The challenge requires that IDs are **globally unique** across all nodes in the cluster, even when Maelstrom cuts off network communication between nodes (`--nemesis partition`).

### How Uniqueness is Guaranteed

Because nodes cannot talk to each other during a network partition, we cannot use centralized coordination (like a database or Raft leader). Instead, we generate IDs completely locally:

- **Node Identity (`n.ID()`):** Maelstrom assigns each node a unique ID (`n0`, `n1`, `n2`) during initialization. Prepending the node ID guarantees **zero collisions across different nodes**.
- **Mutex + Atomic Counter:** Each node protects an integer counter with a `sync.Mutex` that increments on every request. This guarantees **zero collisions across concurrent goroutines on the same node**.
- **Unix Timestamp:** Appends `time.Now().Unix()` for safety across process restarts.

Format: `<node_id>-<counter>-<timestamp>` (e.g. `n0-1-1790892627`).

Because zero cross-node communication is needed, the service provides **100% availability** even under heavy network partitions.

## Build & Run

Using the included `Makefile`:

### On Linux / Ubuntu WSL:

```bash
make run-linux
```

Or manually:

```bash
go build -o unique-ids .
../../maelstrom/maelstrom test -w unique-ids --bin ./unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```

### On Windows:

```powershell
make run-win
```

Or manually:

```powershell
go build -o unique-ids.exe .
java -jar ..\..\maelstrom\lib\maelstrom.jar test -w unique-ids --bin .\unique-ids.exe --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition
```

## Running on Windows vs. Linux

- **Windows:** Maelstrom creates a `store/current` symlink. On Windows, you need either **Developer Mode** enabled in Windows Settings or run your terminal as Administrator. Also, if `gnuplot` isn't installed on Windows, Maelstrom will say `Errors occurred during analysis, but no anomalies found` (meaning the test passed, but plotting failed).
- **Linux / WSL:** Works smoothly out of the box. Headless Java needs `-Djava.awt.headless=true` if no display is attached. With `gnuplot` installed, it generates the latency and rate charts.

## Results

```clojure
:perf {:latency-graph {:valid? true},
       :rate-graph {:valid? true},
       :valid? true},
:timeline {:valid? true},
:exceptions {:valid? true},
:stats {:valid? true,
        :count 18507,
        :ok-count 18507,
        :fail-count 0,
        :info-count 0,
        :by-f {:generate {:valid? true,
                          :count 18507,
                          :ok-count 18507,
                          :fail-count 0,
                          :info-count 0}}},
:availability {:valid? true, :ok-fraction 1.0},
:workload {:valid? true,
           :attempted-count 18507,
           :acknowledged-count 18507,
           :duplicated-count 0,
           :duplicated {},
           :range ["n0-1-1790892627" "n2-999-1790892634"]},
:valid? true
```

```text
Everything looks good! ヽ(‘ー`)ノ
```

Everything passed. **18,507 requests**, **100% availability**, and **0 duplicate IDs** across 3 nodes hammering 1,000 req/sec under continuous network partitions.
