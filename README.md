## Build your own Redis

Programming exercise to brush up Golang skills and get basic understanding of Redis internals.

This follows [Build your own Redis](https://build-your-own.org/redis) guide religiously with few caveats.

- Uses `Go` instead of `C`
- Uses `net` package for socket programming instead of raw `syscalls`
- Uses `goroutines` instead of `epoll` for implementing concurrency

To run `server` & `client` use

```
> go run build/byor-server/main.go -port=:6060

> go run build/byor-client/main.go -port=:6060
```
