package connrt

import (
	"time"

	"github.com/gookit/event"
)

type Ticker struct {
	Node

	wait int
	exit chan interface{}
}

func NewTicker(eventManager *event.Manager, wait int) *Ticker {
	return &Ticker{
		Node: Node{eventManager: eventManager},
		wait: wait,
		exit: make(chan interface{}),
	}
}

func (t *Ticker) Run() {
	ticker := time.Tick(time.Duration(t.wait) * time.Second)

tickLoop:
	for {
		select {
		case tm := <-ticker:
			event.Fire("tick", event.M{"time": tm})
			break
		case <-t.exit:
			event.Fire("exit", event.M{})
			break tickLoop
		}
	}
}

func (t *Ticker) Stop() {
	t.exit <- nil
}
