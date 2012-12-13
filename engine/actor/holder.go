package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"sync"
)

var (
	MsgAddHeld    = message.NewKind("AddHeld")
	MsgRemoveHeld = message.NewKind("RemoveHeld")
)

type AddHeld struct {
	*Actor
}

// Returns MsgAddHeld.
func (AddHeld) Kind() message.Kind {
	return MsgAddHeld
}

type RemoveHeld struct {
	*Actor
}

// Returns MsgRemoveHeld.
func (RemoveHeld) Kind() message.Kind {
	return MsgRemoveHeld
}

type Holder struct {
	Actor
	held map[*Actor]bool
	lock sync.RWMutex
}

func (h *Holder) Initialize() (messages message.Receiver, broadcast func(message.Message)) {
	msgIn, broadcast := h.Actor.Initialize()

	messages_ := make(chan message.Message)
	messages = messages_

	h.held = make(map[*Actor]bool)

	go h.dispatch(msgIn, messages_, broadcast)

	return
}

func (h *Holder) dispatch(msgIn message.Receiver, messages message.Sender, broadcast func(message.Message)) {
	for {
		select {
		case msg, ok := <-msgIn:
			if !ok {
				close(messages)
				return
			}
			switch m := msg.(type) {
			case AddHeld:
				h.lock.Lock()
				if !h.held[m.Actor] {
					h.held[m.Actor] = true
					go broadcast(m)
				}
				h.lock.Unlock()

			case RemoveHeld:
				h.lock.Lock()
				if h.held[m.Actor] {
					delete(h.held, m.Actor)
					go broadcast(m)
				}
				h.lock.Unlock()

			default:
				messages <- m
			}
		}
	}
}

func (h *Holder) GetHeld() []*Actor {
	h.lock.RLock()
	defer h.lock.RUnlock()

	held := make([]*Actor, 0, len(h.held))
	for a := range h.held {
		held = append(held, a)
	}

	return held
}

func (h *Holder) EachHeld(f func(*Actor)) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for a := range h.held {
		f(a)
	}
}
