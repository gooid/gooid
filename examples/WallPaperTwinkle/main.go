package main

import (
	"log"

	"github.com/gooid/gooid"
	"github.com/gooid/gooid/input"
)

func main() {
	context := app.Callbacks{
		Event:              event,
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
