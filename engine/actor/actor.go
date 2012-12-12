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
	Initialize() (message.Receiver, message.Sender)
},) {
	bottom.tag = tag
	go bottom.topLevel(top.Initialize())
}

var msgAlivePing = message.NewKind("alivePing")

type alivePing struct{}

func (alivePing) Kind() message.Kind { return msgAlivePing }

func (a *Actor) checkAlive(isAlive chan struct{}) {
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
}

func (a *Actor) topLevel(messages message.Receiver, _ message.Sender) {
	isAlive := make(chan struct{}, 1)
	go a.checkAlive(isAlive)
	for msg := range messages {
		if _, ok := msg.(alivePing); ok {
			isAlive <- struct{}{}
			continue
		}
		log.Printf("unhandled message: %s", a.tag)
		panic(msg)
	}
	//log.Printf("actor removed: %s", a.tag)
}
