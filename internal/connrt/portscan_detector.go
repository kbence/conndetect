package connrt

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/utils"
)

type ExpiringConnection struct {
	connlib.DirectionalConnection

	ExpiresAt time.Time
}

type ExpiringConnectionList []ExpiringConnection

func UniqueConnectionKey(c connlib.DirectionalConnection) string {
	return fmt.Sprintf("%s>%s", c.Source.IP.String(), c.Destination.IP.String())
}

type PortMap map[uint16]ExpiringConnection

func (m *PortMap) GetSortedPorts() []int {
	ports := []int{}

	for port := range *m {
		ports = append(ports, int(port))
	}

	sort.Ints(ports)
	return ports
}

type PortscanSettings struct {
	MaxPorts int
	Interval time.Duration
}

func NewPortscanSettings(maxPorts int, period time.Duration) *PortscanSettings {
	return &PortscanSettings{
		MaxPorts: maxPorts,
		Interval: period,
	}
}

type PortscanDetector struct {
	Node

	printer utils.Printer
	time    utils.Time

	settings      *PortscanSettings
	connectionLog map[string]PortMap
	scanReported  map[string]interface{}
}

func NewPortscanDetector(eventManager event.ManagerFace, settings *PortscanSettings) *PortscanDetector {
	detector := &PortscanDetector{
		Node: Node{eventManager: eventManager},

		printer: utils.NewPrinter(),
		time:    utils.NewTime(),

		settings:      settings,
		connectionLog: map[string]PortMap{},
		scanReported:  map[string]interface{}{},
	}

	eventManager.On("newConnection", event.ListenerFunc(detector.Handle))

	return detector
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

	d.saveConnection(connection)
	portMap := d.findPortMap(connection)

	if len(*portMap) >= d.settings.MaxPorts {
		d.reportScan(connection, portMap)
	}

	return nil
}

func (d *PortscanDetector) reportScan(connection *connlib.DirectionalConnection, portMap *PortMap) {
	connectionKey := UniqueConnectionKey(*connection)

	if _, found := d.scanReported[connectionKey]; found {
		return
	}

	ports := []string{}
	for _, port := range portMap.GetSortedPorts() {
		ports = append(ports, strconv.Itoa(port))
	}

	d.printer.Printf(
		"%s: Port scan detected: %s -> %s on ports %s\n",
		d.time.Now().Format(TIME_FORMAT),
		connection.Source.IP.String(),
		connection.Destination.IP.String(),
		strings.Join(ports, ","),
	)

	d.scanReported[connectionKey] = nil
}

func (d *PortscanDetector) findPortMap(connection *connlib.DirectionalConnection) *PortMap {
	if portMap, found := d.connectionLog[UniqueConnectionKey(*connection)]; found {
		return &portMap
	}

	return nil
}

func (d *PortscanDetector) saveConnection(connection *connlib.DirectionalConnection) {
	expiringConn := ExpiringConnection{
		DirectionalConnection: *connection,
		ExpiresAt:             d.time.Now().Add(d.settings.Interval),
	}
	connectionKey := UniqueConnectionKey(*connection)

	if _, found := d.connectionLog[connectionKey]; !found {
		d.connectionLog[connectionKey] = PortMap{}
	}

	// let's do some lazy cleanup
	d.cleanUpConnection(connectionKey)

	d.connectionLog[connectionKey][expiringConn.Destination.Port] = expiringConn
}

func (d *PortscanDetector) cleanUpConnection(connectionKey string) {
	newPortMap := PortMap{}

	for _, conn := range d.connectionLog[connectionKey] {
		if d.time.Now().Before(conn.ExpiresAt) {
			newPortMap[conn.Destination.Port] = conn
		}
	}

	d.connectionLog[connectionKey] = newPortMap
}
