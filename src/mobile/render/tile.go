package render

import (
  "log"
  "image"
  "math/rand"

  _ "image/png" // gotta do this or else u cant decode pngs

  "golang.org/x/mobile/asset"
  "golang.org/x/mobile/exp/f32"
  "golang.org/x/mobile/exp/sprite"
  "golang.org/x/mobile/exp/sprite/clock"
)

var spriteCache = make(map[string]*sprite.SubTex)
const tileWidth, tileHeight float32 = 49.0, 50.0
const tileOffsetHorizX, tileOffsetHorizY float32 = 32.0, 8.0
const tileOffsetVertX, tileOffsetVertY float32 = 16.0, 16.0
type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)
func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

func DrawTile (eng sprite.Engine, scene *sprite.Node, filen string, x int, y int, z int) {
  if rand.Intn(25) < 7 {
    filen += "1.png"
  } else {
    filen += "0.png"
  }

  tex := loadTexture(eng, filen)

  var tileArranger arrangerFunc = func(eng sprite.Engine, n *sprite.Node, t clock.Time) {
    eng.SetSubTex(n, *tex)

    /*
     * The top-left of the tiles you can see doesn't line up with the top-left
     * of the screen, because of the way the perspective works out. As you go down
     * and stay in the same column, the position of the tile on the screen moves
     * left by a factor of 16/49 of the tile's width. Therefore, we need to offset
     * the whole tile cluster by that factor
    **/

    const clusterWidth float32 = 10.0
    const clusterHeight float32 = 10.0

    var scrollX float32 = clusterHeight * tileWidth * 16.0 / 49.0
    var scrollY float32 = 16

    // Moving right 1 tile is equivalent to moving 32 pixels right + 8 pixels down
    // Moving down 1 tile is equivalent to moving 16 pixels left + 16 pixels down
    // Therefore your position is given by:
    var tileOffsetX = float32(x)*tileOffsetHorizX - float32(y)*tileOffsetVertX
    var tileOffsetY = float32(x)*tileOffsetHorizY + float32(y)*tileOffsetVertY
    eng.SetTransform(n, f32.Affine{{tileWidth, 0, scrollX + tileOffsetX},{0, tileHeight, scrollY + tileOffsetY},})
  }

  n := &sprite.Node{Arranger: tileArranger}
  eng.Register(n)
  scene.AppendChild(n)
}


func InitScene (eng sprite.Engine) *sprite.Node {
  scene := &sprite.Node{}
  eng.Register(scene)
  eng.SetTransform(scene, f32.Affine{
    {1, 0, 0},
    {0, 1, 0},
  })

  // Call DrawTile a bunch of times to initialize the scene
  for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
      DrawTile(eng, scene, "sprites/terrain/plains/", i, j, 0)
    }
  }

  return scene
}

func loadTexture (eng sprite.Engine, filen string) *sprite.SubTex {
  var tex *sprite.SubTex = nil
  var loaded bool = false

  if tex, loaded = spriteCache[filen]; !loaded {
    a, err := asset.Open(filen)
    if err != nil {
      log.Fatal(err)
    }
    defer a.Close()

    m, _, err := image.Decode(a)
    if err != nil {
      log.Fatal(err)
    }

    t, err := eng.LoadTexture(m)
    if err != nil {
      log.Fatal(err)
    }

    spriteCache[filen] = &sprite.SubTex{t, image.Rect(0, 0, int(tileWidth), int(tileHeight))}
    tex = spriteCache[filen]
  }

  return tex
}
