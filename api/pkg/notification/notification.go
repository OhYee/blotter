package notification

import (
	"sync"
)

type WritePackage struct {
	MessageType int
	MessageData []byte
}

type Type struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type mutexHub struct {
	mutex   *sync.Mutex
	channel map[string]chan WritePackage
}

// Hub of websocket
var Hub = mutexHub{
	mutex:   new(sync.Mutex),
	channel: make(map[string]chan WritePackage),
}

func (hub *mutexHub) Set(id string, channel chan WritePackage) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.channel[id] = channel
}

func (hub *mutexHub) Get(id string) chan WritePackage {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	return hub.channel[id]
}

func (hub *mutexHub) Remove(id string) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	delete(hub.channel, id)
}
