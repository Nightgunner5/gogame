// +build profile,!noclient

package main

import (
	clientpkg "github.com/Nightgunner5/gogame/client"
)

func init() {
	oldDisconnected := clientpkg.Disconnected
	clientpkg.Disconnected = func() {
		profileCleanup <- nil
		oldDisconnected()
	}
}
