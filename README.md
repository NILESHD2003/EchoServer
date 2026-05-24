# Echo Server

A TCP echo server written in Go to explore:

- TCP networking
- synchronous vs concurrent connection handling
- goroutines
- race conditions
- atomic operations
- Go's `net` package

The project starts two TCP servers simultaneously:

| Server Type | Default Port |
| --- | ---: |
| Synchronous TCP Server | `7379` |
| Concurrent TCP Server | `7380` |

Both servers accept TCP client messages and echo the same message back to the client.

---

# Project Structure

```text
echo-server/
├── config/
│   └── config.go
├── server/
│   ├── asynchronous_tcp_server.go
│   ├── synchronous_tcp_server.go
│   └── utils.go
├── go.mod
├── main.go
└── README.md
```

---

# Technical Overview

This project implements TCP servers directly using Go's standard `net` package without relying on higher-level HTTP abstractions.

The server lifecycle is:

```text
Client Connection
↓
Accept TCP Socket
↓
Read Incoming Bytes
↓
Process Message
↓
Write Response Back
```

The project contains two implementations:

## 1. Synchronous TCP Server

- Handles one accepted client connection path at a time
- Uses a single execution flow
- Easier to reason about
- Limited concurrency

## 2. Concurrent TCP Server

- Spawns one goroutine per client connection
- Handles multiple clients simultaneously
- Uses atomic operations for shared state synchronization

---

# Configuration

Configuration is defined in:

```text
config/config.go
```

```go
type Config struct {
	Host string
	Port int
}
```

The application supports the following CLI flags:

| Flag | Default | Description |
| --- | ---: | --- |
| `-sync-host` | `0.0.0.0` | Host for synchronous TCP server |
| `-sync-port` | `7379` | Port for synchronous TCP server |
| `-async-host` | `0.0.0.0` | Host for concurrent TCP server |
| `-async-port` | `7380` | Port for concurrent TCP server |

Example:

```sh
go run . \
  -sync-port 9001 \
  -async-port 9002
```

`0.0.0.0` means the server listens on all available network interfaces.

For local testing:

```text
localhost
```

---

# Application Entrypoint

The application starts from:

```text
main.go
```

Both TCP servers are started in separate goroutines:

```go
go server.RunSynchronousTCPServer(cfg_sync)
go server.RunAsynchronousTCPServer(cfg_async)

select {}
```

The final:

```go
select {}
```

blocks forever and prevents the main goroutine from exiting.

This keeps the process alive while both TCP servers continue running.

---

# TCP Read/Write Utilities

Common TCP read/write logic is implemented in:

```text
server/utils.go
```

---

## Reading Client Input

```go
func readIncomingCommand(c net.Conn) (string, error)
```

This function:

1. Allocates a `512` byte buffer
2. Reads bytes from the TCP socket
3. Converts bytes into a string
4. Returns the received message

---

## Responding to the Client

```go
func respondToClient(c net.Conn, cmd string) error
```

This function writes the same message back to the client.

Example:

Client sends:

```text
hello server
```

Server responds:

```text
hello server
```

---

# Synchronous TCP Server

Implemented in:

```text
server/synchronous_tcp_server.go
```

Function:

```go
RunSynchronousTCPServer(config config.Config)
```

---

## How It Works

The synchronous server:

1. Starts a TCP listener
2. Accepts a client connection
3. Handles the client in the same goroutine
4. Reads messages
5. Echoes responses
6. Waits for disconnect before handling next client

Default address:

```text
0.0.0.0:7379
```

Local testing:

```text
localhost:7379
```

---

## Synchronous Behavior

The synchronous server processes client connections sequentially.

Because connection handling occurs in a single execution flow:

- one slow client can delay others
- later clients wait until the current client disconnects

The connected client counter is:

```go
var concurrent_clients int = 0
```

This is safe because only one execution path updates the value.

---

# Concurrent TCP Server

Implemented in:

```text
server/asynchronous_tcp_server.go
```

Function:

```go
RunAsynchronousTCPServer(config config.Config)
```

---

## How It Works

The concurrent server:

1. Starts a TCP listener
2. Accepts incoming TCP connections continuously
3. Spawns a new goroutine per client
4. Handles each client independently

Each client connection is handled using:

```go
go handleClientConnection(c, &concurrent_clients)
```

Default address:

```text
0.0.0.0:7380
```

Local testing:

```text
localhost:7380
```

---

## Concurrent Behavior

Each connected client runs in its own goroutine.

This means:

- multiple clients can connect simultaneously
- one blocked client does not stop others
- connection handling becomes concurrent

Go runtime efficiently schedules many goroutines across a smaller number of OS threads.

---

# Race Condition Handling

The concurrent server contains shared mutable state:

```go
var concurrent_clients int64 = 0
```

This counter is accessed by multiple goroutines simultaneously.

Without synchronization, operations like:

```go
concurrent_clients += 1
```

are unsafe.

Internally this operation becomes:

```text
Read Current Value
↓
Modify Value
↓
Write Updated Value
```

Multiple goroutines performing this simultaneously can corrupt the final value.

This is called a:

```text
Race Condition
```

---

# Atomic Operations

The concurrent server uses Go's:

```go
sync/atomic
```

package for synchronization.

---

## Increment Counter

```go
atomic.AddInt64(&concurrent_clients, 1)
```

---

## Decrement Counter

```go
atomic.AddInt64(concurrent_clients, -1)
```

---

## Read Counter Safely

```go
atomic.LoadInt64(&concurrent_clients)
```

Atomic operations ensure that read-modify-write operations happen indivisibly without unsafe concurrent interleaving between goroutines.

---

# EOF Handling

When a TCP client disconnects gracefully:

```go
io.EOF
```

is returned from:

```go
c.Read(...)
```

This indicates:

```text
The remote side closed the TCP connection
```

The connection handler then exits and the goroutine terminates naturally.

---

# Running the Server

Run from project root:

```sh
go run .
```

This starts both TCP servers.

---

# Running with Race Detection

Go provides a built-in race detector.

Run:

```sh
go run -race .
```

This helps detect unsafe concurrent memory access.

---

# Testing with Netcat

## Test Synchronous Server

```sh
nc localhost 7379
```

Example:

```text
hello sync server
```

Expected response:

```text
hello sync server
```

---

## Test Concurrent Server

```sh
nc localhost 7380
```

Open multiple terminals to test concurrent client handling.

---

# Summary

This project demonstrates:

- low-level TCP networking in Go
- synchronous vs concurrent server architectures
- goroutines
- shared mutable state
- race conditions
- atomic synchronization
- TCP connection lifecycle handling

The concurrent server architecture is significantly better suited for handling multiple simultaneous TCP clients.
