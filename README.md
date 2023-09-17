Multi-threaded TCP Server
===

This codebase demonstartes multi-threaded tcp server in Golang.
For detailed explanation please refer to the following video

- https://youtu.be/f9gUFy-9uCM

## How to Run

```
$ go run main.go
```

Fire the following commands on another terminal to simulate
multiple concurrent requests.

```
$ curl http://localhost:1729 &
$ curl http://localhost:1729 &
$ curl http://localhost:1729 &
```

## Improvements suggested by arpit at the end of video

1. Limiting the number of threads
2. Add threadpool to save on thread creation time
3. Connection timeout
4. Tcp backlog queue configuration
