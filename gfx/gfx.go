package gfx

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//go:embed shaders/*
var shaderFS embed.FS

// InitGlfw initializes glfw and returns a Window to use.
func InitGlfw(width, height int, title string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	return window
}

// InitOpenGL initializes OpenGL and returns an initialized program.
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	// 1. Read the Vertex Shader
	vertSrc, err := shaderFS.ReadFile("shaders/main.vert")
	if err != nil {
		panic(err)
	}

	// 2. Read the Fragment Shader
	fragSrc, err := shaderFS.ReadFile("shaders/main.frag")
	if err != nil {
		panic(err)
	}

	// 3. Compile them (We must append the null terminator "\x00")
	vertexShader, err := compileShader(string(vertSrc)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(string(fragSrc)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

// CreateBatchObjects creates a VAO and VBO meant for dynamic updates
func CreateBatchObjects() (uint32, uint32) {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao, vbo
}

// UpdateBatch uploads the new vertex data to the GPU
func UpdateBatch(vbo uint32, points []float32) {
	if len(points) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// GL_DYNAMIC_DRAW tells the driver we will change this data often
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.DYNAMIC_DRAW)
}

// DrawBatch performs the single draw call
func DrawBatch(vao uint32, vertexCount int32) {
	if vertexCount == 0 {
		return
	}
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
