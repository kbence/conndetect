package connlib

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var fieldSeparator = regexp.MustCompile("\\s+")

// ParseEndpoint - Parses and hexadecimal IPv4 endpoint (eg.: "0100007F:1F90")
func ParseEndpoint(endpoint string) Endpoint {
	parts := strings.Split(endpoint, ":")

	ip, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		panic(err)
	}

	port, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		panic(err)
	}

	return Endpoint{
		IP:   IPv4Address{byte(ip & 0xff), byte(ip >> 8), byte(ip >> 16), byte(ip >> 24)},
		Port: uint16(port),
	}
}

// ParseTCPFile - Parses a file that's in the format of /proc/net/tcp
func ParseTCPFile(r *bufio.Reader) ConnectionList {
	var err error
	var line []byte

	_, _, err = r.ReadLine()
	if err != nil {
		panic(err)
	}

	connections := ConnectionList{}

	for err == nil {
		line, _, err = r.ReadLine()

		switch err {
		case nil:
			lineStr := string(line)
			fields := fieldSeparator.Split(lineStr, 15)
			src := ParseEndpoint(fields[2])
			dst := ParseEndpoint(fields[3])

			connections = append(connections, Connection{Local: src, Remote: dst})
			break

		case io.EOF:
			break

		default:
			panic(err)
		}
	}

	return connections
}

// ReadEstabilishedTCPConnections - Reads /proc/net/tcp and returns
// the slice of parsed connections
func ReadEstabilishedTCPConnections() ConnectionList {
	file, err := os.Open("/proc/net/tcp")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return filterConnections(
		ParseTCPFile(bufio.NewReader(file)),
		func(c Connection) bool {
			return !c.Remote.IsUnbound()
		},
	)
}
