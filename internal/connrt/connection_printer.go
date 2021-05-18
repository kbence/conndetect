package connrt

import (
	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/utils"
)

type ConnectionPrinter struct {
	Node

	printer utils.Printer
	time    utils.Time
}

func NewConnectionPrinter(eventManager *event.Manager) *ConnectionPrinter {
	printer := &ConnectionPrinter{
		Node: Node{eventManager: eventManager},

		printer: utils.NewPrinter(),
		time:    utils.NewTime(),
	}

	event.On("newConnection", event.ListenerFunc(printer.Handle))

	return printer
}

func (p *ConnectionPrinter) Handle(e event.Event) error {
	var connection *connlib.DirectionalConnection = nil

	if connObj := e.Get("connection"); connObj != nil {
		switch conn := connObj.(type) {
		case connlib.DirectionalConnection:
			connection = &conn
		}
	}

	// Swallow the error now
	// TODO: handle this error more gracefully, by eg. logging it
	// or sending an error event! :o
	if connection == nil {
		return nil
	}

	p.printer.Printf(
		"%s: New connection: %s -> %s\n",
		p.time.Now().Format(TIME_FORMAT),
		connection.Source.String(),
		connection.Destination.String(),
	)

	return nil
}
