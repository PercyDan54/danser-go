//go:build windows

package gcontext

import (
	"github.com/TheTitanrain/w32"
	"github.com/Zyko0/go-sdl3/sdl"
)

var winInterval = 1

func SetSwapInterval(interval int) {
	winInterval = interval

	if !fullscreen {
		interval = 0
	}

	if err := sdl.GL_SetSwapInterval(int32(interval)); err != nil {
		panic(err)
	}
}

func SwapBuffers() {
	if !fullscreen {
		for range winInterval {
			w32.DwmFlush() // For some reason SDL doesn't do this so...
		}
	}

	if err := sdl.GL_SwapWindow(sdlWindow); err != nil {
		panic(err)
	}
}
