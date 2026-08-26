# maelstrom-echo

a solution to [Challenge #1: Echo](https://fly.io/dist-sys/1/) from the Fly.io Distributed Systems series. Just messing around with Maelstrom and learning how distributed systems work.

## What it does

Maelstrom sends your program JSON messages over stdin, and your program responds on stdout. The echo workload just wants you to send back what you received:

```
in:  {"type": "echo", "msg_id": 1, "echo": "Please echo 35"}
out: {"type": "echo_ok", "msg_id": 1, "in_reply_to": 1, "echo": "Please echo 35"}
```

Maelstrom then spins this up, hits it with network faults, and checks everything still works correctly.

## Build & Run

```bash
cd maelstrom-echo
go build -o maelstrom-echo .
```

Then from wherever you extracted Maelstrom:

```bash
./maelstrom test -w echo --bin /path/to/maelstrom-echo --node-count 1 --time-limit 10
```

## Running on Windows

it does Work fine, just a couple things you should look out for,

- Build with `go build -o maelstrom-echo.exe .`
- The `./maelstrom` bash script won't work on Windows, but you can run the JAR directly instead:
  ```
  java -jar <maelstrom-path>\lib\maelstrom.jar test -w echo --bin ./maelstrom-echo.exe --node-count 1 --time-limit 10
  ```
- Maelstrom creates symlinks which need Developer Mode on, or you can run your terminal as admin and it would work just fine.

I'd recommend running this on Linux tho, it saves you all that hassle. 

## Results

```
:timeline {:valid? true}
:stats    {:valid? true, :count 51, :ok-count 51, :fail-count 0}
:availability {:valid? true, :ok-fraction 1.0}
:workload {:valid? true, :errors ()}
```

Everything passed. 51 requests, all ok, zero failures.

## Running the server

When you run Maelstrom it will spin up your binary, send it an `init` message to register the node, then start hammering it with `echo` messages for the time limit you set. You can watch the whole thing in real time as it logs each request and response.

Here's what it looks like when its running:

![Maelstrom echo test output](screenshot/Screenshot%20(245).png)

You can see Maelstrom sending echo requests like "Please echo 35" and getting back `echo_ok` responses. After the time limit it tears everything down and runs its analysis to check if your node handled everything correctly.
