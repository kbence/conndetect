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
func ParseEndpoint(endpoint string) (*Endpoint, error) {
	parts := strings.Split(endpoint, ":")

	ip, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		IP:   IPv4Address{byte(ip & 0xff), byte(ip >> 8), byte(ip >> 16), byte(ip >> 24)},
		Port: uint16(port),
	}, nil
}

// ParseTCPFile - Parses a file that's in the format of /proc/net/tcp
func ParseTCPFile(r io.Reader) (ConnectionList, error) {
	var err error
	var line []byte

	reader := bufio.NewReader(r)

	_, _, err = reader.ReadLine()
	if err != nil {
		return nil, err
	}

	connections := ConnectionList{}

	for err == nil {
		line, _, err = reader.ReadLine()

		switch err {
		case nil:
			lineStr := string(line)
			entryStr := strings.SplitN(lineStr, ":", 2)[1]
			fields := fieldSeparator.Split(entryStr, 15)

			var src, dst *Endpoint

			if src, err = ParseEndpoint(fields[1]); err != nil {
				return nil, err
			}

			if dst, err = ParseEndpoint(fields[2]); err != nil {
				return nil, err
			}

			connections = append(connections, Connection{Local: *src, Remote: *dst})
			break

		case io.EOF:
			break

		default:
			return nil, err
		}
	}

	return connections, nil
}

// ReadEstabilishedTCPConnections - Reads /proc/net/tcp and returns
// the slice of parsed connections
func ReadEstabilishedTCPConnections() (ConnectionList, error) {
	var file *os.File
	var err error

	if file, err = os.Open("/proc/net/tcp"); err != nil {
		return nil, err
	}
	defer file.Close()

	connections, err := ParseTCPFile(file)

	return filterConnections(
		connections,
		func(c Connection) bool {
			return !c.Remote.IsUnbound()
		},
	), nil
}
