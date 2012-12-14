package main

import (
	"flag"
	"log"
	"net"
	"os"
)

var (
	server = flag.Bool("server", DefaultServer, "Start in server mode")
	addr   = flag.String("addr", DefaultAddr, "Address to listen on or connect to")
	user   = flag.String("user", os.Getenv("USER"), "Username (ignored in server mode)")
)

type Handshake struct {
	Monkey bool
}

func main() {
	flag.Parse()
	if *addr == "" {
		log.Fatalf("error: must specify -addr. see -help for arguments.")
	}

	if *server {
		listenAndServe(*addr)
	} else {
		connectTo(*user, *addr)
	}
}

func listenAndServe(addr string) {
	if !canServe {
		serve("", nil)
		return
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen(%q): %s", addr, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			if err.(net.Error).Temporary() {
				log.Printf("accept(): %s", err)
			} else {
				log.Fatalf("accept(): %s", err)
			}
		}

		go serve(conn.RemoteAddr().String(), conn)
	}
}

func connectTo(username, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial(%q): %s", addr, err)
	}

	client(username, conn)
}
