// +build noclient,noserver,monkey

package main

import (
	"bufio"
	"encoding/gob"
	"github.com/Nightgunner5/gogame/shared/layout"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	DefaultAddr   = "127.0.0.1:7031"
	DefaultServer = false
)

func init() {
	monkey = true
}

func client(username string, server io.ReadWriteCloser) {
	log.Print("Connected to server")
	defer server.Close()
	defer log.Print("Server disconnected")

	buf := bufio.NewWriter(server)
	encode := gob.NewEncoder(buf)

	var handshake Handshake

	handshake.Monkey = true

	if err := encode.Encode(handshake); err != nil {
		log.Printf("Error while sending handshake: %s", err)
		return
	}

	send := make(chan *packet.Packet)
	defer close(send)

	go clientEncode(encode, buf, send)
	go io.Copy(ioutil.Discard, server)

	for {
		send <- &packet.Packet{
			Location: &packet.Location{
				Coord: layout.Coord{
					rand.Intn(3) - 1,
					rand.Intn(3) - 1,
				},
			},
		}

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func clientEncode(encode *gob.Encoder, buf *bufio.Writer, send <-chan *packet.Packet) {
	for msg := range send {
		if err := encode.Encode(msg); err != nil {
			log.Printf("Error encoding packet: %s", err)
			os.Exit(1)
		}
		if err := buf.Flush(); err != nil {
			log.Printf("Error sending buffer: %s", err)
			os.Exit(1)
		}
	}
}
