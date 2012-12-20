// +build !noserver

package main

import (
	"bufio"
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

	bufclient := bufio.NewWriter(client)
	encode, decode := gob.NewEncoder(bufclient), gob.NewDecoder(client)

	var handshake Handshake
	if err := decode.Decode(&handshake); err != nil {
		log.Printf("[client:%s] decoding handshake: %s", id, err)
		return
	}
	log.Printf("[client:%s] %#v", id, handshake)

	send := make(chan *packet.Packet)

	go serverEncode(id, bufclient, encode, client, send)

	defer close(send)
	if !addUser(send, handshake) {
		return
	}
	defer delUser(send)

	var flags uint32
	if handshake.Monkey {
		flags |= packet.FlagMonkey
	} else {
		flags |= packet.FlagSuit
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

func serverEncode(id string, bufclient *bufio.Writer, encode *gob.Encoder, client net.Conn, send <-chan *packet.Packet) {
	tick := time.NewTicker(time.Second / 20)
	defer tick.Stop()

	for {
		select {
		case msg, ok := <-send:
			if !ok {
				return
			}
			if err := encode.Encode(msg); err != nil {
				log.Printf("[client:%s] encoding packet: %s", id, err)
			}

		case <-tick.C:
			client.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if err := bufclient.Flush(); err != nil {
				log.Printf("[client:%s] sending buffer: %s", id, err)
			}
		}
	}
}
