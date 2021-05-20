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
cd conndetect
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

### Running in Docker

First, build the docker image:

```shell
docker build -t conndetect .
```

Then run the container with the following command. Make sure to pass `--network host` to get access to all the network activity.

```shell
docker run -it --rm --network host conndetect
```
