package main

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/gooid/gooid"
	"github.com/gooid/gooid/sensor"
)

type sensorInfo struct {
	Type         sensor.TYPE
	Name, Vender string
	Resolution   float32
	MinDelay     time.Duration
	sensor       *sensor.Sensor
}

func SensorInfo(act *app.Activity) []sensorInfo {
	sensorManager := sensor.ManagerInstance()
	sensorList := sensorManager.GetSensorList()
	log.Println("SensorInfo:", len(sensorList))
	sis := make([]sensorInfo, len(sensorList))
	for i, sensor := range sensorList {
		sis[i].Name = sensor.GetName()
		sis[i].Vender = sensor.GetVendor()
		sis[i].Type = sensor.GetType()
		sis[i].Resolution = sensor.GetResolution()
		sis[i].MinDelay = sensor.GetMinDelay()
		sis[i].sensor = sensor
	}
	// 使 sensor 有效
	//act.Context().EnableSensor(sensorManager.GetDefaultSensor(sensor.TYPE_ACCELEROMETER))
	//act.Context().EnableSensor(sensorManager.GetDefaultSensor(sensor.TYPE_GYROSCOPE))
	//act.Context().EnableSensor(sensorManager.GetDefaultSensor(sensor.TYPE_GRAVITY))
	return sis
}

var enableSensors = map[*sensor.Sensor]bool{}

func EnableSensor(act *app.Activity, s *sensor.Sensor, t time.Duration) {
	if b, _ := enableSensors[s]; !b {
		s.Enable(act)
		s.SetEventRate(act, t)
		enableSensors[s] = true
	}
}

func DisableSensor(act *app.Activity, s *sensor.Sensor) {
	if b, _ := enableSensors[s]; b {
		s.Disable(act)
		enableSensors[s] = false
	}
}

func SetEventRate(act *app.Activity, s *sensor.Sensor, t time.Duration) {
	if b, _ := enableSensors[s]; b {
		s.SetEventRate(act, t)
	}
}

const threshold = float64(sensor.STANDARD_GRAVITY / 3)

var lastVec sensor.Vector
var vertical = float32(math.Sqrt(sensor.STANDARD_GRAVITY * sensor.STANDARD_GRAVITY / 2))
var gRotation = ROTATION0

var eventMap = map[sensor.TYPE]map[int]sensor.Event{}

func sensorEevent(act *app.Activity, events []sensor.Event) {
	n := len(events) - 1
	log.Println("sensor:", events[n].GetType(), events[n].GetTimestamp())
	var vec sensor.Vector
	for _, event := range events {
		if _, ok := eventMap[event.GetType()]; !ok {
			eventMap[event.GetType()] = map[int]sensor.Event{}
		}
		eventMap[event.GetType()][event.GetSensor()] = event

		switch event.GetType() {
		case sensor.TYPE_ACCELEROMETER:
			event.GetData(&vec)
			if math.Abs(float64(lastVec.X-vec.X)) > threshold ||
				math.Abs(float64(lastVec.Y-vec.Y)) > threshold ||
				math.Abs(float64(lastVec.Z-vec.Z)) > threshold {

				lastVec = vec

				if isVertical(&vec) {
					if math.Abs(float64(vec.X)) > math.Abs(float64(vec.Y)) {
						if vec.X > 0 {
							gRotation = ROTATION270
						} else {
							gRotation = ROTATION90
						}
					} else {
						if vec.Y > 0 {
							gRotation = ROTATION0
						} else {
							gRotation = ROTATION180
						}
					}
				}

				log.Printf("Rotation: %v %+v\n", gRotation, vec)
			}
		}
	}
}

func isVertical(vec *sensor.Vector) bool {
	if vec.Z >= -vertical &&
		vec.Z <= vertical {
		return true
	}
	return false
}

type Rotation int

const (
	ROTATION0 Rotation = iota
	ROTATION90
	ROTATION180
	ROTATION270
)

func (r Rotation) String() string {
	return []string{"0°", "90°", "180°", "270°"}[int(r)]
}

// TYPE 和 值的对应关系
// 参考 https://developer.android.com/guide/topics/sensors/sensors_overview
func eventString(event sensor.Event) string {
	var value interface{}
	switch event.GetType() {
	case sensor.TYPE_ACCELEROMETER, sensor.TYPE_GRAVITY, sensor.TYPE_GYROSCOPE,
		sensor.TYPE_LINEAR_ACCELERATION, sensor.TYPE_MAGNETIC_FIELD:
		var v sensor.Vector
		value = &v

	case sensor.TYPE_ORIENTATION:
		var v sensor.MagneticVector
		value = &v

	case sensor.TYPE_AMBIENT_TEMPERATURE, sensor.TYPE_LIGHT, sensor.TYPE_PRESSURE:
		var v float32
		value = &v

	case sensor.TYPE_STEP_COUNTER:
		var v int64
		value = &v

	default:
		var v [4]float32
		value = v[:]
	}
	event.GetData(value)
	return fmt.Sprintf("%v: %v\n t=%v v=%+v", event.GetType(), event.GetSensor(), event.GetTimestamp(), reflect.Indirect(reflect.ValueOf(value)))
}
