package render

import (
	"log"
	"time"

	"github.com/gooid/gl/es2"
)

const (
	//Rotation
	ROTATION0    = 0
	ROTATION90   = 1
	ROTATION180  = 2
	ROTATION270  = 3
	ROTATIONMASK = 3

	//Flip Horizontal
	FLIPHOR = 4
	//Flip Vertical
	FLIPVER = 8
)

type Render interface {
	// Init render, userdata is ignore
	Init()
	Draw(pixels interface{})
	Release()

	// SetProperty
	// wW, wH is windows/client width, height
	// iW, iH is image width, height
	// x, y, w, h is draw image to rect
	// op is ROTATION?? | FLIPHOR | FLIPVER
	SetProperty(wW, wH int, iW, iH int, x, y, w, h int, op int)

	// 验证 pixels 是否符合指定 width、height
	Validate(width, height int, pixels interface{}) bool
}

type BaseRender struct {
	wW, wH int32
	iW, iH int32

	x, y, w, h, s float32
	op            int
	restore       func()
}

func (r *BaseRender) SetProperty(wW, wH int, iW, iH int, x, y, rw, rh int, op int) {
	r.wW, r.wH = int32(wW), int32(wH)
	r.iW, r.iH = int32(iW), int32(iH)
	r.op = op

	rotation := r.op & ROTATIONMASK

	iw, ih := r.iW, r.iH
	if rotation == ROTATION90 || rotation == ROTATION270 {
		iw, ih = ih, iw
	}

	var s float32
	if iw*100/int32(rw) > ih*100/int32(rh) {
		s = float32(rw) / float32(iw)
	} else {
		s = float32(rh) / float32(ih)
	}
	r.s = s
	r.x, r.y, r.w, r.h = float32(x), float32(y), float32(iw)*s, float32(ih)*s
}

func (r *BaseRender) Vertices() []float32 {
	return []float32{
		0, 0, 0,
		r.w, 0, 0,
		r.w, r.h, 0,
		0, r.h, 0,
	}
}

// 材质四边形坐标(0-1)
var coords = []float32{
	1.0, 0.0,
	0.0, 0.0,
	0.0, 1.0,
	1.0, 1.0,
	1.0, 0.0,
	0.0, 0.0,
	0.0, 1.0,
}

func (r *BaseRender) TexCoords() []float32 {
	return coords[2*(r.op&ROTATIONMASK):]
}

func (r *BaseRender) OrthoProjection() []float32 {
	var orthoProjection = [16]float32{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, -1.0, 0.0,
		-1.0, -1.0, 0.0, 1.0,
	}

	//	orthoProjection[0], orthoProjection[12] = -2.0/float32(r.wW), 1.0+r.rX/float32(r.wW) // 缩放比,X偏移
	//	if r.op&FLIPVER == FLIPVER {
	//		orthoProjection[0], orthoProjection[12] = 2.0/float32(r.wW), -1.0-2*(float32(r.wW)-r.rW-r.rX)/float32(r.wW) // 垂直翻转
	//	}
	//	orthoProjection[5], orthoProjection[13] = -2.0/float32(r.wH), 1.0-r.rY/float32(r.wH) // 缩放比,Y偏移
	//	if r.op&FLIPHOR == FLIPHOR {
	//		orthoProjection[5], orthoProjection[13] = 2.0/float32(r.wH), -1.0+2*(float32(r.wH)-r.rH-r.rY)/float32(r.wH) // 水平翻转
	//	}

	orthoProjection[0], orthoProjection[12] = -2.0/float32(r.wW), 1.0-2*(float32(r.wW)-r.w-r.x)/float32(r.wW) // 缩放比,X偏移
	if r.op&FLIPVER == FLIPVER {
		orthoProjection[0], orthoProjection[12] = 2.0/float32(r.wW), -1.0+r.x/float32(r.wW) // 垂直翻转
	}
	orthoProjection[5], orthoProjection[13] = -2.0/float32(r.wH), 1.0-r.y/float32(r.wH) // 缩放比,Y偏移
	if r.op&FLIPHOR == FLIPHOR {
		orthoProjection[5], orthoProjection[13] = 2.0/float32(r.wH), -1.0+2*(float32(r.wH)-r.h-r.y)/float32(r.wH) // 水平翻转
	}
	return orthoProjection[:]
}

func (r *BaseRender) backupGLState() {
	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastArrayBuffer int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	//var lastVertexArray int32
	//gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(gl.BLEND)
	lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)
	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Println(" Draw:", glerr)
	}

	r.restore = func() {
		// Restore modified GL state
		gl.UseProgram(uint32(lastProgram))
		gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
		gl.ActiveTexture(uint32(lastActiveTexture))
		//gl.BindVertexArray(uint32(lastVertexArray))
		gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
		gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
		gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
		if lastEnableBlend {
			gl.Enable(gl.BLEND)
		} else {
			gl.Disable(gl.BLEND)
		}
		if lastEnableCullFace {
			gl.Enable(gl.CULL_FACE)
		} else {
			gl.Disable(gl.CULL_FACE)
		}
		if lastEnableDepthTest {
			gl.Enable(gl.DEPTH_TEST)
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
		if lastEnableScissorTest {
			gl.Enable(gl.SCISSOR_TEST)
		} else {
			gl.Disable(gl.SCISSOR_TEST)
		}
		gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
		gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
	}
}

func (r *BaseRender) restoreGLState() {
	if r.restore != nil {
		r.restore()
		r.restore = nil
	}
}

func (r *BaseRender) setGL() {
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Enable(gl.STENCIL_TEST)
}

type C3Render struct {
	// Handle to a program object
	program gl.Program

	// Attribute locations
	positionLoc gl.Attrib
	texCoordLoc gl.Attrib

	// Projection Matrix
	projMtxLoc gl.Uniform

	// Sampler location
	samplerLoc gl.Uniform

	// Texture handle
	textureId gl.Texture

	BaseRender
}

func (r *C3Render) Init(userData interface{}) {
	var err error
	var vertex_shader = "uniform mat4 ProjMtx;\n" +
		"attribute vec3 Position;\n" +
		"attribute vec2 TexCoord;\n" +
		"varying vec2 vTexCoord;\n" +
		"void main()\n" +
		"{\n" +
		"	vTexCoord = TexCoord;\n" +
		"	gl_Position = ProjMtx * vec4(Position.xy,0,1);\n" +
		"}\n"

	var fragment_shader = "precision mediump float;\n" +
		"uniform sampler2D Texture;\n" +
		"varying vec2 vTexCoord;\n" +
		"void main()\n" +
		"{\n" +
		"   vec4 color = texture2D(Texture, vTexCoord);          \n" +
		//We finally set the RGB color of our pixel
		"   gl_FragColor = vec4(color.b, color.g, color.r, 1.0); \n" +
		"}\n"

	r.program, err = gl.NewProgram([]string{vertex_shader}, []string{fragment_shader})
	if !r.program.IsValid() || err != nil {
		log.Panicln("CreateProgram.err:", err)
		return
	}

	r.projMtxLoc = r.program.GetUniformLocation("ProjMtx")
	r.samplerLoc = r.program.GetUniformLocation("Texture")

	r.positionLoc = r.program.GetAttribLocation("Position")
	r.texCoordLoc = r.program.GetAttribLocation("TexCoord")

	// Use tightly packed data
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	// Generate a texture object
	r.textureId = gl.CreateTexture()
}

func (r *C3Render) SetData(pixels interface{}) {
	// Active the texture
	gl.ActiveTexture(gl.TEXTURE0)

	// Bind the texture object
	r.textureId.Bind(gl.TEXTURE_2D)

	// Boader
	boader := r.iW % 8

	// Load the texture
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, r.iW, r.iH, boader, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	r.textureId.Unbind(gl.TEXTURE_2D)
}

func (r *C3Render) Release() {
	r.program.Delete()
	r.textureId.Delete()
}

func (r *C3Render) drawBGR(pixels interface{}, x0, y0, x1, y1 float32, op int) {
	entry := time.Now()

	var indices = []uint16{0, 1, 2, 3}

	r.setGL()

	// Use the program object
	r.program.Use()

	r.projMtxLoc.Matrix4fv(r.OrthoProjection()[:])

	// 绘制四边形
	r.texCoordLoc.Enable()
	r.texCoordLoc.Pointer(2, gl.FLOAT, false, 0, r.TexCoords())

	r.positionLoc.Enable()
	r.positionLoc.Pointer(3, gl.FLOAT, false, 0, r.Vertices())

	// Active the texture
	gl.ActiveTexture(gl.TEXTURE0)

	// Bind the texture object
	r.textureId.Bind(gl.TEXTURE_2D)

	if pixels != nil {
		// Boader
		boader := r.iW % 8

		// Load the texture
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, r.iW, r.iH, boader, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	}

	// Set the sampler texture unit to 0
	r.samplerLoc.I(0)

	// Draw
	gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_SHORT, gl.Ptr(indices))

	r.program.Unuse()
	r.textureId.Unbind(gl.TEXTURE_2D)
	log.Println("\tBGR FPS:", 1.0/time.Now().Sub(entry).Seconds())
}

func (r *C3Render) Draw(pixels interface{}) {
	r.drawBGR(pixels, r.x, r.y, r.w, r.h, r.op)
}
