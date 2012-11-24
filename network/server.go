package network

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net"
	"net/http"
	"sync"

	_ "github.com/Nightgunner5/gogame/client"
)

func init() {
	http.Handle("/socket", websocket.Handler(handleSocket))
}

var (
	connections    = make(map[*websocket.Conn]chan<- Packet)
	connectionLock sync.RWMutex
)

func Broadcast(packet Packet, ensureSend bool) {
	connectionLock.RLock()
	defer connectionLock.RUnlock()

	for _, stream := range connections {
		select {
		case stream <- packet:
		default:
			if ensureSend {
				go func(stream chan<- Packet) {
					stream <- packet
				}(stream)
			}
		}
	}
}

type remoteAddr string

func (r remoteAddr) String() string {
	return string(r)
}

func (remoteAddr) Network() string {
	return "websocket"
}

func handleSocket(conn *websocket.Conn) {
	in, out := DecodeStream(conn), EncodeStream(conn)

	connectionLock.Lock()
	connections[conn] = out
	connectionLock.Unlock()

	addr := remoteAddr(conn.Request().RemoteAddr)

	startupListener(out, addr)

	for p := range in {
		dispatchPacket(p, out, addr, conn)
	}

	connectionLock.Lock()
	delete(connections, conn)
	connectionLock.Unlock()
}

func dispatchPacket(p Packet, send chan<- Packet, addr net.Addr, conn *websocket.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic when handling packet %d for client %v: %v", p.ID(), addr, err)
		}
	}()

	if handler, ok := registered[p.ID()]; ok {
		handler(p, send, addr)
	} else {
		conn.Close()
		log.Panicf("No handler found for packet ID %d", p.ID())
	}
}

var startupListener func(chan<- Packet, net.Addr) = func(chan<- Packet, net.Addr) {}
var registered = make(map[PacketID]func(Packet, chan<- Packet, net.Addr))

// Must be called from an init() function.
func RegisterStartupListener(f func(chan<- Packet, net.Addr)) {
	oldListener := startupListener

	startupListener = func(send chan<- Packet, addr net.Addr) {
		oldListener(send, addr)
		f(send, addr)
	}
}

// Must be called from an init() function.
func RegisterHandler(id PacketID, f func(Packet, chan<- Packet, net.Addr)) {
	if _, exists := registered[id]; exists {
		log.Panicf("Duplicate registration for packet %d", id)
	}
	registered[id] = f
}

func init() {
	RegisterHandler(debugEcho, func(p Packet, reply chan<- Packet, addr net.Addr) {
		reply <- p
	})
}
