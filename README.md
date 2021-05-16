# conndetect

Small utility to track connections via `/proc/net/tcp`.

This is the first, naive implementation that considers the direction of connection to always be REMOTE -> LOCAL (which is obviously not true). I'll extend it later to use the listening ports to identify in what direction the socket was opened.

## Requirements

- [Go 1.12+](https://golang.org/doc/install)

## Installation

```shell
go get github.com/kbence/conndetect
```

or

```shell
git clone github.com:kbence/conndetect.git
cd connetect
go get .
```

## How to run the application?

If you have installed it with `go get` (given you have `$GOPATH/bin` on your `$PATH`), just execute:

```shell
conndetect
```

Or you can also run it directly using `go` (from the root directory):

```shell
go run .
```

## Tasks remaining

### Level 1

- [x] Currently the program assumes that all connections are incoming, let's use the listening sockets that we drop now to identify the incoming ones. Note: 0.0.0.0 is a wildcard, we have to only match those by ports.
- [x] Make it a bit more configurable (low hanging fruit, but it'll be useful for testing)
