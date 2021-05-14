package connlib

func filterConnections(connections []Connection, filter func(Connection) bool) []Connection {
	filteredConnections := []Connection{}

	for _, conn := range connections {
		if filter(conn) {
			filteredConnections = append(filteredConnections, conn)
		}
	}

	return filteredConnections
}
