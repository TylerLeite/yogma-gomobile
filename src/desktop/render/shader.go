package render

import (
  "fmt"
  "strings"

  "github.com/go-gl/gl/v4.1-core/gl"
)

func CompileShader (source string, shaderType uint32) (uint32, error) {
  // Create an empty shader
  shader := gl.CreateShader(shaderType)

  // Load it up with the source
  sourceStrings, freeFn := gl.Strs(source)
  gl.ShaderSource(shader, 1, sourceStrings, nil)
  freeFn()

  // Compile the shader
  gl.CompileShader(shader)

  var status int32
  gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
  if status == gl.FALSE {
    // Compilation errors occurred

    // Find the length of the error log
    var logLength int32
    gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

    // Request that many lines
    log := strings.Repeat("\x00", int(logLength+1))
    gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

    // Print
    return 0, fmt.Errorf("failed to compile %v: %v", source, log)
  }

  // Else success
  return shader, nil
}
