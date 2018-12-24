package render

import (
	"log"
	"time"

	"github.com/gooid/gl/es2"
)

type YuvRender struct {
	//width, height                    int
	program                          gl.Program
	yTextureId, uvTextureId          gl.Texture
	yTexture, uvTexture, projMmtxLoc gl.Uniform
	positionLoc, texCoordLoc         gl.Attrib

	BaseRender
}

func (r *YuvRender) Draw(data interface{}) {
	yuv420sp := data.([]byte)
	indices := []uint16{0, 1, 2, 3}

	entry := time.Now()
	boader := r.iW % 8

	r.backupGLState()
	defer r.restoreGLState()

	r.setGL()

	gl.ActiveTexture(gl.TEXTURE0)
	r.yTextureId.Bind(gl.TEXTURE_2D)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.LUMINANCE, r.iW-boader, r.iH, boader, gl.LUMINANCE, gl.UNSIGNED_BYTE, gl.Ptr(yuv420sp))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.ActiveTexture(gl.TEXTURE1)
	r.uvTextureId.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.LUMINANCE_ALPHA, (r.iW-boader)/2, r.iH/2, boader/2, gl.LUMINANCE_ALPHA, gl.UNSIGNED_BYTE, gl.Ptr(yuv420sp[r.iW*r.iH:]))

	r.program.Use()
	r.yTexture.I(0)
	r.uvTexture.I(1)

	r.projMmtxLoc.Matrix4fv(r.OrthoProjection())

	// 绘制四边形
	r.texCoordLoc.Enable()
	r.texCoordLoc.Pointer(2, gl.FLOAT, false, 0, r.TexCoords())

	r.positionLoc.Enable()
	r.positionLoc.Pointer(3, gl.FLOAT, false, 0, r.Vertices())

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_SHORT, gl.Ptr(indices))
	r.yTextureId.Unbind(gl.TEXTURE_2D)
	r.uvTextureId.Unbind(gl.TEXTURE_2D)

	r.program.Unuse()

	if glerr := gl.GetError(); glerr != gl.NO_ERROR {
		log.Println("YuvRender.Draw", ":", glerr)
	}
	log.Println("\tFPS:", 1.0/time.Now().Sub(entry).Seconds())
}

func (r *YuvRender) Validate(w, h int, pixels interface{}) bool {
	if bs, ok := pixels.([]byte); ok {
		const ALIGNBYTES = 8
		// stride
		stride := ALIGNBYTES * ((w + ALIGNBYTES - 1) / ALIGNBYTES)
		return stride*h*3/2 == len(bs)
	}
	return false
}

func (r *YuvRender) Release() {
	if r.program != 0 {
		r.program.Delete()
	}
	if r.yTextureId != 0 {
		r.yTextureId.Delete()
	}
	if r.uvTextureId != 0 {
		r.uvTextureId.Delete()
	}
	r.program = 0
	r.yTextureId = 0
	r.uvTextureId = 0
}

func (r *YuvRender) Init() {
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
		"uniform sampler2D yTexture;\n" +
		"uniform sampler2D uvTexture;\n" +
		"varying vec2 vTexCoord;\n" +
		"void main()\n" +
		"{\n" +
		"	float r, g, b, y, u, v;\n" +
		"	vec4 color;\n" +
		//We had put the Y values of each pixel to the R,G,B components by GL_LUMINANCE,
		//that's why we're pulling it from the R component, we could also use G or B
		"   y = texture2D(yTexture, vTexCoord).r;         \n" +

		//We had put the U and V values of each pixel to the A and R,G,B components of the
		//texture respectively using GL_LUMINANCE_ALPHA. Since U,V bytes are interspread
		//in the texture, this is probably the fastest way to use them in the shader
		"   color = texture2D(uvTexture, vTexCoord);  \n" +
		"   u = color.a - 0.5;  \n" +
		"   v = color.r - 0.5;  \n" +

		//The numbers are just YUV to RGB conversion constants
		"   r = y + 1.13983*v;                              \n" +
		"   g = y - 0.39465*u - 0.58060*v;                  \n" +
		"   b = y + 2.03211*u;                              \n" +

		//We finally set the RGB color of our pixel
		"   gl_FragColor = vec4(r, g, b, 1.0);              \n" +
		"}\n"

	r.program, err = gl.NewProgram([]string{vertex_shader}, []string{fragment_shader})
	if err != nil {
		log.Panicln("CreateProgram.err:", err)
		return
	}

	r.yTexture = r.program.GetUniformLocation("yTexture")
	r.uvTexture = r.program.GetUniformLocation("uvTexture")
	r.projMmtxLoc = r.program.GetUniformLocation("ProjMtx")
	r.positionLoc = r.program.GetAttribLocation("Position")
	r.texCoordLoc = r.program.GetAttribLocation("TexCoord")

	r.yTextureId = gl.CreateTexture()
	r.uvTextureId = gl.CreateTexture()
}
