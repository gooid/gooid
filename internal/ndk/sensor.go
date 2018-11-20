package app

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
	"time"
	"unsafe"
)

/*
 * Structures and functions to receive and process sensor events in
 * native code.
 *
 */

/*
 * Sensor types
 * (keep in sync with hardware/sensor.h)
 */

type SENSOR_TYPE int
type SENSOR_STATUS int

const (
	SENSOR_TYPE_ACCELEROMETER  = SENSOR_TYPE(C.ASENSOR_TYPE_ACCELEROMETER)
	SENSOR_TYPE_MAGNETIC_FIELD = SENSOR_TYPE(C.ASENSOR_TYPE_MAGNETIC_FIELD)
	SENSOR_TYPE_GYROSCOPE      = SENSOR_TYPE(C.ASENSOR_TYPE_GYROSCOPE)
	SENSOR_TYPE_LIGHT          = SENSOR_TYPE(C.ASENSOR_TYPE_LIGHT)
	SENSOR_TYPE_PROXIMITY      = SENSOR_TYPE(C.ASENSOR_TYPE_PROXIMITY)

	/*
	 * Sensor accuracy measure
	 */

	SENSOR_STATUS_UNRELIABLE      = SENSOR_STATUS(C.ASENSOR_STATUS_UNRELIABLE)
	SENSOR_STATUS_ACCURACY_LOW    = SENSOR_STATUS(C.ASENSOR_STATUS_ACCURACY_LOW)
	SENSOR_STATUS_ACCURACY_MEDIUM = SENSOR_STATUS(C.ASENSOR_STATUS_ACCURACY_MEDIUM)
	SENSOR_STATUS_ACCURACY_HIGH   = SENSOR_STATUS(C.ASENSOR_STATUS_ACCURACY_HIGH)

	/*
	 * A few useful constants
	 */

	/* Earth's gravity in m/s^2 */
	SENSOR_STANDARD_GRAVITY = 9.80665
	/* Maximum magnetic field on Earth's surface in uT */
	SENSOR_MAGNETIC_FIELD_EARTH_MAX = 60.0
	/* Minimum magnetic field on Earth's surface in uT*/
	SENSOR_MAGNETIC_FIELD_EARTH_MIN = 30.0
)

/*
 * A sensor event.
 */

/* NOTE: Must match hardware/sensors.h */
type SensorVector struct {
	X, Y, Z float32 /* azimuth; pitch; roll; */
	Status  int8
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

/* NOTE: Must match hardware/sensors.h */
type SensorEvent C.ASensorEvent

type SensorManager C.ASensorManager

func (manager *SensorManager) cptr() *C.ASensorManager {
	return (*C.ASensorManager)(manager)
}

type SensorEventQueue C.ASensorEventQueue

func (queue *SensorEventQueue) cptr() *C.ASensorEventQueue {
	return (*C.ASensorEventQueue)(queue)
}

type Sensor C.ASensor

func (sensor *Sensor) cptr() *C.ASensor {
	return (*C.ASensor)(sensor)
}

/*****************************************************************************/

/*
 * Get a reference to the sensor manager. ASensorManager is a singleton.
 *
 * Example:
 *
 *     ASensorManager* sensorManager = ASensorManager_getInstance();
 *
 */
//ASensorManager* ASensorManager_getInstance();
func SensorManagerInstance() *SensorManager {
	return (*SensorManager)(C.ASensorManager_getInstance())
}

/*
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

/*
 * Returns the default sensor for the given type, or NULL if no sensor
 * of that type exist.
 */
//ASensor const* ASensorManager_getDefaultSensor(ASensorManager* manager, int type);
func (manager *SensorManager) GetDefaultSensor(typ SENSOR_TYPE) *Sensor {
	return (*Sensor)(C.ASensorManager_getDefaultSensor(manager.cptr(), C.int(typ)))
}

/*
 * Creates a new sensor event queue and associate it with a looper.
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

/*
 * Destroys the event queue and free all resources associated to it.
 */
//int ASensorManager_destroyEventQueue(ASensorManager* manager, ASensorEventQueue* queue);
func (manager *SensorManager) destroy(queue *SensorEventQueue) int {
	return int(C.ASensorManager_destroyEventQueue(manager.cptr(), queue.cptr()))
}

/*****************************************************************************/

/*
 * Enable the selected sensor. Returns a negative error code on failure.
 */
//int ASensorEventQueue_enableSensor(ASensorEventQueue* queue, ASensor const* sensor);
func (queue *SensorEventQueue) enableSensor(sensor *Sensor) int {
	return int(C.ASensorEventQueue_enableSensor(queue.cptr(), sensor.cptr()))
}

/*
 * Disable the selected sensor. Returns a negative error code on failure.
 */
//int ASensorEventQueue_disableSensor(ASensorEventQueue* queue, ASensor const* sensor);
func (queue *SensorEventQueue) disableSensor(sensor *Sensor) int {
	return int(C.ASensorEventQueue_disableSensor(queue.cptr(), sensor.cptr()))
}

/*
 * Sets the delivery rate of events in microseconds for the given sensor.
 * Note that this is a hint only, generally event will arrive at a higher
 * rate. It is an error to set a rate inferior to the value returned by
 * ASensor_getMinDelay().
 * Returns a negative error code on failure.
 */
//int ASensorEventQueue_setEventRate(ASensorEventQueue* queue, ASensor const* sensor, int32_t usec);
func (queue *SensorEventQueue) setEventRate(sensor *Sensor, t time.Duration) int {
	usec := t / time.Microsecond
	return int(C.ASensorEventQueue_setEventRate(queue.cptr(), sensor.cptr(), C.int32_t(usec)))
}

/*
 * Returns true if there are one or more events available in the
 * sensor queue.  Returns 1 if the queue has events; 0 if
 * it does not have events; and a negative value if there is an error.
 */
//int ASensorEventQueue_hasEvents(ASensorEventQueue* queue);
func (queue *SensorEventQueue) hasEvents() int {
	return int(C.ASensorEventQueue_hasEvents(queue.cptr()))
}

/*
 * Returns the next available events from the queue.  Returns a negative
 * value if no events are available or an error has occurred, otherwise
 * the number of events returned.
 *
 * Examples:
 *   ASensorEvent event;
 *   ssize_t numEvent = ASensorEventQueue_getEvents(queue, &event, 1);
 *
 *   ASensorEvent eventBuffer[8];
 *   ssize_t numEvent = ASensorEventQueue_getEvents(queue, eventBuffer, 8);
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

/*
 * Returns this sensor's name (non localized)
 */
//const char* ASensor_getName(ASensor const* sensor);
func (sensor *Sensor) GetName() string {
	return C.GoString(C.ASensor_getName(sensor.cptr()))
}

/*
 * Returns this sensor's vendor's name (non localized)
 */
//const char* ASensor_getVendor(ASensor const* sensor);
func (sensor *Sensor) GetVendor() string {
	return C.GoString(C.ASensor_getVendor(sensor.cptr()))
}

/*
 * Return this sensor's type
 */
//int ASensor_getType(ASensor const* sensor);
func (sensor *Sensor) GetType() SENSOR_TYPE {
	return SENSOR_TYPE(C.ASensor_getType(sensor.cptr()))
}

/*
 * Returns this sensors's resolution
 */
//float ASensor_getResolution(ASensor const* sensor) __NDK_FPABI__;
func (sensor *Sensor) GetResolution() float32 {
	return float32(C.ASensor_getResolution(sensor.cptr()))
}

/*
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
		return "ACCELEROMETER"
	case SENSOR_TYPE_MAGNETIC_FIELD:
		return "MAGNETIC_FIELD"
	case SENSOR_TYPE_GYROSCOPE:
		return "GYROSCOPE"
	case SENSOR_TYPE_LIGHT:
		return "LIGHT"
	case SENSOR_TYPE_PROXIMITY:
		return "PROXIMITY"
	default:
		return fmt.Sprintf("UNKNOW:%v", int(t))
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

func (event *SensorEvent) GetTimestamp() int64 {
	return int64(event.timestamp)
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
	case *MetaDataEvent:
	case *UncalibratedEvent:
	case *HeartRateEvent:
	case *float32, []float32:
	case *uint64, []uint64:
	case *int64, []int64:
	default:
		assert(false, "SensorEvent.GetData error format.")
		return
	}
	binary.Read(bytes.NewBuffer(event.getDatas()), binary.LittleEndian, data)
}
