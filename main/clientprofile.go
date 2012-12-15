// +build profile,!noclient

package main

import (
	"os"

	clientpkg "github.com/Nightgunner5/gogame/client"
)

func init() {
	oldDisconnected := clientpkg.Disconnected
	clientpkg.Disconnected = func() {
		profileCleanup <- 0
		oldDisconnected()
	}
}
