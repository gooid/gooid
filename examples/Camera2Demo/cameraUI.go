package main

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/gooid/gooid/examples/CameraDemo/render"
	"github.com/gooid/imgui"
)

type CameraI interface {
	GetIds() []string
	GetName(id string) string
	GetSupportPixels() [][2]int // [2]int{w, h}
	GetSupportFormat() []string // { "yuv420sp", ... }
	StartPreview(w, h, formatIdx int, cb func(w, h, formatIdx int, ds ...[]byte))
	StopPreview()
	Capture(w, h, formatIdx int, cb func(w, h, formatIdx int, i []byte))
	GetSensorOrientation() (bool, int)
	Close()
}

type cameraIInfo struct {
	o func(int) (CameraI, error)
	n string
	p int // priority; 值越小优先级越高
}

var cameraIs = []*cameraIInfo{}

func registerCamera(dev string, prority int, open func(int) (CameraI, error)) {
	cameraIs = append(cameraIs, &cameraIInfo{o: open, n: dev, p: prority})
}

func openCamera(dev string, id int) CameraI {
	var get func(d string) *cameraIInfo
	get = func(d string) *cameraIInfo {
		var ret *cameraIInfo
		if d == "" {
			p := math.MaxInt32
			for _, i := range cameraIs {
				if i.p <= p {
					ret = i
					p = i.p
				}
			}
		} else {
			for _, i := range cameraIs {
				if i.n == d {
					ret = i
				}
			}
		}
		if ret == nil {
			return get(VIRTUALCAMERA)
		}
		return ret
	}

	ret := get(dev)
	if ret != nil {
		if cam, err := ret.o(id); err == nil {
			return cam
		} else {
			cam, _ = get(VIRTUALCAMERA).o(id)
			return cam
		}
	}
	return nil
}

// 为了避免数据拷贝，防止数据在使用过程中被改写
type DataUnit struct {
	rwlock sync.RWMutex
	data   []byte
}

func (d *DataUnit) Use(fn func(data []byte)) {
	d.rwlock.RLock()
	fn(d.data)
	d.rwlock.RUnlock()
}

func (d *DataUnit) Update(in ...[]byte) {
	d.rwlock.Lock()

	d.data = d.data[:0]
	for _, bs := range in {
		d.data = append(d.data, bs...)
	}

	d.rwlock.Unlock()
}

type propCombo struct {
	idx   int
	title string
	end   string
	items []string
	cmd   func(int)
}

type cameraUI struct {
	dev CameraI

	props []propCombo

	rotation       int
	ww, wh, iw, ih int
	data           DataUnit
	frameRender    render.Render
	wakeup         func()
}

// TODO ???
//func (ui *cameraUI) Init(act *app.Activity) {
//	faceDetect.LoadOpenCV(act)
//}

func (ui *cameraUI) DrawFrame() {
	if ui.frameRender != nil {
		ui.data.Use(func(data []byte) {
			if data != nil && len(data) > 0 {
				ui.frameRender.Draw(data)
			}
		})
	}
}

func (ui *cameraUI) DrawFrameMark(ds imgui.DrawList) {
	if ui.frameRender != nil {
		// 绘制检测到的人脸所在区域
		c := imgui.NewColor(127, 255, 255).U32()
		for _, r := range faceDetect.LastRects() {
			x1, y1 := ui.frameRender.Pos(r.X, r.Y)
			x2, y2 := ui.frameRender.Pos(r.X+r.Width, r.Y+r.Height)
			if x1 > x2 {
				x1, x2 = x2, x1
			}
			if y1 > y2 {
				y1, y2 = y2, y1
			}
			t1 := imgui.NewVec2(x1, y1)
			defer t1.Delete()
			t2 := imgui.NewVec2(x2, y2)
			defer t2.Delete()
			ds.AddRect(t1, t2, c, 20, imgui.DrawCornerFlags_All, 6)
		}
	}
}

func (ui *cameraUI) DrawUI() {
	for i := 0; i < len(ui.props); i++ {
		prop := &ui.props[i]
		if imgui.BeginCombo(prop.title, prop.items[prop.idx]) {
			for n, s := range prop.items {
				isSelected := n == prop.idx
				if imgui.Selectable(s, isSelected) {
					prop.idx = n
					log.Println("  ", prop.title, n, prop.items[n])
					if prop.cmd != nil {
						prop.cmd(prop.idx)
					}
				}
				if isSelected {
					imgui.SetItemDefaultFocus()
				}
			}
			imgui.EndCombo()
		}
		imgui.Separator()
	}

	if imgui.Button("Take") {
		ui.dev.Capture(1280, 720, 1, func(w, h, formatIdx int, i []byte) {
			log.Println("OnTake:", w, h, formatIdx, len(i))
		})
	}
}

func (ui *cameraUI) ResetProperty() {
	log.Println("ResetProperty:", ui.iw, ui.ih)
	ui.frameRender.SetProperty(ui.ww, ui.wh, ui.iw, ui.ih, 0, 0, ui.ww, ui.wh, ui.rotation)
}

func (ui *cameraUI) startPreview(w, h int) {
	ui.dev.StopPreview()
	ui.dev.StartPreview(w, h, -1,
		func(w, h, f int, data ...[]byte) {
			if w != ui.iw || h != ui.ih {
				ui.iw, ui.ih = w, h
				ui.data.Update([]byte{}) // disable draw frame
				ui.ResetProperty()
			}
			ui.data.Update(data...)

			go func() {
				faceDetect.PrepareYuv(w, h, ui.rotation, ui.data)
			}()

			if ui.wakeup != nil {
				ui.wakeup()
			}
		})
}

func (ui *cameraUI) SetCamera(cam CameraI) {
	ui.dev = cam
	ui.props = nil

	ids := ui.dev.GetIds()
	ns := []string{}
	for _, id := range ids {
		ns = append(ns, ui.dev.GetName(id))
	}
	ui.props = append(ui.props, propCombo{
		title: "camera",
		items: ns,
		cmd: func(idx int) {
			if ui.dev != nil {
				ui.dev.Close()
			}

			// disable draw frame
			ui.frameRender = nil
			ui.data.Update([]byte{})

			ui.SetCamera(openCamera("", idx))
		},
	})

	items := []string{}
	pixels := ui.dev.GetSupportPixels()
	for _, wh := range pixels {
		items = append(items, fmt.Sprint(wh[0], "x", wh[1]))
	}

	ui.props = append(ui.props, propCombo{
		title: "pixels",
		items: items,
		cmd: func(idx int) {
			ui.startPreview(pixels[idx][0], pixels[idx][1])
		},
	})

	_, d := ui.dev.GetSensorOrientation()
	ui.rotation = d / 90
	ui.props = append(ui.props, propCombo{
		idx:   d / 90,
		title: "degree",
		items: []string{"0", "90", "180", "270"},
		cmd: func(idx int) {
			ui.rotation = idx
			ui.ResetProperty()
		},
	})

	ui.iw, ui.ih = pixels[0][0], pixels[0][1]
	ui.frameRender = &render.YuvRender{}
	ui.frameRender.Init()
	ui.ResetProperty()
	ui.startPreview(ui.iw, ui.ih)
}

func (ui *cameraUI) SetWindow(ww, wh int) {
	ui.ww, ui.wh = ww, wh
}

func (ui *cameraUI) SetWakeup(fn func()) {
	ui.wakeup = fn
}

func (ui *cameraUI) Close() {
	if ui.dev != nil {
		ui.dev.Close()
		ui.dev = nil

		ui.props = nil
	}
}

var gCUI *cameraUI

func CameraUI() *cameraUI {
	if gCUI == nil {
		gCUI = &cameraUI{}
	}
	return gCUI
}
