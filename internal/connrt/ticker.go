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

func NewTicker(eventManager event.ManagerFace, wait int) *Ticker {
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
			t.eventManager.Fire(eventTick, event.M{"time": tm})
			break

		case <-t.exit:
			t.eventManager.Fire(eventExit, event.M{})
			break tickLoop
		}
	}
}

func (t *Ticker) Stop() {
	t.exit <- nil
}
