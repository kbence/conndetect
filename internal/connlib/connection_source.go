package connlib

import "os"

type ConnectionSource interface {
	ReadEstablishedTCPConnections(fileName string) (*CategorizedConnections, error)
}

type ConnectionSourceImpl struct{}

func NewConnectionSource() ConnectionSource {
	return &ConnectionSourceImpl{}
}

// ReadEstabilishedTCPConnections - Reads /proc/net/tcp and returns
// the slice of parsed connections
func (s *ConnectionSourceImpl) ReadEstablishedTCPConnections(fileName string) (*CategorizedConnections, error) {
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
