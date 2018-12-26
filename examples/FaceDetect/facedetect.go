package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/gooid/gocv/opencv3/core"
	"github.com/gooid/gocv/opencv3/objdetect"
	"github.com/gooid/gooid"
	"github.com/gooid/gooid/examples/CameraDemo/render"
	"github.com/gooid/gooid/storage"
)

type faceDetectStat struct {
	cascade *objdetect.CascadeClassifier
	result  []*core.Rect

	lock     int32 // atomic
	rotation int
	ch       chan *core.Mat

	// 统计时间消耗
	tPrepare, sPrepare,
	tDetect, sDetect time.Duration
}

var faceDetect faceDetectStat

func (s *faceDetectStat) LoadOpenCV(act *app.Activity) {
	ls := app.FindMatchLibrary("libopencv_*.so")
	if len(ls) > 0 {
		log.Println("LoadCvLib:", ls[0])
		err := core.LoadCvLib(ls[0])
		if err != nil {
			panic(err)
		}
	}
	log.Println(core.GetBuildInformation())

	s.cascade = loadCascade(act)
	if s.cascade != nil {
		s.ch = make(chan *core.Mat)
		// 用于识别的 goroutine
		go func() {
			for gray := range s.ch {
				s.detect(gray)
			}
		}()
	}
}

// 加载 CascadeClassifier，但因为只有加载文件的接中，
// 因此只能从asset中读出，生成临时文件的方式加载。
func loadCascade(act *app.Activity) *objdetect.CascadeClassifier {
	assetMgr := act.AssetManager()
	altname := "haarcascade_frontalface_alt2.xml"
	tmpPath := filepath.Join("/data/data", act.Context().Package(), "files")
	os.MkdirAll(tmpPath, 0666)
	fw, err := os.Create(filepath.Join(tmpPath, altname))
	if err != nil {
		tmpPath = `/sdcard`
		fw, err = os.Create(filepath.Join(tmpPath, altname))
		if err != nil {
			panic(err)
		}
	}

	r := assetMgr.Open(altname, storage.ASSET_MODE_STREAMING)
	if r == nil {
		panic(fmt.Errorf("open ", altname, " fail"))
	}
	io.Copy(fw, r)
	r.Close()
	fw.Close()

	cascade := objdetect.NewCascadeClassifier2(filepath.Join(tmpPath, altname))
	log.Println("CascadeClassifier:", cascade, tmpPath)
	return cascade
}

// 仅仅适用于 yuv420 数据
// 它的前 stride*h 个数据就是一个灰度图像数据
// 这里是传入 stride，而不是 width，需调用者特别注意。
func (s *faceDetectStat) PrepareYuv(stride, h, r int, d DataUnit) bool {
	if s.cascade != nil && s.ch != nil &&
		atomic.CompareAndSwapInt32(&s.lock, 0, 1) {
		if d.Lock() {
			t1 := time.Now()

			s.rotation = r
			gray := core.NewMat3(h, stride, core.CvTypeCV_8UC1)
			gray.PutB(0, 0, d.Data()[:stride*h])

			t2 := time.Now()
			defer d.Unlock()

			s.tPrepare = t2.Sub(t1)
			s.sPrepare += s.tPrepare
			s.ch <- gray
			return true
		}
	}
	return false
}

func (s *faceDetectStat) detect(gray *core.Mat) {
	t1 := time.Now()
	coreR := core.CoreROTATE_90_CLOCKWISE
	revertR := render.ROTATION270
	switch s.rotation {
	case render.ROTATION0:
	case render.ROTATION90:
		coreR = core.CoreROTATE_90_CLOCKWISE
		revertR = render.ROTATION270

	case render.ROTATION180:
		coreR = core.CoreROTATE_180
		revertR = render.ROTATION180

	case render.ROTATION270:
		coreR = core.CoreROTATE_90_COUNTERCLOCKWISE
		revertR = render.ROTATION90
	}

	if s.rotation != render.ROTATION0 {
		core.Rotate(gray, gray, coreR)
	}

	rv := core.NewMatOfRect()
	s.cascade.DetectMultiScale2(gray, rv)

	rects := rv.ToArray()
	// 因为结果是基于旋转后的图像，因此需把结果坐标转换成基于原始图像的坐标
	if s.rotation != render.ROTATION0 {
		r := render.BaseRender{}
		iw, ih := gray.Cols(), gray.Rows()
		cw, ch := iw, ih
		if revertR == render.ROTATION90 || revertR == render.ROTATION270 {
			cw, ch = ch, cw
		}
		r.SetProperty(cw, ch, iw, ih, 0, 0, cw, ch, revertR)

		s.result = nil
		for _, rc := range rects {
			x1, y1 := r.Pos(rc.X, rc.Y)
			x2, y2 := r.Pos(rc.X+rc.Width, rc.Y+rc.Height)

			s.result = append(s.result,
				core.NewRect3(core.NewPoint(float64(x1), float64(y1)),
					core.NewPoint(float64(x2), float64(y2))))
		}
	} else {
		s.result = rects
	}

	t2 := time.Now()
	s.tDetect = t2.Sub(t1)
	s.sDetect += s.tDetect

	if !atomic.CompareAndSwapInt32(&s.lock, 1, 0) {
		panic(fmt.Errorf("faceDetect lock abnormal"))
	}
	log.Println("DetectYuv:", s.tPrepare, s.tDetect, s.sPrepare, s.sDetect)
}

func (s *faceDetectStat) LastRects() []*core.Rect {
	return s.result
}

func (s *faceDetectStat) Release() {
	close(s.ch)
}
