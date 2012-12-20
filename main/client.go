// +build !noclient

package main

import (
	"bufio"
	"encoding/gob"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"
	"time"

	clientpkg "github.com/Nightgunner5/gogame/client"
)

const (
	DefaultAddr   = "nightgunner5.is-a-geek.net:7031"
	DefaultServer = false
)

func client(username string, server io.ReadWriteCloser) {
	log.Print("Connected to server")
	defer server.Close()
	defer log.Print("Server disconnected")

	buf := bufio.NewWriter(server)
	encode, decode := gob.NewEncoder(buf), gob.NewDecoder(server)

	var handshake Handshake
	if err := encode.Encode(handshake); err != nil {
		log.Printf("Error while sending handshake: %s", err)
		return
	}

	send := make(chan *packet.Packet)
	defer close(send)

	clientpkg.Network = send

	go clientEncode(encode, buf, send)
	go clientDecode(decode)

	go keepAlive(send)

	clientpkg.Main()
}

func keepAlive(c chan<- *packet.Packet) {
	var keepAlive packet.Packet
	for {
		c <- &keepAlive
		time.Sleep(1 * time.Minute)
	}
}

func clientEncode(encode *gob.Encoder, buf *bufio.Writer, send <-chan *packet.Packet) {
	for msg := range send {
		if err := encode.Encode(msg); err != nil {
			log.Printf("Error encoding packet: %s", err)
		}
		if err := buf.Flush(); err != nil {
			log.Printf("Error writing buffer: %s", err)
		}
	}
}

func clientDecode(decode *gob.Decoder) {
	defer clientpkg.Disconnected()
	for {
		msg := new(packet.Packet)

		if err := decode.Decode(msg); err != nil {
			if err == io.EOF {
				return
			}

			log.Printf("Error decoding packet: %s", err)
			continue
		}

		clientpkg.Handle(msg)
	}
}
