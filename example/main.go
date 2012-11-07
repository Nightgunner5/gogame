package main

import (
	"flag"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"log"
	"math/rand"
	"runtime"
	_ "net/http/pprof"
	"net/http"
)

const (
	maxHealth     = 100
	maxMana       = 20
	manaPerSecond = 1
	initialHealth = 100
	initialMana   = 10

	manaForDamageSpell = 10
	damageCastTime     = 1
	spellDamage        = 75

	manaForHealingSpell = 3
	healCastTime        = 4
	spellHealing        = 5
)

var mageCount = flag.Int("mages", 4, "The number of mages at the start")

func main() {
	flag.Parse()

	go http.ListenAndServe("localhost:6060", nil)

	go func() {
		for i := 0; i < *mageCount; i++ {
			entity.Spawn(&mage{
				health: initialHealth * variance(),
				mana:   initialMana * variance(),
				x:      rand.Float64()*20 - 10,
				y:      rand.Float64()*20 - 10,
			})
		}
	}()

	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}

	defer glfw.Terminate()

	if err := glfw.OpenWindow(256, 256, 8, 8, 8, 0, 0, 0, glfw.Windowed); err != nil {
		log.Fatal(err)
	}

	defer glfw.CloseWindow()

	glfw.SetWindowTitle("GoGame")
	glfw.SetSwapInterval(1)
	glfw.SetWindowSizeCallback(windowSize)

	gl.ClearColor(1, 1, 1, 1)

	for glfw.WindowParam(glfw.Opened) == 1 {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		render()

		glfw.SwapBuffers()

		runtime.Gosched()
	}
}

func windowSize(width, height int) {
	if width == 0 || height == 0 {
		return
	}

	aspect := float64(width) / float64(height)

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	if aspect < 1 {
		gl.Ortho(-10, 10, -10/aspect, 10/aspect, -10, 10)
	} else {
		gl.Ortho(-10*aspect, 10*aspect, -10, 10, -10, 10)
	}
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
