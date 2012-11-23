package network

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

const MaxBufferedPackets = 32

func DecodeStream(r net.Conn) <-chan Packet {
	ch := make(chan Packet)
	decoder := json.NewDecoder(r)

	go func() {
		var p packetTransmit
		for {
			r.SetReadDeadline(time.Now().Add(time.Minute))
			if err := decoder.Decode(&p); err != nil {
				if err != io.EOF {
					log.Print("Error decoding packet: ", err)
				}
				r.Close()
				close(ch)
				return
			}
			ch <- p.toPacket()
		}
	}()

	return ch
}

func EncodeStream(w net.Conn) chan<- Packet {
	ch := make(chan Packet, MaxBufferedPackets)
	encoder := json.NewEncoder(w)

	go func() {
		var p packetTransmit
		for packet := range ch {
			p.fromPacket(packet)
			w.SetWriteDeadline(time.Now().Add(time.Minute))
			if err := encoder.Encode(p); err != nil {
				log.Print("Error encoding packet: ", err)
				w.Close()
				return
			}
		}
		w.Close()
	}()

	return ch
}
