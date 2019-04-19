// +build android

package main

import (
	"container/list"
	"fmt"
	"log"
	"sync"

	app "github.com/gooid/gooid"
	camera "github.com/gooid/gooid/camera24"
	media "github.com/gooid/gooid/media24"
	"github.com/gooid/util"
)

func winRedraw(act *app.Activity, win *app.Window) {
	manager := CameraManagerInstance()
	ids, status := manager.GetCameraIdList()
	log.Println("GetCameraIdList>>", ids, status)

	metadata, status := manager.GetCameraCharacteristics("0")
	defer metadata.Free()
	log.Println("GetCameraCharacteristics>>", status)
	if status == nil {
		dumpMetadata(metadata)
	}
}

func dumpMetadata(metadata *camera.Metadata) {
	tags, status := metadata.GetAllTags()
	log.Println("  dumpMetadata>>", len(tags), status)
	for _, tag := range tags {
		entry, status := metadata.GetConstEntry(tag)
		if status == nil {
			//log.Println("metadata.GetConstEntry>>", status)
			log.Println("    .GetConstEntry>> {", entry.Tag(), entry.Type(), entry.Count(), entry.Data(), "}", status)
		} else {
			log.Println("    .GetConstEntry>>", tag, status)
		}
	}
}

type ndkCamera struct {
	dev *NDKCamera
	yuvReader,
	jpgReader *media.ImageReader
	imgLock sync.Mutex
	bStart  bool

	// 为了适应 camera2 的 planes，这里允许以多份数据方式传入
	previewCB func(w, h, formatIdx int, ds ...[]byte)
}

func (cam *ndkCamera) OnImage(r *media.ImageReader) {
	cam.imgLock.Lock()
	defer cam.imgLock.Unlock()

	//format, _ := r.GetFormat()
	//log.Println(" ImageListener: ", format)
	image, err := r.AcquireNextImage()
	if err == nil && image != nil {
		defer image.Delete()

		w, _ := image.GetWidth()
		h, _ := image.GetHeight()
		t, _ := image.GetTimestamp()
		format, _ := image.GetFormat()
		planes, _ := image.GetNumberOfPlanes()

		pixelStride, _ := image.GetPlanePixelStride(1)
		rowStride, _ := image.GetPlaneRowStride(1)
		//log.Println("\t", w, h, planes, pixelStride, rowStride, t)
		if format == media.FORMAT_YUV_420_888 {
			dataY, _ := image.GetPlaneData(0)
			dataUV, _ := image.GetPlaneData(2)

			cam.previewCB(w, h, 0, dataY, dataUV)
		} else {
			log.Println(" ImageListener: ", format, w, h, planes, pixelStride, rowStride, t)
			// 保存照片
		}
	} else {
		log.Println("   ImageListener:", err)
	}
}

func (cam *ndkCamera) initCamera(id string) {
	if cam.dev == nil {
		cam.imgLock.Lock()
		defer cam.imgLock.Unlock()

		var err error
		cam.dev, err = CameraManagerInstance().OpenCamera(id)
		if err == nil && cam.dev != nil {
			var err error
			cam.jpgReader, err = media.NewImageReader(1280, 720, media.FORMAT_JPEG, 1)
			util.Assert(err)
			cam.jpgReader.SetImageListener(cam.OnImage)
			win, _ := cam.jpgReader.GetWindow()
			cam.dev.resetTakePhoto(win, 0)
		} else {
			log.Println("OpenCamera:", err)
		}
	}
	log.Printf("initCamera: %+v\n", cam.dev)
}

func (cam *ndkCamera) Close() {
	cam.StopPreview()
	cam.dev.Close()

	cam.imgLock.Lock()
	defer cam.imgLock.Unlock()

	cam.yuvReader.Delete()
	cam.jpgReader.Delete()
	cam.dev = nil
	cam.yuvReader = nil
	cam.jpgReader = nil
}

func (cam *ndkCamera) StartPreview(w, h, formatIdx int, cb func(w, h, formatIdx int, ds ...[]byte)) {
	if !cam.bStart {
		cam.imgLock.Lock()
		defer cam.imgLock.Unlock()

		cam.bStart = true
		cam.dev.releasetPreview()
		if cam.yuvReader != nil {
			cam.yuvReader.Delete()
		}

		var err error
		cam.yuvReader, err = media.NewImageReader(w, h, media.FORMAT_YUV_420_888, 4)
		util.Assert(err)
		cam.yuvReader.SetImageListener(cam.OnImage)
		win, _ := cam.yuvReader.GetWindow()
		cam.dev.resetPreview(win, 0)

		cam.dev.StartPreview(true)
		cam.previewCB = cb
	}
}

func (cam *ndkCamera) StopPreview() {
	if cam.bStart {
		cam.bStart = false
		cam.dev.StartPreview(false)
	}
}

type CameraState struct {
	bAvailable bool
	devices    *list.List
}

// for vcamera

func (cam *ndkCamera) GetIds() []string {
	ids, err := CameraManagerInstance().GetCameraIdList()
	if err != nil {
		return nil
	}
	return ids
}

func (cam *ndkCamera) GetName(id string) string {
	f, _ := CameraManagerInstance().GetSensorOrientation(id)
	return fmt.Sprint(f, "-", id)
}

func (cam *ndkCamera) GetSupportPixels() [][2]int {
	return CameraManagerInstance().GetSupportPixels(cam.dev.device.GetId())
}

func (cam *ndkCamera) GetSupportFormat() []string { return []string{"yuv420sp"} }
func (cam *ndkCamera) GetSensorOrientation() (bool, int) {
	f, r := CameraManagerInstance().GetSensorOrientation(cam.dev.device.GetId())
	return f == camera.LENS_FACING_FRONT, r
}

func (cam *ndkCamera) Capture(w, h, formatIdx int, cb func(w, h, formatIdx int, i []byte)) {
	cam.dev.TakePhoto()
}

func openNDKCamera(idx int) (CameraI, error) {
	cam := &ndkCamera{}
	cam.initCamera(fmt.Sprint(idx))
	return cam, nil
}

func init() {
	registerCamera("Camera2", 1, openNDKCamera)
}
