package connrt

import (
	"strconv"
	"strings"

	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/utils"
)

type TimedConnection struct {
	connlib.DirectionalConnection

	Time utils.Time
}

type PortscanDetector struct {
	Node

	printer utils.Printer
	time    utils.Time

	connectionLog []connlib.DirectionalConnection
}

func NewPortscanDetector(eventManager event.ManagerFace) *PortscanDetector {
	return &PortscanDetector{
		Node: Node{eventManager: eventManager},

		printer: utils.NewPrinter(),
		time:    utils.NewTime(),
	}
}

func (d *PortscanDetector) Handle(e event.Event) error {
	var connection *connlib.DirectionalConnection = nil

	if connObj := e.Get("connection"); connObj != nil {
		switch conn := connObj.(type) {
		case connlib.DirectionalConnection:
			connection = &conn
		}
	}

	// Do nothing
	if connection == nil {
		return nil
	}

	d.connectionLog = append(d.connectionLog, *connection)

	if len(d.connectionLog) >= 3 {
		ports := []string{}

		for _, conn := range d.connectionLog {
			ports = append(ports, strconv.Itoa(int(conn.Destination.Port)))
		}

		d.printer.Printf(
			"%s: Port scan detected: %s -> %s on ports %s",
			d.time.Now().Format(TIME_FORMAT),
			connection.Source.IP,
			connection.Destination.IP,
			strings.Join(ports, ","),
		)
	}

	return nil
}
