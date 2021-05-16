package connlib

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
)

func filterConnections(connections []Connection, filter func(Connection) bool) []Connection {
	filteredConnections := []Connection{}

	for _, conn := range connections {
		if filter(conn) {
			filteredConnections = append(filteredConnections, conn)
		}
	}

	return filteredConnections
}

var ephemeralPortRange []uint16 = nil

func isEphemeralPort(port uint16) bool {
	if ephemeralPortRange == nil {
		// set default values to Linux defaults
		ephemeralPortRange = []uint16{32768, 60999}

		// read local port range (just stop on an error for now)
		if file, err := os.Open("/proc/sys/net/ipv4/ip_local_port_range"); err == nil {
			if line, _, err := bufio.NewReader(file).ReadLine(); err == nil {
				re := regexp.MustCompile("\\s+")
				ports := re.Split(string(line), 2)

				if port, err := strconv.Atoi(ports[0]); err != nil {
					ephemeralPortRange[0] = uint16(port)
				}

				if port, err := strconv.Atoi(ports[1]); err != nil {
					ephemeralPortRange[1] = uint16(port)
				}
			}
		}
	}

	return port >= ephemeralPortRange[0] && port <= ephemeralPortRange[1]
}
