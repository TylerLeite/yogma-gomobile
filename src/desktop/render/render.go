package render

import (
  "log"

  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/glfw/v3.2/glfw"
)

const (
  texVertexShaderSource = `
    #version 410 core

    layout (location = 0) in vec3 position;
    layout (location = 1) in vec2 texCoord;

    out vec2 TexCoord;

    void main () {
      gl_Position = vec4(position, 1.0);
      TexCoord = texCoord;
    }
  ` + "\x00"

  texFragmentShaderSource = `
    #version 410 core

    in vec2 TexCoord;

    uniform sampler2D ourTexture0;

    out vec4 FragColor;

    void main () {
      FragColor = texture(ourTexture0, TexCoord);
    }
  ` + "\x00"
)

// var (
//   square = []float32{
//     0.25, 0.5, 0, // top right
//     1.0, 0.0,
//     0.25, -0.5, 0, // bottom right
//     1.0, 1.0,
//     -0.25, -0.5, 0, // bottom left
//     0.0, 1.0,
//
//     0.25, 0.5, 0, //top right
//     1.0, 0.0,
//     -0.25, 0.5, 0, // top left
//     0.0, 0.0,
//     -0.25, -0.5, 0, // bottom left
//     0.0, 1.0,
//   }
// )

const (
  width = 950
  height = 450
)
func RectAt (x, y, wd, hg float32) []float32 {
  // Need to take into account image size & screen size
  x = x/width
  y = y/height
  wd = wd/width
  hg = hg/height
  return rectAtGLCoords(x, y, wd, hg)
}

func rectAtGLCoords (x, y, wd, hg float32) []float32 {
  return []float32{
    x+wd, y+hg, 0, // top right
    1.0, 0.0,
    x+wd, y-hg, 0, // bottom right
    1.0, 1.0,
    x-wd, y-hg, 0, // bottom left
    0.0, 1.0,

    x+wd, y+hg, 0, //top right
    1.0, 0.0,
    x-wd, y+hg, 0, // top left
    0.0, 0.0,
    x-wd, y-hg, 0, // bottom left
    0.0, 1.0,
  }
}


// Init gl and glfw, return the gl program and glfw window
func InitScene (width, height int) (uint32, *glfw.Window) {
  // init glfw
  if err := glfw.Init(); err != nil {
    log.Printf("Error initiating glfw")
    panic(err)
  }

  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
  glfw.WindowHint(glfw.ContextVersionMinor, 1) // would like to use 4.6, but glfw wont let me for some reason
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

  window, err := glfw.CreateWindow(width, height, "Yogma", nil, nil)
  if err != nil {
    log.Printf("Error creating window")
    panic(err)
  }
  window.MakeContextCurrent()

  // init opengl
  // note: this must happen 2nd to avoid locking issues with cgo calls
  //       or something
  if err := gl.Init(); err != nil {
    log.Printf("Error initiating gl")
    panic(err)
  }

  var program uint32 = gl.CreateProgram()

  vertexShader, err := CompileShader(texVertexShaderSource, gl.VERTEX_SHADER)
  if err != nil {
    panic(err)
  }

  fragmentShader, err := CompileShader(texFragmentShaderSource, gl.FRAGMENT_SHADER)
  if err != nil {
    panic(err)
  }

  gl.AttachShader(program, vertexShader)
  gl.AttachShader(program, fragmentShader)

  gl.LinkProgram(program)

  // gl.ClearColor(1,1,1,1)

  // well the good news is we dont need to return any errors
  // the bad news is thats because all possible errors here are fatal
  return program, window
}

func makeVertexArrayObject (points []float32) uint32 {
  var vertexArrayObject uint32
  gl.GenVertexArrays(1, &vertexArrayObject)
  gl.BindVertexArray(vertexArrayObject)

  var vertexBufferObject uint32
  gl.GenBuffers(1, &vertexBufferObject)
  gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
  gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

  // var elementBufferObject uint32
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, elementBufferObject)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(triangles)*4, gl.Ptr(triangles), gl.STATIC_DRAW)

  var pointDataSize int32 = 3*4 + 2*4
	var pointDataPosition int = 0

  // position
  gl.VertexAttribPointer(0, 3, gl.FLOAT, false, pointDataSize, gl.PtrOffset(pointDataPosition))
  gl.EnableVertexAttribArray(0)
  pointDataPosition += 3*4

  // texture
  gl.VertexAttribPointer(1, 2, gl.FLOAT, false, pointDataSize, gl.PtrOffset(pointDataPosition))
	gl.EnableVertexAttribArray(1)
	pointDataPosition += 2*4

  gl.BindVertexArray(0)

  return vertexArrayObject
}

func OneFrame (program uint32, window *glfw.Window) {
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
  gl.UseProgram(program)

  var err error

  var plainPlains *Texture
  plainPlains, err = LoadTexture("assets/sprites/terrain/plains/0.png")
  if err != nil {
    panic(err)
  }

  var flowerPlains *Texture
  flowerPlains, err = LoadTexture("assets/sprites/terrain/plains/1.png")
  if err != nil {
    panic(err)
  }

  // set texture0 to uniform0 in the fragment shader
  plainPlains.Bind(gl.TEXTURE0)
  plainPlains.SetUniform(gl.GetUniformLocation(program, gl.Str("ourTexture0" + "\x00")))

  flowerPlains.Bind(gl.TEXTURE1)
  flowerPlains.SetUniform(gl.GetUniformLocation(program, gl.Str("ourTexture1" + "\x00")))

  const tileWidth, tileHeight float32 = 49.0, 50.0
  const tileOffsetHorizX, tileOffsetHorizY float32 = 62.0, -16.0
  const tileOffsetVertX, tileOffsetVertY float32 = -32.0, -32.0

  var tiles []float32
  var t1 []float32 = RectAt(-300, 0, tileWidth, tileHeight)
  var t2 []float32 = RectAt(-300 + tileOffsetHorizX, 0 + tileOffsetHorizY, tileWidth, tileHeight)
  var t3 []float32 = RectAt(-300 + tileOffsetVertX, 0 + tileOffsetVertY, tileWidth, tileHeight)
  tiles = append(tiles, t1...)
  tiles = append(tiles, t2...)
  tiles = append(tiles, t3...)

  vao := makeVertexArrayObject(tiles)
  gl.BindVertexArray(vao)
  gl.DrawArrays(gl.TRIANGLES, 0, int32(len(tiles) / 3))

  plainPlains.UnBind()
  flowerPlains.UnBind()

  glfw.PollEvents()
  window.SwapBuffers()
}
