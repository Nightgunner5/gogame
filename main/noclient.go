// +build noclient,!noserver

package main

import (
	"io"
	"log"
)

func client(string, io.ReadWriteCloser) {
	log.Fatal("client mode is not available in this binary")
}
