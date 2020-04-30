package notification

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/OhYee/blotter/output"
)

type WritePackage struct {
	MessageType int
	MessageData []byte
}

type Type struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Socket struct {
	ID      string            `json:"id"`
	Channel chan WritePackage `json:"channel"`
}
type mutexHub struct {
	mutex   *sync.Mutex
	channel map[string][]Socket
}

func generateID(name, token string) string {
	return fmt.Sprintf("%s|%s", name, token)
}

// Hub of websocket
var Hub = mutexHub{
	mutex:   new(sync.Mutex),
	channel: make(map[string][]Socket),
}

func (hub *mutexHub) Set(name string, token string, channel chan WritePackage) string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	id := generateID(name, token)

	_, existed := hub.channel[id]
	if !existed {
		hub.channel[id] = make([]Socket, 0)
	}
	chanID := fmt.Sprintf("%d|%d", time.Now().Unix(), rand.Int31)
	hub.channel[id] = append(hub.channel[id], Socket{ID: chanID, Channel: channel})

	output.Debug("set %+v", hub)
	return chanID
}

func (hub *mutexHub) Get(name, token string) []Socket {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	output.Debug("get %+v", hub)

	id := generateID(name, token)
	return hub.channel[id]
}

func (hub *mutexHub) Remove(name, token string, chanID string) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	id := generateID(name, token)

	channels, existed := hub.channel[id]
	if !existed {
		return
	}
	pos := -1
	for idx, t := range channels {
		if t.ID == chanID {
			pos = idx
		}
	}
	if pos != -1 {
		channels = append(channels[:pos], channels[pos+1:]...)
	}
	if len(channels) == 0 {
		delete(hub.channel, id)
	} else {
		hub.channel[id] = channels
	}

	delete(hub.channel, token)
	output.Debug("remove %+v", hub)

}
