package actor

import (
	"github.com/Nightgunner5/gogame/engine/message"
	"log"
	"time"
)

type Actor struct {
	Send message.Sender
	closed chan struct{}
	tag  string
}

func (a *Actor) Initialize() (messages message.Receiver, broadcast message.Sender) {
	send_ := make(chan message.Message)
	a.Send = send_

	messages_ := make(chan message.Message)
	messages = messages_

	broadcast_ := make(chan message.Message)
	broadcast = broadcast_

	a.closed = make(chan struct{})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Print("panic in ", a.tag)
				panic(r)
			}
		}()

		subscribers := make(map[Subscriber]bool)

		for {
			select {
			case msg, ok := <-send_:
				if !ok {
					close(a.closed)
					close(messages_)
					return
				}
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

func Init(tag string, bottom *Actor, top interface {
	Initialize() (message.Receiver, message.Sender)
},) {
	bottom.tag = tag
	go topLevel(top.Initialize())(bottom)
}

var msgAlivePing = message.NewKind("alivePing")

type alivePing struct{}

func (alivePing) Kind() message.Kind { return msgAlivePing }

func topLevel(messages message.Receiver, _ message.Sender) func(*Actor) {
	return func(a *Actor) {
		isAlive := make(chan bool, 1)
		go func() {
			isAlive <- true
			for {
				select {
				case _, ok := <-isAlive:
					if !ok {
						return
					}
				default:
					panic("actor is frozen: " + a.tag)
				}

				select {
				case _, _ = <-a.closed:
					return
				}
				a.Send <- alivePing{}
				time.Sleep(time.Second)
			}
		}()
		for msg := range messages {
			if _, ok := msg.(alivePing); ok {
				isAlive <- true
				continue
			}
			log.Printf("unhandled message: %s", a.tag)
			panic(msg)
		}
		//log.Printf("actor removed: %s", a.tag)
	}
}
