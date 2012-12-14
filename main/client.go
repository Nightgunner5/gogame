// +build !noclient

package main

import (
	"encoding/gob"
	"github.com/Nightgunner5/gogame/shared/packet"
	"io"
	"log"

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

	encode, decode := gob.NewEncoder(server), gob.NewDecoder(server)

	var handshake Handshake
	if err := encode.Encode(handshake); err != nil {
		log.Printf("Error while sending handshake: %s", err)
		return
	}

	send := make(chan *packet.Packet)
	defer close(send)

	clientpkg.Network = send

	go clientEncode(encode, send)
	go clientDecode(decode)

	clientpkg.Main()
}

func clientEncode(encode *gob.Encoder, send <-chan *packet.Packet) {
	for msg := range send {
		if err := encode.Encode(msg); err != nil {
			log.Printf("Error encoding packet: %s", err)
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
