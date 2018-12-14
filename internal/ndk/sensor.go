package app

/*
 * Copyright (C) 2010 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
#include <android/sensor.h>
extern int cgoCallback(int fd, int events, void* data);
static void* ASensorEvent_getDatas(ASensorEvent* event) {
	return &(event->data[0]);
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"
	"unsafe"
)

/******************************************************************
 *
 * IMPORTANT NOTICE:
 *
 *   This file is part of Android's set of stable system headers
 *   exposed by the Android NDK (Native Development Kit).
 *
 *   Third-party source AND binary code relies on the definitions
 *   here to be FROZEN ON ALL UPCOMING PLATFORM RELEASES.
 *
 *   - DO NOT MODIFY ENUMS (EXCEPT IF YOU ADD NEW 32-BIT VALUES)
 *   - DO NOT MODIFY CONSTANTS OR FUNCTIONAL MACROS
 *   - DO NOT CHANGE THE SIGNATURE OF FUNCTIONS IN ANY WAY
 *   - DO NOT CHANGE THE LAYOUT OR SIZE OF STRUCTURES
 */

/**
 * Structures and functions to receive and process sensor events in
 * native code.
 *
 */

type HardwareBuffer C.AHardwareBuffer

//const RESOLUTION_INVALID = C.ASENSOR_RESOLUTION_INVALID //(nanf(""))
const SENSOR_FIFO_COUNT_INVALID = C.ASENSOR_FIFO_COUNT_INVALID //(-1)
const SENSOR_DELAY_INVALID = C.ASENSOR_DELAY_INVALID           //INT32_MIN

/**
 * Sensor types.
 * (keep in sync with hardware/sensors.h)
 */
type SENSOR_TYPE int32

const (
	/**
	 * Invalid sensor type. Returned by {@link ASensor_getType} as error value.
	 */
	SENSOR_TYPE_INVALID = SENSOR_TYPE(C.ASENSOR_TYPE_INVALID)
	/**
	 * {@link ASENSOR_TYPE_ACCELEROMETER}
	 * reporting-mode: continuous
	 *
	 *  All values are in SI units (m/s^2) and measure the acceleration of the
	 *  device minus the force of gravity.
	 */
	SENSOR_TYPE_ACCELEROMETER = SENSOR_TYPE(C.ASENSOR_TYPE_ACCELEROMETER)
	/**
	 * {@link ASENSOR_TYPE_MAGNETIC_FIELD}
	 * reporting-mode: continuous
	 *
	 *  All values are in micro-Tesla (uT) and measure the geomagnetic
	 *  field in the X, Y and Z axis.
	 */
	SENSOR_TYPE_MAGNETIC_FIELD = SENSOR_TYPE(C.ASENSOR_TYPE_MAGNETIC_FIELD)
	/**
	 * {@link ASENSOR_TYPE_GYROSCOPE}
	 * reporting-mode: continuous
	 *
	 *  All values are in radians/second and measure the rate of rotation
	 *  around the X, Y and Z axis.
	 */
	SENSOR_TYPE_GYROSCOPE = SENSOR_TYPE(C.ASENSOR_TYPE_GYROSCOPE)
	/**
	 * {@link ASENSOR_TYPE_LIGHT}
	 * reporting-mode: on-change
	 *
	 * The light sensor value is returned in SI lux units.
	 */
	SENSOR_TYPE_LIGHT = SENSOR_TYPE(C.ASENSOR_TYPE_LIGHT)
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
	SENSOR_TYPE_PROXIMITY = SENSOR_TYPE(C.ASENSOR_TYPE_PROXIMITY)
	/**
	 * {@link ASENSOR_TYPE_LINEAR_ACCELERATION}
	 * reporting-mode: continuous
	 *
	 *  All values are in SI units (m/s^2) and measure the acceleration of the
	 *  device not including the force of gravity.
	 */
	SENSOR_TYPE_LINEAR_ACCELERATION = SENSOR_TYPE(C.ASENSOR_TYPE_LINEAR_ACCELERATION)

	// java
	SENSOR_TYPE_ORIENTATION                 = SENSOR_TYPE(3)
	SENSOR_TYPE_PRESSURE                    = SENSOR_TYPE(6)
	SENSOR_TYPE_TEMPERATURE                 = SENSOR_TYPE(7)
	SENSOR_TYPE_GRAVITY                     = SENSOR_TYPE(9)
	SENSOR_TYPE_ROTATION_VECTOR             = SENSOR_TYPE(11)
	SENSOR_TYPE_RELATIVE_HUMIDITY           = SENSOR_TYPE(12)
	SENSOR_TYPE_AMBIENT_TEMPERATURE         = SENSOR_TYPE(13)
	SENSOR_TYPE_MAGNETIC_FIELD_UNCALIBRATED = SENSOR_TYPE(14)
	SENSOR_TYPE_GAME_ROTATION_VECTOR        = SENSOR_TYPE(15)
	SENSOR_TYPE_GYROSCOPE_UNCALIBRATED      = SENSOR_TYPE(16)
	SENSOR_TYPE_SIGNIFICANT_MOTION          = SENSOR_TYPE(17)
	SENSOR_TYPE_STEP_DETECTOR               = SENSOR_TYPE(18)
	SENSOR_TYPE_STEP_COUNTER                = SENSOR_TYPE(19)
	SENSOR_TYPE_GEOMAGNETIC_ROTATION_VECTOR = SENSOR_TYPE(20)
	SENSOR_TYPE_HEART_RATE                  = SENSOR_TYPE(21)
	SENSOR_TYPE_TILT_DETECTOR               = SENSOR_TYPE(22)
	SENSOR_TYPE_WAKE_GESTURE                = SENSOR_TYPE(23)
	SENSOR_TYPE_GLANCE_GESTURE              = SENSOR_TYPE(24)
	SENSOR_TYPE_PICK_UP_GESTURE             = SENSOR_TYPE(25)
	SENSOR_TYPE_WRIST_TILT_GESTURE          = SENSOR_TYPE(26)
	SENSOR_TYPE_DEVICE_ORIENTATION          = SENSOR_TYPE(27)
	SENSOR_TYPE_POSE_6DOF                   = SENSOR_TYPE(28)
	SENSOR_TYPE_STATIONARY_DETECT           = SENSOR_TYPE(29)
	SENSOR_TYPE_MOTION_DETECT               = SENSOR_TYPE(30)
	SENSOR_TYPE_HEART_BEAT                  = SENSOR_TYPE(31)
	SENSOR_TYPE_DYNAMIC_SENSOR_META         = SENSOR_TYPE(32)
	SENSOR_TYPE_LOW_LATENCY_OFFBODY_DETECT  = SENSOR_TYPE(34)
	SENSOR_TYPE_ACCELEROMETER_UNCALIBRATED  = SENSOR_TYPE(35)
)

/**
 * Sensor accuracy measure.
 */
const (
	/** no contact */
	SENSOR_STATUS_NO_CONTACT = C.ASENSOR_STATUS_NO_CONTACT
	/** unreliable */
	SENSOR_STATUS_UNRELIABLE = C.ASENSOR_STATUS_UNRELIABLE
	/** low accuracy */
	SENSOR_STATUS_ACCURACY_LOW = C.ASENSOR_STATUS_ACCURACY_LOW
	/** medium accuracy */
	SENSOR_STATUS_ACCURACY_MEDIUM = C.ASENSOR_STATUS_ACCURACY_MEDIUM
	/** high accuracy */
	SENSOR_STATUS_ACCURACY_HIGH = C.ASENSOR_STATUS_ACCURACY_HIGH
)

/**
 * Sensor Reporting Modes.
 */
const (
	/** invalid reporting mode */
	REPORTING_MODE_INVALID = C.AREPORTING_MODE_INVALID
	/** continuous reporting */
	REPORTING_MODE_CONTINUOUS = C.AREPORTING_MODE_CONTINUOUS
	/** reporting on change */
	REPORTING_MODE_ON_CHANGE = C.AREPORTING_MODE_ON_CHANGE
	/** on shot reporting */
	REPORTING_MODE_ONE_SHOT = C.AREPORTING_MODE_ONE_SHOT
	/** special trigger reporting */
	REPORTING_MODE_SPECIAL_TRIGGER = C.AREPORTING_MODE_SPECIAL_TRIGGER
)

/**
 * Sensor Direct Report Rates.
 */
const (
	/** stopped */
	SENSOR_DIRECT_RATE_STOP = C.ASENSOR_DIRECT_RATE_STOP
	/** nominal 50Hz */
	SENSOR_DIRECT_RATE_NORMAL = C.ASENSOR_DIRECT_RATE_NORMAL
	/** nominal 200Hz */
	SENSOR_DIRECT_RATE_FAST = C.ASENSOR_DIRECT_RATE_FAST
	/** nominal 800Hz */
	SENSOR_DIRECT_RATE_VERY_FAST = C.ASENSOR_DIRECT_RATE_VERY_FAST
)

/**
 * Sensor Direct Channel Type.
 */
const (
	/** shared memory created by ASharedMemory_create */
	SENSOR_DIRECT_CHANNEL_TYPE_SHARED_MEMORY = C.ASENSOR_DIRECT_CHANNEL_TYPE_SHARED_MEMORY
	/** AHardwareBuffer */
	SENSOR_DIRECT_CHANNEL_TYPE_HARDWARE_BUFFER = C.ASENSOR_DIRECT_CHANNEL_TYPE_HARDWARE_BUFFER
)

/*
 * A few useful constants
 */

/** Earth's gravity in m/s^2 */
const SENSOR_STANDARD_GRAVITY = C.ASENSOR_STANDARD_GRAVITY //(9.80665f)
/** Maximum magnetic field on Earth's surface in uT */
const SENSOR_MAGNETIC_FIELD_EARTH_MAX = C.ASENSOR_MAGNETIC_FIELD_EARTH_MAX //(60.0f)
/** Minimum magnetic field on Earth's surface in uT*/
const SENSOR_MAGNETIC_FIELD_EARTH_MIN = C.ASENSOR_MAGNETIC_FIELD_EARTH_MIN //(30.0f)

/**
 * A sensor event.
 */

/* NOTE: Must match hardware/sensors.h */
type SensorVector struct {
	X, Y, Z float32
	Status  int8
}
type MagneticVector struct {
	Azimuth, Pitch, Roll float32
	Status               int8
}

type MetaDataEvent struct {
	What, Sensor int32
}

type UncalibratedEvent struct {
	XUncalib, YUncalib, ZUncalib float32
	XBias, YBias, ZBias          float32
}

type HeartRateEvent struct {
	Bpm    float32
	Status int8
}

type DynamicSensorEvent struct {
	Connected, Handle int32
}

type AdditionalInfoEvent struct {
	Type, Serial int32
	dataInt32    [14]int32
	//DataInt32 [14]int32
	//DataFloat [14]float32
}

/* NOTE: Must match hardware/sensors.h */
type SensorEvent C.ASensorEvent

/*
typedef struct ASensorEvent {
    int32_t version; // sizeof(struct ASensorEvent)
    int32_t sensor;
    int32_t type;
    int32_t reserved0;
    int64_t timestamp;
    union {
        union {
            float           data[16];
            ASensorVector   vector;
            ASensorVector   acceleration;
            ASensorVector   magnetic;
            float           temperature;
            float           distance;
            float           light;
            float           pressure;
            float           relative_humidity;
            AUncalibratedEvent uncalibrated_gyro;
            AUncalibratedEvent uncalibrated_magnetic;
            AMetaDataEvent meta_data;
            AHeartRateEvent heart_rate;
            ADynamicSensorEvent dynamic_sensor_meta;
            AAdditionalInfoEvent additional_info;
        };
        union {
            uint64_t        data[8];
            uint64_t        step_counter;
        } u64;
    };

    uint32_t flags;
    int32_t reserved1[3];
} ASensorEvent;
*/

/**
 * {@link ASensorManager} is an opaque type to manage sensors and
 * events queues.
 *
 * {@link ASensorManager} is a singleton that can be obtained using
 * ASensorManager_getInstance().
 *
 * This file provides a set of functions that uses {@link
 * ASensorManager} to access and list hardware sensors, and
 * create and destroy event queues:
 * - ASensorManager_getSensorList()
 * - ASensorManager_getDefaultSensor()
 * - ASensorManager_getDefaultSensorEx()
 * - ASensorManager_createEventQueue()
 * - ASensorManager_destroyEventQueue()
 */
type SensorManager C.ASensorManager

func (manager *SensorManager) cptr() *C.ASensorManager {
	return (*C.ASensorManager)(manager)
}

/**
 * {@link ASensorEventQueue} is an opaque type that provides access to
 * {@link ASensorEvent} from hardware sensors.
 *
 * A new {@link ASensorEventQueue} can be obtained using ASensorManager_createEventQueue().
 *
 * This file provides a set of functions to enable and disable
 * sensors, check and get events, and set event rates on a {@link
 * ASensorEventQueue}.
 * - ASensorEventQueue_enableSensor()
 * - ASensorEventQueue_disableSensor()
 * - ASensorEventQueue_hasEvents()
 * - ASensorEventQueue_getEvents()
 * - ASensorEventQueue_setEventRate()
 */
type SensorEventQueue C.ASensorEventQueue

func (queue *SensorEventQueue) cptr() *C.ASensorEventQueue {
	return (*C.ASensorEventQueue)(queue)
}

/**
 * {@link ASensor} is an opaque type that provides information about
 * an hardware sensors.
 *
 * A {@link ASensor} pointer can be obtained using
 * ASensorManager_getDefaultSensor(),
 * ASensorManager_getDefaultSensorEx() or from a {@link ASensorList}.
 *
 * This file provides a set of functions to access properties of a
 * {@link ASensor}:
 * - ASensor_getName()
 * - ASensor_getVendor()
 * - ASensor_getType()
 * - ASensor_getResolution()
 * - ASensor_getMinDelay()
 * - ASensor_getFifoMaxEventCount()
 * - ASensor_getFifoReservedEventCount()
 * - ASensor_getStringType()
 * - ASensor_getReportingMode()
 * - ASensor_isWakeUpSensor()
 */
type Sensor C.ASensor

func (sensor *Sensor) cptr() *C.ASensor {
	return (*C.ASensor)(sensor)
}

/**
 * {@link ASensorRef} is a type for constant pointers to {@link ASensor}.
 *
 * This is used to define entry in {@link ASensorList} arrays.
 */
//typedef ASensor const* ASensorRef;
/**
 * {@link ASensorList} is an array of reference to {@link ASensor}.
 *
 * A {@link ASensorList} can be initialized using ASensorManager_getSensorList().
 */
//typedef ASensorRef const* ASensorList;

/*****************************************************************************/

/**
 * Get a reference to the sensor manager. ASensorManager is a singleton
 * per package as different packages may have access to different sensors.
 *
 * Deprecated: Use ASensorManager_getInstanceForPackage(const char*) instead.
 *
 * Example:
 *
 *     ASensorManager* sensorManager = ASensorManager_getInstance();
 *
 */
//#if __ANDROID_API__ >= __ANDROID_API_O__
//__attribute__ ((deprecated)) ASensorManager* ASensorManager_getInstance();
//#endif
func SensorManagerInstance() *SensorManager {
	return (*SensorManager)(C.ASensorManager_getInstance())
}

/**
 * Returns the list of available sensors.
 */
//int ASensorManager_getSensorList(ASensorManager* manager, ASensorList* list);
func (manager *SensorManager) GetSensorList() []*Sensor {
	var list C.ASensorList
	n := C.ASensorManager_getSensorList(manager.cptr(), &list)
	if n > 0 {
		return (*[1 << 26]*Sensor)(unsafe.Pointer(list))[:n]
	}
	return nil
}

/**
 * Returns the default sensor for the given type, or NULL if no sensor
 * of that type exists.
 */
//ASensor const* ASensorManager_getDefaultSensor(ASensorManager* manager, int type);
func (manager *SensorManager) GetDefaultSensor(typ SENSOR_TYPE) *Sensor {
	return (*Sensor)(C.ASensorManager_getDefaultSensor(manager.cptr(), C.int(typ)))
}

/**
 * Creates a new sensor event queue and associate it with a looper.
 *
 * "ident" is a identifier for the events that will be returned when
 * calling ALooper_pollOnce(). The identifier must be >= 0, or
 * ALOOPER_POLL_CALLBACK if providing a non-NULL callback.
 */
//ASensorEventQueue* ASensorManager_createEventQueue(ASensorManager* manager,
//        ALooper* looper, int ident, ALooper_callbackFunc callback, void* data);
func (manager *SensorManager) createEventQueue(looper *Looper, ident int, callback LooperCallback, data unsafe.Pointer) *SensorEventQueue {
	if callback != nil {
		return (*SensorEventQueue)(C.ASensorManager_createEventQueue(manager.cptr(), looper.cptr(),
			C.int(ident), (C.ALooper_callbackFunc)(C.cgoCallback),
			unsafe.Pointer(&callbackParam{callback, data})))
	} else {
		return (*SensorEventQueue)(C.ASensorManager_createEventQueue(manager.cptr(), looper.cptr(),
			C.int(ident), nil, data))
	}
}

/**
 * Destroys the event queue and free all resources associated to it.
 */
//int ASensorManager_destroyEventQueue(ASensorManager* manager, ASensorEventQueue* queue);
func (manager *SensorManager) destroy(queue *SensorEventQueue) int {
	return int(C.ASensorManager_destroyEventQueue(manager.cptr(), queue.cptr()))
}

/*****************************************************************************/

/**
 * Enable the selected sensor with sampling and report parameters
 *
 * Enable the selected sensor at a specified sampling period and max batch report latency.
 * To disable  sensor, use {@link ASensorEventQueue_disableSensor}.
 *
 * \param queue {@link ASensorEventQueue} for sensor event to be report to.
 * \param sensor {@link ASensor} to be enabled.
 * \param samplingPeriodUs sampling period of sensor in microseconds.
 * \param maxBatchReportLatencyus maximum time interval between two batch of sensor events are
 *                                delievered in microseconds. For sensor streaming, set to 0.
 * \return 0 on success or a negative error code on failure.
 */
//int ASensorEventQueue_registerSensor(ASensorEventQueue* queue, ASensor const* sensor,
//        int32_t samplingPeriodUs, int64_t maxBatchReportLatencyUs);
/*func (queue *SensorEventQueue) registerSensor(sensor *Sensor,
	samplingPeriodUs int32, maxBatchReportLatencyUs int64) int {
	return int(C.ASensorEventQueue_registerSensor(queue.cptr(), sensor.cptr(),
		C.int32_t(samplingPeriodUs), C.int64_t(maxBatchReportLatencyUs)))
}*/

/**
 * Enable the selected sensor at default sampling rate.
 *
 * Start event reports of a sensor to specified sensor event queue at a default rate.
 *
 * \param queue {@link ASensorEventQueue} for sensor event to be report to.
 * \param sensor {@link ASensor} to be enabled.
 *
 * \return 0 on success or a negative error code on failure.
 */
//int ASensorEventQueue_enableSensor(ASensorEventQueue* queue, ASensor const* sensor);
func (queue *SensorEventQueue) enableSensor(sensor *Sensor) int {
	return int(C.ASensorEventQueue_enableSensor(queue.cptr(), sensor.cptr()))
}

/**
 * Disable the selected sensor.
 *
 * Stop event reports from the sensor to specified sensor event queue.
 *
 * \param queue {@link ASensorEventQueue} to be changed
 * \param sensor {@link ASensor} to be disabled
 * \return 0 on success or a negative error code on failure.
 */
//int ASensorEventQueue_disableSensor(ASensorEventQueue* queue, ASensor const* sensor);
func (queue *SensorEventQueue) disableSensor(sensor *Sensor) int {
	return int(C.ASensorEventQueue_disableSensor(queue.cptr(), sensor.cptr()))
}

/**
 * Sets the delivery rate of events in microseconds for the given sensor.
 *
 * This function has to be called after {@link ASensorEventQueue_enableSensor}.
 * Note that this is a hint only, generally event will arrive at a higher
 * rate. It is an error to set a rate inferior to the value returned by
 * ASensor_getMinDelay().
 *
 * \param queue {@link ASensorEventQueue} to which sensor event is delivered.
 * \param sensor {@link ASensor} of which sampling rate to be updated.
 * \param usec sensor sampling period (1/sampling rate) in microseconds
 * \return 0 on sucess or a negative error code on failure.
 */
//int ASensorEventQueue_setEventRate(ASensorEventQueue* queue, ASensor const* sensor, int32_t usec);
func (queue *SensorEventQueue) setEventRate(sensor *Sensor, t time.Duration) int {
	usec := t / time.Microsecond
	if usec >= math.MaxInt32 {
		usec = math.MaxInt32
	}
	return int(C.ASensorEventQueue_setEventRate(queue.cptr(), sensor.cptr(), C.int32_t(usec)))
}

/**
 * Determine if a sensor event queue has pending event to be processed.
 *
 * \param queue {@link ASensorEventQueue} to be queried
 * \return 1 if the queue has events; 0 if it does not have events;
 *         or a negative value if there is an error.
 */
//int ASensorEventQueue_hasEvents(ASensorEventQueue* queue);
func (queue *SensorEventQueue) hasEvents() int {
	return int(C.ASensorEventQueue_hasEvents(queue.cptr()))
}

/**
 * Retrieve pending events in sensor event queue
 *
 * Retrieve next available events from the queue to a specified event array.
 *
 * \param queue {@link ASensorEventQueue} to get events from
 * \param events pointer to an array of {@link ASensorEvents}.
 * \param count max number of event that can be filled into array event.
 * \return number of events returned on success; negative error code when
 *         no events are pending or an error has occurred.
 *
 * Examples:
 *
 *     ASensorEvent event;
 *     ssize_t numEvent = ASensorEventQueue_getEvents(queue, &event, 1);
 *
 *     ASensorEvent eventBuffer[8];
 *     ssize_t numEvent = ASensorEventQueue_getEvents(queue, eventBuffer, 8);
 *
 */
//ssize_t ASensorEventQueue_getEvents(ASensorEventQueue* queue,
//                ASensorEvent* events, size_t count);
func (queue *SensorEventQueue) getEvents(count int) ([]SensorEvent, int) {
	if count == 0 {
		return nil, 0
	}
	evts := make([]SensorEvent, count)
	ret := C.ASensorEventQueue_getEvents(queue.cptr(), (*C.ASensorEvent)(unsafe.Pointer(&evts[0])), C.size_t(count))
	if ret < 0 {
		return nil, int(ret)
	}
	return evts[:ret], int(ret)
}

/*****************************************************************************/

/**
 * Returns this sensor's name (non localized)
 */
//const char* ASensor_getName(ASensor const* sensor);
func (sensor *Sensor) GetName() string {
	return C.GoString(C.ASensor_getName(sensor.cptr()))
}

/**
 * Returns this sensor's vendor's name (non localized)
 */
//const char* ASensor_getVendor(ASensor const* sensor);
func (sensor *Sensor) GetVendor() string {
	return C.GoString(C.ASensor_getVendor(sensor.cptr()))
}

/**
 * Return this sensor's type
 */
//int ASensor_getType(ASensor const* sensor);
func (sensor *Sensor) GetType() SENSOR_TYPE {
	return SENSOR_TYPE(C.ASensor_getType(sensor.cptr()))
}

/**
 * Returns this sensors's resolution
 */
//float ASensor_getResolution(ASensor const* sensor) __NDK_FPABI__;
func (sensor *Sensor) GetResolution() float32 {
	return float32(C.ASensor_getResolution(sensor.cptr()))
}

/**
 * Returns the minimum delay allowed between events in microseconds.
 * A value of zero means that this sensor doesn't report events at a
 * constant rate, but rather only when a new data is available.
 */
//int ASensor_getMinDelay(ASensor const* sensor);
func (sensor *Sensor) GetMinDelay() time.Duration {
	return time.Duration(int(C.ASensor_getMinDelay(sensor.cptr()))) * time.Microsecond
}

// SENSOR_TYPE
func (t SENSOR_TYPE) String() string {
	switch t {
	case SENSOR_TYPE_ACCELEROMETER:
		return "sensor.accelerometer"
	case SENSOR_TYPE_MAGNETIC_FIELD:
		return "sensor.magnetic_field"
	case SENSOR_TYPE_ORIENTATION:
		return "sensor.orientation"
	case SENSOR_TYPE_GYROSCOPE:
		return "sensor.gyroscope"
	case SENSOR_TYPE_LIGHT:
		return "sensor.light"
	case SENSOR_TYPE_PRESSURE:
		return "sensor.pressure"
	case SENSOR_TYPE_TEMPERATURE:
		return "sensor.temperature"
	case SENSOR_TYPE_PROXIMITY:
		return "sensor.proximity"
	case SENSOR_TYPE_GRAVITY:
		return "sensor.gravity"
	case SENSOR_TYPE_LINEAR_ACCELERATION:
		return "sensor.linear_acceleration"
	case SENSOR_TYPE_ROTATION_VECTOR:
		return "sensor.rotation_vector"
	case SENSOR_TYPE_RELATIVE_HUMIDITY:
		return "sensor.relative_humidity"
	case SENSOR_TYPE_AMBIENT_TEMPERATURE:
		return "sensor.ambient_temperature"
	case SENSOR_TYPE_MAGNETIC_FIELD_UNCALIBRATED:
		return "sensor.magnetic_field_uncalibrated"
	case SENSOR_TYPE_GAME_ROTATION_VECTOR:
		return "sensor.game_rotation_vector"
	case SENSOR_TYPE_GYROSCOPE_UNCALIBRATED:
		return "sensor.gyroscope_uncalibrated"
	case SENSOR_TYPE_SIGNIFICANT_MOTION:
		return "sensor.significant_motion"
	case SENSOR_TYPE_STEP_DETECTOR:
		return "sensor.step_detector"
	case SENSOR_TYPE_STEP_COUNTER:
		return "sensor.step_counter"
	case SENSOR_TYPE_GEOMAGNETIC_ROTATION_VECTOR:
		return "sensor.geomagnetic_rotation_vector"
	case SENSOR_TYPE_HEART_RATE:
		return "sensor.heart_rate"
	case SENSOR_TYPE_WAKE_GESTURE:
		return "sensor.wake_gesture"
	case SENSOR_TYPE_GLANCE_GESTURE:
		return "sensor.glance_gesture"
	case SENSOR_TYPE_PICK_UP_GESTURE:
		return "sensor.pick_up_gesture"
	case SENSOR_TYPE_WRIST_TILT_GESTURE:
		return "sensor.wrist_tilt_gesture"
	case SENSOR_TYPE_DEVICE_ORIENTATION:
		return "sensor.device_orientation"
	case SENSOR_TYPE_POSE_6DOF:
		return "sensor.pose_6dof"
	case SENSOR_TYPE_STATIONARY_DETECT:
		return "sensor.stationary_detect"
	case SENSOR_TYPE_MOTION_DETECT:
		return "sensor.motion_detect"
	case SENSOR_TYPE_HEART_BEAT:
		return "sensor.heart_beat"
	case SENSOR_TYPE_DYNAMIC_SENSOR_META:
		return "sensor.dynamic_sensor_meta"
	case SENSOR_TYPE_LOW_LATENCY_OFFBODY_DETECT:
		return "sensor.low_latency_offbody_detect"
	case SENSOR_TYPE_ACCELEROMETER_UNCALIBRATED:
		return "sensor.accelerometer_uncalibrated"
	default:
		return fmt.Sprintf("sensor.type_%v", int(t))
	}
}

// SensorEvent
func (event *SensorEvent) cptr() *C.ASensorEvent {
	return (*C.ASensorEvent)(event)
}
func (event *SensorEvent) GetSensor() int {
	return int(event.sensor)
}

func (event *SensorEvent) GetType() SENSOR_TYPE {
	return SENSOR_TYPE(event._type)
}

func (event *SensorEvent) GetTimestamp() time.Duration {
	return time.Duration(event.timestamp) * time.Nanosecond
}

func (event *SensorEvent) getDatas() []byte {
	return (*[1 << 28]byte)(C.ASensorEvent_getDatas(event.cptr()))[:64]
}

func (event *SensorEvent) getReserveds() []uint32 {
	return (*[1 << 26]uint32)(unsafe.Pointer(&event.reserved1[0]))[:4]
}

func (event *SensorEvent) GetData(data interface{}) {
	switch data.(type) {
	case *SensorVector:
	case *MagneticVector:
	case *MetaDataEvent:
	case *UncalibratedEvent:
	case *HeartRateEvent:
	case *DynamicSensorEvent:
	case *AdditionalInfoEvent:
	case *float32, []float32:
	case *uint64, []uint64:
	case *int64, []int64:
	default:
		assert(false, "SensorEvent.GetData error format.")
		return
	}
	binary.Read(bytes.NewBuffer(event.getDatas()), binary.LittleEndian, data)
}
