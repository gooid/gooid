package main

import (
	"math"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gooid/gocv/opencv3/core"
	"github.com/gooid/gocv/opencv3/imgproc"
)

type demoCamera int

func (d demoCamera) GetIds() []string                  { return []string{"0", "1"} }
func (d demoCamera) GetName(id string) string          { return "camera" + id }
func (d demoCamera) GetSupportPixels() [][2]int        { return [][2]int{{640, 480}, {1280, 960}} }
func (d demoCamera) GetSupportFormat() []string        { return []string{"yuv420sp"} }
func (d demoCamera) GetSensorOrientation() (bool, int) { return false, 90 }
func (d demoCamera) Close()                            { d.StopPreview() }

func (d demoCamera) Capture(w, h, formatIdx int, cb func(w, h, formatIdx int, i []byte)) {}

func (d demoCamera) StartPreview(w, h, formatIdx int, cb func(w, h, formatIdx int, ds ...[]byte)) {
	if atomic.AddInt32(&vcamPreview, 0) == 0 {
		if atomic.AddInt32(&vcamPreview, 1) == 1 {
			go func() {
				for {
					select {
					case <-stopCh:
						return

					case <-time.After(25 * time.Millisecond):
						img := core.NewMat3(w, h, core.CvTypeCV_8UC3)
						rimg := core.NewMat2()
						yuv := core.NewMat2()

						// eye angle
						a := float64(time.Now().UnixNano()*int64(time.Nanosecond)) * math.Pi / float64(time.Second)
						drawImage(img, (1+math.Sin(a))/2)
						core.Rotate(img, rimg, core.CoreROTATE_90_COUNTERCLOCKWISE)
						BGR2YUVNV21(rimg, yuv)
						img.Release()
						rimg.Release()

						bs := (*[1 << 30]byte)(unsafe.Pointer(uintptr(yuv.DataAddr())))[:yuv.Width()*yuv.Height()]
						cb(w, h, formatIdx, bs)

						yuv.Release()
					}
				}
			}()
		}
	}
}
func (d demoCamera) StopPreview() {
	if atomic.AddInt32(&vcamPreview, 0) == 1 {
		if atomic.AddInt32(&vcamPreview, -1) == 0 {
			stopCh <- true
		}
	}
}

func _openVCamera(idx int) (CameraI, error) { return demoCamera(0), nil }

const VIRTUALCAMERA = "vcamera"

var vcamPreview int32
var stopCh = make(chan bool, 1)

func init() {
	registerCamera(VIRTUALCAMERA, 1000, _openVCamera)
}

func drawImage(org *core.Mat, eyePos float64) {
	var n float64
	w, h := float64(org.Width()), float64(org.Height())
	if w < h {
		n = w / 100
	} else {
		n = h / 100
	}

	clr := core.NewScalar2(127, 127, 127)
	pt0 := core.NewMatOfPoint4([]*core.Point{core.NewPoint(0, 0),
		core.NewPoint(w-1, 0), core.NewPoint(w-1, h-1), core.NewPoint(0, h-1), core.NewPoint(0, 0)})
	imgproc.FillPoly2(org, []*core.MatOfPoint{pt0}, clr)

	clr = core.NewScalar2(0, 196, 196)
	imgproc.Ellipse2(org, core.NewPoint(n*100/2, n*100/2), core.NewSize(n*100/2-1, n*100/2-1),
		360.0, 0.0, 360.0, clr, int(n))

	clr = core.NewScalar2(0, 0, 0)
	// eye
	imgproc.Ellipse2(org, core.NewPoint(32.0*n, 34*n), core.NewSize(12*n, 16*n),
		360.0, 0.0, 360.0, clr, int(n))
	imgproc.Ellipse2(org, core.NewPoint(32.0*n, (40-12*eyePos)*n), core.NewSize(10*n, 10*n),
		360.0, 0.0, 360.0, clr, int(n))
	imgproc.FloodFill2(org, core.NewMat2(), core.NewPoint(32.0*n, (40-12*eyePos)*n), clr)
	clr = core.NewScalar2(196, 240, 255)
	imgproc.FloodFill2(org, core.NewMat2(), core.NewPoint(21.0*n, 34*n), clr)
	clr = core.NewScalar2(0, 0, 0)

	imgproc.Ellipse2(org, core.NewPoint(68.0*n, 34*n), core.NewSize(12*n, 16*n),
		360.0, 0.0, 360.0, clr, int(n))
	imgproc.Ellipse2(org, core.NewPoint(68.0*n, (40-12*eyePos)*n), core.NewSize(10*n, 10*n),
		360.0, 0.0, 360.0, clr, int(n))
	imgproc.FloodFill2(org, core.NewMat2(), core.NewPoint(68.0*n, (40-12*eyePos)*n), clr)
	clr = core.NewScalar2(196, 240, 255)
	imgproc.FloodFill2(org, core.NewMat2(), core.NewPoint(57.0*n, 34*n), clr)
	clr = core.NewScalar2(0, 0, 0)

	mls := [][2]float64{{26, 60}, {25, 63}, {24, 65}, {20, 67}}
	mrs := [][2]float64{{71, 59}, {72, 61}, {75, 64}, {78, 65}}
	drawLine := func(ps [][2]float64, c *core.Scalar) {
		for i := 1; i < len(ps); i++ {
			imgproc.Line2(org, core.NewPoint(ps[i-1][0]*n, ps[i-1][1]*n), core.NewPoint(ps[i][0]*n, ps[i][1]*n), c, int(n))
		}
	}

	drawLine(mls, clr)
	drawLine(mrs, clr)

	imgproc.Ellipse2(org, core.NewPoint(48.0*n, 10*n), core.NewSize(42*n, 68*n),
		180.0, 232.0, 305.0, clr, int(n))

	clr = core.NewScalar2(0, 255, 255)
	imgproc.FloodFill2(org, core.NewMat2(), core.NewPoint(w/2, h/2), clr)
}

func _YuvI420ToNV21(src *core.Mat, w, h int) {
	// I420 => NV21
	uvI420 := make([]byte, w*h/2)
	src.GetB(h, 0, uvI420)

	uvNV21 := make([]byte, w*h/2)
	n := w * h / 4
	for i, v := range uvI420 {
		x := i % n
		x *= 2
		if i < n {
			x++
		}
		uvNV21[x] = v
	}
	src.PutB(h, 0, uvNV21)
}

func BGR2YUVNV21(src, dst *core.Mat) {
	imgproc.CvtColor2(src, dst, imgproc.COLOR_BGR2YUV_I420)
	_YuvI420ToNV21(dst, src.Cols(), src.Rows())
}
