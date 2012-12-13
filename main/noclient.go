// +build noclient,!noserver,!monkey

package main

import (
	"io"
	"log"
)

const (
	DefaultAddr   = ""
	DefaultServer = true
)

func client(string, io.ReadWriteCloser) {
	log.Fatal("client mode is not available in this binary")
}
