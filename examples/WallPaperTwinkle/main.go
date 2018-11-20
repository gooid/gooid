package main

import (
	"log"

	"github.com/gooid/gooid"
	camera "github.com/gooid/gooid/camera24"
	"github.com/gooid/gooid/input"
)

func main() {
	context := app.Callbacks{
		Event:              event,
		WindowCreated:      winCreate,
		WindowRedrawNeeded: winRedraw,
	}
	app.SetMainCB(func(ctx *app.Context) {
		ctx.Run(context)
	})
	for app.Loop() {
	}
	log.Println("done.")
}

var bWallpaper = false

func event(act *app.Activity, e *app.InputEvent) {
	if mot := e.Motion(); mot != nil {
		if mot.GetAction() == input.KEY_EVENT_ACTION_UP {
			log.Println("event:", mot)

			bWallpaper = !bWallpaper
			if bWallpaper {
				act.SetWindowFlags(app.FLAG_SHOW_WALLPAPER, 0)
			} else {
				act.SetWindowFlags(0, app.FLAG_SHOW_WALLPAPER)
			}
		}
	}
}

type cbs struct{}

func (cs *cbs) OnCameraAvailable(id string) {
	log.Println("... onCameraAvailable", id)

	metadata, status := manager.GetCameraCharacteristics(id)
	log.Println("GetCameraCharacteristics>>", metadata, status)
}
func (cs *cbs) OnCameraUnavailable(id string) {
	log.Println("... onCameraUnavailable", id)
}

var manager *camera.CameraManager

func winCreate(act *app.Activity, win *app.Window) {
	log.Println("winCreate...")

	manager = camera.ManagerCreate()

	manager.RegisterAvailabilityCallback(&cbs{})
}

func winRedraw(act *app.Activity, win *app.Window) {
	//defer manager.Delete()
	ids, status := manager.GetCameraIdList()
	log.Println("GetCameraIdList>>", ids, status)

	metadata, status := manager.GetCameraCharacteristics("0")
	log.Println("GetCameraCharacteristics>>", metadata, status)
}
