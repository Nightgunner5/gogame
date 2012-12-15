// +build profile,!noserver,noclient profile,noserver,!noclient

package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"time"
)

func init() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}

	cleanup := make(chan os.Signal, 1)

	go func() {
		tick := time.Tick(time.Minute)
		for {
			select {
			case <-cleanup:
				pprof.StopCPUProfile()
				os.Exit(0)
			case <-tick:
				for _, profile := range pprof.Profiles() {
					f, err := os.Create(profile.Name() + ".prof")
					if err != nil {
						panic(err)
					}
					err = profile.WriteTo(f, map[string]int{
						"threadcreate": 1,
						"goroutine":    2,
					}[profile.Name()])
					if err != nil {
						panic(err)
					}
					f.Close()
				}
			}
		}
	}()

	signal.Notify(cleanup, os.Interrupt)
}
