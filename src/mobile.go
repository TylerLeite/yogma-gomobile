// +build !darwin,!linux

package main

import (
  "math/rand"
  "time"

  "golang.org/x/mobile/app"
  "golang.org/x/mobile/event/key"
  "golang.org/x/mobile/event/lifecycle"
  "golang.org/x/mobile/event/paint"
  "golang.org/x/mobile/event/size"
  "golang.org/x/mobile/event/touch"
  "golang.org/x/mobile/exp/gl/glutil"
  "golang.org/x/mobile/exp/sprite"
  "golang.org/x/mobile/exp/sprite/clock"
  "golang.org/x/mobile/exp/sprite/glsprite"
  "golang.org/x/mobile/geom"
  "golang.org/x/mobile/gl"

  "github.com/TylerLeite/yogma/src/mobile/input"
  "github.com/TylerLeite/yogma/src/mobile/render"
)

func setDimensions (sz size.Event, width, height int) size.Event {
  var ppp float32 = sz.PixelsPerPt
  var orientation size.Orientation = size.OrientationLandscape
  if height > width {
    orientation = size.OrientationPortrait
  }

  var newSize size.Event = size.Event{
    WidthPx: width,
    HeightPx: height,
    WidthPt: geom.Pt(float32(width)/ppp),
    HeightPt: geom.Pt(float32(height)/ppp),
    PixelsPerPt: ppp,
    Orientation: orientation,
  }

  return newSize
}

func main () {
  rand.Seed(time.Now().UnixNano())

  app.Main(func(a app.App) {
    // return
    var glctx gl.Context
    var sz size.Event

    for e := range a.Events() {
      switch e := a.Filter(e).(type) {
      case lifecycle.Event:
        switch e.Crosses(lifecycle.StageVisible) {
        case lifecycle.CrossOn:
          glctx, _ = e.DrawContext.(gl.Context)
          onStart(glctx)
          a.Send(paint.Event{})
        case lifecycle.CrossOff:
          onStop()
          glctx = nil
        }
      case size.Event:
        // sz = setDimensions(e, 2160, 1080)
        sz = e;
      case paint.Event:
        if glctx == nil || e.External {
          continue
        }

        onPaint(glctx, sz)
        a.Publish() // write to screen
        a.Send(paint.Event{}) // call for next frame
      case touch.Event:
        input.RouteTouch(e, scene)
      case key.Event:
        if e.Code != key.CodeSpacebar {
          break
        }
        if down := e.Direction == key.DirPress; down || e.Direction == key.DirRelease {
          // game.Press(down)
        }
      }
    }
  })
}

var (
  startTime = time.Now()
  images    *glutil.Images
  eng       sprite.Engine
  scene     *sprite.Node
  // game      *Game
)

func onStart (glctx gl.Context) {
  images = glutil.NewImages(glctx)
  eng = glsprite.Engine(images)
  // game = NewGame()
  scene = render.InitScene(eng)
}

func onStop () {
  eng.Release()
  images.Release()
  // game = nil
}

func onPaint (glctx gl.Context, sz size.Event) {
  glctx.ClearColor(1, 1, 1, 1)
  glctx.Clear(gl.COLOR_BUFFER_BIT)

  // 60 fps
  now := clock.Time(time.Since(startTime) * 60 / time.Second)
  // game.Update(now)

  eng.Render(scene, now, sz)
}
