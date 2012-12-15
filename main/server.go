// +build !noserver

package main

import (
	"encoding/gob"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"
	"net"
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
		ch <- msg
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

func serve(id string, client net.Conn) {
	log.Printf("[client:%s] connected", id)
	defer client.Close()
	defer log.Printf("[client:%s] disconnected", id)

	encode, decode := gob.NewEncoder(client), gob.NewDecoder(client)

	var handshake Handshake
	if err := decode.Decode(&handshake); err != nil {
		log.Printf("[client:%s] decoding handshake: %s", id, err)
		return
	}
	log.Printf("[client:%s] %#v", id, handshake)

	send := make(chan *packet.Packet)

	if handshake.Monkey {
		go serverMonkeyEncode(send)
	} else {
		go serverEncode(id, encode, client, send)
	}

	defer close(send)
	if !addUser(send, handshake) {
		return
	}
	defer delUser(send)

	var flags uint32
	if handshake.Monkey {
		flags |= packet.FlagMonkey
	}

	// TODO: names
	player := serverpkg.NewPlayer(id, send, flags)
	defer player.Disconnected()

	var fastZero packet.Packet
	for {
		msg := new(packet.Packet)
		client.SetReadDeadline(time.Now().Add(2 * time.Minute))
		if err := decode.Decode(msg); err != nil {
			if err == io.EOF {
				return
			}

			log.Printf("[client:%s] decoding packet: %s", id, err)
			continue
		}
		if *msg != fastZero && !serverpkg.Dispatch(player, msg) {
			return
		}
	}
}

func serverMonkeyEncode(send <-chan *packet.Packet) {
	for _ = range send {
	}
}

func serverEncode(id string, encode *gob.Encoder, client net.Conn, send <-chan *packet.Packet) {
	for msg := range send {
		client.SetWriteDeadline(time.Now().Add(30 * time.Second))
		err := encode.Encode(msg)
		if err != nil {
			log.Printf("[client:%s] encoding packet: %s", id, err)
		}
	}
}
