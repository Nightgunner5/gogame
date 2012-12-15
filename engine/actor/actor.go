package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"log"
)

type Actor struct {
	Send   message.Sender
	closed chan struct{}
	tag    string
}

func (a *Actor) Initialize() (messages message.Receiver, broadcast func(message.Message)) {
	send_ := make(chan message.Message)
	a.Send = send_

	messages_ := make(chan message.Message)
	messages = messages_

	broadcast_ := make(chan message.Message)
	broadcast = func(m message.Message) {
		broadcast_ <- m
	}

	a.closed = make(chan struct{})

	go a.dispatch(send_, messages_, broadcast_)

	return
}

func (a *Actor) dispatch(send, messages, broadcast chan message.Message) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("panic in ", a.tag)
			panic(r)
		}
	}()

	subscribers := make(map[Subscriber]bool)

	for {
		select {
		case msg, ok := <-send:
			if !ok {
				close(a.closed)
				close(messages)
				return
			}

			switch m := msg.(type) {
			case AddSubscriber:
				subscribers[m.Subscriber] = true

			case RemoveSubscriber:
				delete(subscribers, m.Subscriber)

			default:
				messages <- m
			}
		case msg := <-broadcast:
			for s := range subscribers {
				s.offer(msg)
			}
		}
	}
}

func Init(tag string, bottom *Actor, top interface {
	Initialize() (message.Receiver, func(message.Message))
},) {
	bottom.tag = tag

	// TODO: remove temporary hack
	a, b := top.Initialize()
	go bottom.topLevel(a, b)
}

func (a *Actor) topLevel(messages message.Receiver, _ func(message.Message)) {
	for msg := range messages {
		log.Printf("unhandled message: %s", a.tag)
		panic(msg)
	}
	//log.Printf("actor removed: %s", a.tag)
}
