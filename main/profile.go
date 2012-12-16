// +build profile,!noserver,noclient profile,noserver,!noclient

package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"time"
)

var (
	profileCleanup = make(chan os.Signal, 1)
	clientCanExit  = make(chan struct{})
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

	go func() {
		tick := time.Tick(time.Minute)
		for {
			select {
			case signal := <-profileCleanup:
				pprof.StopCPUProfile()
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
				if signal == os.Interrupt {
					os.Exit(0)
				} else {
					clientCanExit <- struct{}{}
				}

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

	signal.Notify(profileCleanup, os.Interrupt)
}
