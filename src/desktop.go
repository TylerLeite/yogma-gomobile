package main

import (
  "math/rand"
  "time"
  // "log"
  "runtime"

  "github.com/go-gl/glfw/v3.2/glfw"

  // "github.com/TylerLeite/yogma/src/desktop/input"
  "github.com/TylerLeite/yogma/src/desktop/render"
)

const (
  width int = 950
  height int = 450
)

var (
  frame int = 0
)

func main () {
  rand.Seed(time.Now().UnixNano())

  // window only available on the thread on which it was created
  runtime.LockOSThread()

  // program uint32
  // window *glfw.Window
  program, window := render.InitScene(width, height)
  defer glfw.Terminate()

  for !window.ShouldClose() {
    // TODO: input
    // TODO: game logic
    // TODO: animation logic

    render.OneFrame(program, window)

    timer := time.NewTimer(16 * time.Millisecond)
    <-timer.C
    frame += 1

    // TODO: if there's spare time at the end of a frame, load shit or autosave
  }
}
