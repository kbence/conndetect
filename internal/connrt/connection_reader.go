package connrt

import (
	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/utils"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

type ConnectionReader struct {
	Node

	fileName    string
	connections *connlib.CategorizedConnections

	connectionSource connlib.ConnectionSource
	printer          utils.Printer
	time             utils.Time
}

func NewConnectionReader(eventManager event.ManagerFace, fileName string, connSrc ...connlib.ConnectionSource) (*ConnectionReader, error) {
	var err error

	reader := &ConnectionReader{
		Node:     Node{eventManager: eventManager},
		fileName: fileName,
		printer:  utils.NewPrinter(),
		time:     utils.NewTime(),
	}

	if len(connSrc) > 0 {
		reader.connectionSource = connSrc[0]
	} else {
		reader.connectionSource = connlib.NewConnectionSource()
	}

	if reader.connections, err = reader.readConnections(); err != nil {
		return nil, err
	}

	event.On("tick", event.ListenerFunc(reader.Handle))

	return reader, nil
}

func (r *ConnectionReader) readConnections() (*connlib.CategorizedConnections, error) {
	return r.connectionSource.ReadEstablishedTCPConnections(r.fileName)
}

func (r *ConnectionReader) Handle(event.Event) error {
	var connections *connlib.CategorizedConnections
	var err error

	if connections, err = r.readConnections(); err != nil {
		return err
	}

	oldConnectionsMap := r.connections.Established.ToMap()

	for _, conn := range connections.Established {
		if _, found := oldConnectionsMap[conn]; !found {
			dirConn := connlib.CalculateDirection(connections.Listening, conn)
			event.Fire("newConnection", event.M{"connection": dirConn})
		}
	}

	r.connections = connections

	return nil
}
