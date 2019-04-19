// +build !android

package main

func preCreate(ctx interface{}) {

}

func postCreate(ctx interface{}) {
	faceDetect.LoadOpenCV("", "basic/assets/haarcascade_frontalface_alt2.xml")
}
