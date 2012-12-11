// +build !noserver

package main

import (
	"io"
	"log"
	"sync"

	"github.com/Nightgunner5/gogame/shared/packet"
	"github.com/kylelemons/fatchan"
)

const canServe = true

var users = struct {
	sync.RWMutex
	chans map[chan packet.Packet]string
	taken map[string]bool
}{
	chans: make(map[chan packet.Packet]string),
	taken: map[string]bool{"SERVER": true},
}

func addUser(username string, channel chan packet.Packet) bool {
	users.Lock()
	defer users.Unlock()

	if users.taken[username] {
		log.Printf("Username taken: %q", username)
		return false
	}

	log.Printf("New user %q", username)
	users.chans[channel] = username
	users.taken[username] = true
	return true
}
func delUser(channel chan packet.Packet) {
	users.Lock()
	defer users.Unlock()

	log.Printf("User signoff %q", users.chans[channel])
	delete(users.taken, users.chans[channel])
	delete(users.chans, channel)
}
func sendAll(msg packet.Packet) {
	users.RLock()
	defer users.RUnlock()
	for ch := range users.chans {
		ch <- msg
	}
}

func serve(id string, client io.ReadWriteCloser) {
	log.Printf("Client %q connected", id)
	defer log.Printf("Client %q disconnected", id)

	xport := fatchan.New(client, nil)
	login := make(chan Handshake)
	xport.ToChan(login)

	user := <-login
	log.Printf("Client %q registered as %q", id, user.User)
	if !addUser(user.User, user.Recv) {
		user.Recv <- packet.Packet{
			Chat: &packet.Chat{
				User:    "SERVER",
				Message: "That username is already taken.",
			},
		}
		close(user.Recv)
		return
	}
	defer delUser(user.Recv)
	defer close(user.Recv)

	for msg := range user.Send {
		switch {
		case msg.Chat != nil:
			msg.Chat.User = user.User
			go sendAll(packet.Packet{
				Chat: msg.Chat,
			})
		default:
			log.Printf("Client %q sent unknown packet %#v", id, msg)
			return
		}
	}
}