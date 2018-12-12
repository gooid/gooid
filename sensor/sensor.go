// Copyright 2018 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sensor

import (
	"github.com/gooid/gooid/internal/ndk"
)

type TYPE = app.SENSOR_TYPE

type AdditionalInfoEvent = app.AdditionalInfoEvent
type DynamicSensorEvent = app.DynamicSensorEvent
type HeartRateEvent = app.HeartRateEvent
type MetaDataEvent = app.MetaDataEvent
type UncalibratedEvent = app.UncalibratedEvent
type Vector = app.SensorVector
type MagneticVector = app.MagneticVector

type Sensor = app.Sensor
type Event = app.SensorEvent
type Manager = app.SensorManager

const (
	FIFO_COUNT_INVALID = app.SENSOR_FIFO_COUNT_INVALID
	DELAY_INVALID      = app.SENSOR_DELAY_INVALID

	/**
	 * Invalid sensor type. Returned by {@link ASensor_getType} as error value.
	 */
	TYPE_INVALID = app.SENSOR_TYPE_INVALID
	/**
	 * {@link ASENSOR_TYPE_ACCELEROMETER}
	 * reporting-mode: continuous
	 *
	 *  All values are in SI units (m/s^2) and measure the acceleration of the
	 *  device minus the force of gravity.
	 */
	TYPE_ACCELEROMETER = app.SENSOR_TYPE_ACCELEROMETER
	/**
	 * {@link ASENSOR_TYPE_MAGNETIC_FIELD}
	 * reporting-mode: continuous
	 *
	 *  All values are in micro-Tesla (uT) and measure the geomagnetic
	 *  field in the X, Y and Z axis.
	 */
	TYPE_MAGNETIC_FIELD = app.SENSOR_TYPE_MAGNETIC_FIELD
	/**
	 * {@link ASENSOR_TYPE_GYROSCOPE}
	 * reporting-mode: continuous
	 *
	 *  All values are in radians/second and measure the rate of rotation
	 *  around the X, Y and Z axis.
	 */
	TYPE_GYROSCOPE = app.SENSOR_TYPE_GYROSCOPE
	/**
	 * {@link ASENSOR_TYPE_LIGHT}
	 * reporting-mode: on-change
	 *
	 * The light sensor value is returned in SI lux units.
	 */
	TYPE_LIGHT = app.SENSOR_TYPE_LIGHT
	/**
	 * {@link ASENSOR_TYPE_PROXIMITY}
	 * reporting-mode: on-change
	 *
	 * The proximity sensor which turns the screen off and back on during calls is the
	 * wake-up proximity sensor. Implement wake-up proximity sensor before implementing
	 * a non wake-up proximity sensor. For the wake-up proximity sensor set the flag
	 * SENSOR_FLAG_WAKE_UP.
	 * The value corresponds to the distance to the nearest object in centimeters.
	 */
	TYPE_PROXIMITY = app.SENSOR_TYPE_PROXIMITY
	/**
	 * {@link ASENSOR_TYPE_LINEAR_ACCELERATION}
	 * reporting-mode: continuous
	 *
	 *  All values are in SI units (m/s^2) and measure the acceleration of the
	 *  device not including the force of gravity.
	 */
	TYPE_LINEAR_ACCELERATION = app.SENSOR_TYPE_LINEAR_ACCELERATION

	// java
	TYPE_ORIENTATION                 = app.SENSOR_TYPE_ORIENTATION
	TYPE_PRESSURE                    = app.SENSOR_TYPE_PRESSURE
	TYPE_TEMPERATURE                 = app.SENSOR_TYPE_TEMPERATURE
	TYPE_GRAVITY                     = app.SENSOR_TYPE_GRAVITY
	TYPE_ROTATION_VECTOR             = app.SENSOR_TYPE_ROTATION_VECTOR
	TYPE_RELATIVE_HUMIDITY           = app.SENSOR_TYPE_RELATIVE_HUMIDITY
	TYPE_AMBIENT_TEMPERATURE         = app.SENSOR_TYPE_AMBIENT_TEMPERATURE
	TYPE_MAGNETIC_FIELD_UNCALIBRATED = app.SENSOR_TYPE_MAGNETIC_FIELD_UNCALIBRATED
	TYPE_GAME_ROTATION_VECTOR        = app.SENSOR_TYPE_GAME_ROTATION_VECTOR
	TYPE_GYROSCOPE_UNCALIBRATED      = app.SENSOR_TYPE_GYROSCOPE_UNCALIBRATED
	TYPE_SIGNIFICANT_MOTION          = app.SENSOR_TYPE_SIGNIFICANT_MOTION
	TYPE_STEP_DETECTOR               = app.SENSOR_TYPE_STEP_DETECTOR
	TYPE_STEP_COUNTER                = app.SENSOR_TYPE_STEP_COUNTER
	TYPE_GEOMAGNETIC_ROTATION_VECTOR = app.SENSOR_TYPE_GEOMAGNETIC_ROTATION_VECTOR
	TYPE_HEART_RATE                  = app.SENSOR_TYPE_HEART_RATE
	TYPE_TILT_DETECTOR               = app.SENSOR_TYPE_TILT_DETECTOR
	TYPE_WAKE_GESTURE                = app.SENSOR_TYPE_WAKE_GESTURE
	TYPE_GLANCE_GESTURE              = app.SENSOR_TYPE_GLANCE_GESTURE
	TYPE_PICK_UP_GESTURE             = app.SENSOR_TYPE_PICK_UP_GESTURE
	TYPE_WRIST_TILT_GESTURE          = app.SENSOR_TYPE_WRIST_TILT_GESTURE
	TYPE_DEVICE_ORIENTATION          = app.SENSOR_TYPE_DEVICE_ORIENTATION
	TYPE_POSE_6DOF                   = app.SENSOR_TYPE_POSE_6DOF
	TYPE_STATIONARY_DETECT           = app.SENSOR_TYPE_STATIONARY_DETECT
	TYPE_MOTION_DETECT               = app.SENSOR_TYPE_MOTION_DETECT
	TYPE_HEART_BEAT                  = app.SENSOR_TYPE_HEART_BEAT
	TYPE_DYNAMIC_SENSOR_META         = app.SENSOR_TYPE_DYNAMIC_SENSOR_META
	TYPE_LOW_LATENCY_OFFBODY_DETECT  = app.SENSOR_TYPE_LOW_LATENCY_OFFBODY_DETECT
	TYPE_ACCELEROMETER_UNCALIBRATED  = app.SENSOR_TYPE_ACCELEROMETER_UNCALIBRATED

	/** no contact */
	STATUS_NO_CONTACT = app.SENSOR_STATUS_NO_CONTACT
	/** unreliable */
	STATUS_UNRELIABLE = app.SENSOR_STATUS_UNRELIABLE
	/** low accuracy */
	STATUS_ACCURACY_LOW = app.SENSOR_STATUS_ACCURACY_LOW
	/** medium accuracy */
	STATUS_ACCURACY_MEDIUM = app.SENSOR_STATUS_ACCURACY_MEDIUM
	/** high accuracy */
	STATUS_ACCURACY_HIGH = app.SENSOR_STATUS_ACCURACY_HIGH

	/** invalid reporting mode */
	REPORTING_MODE_INVALID = app.REPORTING_MODE_INVALID
	/** continuous reporting */
	REPORTING_MODE_CONTINUOUS = app.REPORTING_MODE_CONTINUOUS
	/** reporting on change */
	REPORTING_MODE_ON_CHANGE = app.REPORTING_MODE_ON_CHANGE
	/** on shot reporting */
	REPORTING_MODE_ONE_SHOT = app.REPORTING_MODE_ONE_SHOT
	/** special trigger reporting */
	REPORTING_MODE_SPECIAL_TRIGGER = app.REPORTING_MODE_SPECIAL_TRIGGER

	/** stopped */
	DIRECT_RATE_STOP = app.SENSOR_DIRECT_RATE_STOP
	/** nominal 50Hz */
	DIRECT_RATE_NORMAL = app.SENSOR_DIRECT_RATE_NORMAL
	/** nominal 200Hz */
	DIRECT_RATE_FAST = app.SENSOR_DIRECT_RATE_FAST
	/** nominal 800Hz */
	DIRECT_RATE_VERY_FAST = app.SENSOR_DIRECT_RATE_VERY_FAST

	/** shared memory created by ASharedMemory_create */
	DIRECT_CHANNEL_TYPE_SHARED_MEMORY = app.SENSOR_DIRECT_CHANNEL_TYPE_SHARED_MEMORY
	/** AHardwareBuffer */
	DIRECT_CHANNEL_TYPE_HARDWARE_BUFFER = app.SENSOR_DIRECT_CHANNEL_TYPE_HARDWARE_BUFFER

	/** Earth's gravity in m/s^2 */
	STANDARD_GRAVITY = app.SENSOR_STANDARD_GRAVITY
	/** Maximum magnetic field on Earth's surface in uT */
	MAGNETIC_FIELD_EARTH_MAX = app.SENSOR_MAGNETIC_FIELD_EARTH_MAX
	/** Minimum magnetic field on Earth's surface in uT*/
	MAGNETIC_FIELD_EARTH_MIN = app.SENSOR_MAGNETIC_FIELD_EARTH_MIN
)

func ManagerInstance() *Manager {
	return app.SensorManagerInstance()
}
