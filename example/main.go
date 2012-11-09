package main

import (
	"flag"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
	"log"
	"math/rand"
	"runtime"
)

const (
	maxHealth     = 100
	maxMana       = 20
	manaPerSecond = 1

	manaForDamageSpell = 10
	damageCastTime     = 1
	spellDamage        = 75

	manaForHealingSpell = 3
	healCastTime        = 4
	spellHealing        = 10

	visWidth  = 10
	visHeight = 10
)

var mageCount = flag.Int("mages", 4, "The number of mages at the start")

func main() {
	flag.Parse()

	go func() {
		for i := 0; i < *mageCount; i++ {
			entity.Spawn(&mage{
				BaseHealth:   entity.BaseHealth{Max: maxHealth},
				BaseResource: entity.BaseResource{Max: maxMana},
				x:            rand.Float64()*2*visWidth - visWidth,
				y:            rand.Float64()*2*visHeight - visHeight,
			})
		}
		entity.Spawn(&spawner{
			BaseResource: entity.BaseResource{Max: 10},
		})
	}()

	runtime.LockOSThread() // OpenGL doesn't like thread switches

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
		gl.Ortho(-visWidth, visWidth, -visHeight/aspect, visHeight/aspect, -10, 10)
	} else {
		gl.Ortho(-visWidth*aspect, visWidth*aspect, -visHeight, visHeight, -10, 10)
	}
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
