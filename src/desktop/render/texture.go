package render

import (
  "os"

  "image"
  "image/draw"
  _ "image/png"

  "github.com/go-gl/gl/v4.1-core/gl"
)

type Texture struct {
  handle uint32
  target uint32
  texUnit uint32
}

func (tex *Texture) Bind (texUnit uint32) {
  gl.ActiveTexture(texUnit)
  gl.BindTexture(tex.target, tex.handle)
  tex.texUnit = texUnit
}

func (tex *Texture) UnBind() {
  tex.texUnit = 0
  gl.BindTexture(tex.target, 0)
}

func (tex *Texture) SetUniform(uniformLoc int32) error {
  gl.Uniform1i(uniformLoc, int32(tex.texUnit - gl.TEXTURE0))
  return nil
}

// loading a texture is expensive, want to do it only once per file
var spriteCache = make(map[string]*Texture)

func LoadTexture (filen string) (*Texture, error) {
  // make sure the texture has never been loaded before loading it
  var err error
  if _, loaded := spriteCache[filen]; !loaded {
    spriteCache[filen], err = loadTextureFile(filen) // and cache it once its loaded
  }

  if err != nil {
    return nil, err
  } else {
    return spriteCache[filen], nil
  }
}

func loadTextureFile (filen string) (*Texture, error) {
  imgFile, err := os.Open(filen)
  if err != nil {
    return nil, err
  }

  defer imgFile.Close()

  img, _, err := image.Decode(imgFile)
  if err != nil {
    return nil, err
  }

  return NewTexture(img)
}

func NewTexture (img image.Image) (*Texture, error) {
  rgba := image.NewRGBA(img.Bounds())
  draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

  var handle uint32
  gl.GenTextures(1, &handle)

  target := uint32(gl.TEXTURE_2D)
  internalFmt := int32(gl.SRGB_ALPHA)
  format := uint32(gl.RGBA)
  width := int32(rgba.Rect.Size().X)
  height := int32(rgba.Rect.Size().Y)
  pixType := uint32(gl.UNSIGNED_BYTE)
  dataPtr := gl.Ptr(rgba.Pix)

  texture := Texture{
    handle: handle,
    target: target,
  }

  texture.Bind(gl.TEXTURE0)
  defer texture.UnBind()

  // define sampling behavior
  // gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
  gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
  gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
  gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)  // minification filter
  gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)  // magnification filter


  gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)
  gl.GenerateMipmap(texture.handle)

  return &texture, nil
}
