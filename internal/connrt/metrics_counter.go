package connrt

import (
	"github.com/gookit/event"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics
var metricConnectionsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "conndetect_connections_total",
	Help: "Total number of connections detected",
})

type MetricsCounter struct {
	Node
}

func NewMetricsCounter(eventManager event.ManagerFace) *MetricsCounter {
	counter := &MetricsCounter{
		Node: Node{eventManager: eventManager},
	}

	eventManager.On("newConnection", event.ListenerFunc(counter.Handle))

	return counter
}

func (c *MetricsCounter) Handle(e event.Event) error {
	switch e.Name() {
	case "newConnection":
		metricConnectionsTotal.Inc()
	}
	return nil
}

func init() {
	prometheus.MustRegister(metricConnectionsTotal)
}
