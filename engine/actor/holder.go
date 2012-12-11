package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
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
	getHeld chan chan []*Actor
}

func (h *Holder) Initialize() (messages message.Receiver, broadcast message.Sender) {
	msgIn, broadcast := h.Actor.Initialize()

	messages_ := make(chan message.Message)
	messages = messages_

	held := make(map[*Actor]bool)
	h.getHeld = make(chan chan []*Actor)

	go func() {
		getHeld := make(chan []*Actor)
		for {
			select {
			case msg, ok := <-msgIn:
				if !ok {
					close(messages_)
					return
				}
				switch m := msg.(type) {
				case AddHeld:
					if !held[m.Actor] {
						held[m.Actor] = true
						broadcast <- m
					}
				case RemoveHeld:
					if held[m.Actor] {
						delete(held, m.Actor)
						broadcast <- m
					}
				default:
					messages_ <- m
				}

			case h.getHeld <- getHeld:
				slice := make([]*Actor, 0, len(held))
				for a := range held {
					slice = append(slice, a)
				}
				getHeld <- slice
			}
		}
	}()

	return
}

func (h *Holder) GetHeld() []*Actor {
	return <-<-h.getHeld
}
