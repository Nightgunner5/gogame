package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
)

type Actor struct {
	Send message.Sender
}

func (a *Actor) Initialize() (messages message.Reciever, broadcast message.Sender) {
	send_ := make(chan message.Message)
	a.Send = send_

	messages_ := make(chan message.Message)
	messages = messages_

	broadcast_ := make(chan message.Message)
	broadcast = broadcast_

	go func() {
		subscribers := make(map[Subscriber]bool)

		for {
			select {
			case msg := <-send_:
				switch m := msg.(type) {
				case AddSubscriber:
					subscribers[m.Subscriber] = true
				case RemoveSubscriber:
					delete(subscribers, m.Subscriber)
				default:
					messages_ <- m
				}

			case msg := <-broadcast_:
				for s := range subscribers {
					s.offer(msg)
				}
			}
		}
	}()

	return
}

func TopLevel(messages message.Reciever, _ message.Sender) {
	go func() {
		// unhandled message
		panic(<-messages)
	}()
}
