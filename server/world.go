package server

import (
	"github.com/Nightgunner5/gogame/engine/actor"
	"github.com/Nightgunner5/gogame/engine/message"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"sync"
)

var (
	SendToAll = make(chan *packet.Packet)
)

type World struct {
	actor.Holder

	onConnect chan<- chan<- *packet.Packet

	idToActor map[uint64]*actor.Actor
	location  map[*actor.Actor]layout.Coord

	doors   map[layout.Coord]*actor.Actor
	doorMtx sync.Mutex
}

func (w *World) Initialize() (message.Receiver, func(message.Message)) {
	msgIn, broadcast := w.Holder.Initialize()

	onConnect := make(chan chan<- *packet.Packet)
	w.onConnect = onConnect

	w.idToActor = make(map[uint64]*actor.Actor)
	w.location = make(map[*actor.Actor]layout.Coord)

	w.doors = make(map[layout.Coord]*actor.Actor)

	messages := make(chan message.Message)

	go w.dispatch(msgIn, messages, broadcast, onConnect)

	return messages, broadcast
}

func (w *World) dispatch(msgIn message.Receiver, messages message.Sender, broadcast func(message.Message), onConnect <-chan chan<- *packet.Packet) {
	for {
		select {
		case msg, ok := <-msgIn:
			if !ok {
				close(messages)
				return
			}

			switch m := msg.(type) {
			case SetLocation:
				w.idToActor[m.ID] = m.Actor
				w.location[m.Actor] = m.Coord

				SendToAll <- &packet.Packet{
					Location: &packet.Location{
						ID:    m.ID,
						Coord: m.Coord,
					},
				}

				go broadcast(m)

			case packet.Despawn:
				a := w.idToActor[m.ID]
				delete(w.idToActor, m.ID)
				delete(w.location, a)
				SendToAll <- &packet.Packet{
					Despawn: &m,
				}

			case AddDoor:
				w.doorMtx.Lock()
				w.doors[m.Coord] = m.Actor
				w.doorMtx.Unlock()

			default:
				messages <- m
			}

		case c := <-onConnect:
			w.EachHeld(func(a *actor.Actor) {
				go sendSendLocation(a, SendLocation(c))
			})
			c <- &packet.Packet{MapOverride: &packet.MapOverride{layout.GetChanges()}}
		}
	}
}

func (w *World) OpenDoor(opener *Player, coord layout.Coord) {
	w.doorMtx.Lock()
	door := w.doors[coord]
	w.doorMtx.Unlock()

	if door != nil {
		door.Send <- OpenDoor{opener}
	}
}

func sendSendLocation(a *actor.Actor, c SendLocation) {
	// This function will only panic if a player disconnects between another
	// player joining and the location sender being recieved.
	defer recover()

	a.Send <- c
}

func init() {
	layout.OnChange = func(c layout.Coord, t layout.MultiTile) {
		SendToAll <- &packet.Packet{MapChange: &packet.MapChange{c, t}}
	}
}

var world = NewWorld()

func NewWorld() (world *World) {
	world = new(World)
	actor.Init("world", &world.Actor, world)
	return
}
