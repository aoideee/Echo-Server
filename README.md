# Improved TCP Echo Server

A concurrent, feature‑rich TCP echo server written in Go.

---

## Video

[Improved Echo Server](https://youtu.be/esPqCw6msZg)

---

## Requirements

- Go 1.18 or newer

---

## Building

```bash
git clone https://github.com/aoideee/Echo‑Server.git
cd tcp‑echo‑server
go build -o echo-server main.go
```

---

## Running
By default the server listens on port 4000:

```bash
./echo-server
```
To choose a different port, use the --port flag:

```bash
./echo-server --port 5000
```

## Most Educationally Enriching Feature
Implementing the 30 second inactivity timeout using ```conn.SetReadDeadline``` and learning to detect and handle ```net.Error``` timeouts without crashing the server deepened my understanding of robust network I/O in Go.

---

## Most Research‑Intensive Feature
Setting up per‑client log files (```<IP>.log```) with concurrent writes taught me about safe file handling in goroutines and how to use ```log.New``` to create separate logger instances for each connection.
