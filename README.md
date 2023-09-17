Multi-threaded TCP Server w/ Improvements
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
$ k6 run simple-test.js
```

## Improvements suggested by arpit at the end of video

- [x] Limiting the number of threads
- [ ] Add threadpool to save on thread creation time
- [ ] Connection timeout
- [ ] Tcp backlog queue configuration
