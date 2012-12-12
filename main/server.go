// +build !noserver

package main

import (
	"github.com/Nightgunner5/fatchan"
	serverpkg "github.com/Nightgunner5/gogame/server"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"
	"sync"
	"time"
)

const canServe = true

var users = struct {
	sync.RWMutex
	chans map[chan packet.Packet]string
	//	taken map[string]bool
}{
	chans: make(map[chan packet.Packet]string),
	//	taken: map[string]bool{"SERVER": true},
}

func addUser(username string, channel chan packet.Packet) bool {
	users.Lock()
	defer users.Unlock()

	//	if users.taken[username] {
	//		log.Printf("Username taken: %q", username)
	//		return false
	//	}

	log.Printf("New user %q", username)
	users.chans[channel] = username
	//	users.taken[username] = true
	return true
}
func delUser(channel chan packet.Packet) {
	users.Lock()
	defer users.Unlock()

	log.Printf("User signoff %q", users.chans[channel])
	//	delete(users.taken, users.chans[channel])
	delete(users.chans, channel)
}
func sendAll(msg packet.Packet) {
	users.RLock()
	defer users.RUnlock()
	for ch := range users.chans {
		select {
		case ch <- msg:
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func broadcastServer() {
	for p := range serverpkg.SendToAll {
		go sendAll(p)
	}
}

func init() {
	go broadcastServer()
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

	player := serverpkg.NewPlayer(id, user.User, user.Recv)
	defer player.Disconnected()

	for msg := range user.Send {
		if !serverpkg.Dispatch(player, msg) {
			return
		}
	}
}
