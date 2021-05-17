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

	if _, _, err = reader.ReadLine(); err != nil {
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
func ReadEstablishedTCPConnections(fileName string) (*CategorizedConnections, error) {
	var file *os.File
	var err error

	if file, err = os.Open(fileName); err != nil {
		return nil, err
	}
	defer file.Close()

	connections, err := ParseTCPFile(file)
	catConnections := &CategorizedConnections{}

	for _, conn := range connections {
		if conn.Remote.IsUnbound() {
			catConnections.Listening = append(catConnections.Listening, conn)
		} else {
			catConnections.Established = append(catConnections.Established, conn)
		}
	}

	return catConnections, nil
}

func CalculateDirection(listeners ConnectionList, conn Connection) DirectionalConnection {
	// Let's see if we have an exact hit first
	for _, l := range listeners {
		if l.Local == conn.Remote {
			return DirectionalConnection{Source: conn.Local, Destination: conn.Remote}
		}

		if l.Local == conn.Local {
			return DirectionalConnection{Source: conn.Remote, Destination: conn.Local}
		}
	}

	// Check if there's a listener on 0.0.0.0 for the same port
	// Note: this is totally not right, since we don't have any
	// information about what IP addresses this specific node has.

	// TODO: parse it out from /proc/net/fib_trie and use that
	// list instead of IP.IsUnspecified()!
	for _, l := range listeners {
		if l.Local.IP.IsUnspecified() && l.Local.Port == conn.Remote.Port {
			return DirectionalConnection{Source: conn.Local, Destination: conn.Remote}
		}

		if l.Local.IP.IsUnspecified() && l.Local.Port == conn.Local.Port {
			return DirectionalConnection{Source: conn.Remote, Destination: conn.Local}
		}
	}

	// Assuming outgoing connection if the local port is ephemeral
	if isEphemeralPort(conn.Local.Port) {
		return DirectionalConnection{Source: conn.Local, Destination: conn.Remote}
	}

	// Assuming incoming connection if nothing else matches
	// TODO: return an error and print a warning about it
	return DirectionalConnection{Source: conn.Remote, Destination: conn.Local}
}
