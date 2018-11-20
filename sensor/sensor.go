// Copyright 2018 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.


package sensor

import (
	"github.com/gooid/gooid/internal/ndk"
)

type TYPE = app.SENSOR_TYPE

type HeartRateEvent = app.HeartRateEvent
type MetaDataEvent = app.MetaDataEvent
type SensorEvent = app.SensorEvent
type UncalibratedEvent = app.UncalibratedEvent

type Sensor = app.Sensor
type Event = app.SensorEvent
type Vector = app.SensorVector
type Manager = app.SensorManager

const (
	TYPE_ACCELEROMETER       = app.SENSOR_TYPE_ACCELEROMETER
	TYPE_MAGNETIC_FIELD      = app.SENSOR_TYPE_MAGNETIC_FIELD
	TYPE_GYROSCOPE           = app.SENSOR_TYPE_GYROSCOPE
	TYPE_LIGHT               = app.SENSOR_TYPE_LIGHT
	TYPE_PROXIMITY           = app.SENSOR_TYPE_PROXIMITY
	STATUS_UNRELIABLE        = app.SENSOR_STATUS_UNRELIABLE
	STATUS_ACCURACY_LOW      = app.SENSOR_STATUS_ACCURACY_LOW
	STATUS_ACCURACY_MEDIUM   = app.SENSOR_STATUS_ACCURACY_MEDIUM
	STATUS_ACCURACY_HIGH     = app.SENSOR_STATUS_ACCURACY_HIGH
	STANDARD_GRAVITY         = app.SENSOR_STANDARD_GRAVITY
	MAGNETIC_FIELD_EARTH_MAX = app.SENSOR_MAGNETIC_FIELD_EARTH_MAX
	MAGNETIC_FIELD_EARTH_MIN = app.SENSOR_MAGNETIC_FIELD_EARTH_MIN
)

func ManagerInstance() *Manager {
	return app.SensorManagerInstance()
}
