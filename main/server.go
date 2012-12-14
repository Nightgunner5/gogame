// +build !noserver

package main

import (
	"encoding/gob"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"
	"sync"
	"time"

	serverpkg "github.com/Nightgunner5/gogame/server"
)

const canServe = true

var users = struct {
	sync.RWMutex
	chans map[chan *packet.Packet]Handshake
}{
	chans: make(map[chan *packet.Packet]Handshake),
}

func addUser(channel chan *packet.Packet, handshake Handshake) bool {
	users.Lock()
	defer users.Unlock()

	users.chans[channel] = handshake
	return true
}
func delUser(channel chan *packet.Packet) {
	users.Lock()
	defer users.Unlock()

	delete(users.chans, channel)
}
func sendAll(msg *packet.Packet) {
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
	defer client.Close()
	defer log.Printf("Client %q disconnected", id)

	encode, decode := gob.NewEncoder(client), gob.NewDecoder(client)

	var handshake Handshake
	if err := decode.Decode(&handshake); err != nil {
		log.Printf("[client:%s] decoding handshake: %s", id, err)
		return
	}

	send := make(chan *packet.Packet)

	go serverEncode(id, encode, send)

	defer close(send)
	if !addUser(send, handshake) {
		return
	}
	defer delUser(send)

	player := serverpkg.NewPlayer(id, "", send) // TODO: names
	defer player.Disconnected()

	for {
		msg := new(packet.Packet)
		if err := decode.Decode(msg); err != nil {
			if err == io.EOF {
				return
			}

			log.Printf("[client:%s] decoding packet: %s", id, err)
			continue
		}
		if !serverpkg.Dispatch(player, msg) {
			return
		}
	}
}

func serverEncode(id string, encode *gob.Encoder, send <-chan *packet.Packet) {
	for msg := range send {
		err := encode.Encode(msg)
		if err != nil {
			log.Printf("[client:%s] encoding packet: %s", id, err)
		}
	}
}
