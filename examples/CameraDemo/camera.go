package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gooid/gooid/camera"
	"github.com/gooid/gooid/examples/CameraDemo/render"
	"github.com/gooid/imgui"
)

type cameraUI struct {
	previewSizes [][2]int
	comboText    string

	resetProp    bool
	previewIndex int
	irender      render.Render
	imgFormat    string
}

type cameraObj struct {
	sync.Mutex
	camera.Camera
	cameraUI

	imageData []byte
	isUpdate  bool
	takeData  int32 // atomic.AddInt32
}

func (cam *cameraObj) setYuv(buffer []byte) {
	if len(cam.imageData) < len(buffer) {
		cam.imageData = make([]byte, len(buffer))
	}
	copy(cam.imageData, buffer)
	cam.isUpdate = true
}

func (cam *cameraObj) Image() ([]byte, bool) {
	cam.Lock()
	defer cam.Unlock()
	return cam.imageData, cam.isUpdate
}

func (cam *cameraObj) Release() {
	if cam != nil && cam.Camera != 0 {
		cam.Disconnect()
	}
	if cam.irender != nil {
		cam.irender.Release()
		cam.irender = nil
	}
}

func cameraInit(id int, usercb func(w, h int, img []byte) bool) *cameraObj {
	cam := &cameraObj{}

	cb := func(buffer []byte) bool {
		// init done
		if cam.takeData > 0 {
			cam.Lock()

			wh := cam.previewSizes[cam.previewIndex]
			if cam.resetProp && cam.irender.Validate(wh[0], wh[1], buffer) {
				cam.ResetProperty()
				cam.resetProp = false
			}

			if !usercb(wh[0], wh[1], buffer) {
				cam.setYuv(buffer)
			}
			cam.Unlock()
			runtime.Gosched()
		} else {
			atomic.AddInt32(&cam.takeData, 1)
		}
		return true
	}

	cam.Camera = camera.Connect(id, cb)
	if cam.Camera == 0 {
		log.Println("Cammera connect fail")
		return nil
	}

	w, h := cam.FrameSize()
	log.Println("camera:", cam, w, h, cam.Fps())
	log.Println("\t", cam.PreviewFormat())
	log.Println("\t", cam.SupportedPreviewSizes())

	cam.cameraUI.Init(cam.Camera)

	// 第一个为默认分辨率
	// 因 camera 还未初始化完，需就异步执行
	go func() {
		runtime.LockOSThread()
		for cam.takeData == 0 {
			time.Sleep(time.Millisecond * 100)
		}
		cam.setPreviewSize(0)
	}()
	return cam
}

func (cam *cameraObj) setPreviewSize(i int) {
	w, h := cam.getPreviewSize(i)
	if cam != nil {
		cam.Lock()
		cam.SetFrameSize(w, h)
		cam.ApplyProperties()
		cam.previewIndex = i
		cam.resetProp = true
		cam.Unlock()
	}
}

func (cam *cameraUI) Init(nativeObj camera.Camera) {
	str := nativeObj.SupportedPreviewSizes()
	for _, s := range strings.Split(str, ",") {
		w, h := 0, 0
		fmt.Sscanf(string(s), "%dx%d", &w, &h)
		cam.previewSizes = append(cam.previewSizes, [2]int{w, h})
	}

	cam.comboText = strings.Replace(str, ",", "\x00", -1) + "\x00"
	cam.imgFormat = nativeObj.PreviewFormat()
}

func (cam *cameraUI) ResetProperty() {
	if cam != nil && cam.irender != nil {
		iw, ih := cam.getPreviewSize(cam.previewIndex)
		cam.irender.SetProperty(width, height, iw, ih, 0, 0, width, height, render.ROTATION90)
	}
}
func (cam *cameraUI) ResetRender() {
	if cam.irender != nil {
		cam.irender.Release()
	}

	switch cam.imgFormat {
	case "yuv420sp":
		cam.irender = &render.YuvRender{}
	default:
		log.Println("not support format:", cam.imgFormat)
		return
	}

	cam.irender.Init()
	cam.ResetProperty()
}

func (cam *cameraUI) getPreviewSize(i int) (w, h int) {
	if i < len(cam.previewSizes) {
		return cam.previewSizes[i][0], cam.previewSizes[i][1]
	}
	return
}

func (cam *cameraObj) Draw() {
	item := int(cam.previewIndex)
	imgui.Text("pixels")
	imgui.SameLine()
	imgui.Combo("", &item, cam.comboText)
	if item != cam.previewIndex {
		cam.setPreviewSize(item)
	}
}

func (cam *cameraObj) DrawImage() {
	if cam != nil {
		data, _ := cam.Image()
		if data != nil && len(data) > 0 {
			cam.irender.Draw(data)
		}
	}
}
