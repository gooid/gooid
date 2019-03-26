// Copyright 2018 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
 * Copyright (C) 2015 The Android Open Source Project
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

/**
 * @addtogroup Camera
 * @{
 */

/**
 * @file NdkCameraManager.h
 */

/*
 * This file defines an NDK API.
 * Do not remove methods.
 * Do not change method signatures.
 * Do not change the value of constants.
 * Do not change the size of any of the classes defined in here.
 * Do not reference types that are not part of the NDK.
 * Do not #include files that aren't part of the NDK.
 */

package camera

import (
	"log"
	"reflect"
	"unsafe"
)

/*
#include <stdbool.h>
#include <stdlib.h>
#include <camera/NdkCameraManager.h>

typedef const char* pcchar;
extern void cgoCameraAvailable(void* context, const char* cameraId);
extern void cgoCameraUnavailable(void* context, const char* cameraId);
extern void cgoCameraDeviceStateCallbacksOnError(void* context, ACameraDevice* device, int error);
extern void cgoCameraDeviceStateCallbacksOnDisconnected(void* context, ACameraDevice* device);

*/
import "C"

/**
 * ACameraManager is opaque type that provides access to camera service.
 *
 * A pointer can be obtained using {@link ACameraManager_create} method.
 */
//typedef struct ACameraManager ACameraManager;
type Manager C.ACameraManager

func (device *Manager) cptr() *C.ACameraManager {
	return (*C.ACameraManager)(device)
}

/**
 * Create ACameraManager instance.
 *
 * <p>The ACameraManager is responsible for
 * detecting, characterizing, and connecting to {@link ACameraDevice}s.</p>
 *
 * <p>The caller must call {@link ACameraManager_delete} to free the resources once it is done
 * using the ACameraManager instance.</p>
 *
 * @return a {@link ACameraManager} instance.
 *
 */
//ACameraManager* ACameraManager_create();
func ManagerCreate() *Manager {
	return (*Manager)(C.ACameraManager_create())
}

/**
 * <p>Delete the {@link ACameraManager} instance and free its resources. </p>
 *
 * @param manager the {@link ACameraManager} instance to be deleted.
 */
//void ACameraManager_delete(ACameraManager* manager);
func (manager *Manager) Delete() {
	C.ACameraManager_delete(manager.cptr())
}

/// Struct to hold list of camera devices
//typedef struct ACameraIdList {
//    int numCameras;          ///< Number of connected camera devices
//    const char** cameraIds;  ///< list of identifier of connected camera devices
//} ACameraIdList;

/**
 * Create a list of currently connected camera devices, including
 * cameras that may be in use by other camera API clients.
 *
 * <p>Non-removable cameras use integers starting at 0 for their
 * identifiers, while removable cameras have a unique identifier for each
 * individual device, even if they are the same model.</p>
 *
 * <p>ACameraManager_getCameraIdList will allocate and return an {@link ACameraIdList}.
 * The caller must call {@link ACameraManager_deleteCameraIdList} to free the memory</p>
 *
 * @param manager the {@link ACameraManager} of interest
 * @param cameraIdList the output {@link ACameraIdList} will be filled in here if the method call
 *        succeeds.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if manager or cameraIdList is NULL.</li>
 *         <li>{@link ACAMERA_ERROR_CAMERA_DISCONNECTED} if connection to camera service fails.</li>
 *         <li>{@link ACAMERA_ERROR_NOT_ENOUGH_MEMORY} if allocating memory fails.</li></ul>
 */
//camera_status_t ACameraManager_getCameraIdList(ACameraManager* manager,
//                                              /*out*/ACameraIdList** cameraIdList);
func (manager *Manager) GetCameraIdList() ([]string, error) {
	var list *C.ACameraIdList
	status := Status(C.ACameraManager_getCameraIdList(manager.cptr(), &list))
	if status == nil {
		log.Println(" GetCameraIdList:", int(list.numCameras))
		ids := make([]string, list.numCameras)
		if list.numCameras > 0 {
			for i, ptr := range (*[1 << 27]*C.char)(unsafe.Pointer(list.cameraIds))[:list.numCameras] {
				ids[i] = C.GoString(ptr)
			}
		}
		C.ACameraManager_deleteCameraIdList(list)
		return ids, status
	}
	return nil, status
}

/**
 * Delete a list of camera devices allocated via {@link ACameraManager_getCameraIdList}.
 *
 * @param cameraIdList the {@link ACameraIdList} to be deleted.
 */
//void ACameraManager_deleteCameraIdList(ACameraIdList* cameraIdList);

/**
 * Definition of camera availability callbacks.
 *
 * @param context The optional application context provided by user in
 *                {@link ACameraManager_AvailabilityCallbacks}.
 * @param cameraId The ID of the camera device whose availability is changing. The memory of this
 *                 argument is owned by camera framework and will become invalid immediately after
 *                 this callback returns.
 */
//typedef void (*ACameraManager_AvailabilityCallback)(void* context, const char* cameraId);

/**
 * A listener for camera devices becoming available or unavailable to open.
 *
 * <p>Cameras become available when they are no longer in use, or when a new
 * removable camera is connected. They become unavailable when some
 * application or service starts using a camera, or when a removable camera
 * is disconnected.</p>
 *
 * @see ACameraManager_registerAvailabilityCallback
 */
//go typedef struct ACameraManager_AvailabilityListener {
/// Optional application context.
//go     void*                               context;
/// Called when a camera becomes available
//go     ACameraManager_AvailabilityCallback onCameraAvailable;
/// Called when a camera becomes unavailable
//go     ACameraManager_AvailabilityCallback onCameraUnavailable;
//go } ACameraManager_AvailabilityCallbacks;
type AvailabilityCallbacks interface {
	OnCameraAvailable(id string)
	OnCameraUnavailable(id string)
}

/**
 * Register camera availability callbacks.
 *
 * <p>onCameraUnavailable will be called whenever a camera device is opened by any camera API client.
 * Other camera API clients may still be able to open such a camera device, evicting the existing
 * client if they have higher priority than the existing client of a camera device.
 * See {@link ACameraManager_openCamera} for more details.</p>
 *
 * <p>The callbacks will be called on a dedicated thread shared among all ACameraManager
 * instances.</p>
 *
 * <p>Since this callback will be registered with the camera service, remember to unregister it
 * once it is no longer needed; otherwise the callback will continue to receive events
 * indefinitely and it may prevent other resources from being released. Specifically, the
 * callbacks will be invoked independently of the general activity lifecycle and independently
 * of the state of individual ACameraManager instances.</p>
 *
 * @param manager the {@link ACameraManager} of interest.
 * @param callback the {@link ACameraManager_AvailabilityCallbacks} to be registered.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if manager or callback is NULL, or
 *                  {ACameraManager_AvailabilityCallbacks#onCameraAvailable} or
 *                  {ACameraManager_AvailabilityCallbacks#onCameraUnavailable} is NULL.</li></ul>
 */
//camera_status_t ACameraManager_registerAvailabilityCallback(
//        ACameraManager* manager, const ACameraManager_AvailabilityCallbacks* callback);
func (manager *Manager) RegisterAvailabilityCallback(cbs AvailabilityCallbacks) error {
	var ccbs C.ACameraManager_AvailabilityCallbacks
	ccbs.context = unsafe.Pointer(^reflect.ValueOf(cbs).Pointer())
	availabilityKeepLives[ccbs.context] = cbs
	ccbs.onCameraAvailable = C.ACameraManager_AvailabilityCallback(C.cgoCameraAvailable)
	ccbs.onCameraUnavailable = C.ACameraManager_AvailabilityCallback(C.cgoCameraUnavailable)
	return Status(C.ACameraManager_registerAvailabilityCallback(manager.cptr(), &ccbs))
}

var availabilityKeepLives = map[unsafe.Pointer]interface{}{}

//export cgoCameraAvailable
func cgoCameraAvailable(context unsafe.Pointer, cameraId C.pcchar) {
	cbs := availabilityKeepLives[context].(AvailabilityCallbacks)
	cbs.OnCameraAvailable(C.GoString(cameraId))
}

//export cgoCameraUnavailable
func cgoCameraUnavailable(context unsafe.Pointer, cameraId C.pcchar) {
	cbs := availabilityKeepLives[context].(AvailabilityCallbacks)
	cbs.OnCameraUnavailable(C.GoString(cameraId))
}

/**
 * Unregister camera availability callbacks.
 *
 * <p>Removing a callback that isn't registered has no effect.</p>
 *
 * @param manager the {@link ACameraManager} of interest.
 * @param callback the {@link ACameraManager_AvailabilityCallbacks} to be unregistered.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if callback,
 *                  {ACameraManager_AvailabilityCallbacks#onCameraAvailable} or
 *                  {ACameraManager_AvailabilityCallbacks#onCameraUnavailable} is NULL.</li></ul>
 */
//camera_status_t ACameraManager_unregisterAvailabilityCallback(
//        ACameraManager* manager, const ACameraManager_AvailabilityCallbacks* callback);
func (manager *Manager) UnregisterAvailabilityCallback(cbs AvailabilityCallbacks) error {
	context := unsafe.Pointer(^reflect.ValueOf(cbs).Pointer())
	if _, ok := availabilityKeepLives[context]; ok {
		delete(availabilityKeepLives, context)
		if len(availabilityKeepLives) == 0 {
			var ccbs C.ACameraManager_AvailabilityCallbacks
			ccbs.context = context
			ccbs.onCameraAvailable = C.ACameraManager_AvailabilityCallback(C.cgoCameraAvailable)
			ccbs.onCameraUnavailable = C.ACameraManager_AvailabilityCallback(C.cgoCameraUnavailable)
			return Status(C.ACameraManager_unregisterAvailabilityCallback(manager.cptr(), &ccbs))
		}
		return nil
	}
	return STATUS_ERROR_INVALID_PARAMETER
}

/**
 * Query the capabilities of a camera device. These capabilities are
 * immutable for a given camera.
 *
 * <p>See {@link ACameraMetadata} document and {@link NdkCameraMetadataTags.h} for more details.</p>
 *
 * <p>The caller must call {@link ACameraMetadata_free} to free the memory of the output
 * characteristics.</p>
 *
 * @param manager the {@link ACameraManager} of interest.
 * @param cameraId the ID string of the camera device of interest.
 * @param characteristics the output {@link ACameraMetadata} will be filled here if the method call
 *        succeeeds.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if manager, cameraId, or characteristics
 *                  is NULL, or cameraId does not match any camera devices connected.</li>
 *         <li>{@link ACAMERA_ERROR_CAMERA_DISCONNECTED} if connection to camera service fails.</li>
 *         <li>{@link ACAMERA_ERROR_NOT_ENOUGH_MEMORY} if allocating memory fails.</li>
 *         <li>{@link ACAMERA_ERROR_UNKNOWN} if the method fails for some other reasons.</li></ul>
 */
//camera_status_t ACameraManager_getCameraCharacteristics(
//        ACameraManager* manager, const char* cameraId,
//        /*out*/ACameraMetadata** characteristics);
func (manager *Manager) GetCameraCharacteristics(cameraId string) (*Metadata, error) {
	ccameraId := C.CString(cameraId)
	defer C.free(unsafe.Pointer(ccameraId))
	var v *C.ACameraMetadata
	ret := Status(C.ACameraManager_getCameraCharacteristics(manager.cptr(), ccameraId, &v))
	return (*Metadata)(v), ret
}

/**
 * Open a connection to a camera with the given ID. The opened camera device will be
 * returned in the `device` parameter.
 *
 * <p>Use {@link ACameraManager_getCameraIdList} to get the list of available camera
 * devices. Note that even if an id is listed, open may fail if the device
 * is disconnected between the calls to {@link ACameraManager_getCameraIdList} and
 * {@link ACameraManager_openCamera}, or if a higher-priority camera API client begins using the
 * camera device.</p>
 *
 * <p>Devices for which the
 * {@link ACameraManager_AvailabilityCallbacks#onCameraUnavailable} callback has been called due to
 * the device being in use by a lower-priority, background camera API client can still potentially
 * be opened by calling this method when the calling camera API client has a higher priority
 * than the current camera API client using this device.  In general, if the top, foreground
 * activity is running within your application process, your process will be given the highest
 * priority when accessing the camera, and this method will succeed even if the camera device is
 * in use by another camera API client. Any lower-priority application that loses control of the
 * camera in this way will receive an
 * {@link ACameraDevice_stateCallbacks#onDisconnected} callback.</p>
 *
 * <p>Once the camera is successfully opened,the ACameraDevice can then be set up
 * for operation by calling {@link ACameraDevice_createCaptureSession} and
 * {@link ACameraDevice_createCaptureRequest}.</p>
 *
 * <p>If the camera becomes disconnected after this function call returns,
 * {@link ACameraDevice_stateCallbacks#onDisconnected} with a
 * ACameraDevice in the disconnected state will be called.</p>
 *
 * <p>If the camera runs into error after this function call returns,
 * {@link ACameraDevice_stateCallbacks#onError} with a
 * ACameraDevice in the error state will be called.</p>
 *
 * @param manager the {@link ACameraManager} of interest.
 * @param cameraId the ID string of the camera device to be opened.
 * @param callback the {@link ACameraDevice_StateCallbacks} associated with the opened camera device.
 * @param device the opened {@link ACameraDevice} will be filled here if the method call succeeds.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if manager, cameraId, callback, or device
 *                  is NULL, or cameraId does not match any camera devices connected.</li>
 *         <li>{@link ACAMERA_ERROR_CAMERA_DISCONNECTED} if connection to camera service fails.</li>
 *         <li>{@link ACAMERA_ERROR_NOT_ENOUGH_MEMORY} if allocating memory fails.</li>
 *         <li>{@link ACAMERA_ERROR_CAMERA_IN_USE} if camera device is being used by a higher
 *                   priority camera API client.</li>
 *         <li>{@link ACAMERA_ERROR_MAX_CAMERA_IN_USE} if the system-wide limit for number of open
 *                   cameras or camera resources has been reached, and more camera devices cannot be
 *                   opened until previous instances are closed.</li>
 *         <li>{@link ACAMERA_ERROR_CAMERA_DISABLED} if the camera is disabled due to a device
 *                   policy, and cannot be opened.</li>
 *         <li>{@link ACAMERA_ERROR_PERMISSION_DENIED} if the application does not have permission
 *                   to open camera.</li>
 *         <li>{@link ACAMERA_ERROR_UNKNOWN} if the method fails for some other reasons.</li></ul>
 */
//camera_status_t ACameraManager_openCamera(
//        ACameraManager* manager, const char* cameraId,
//        ACameraDevice_StateCallbacks* callback,
//        /*out*/ACameraDevice** device);
func (manager *Manager) OpenCamera(cameraId string,
	callbacks DeviceStateCallbacks) (*Device, error) {
	ccameraId := C.CString(cameraId)
	defer C.free(unsafe.Pointer(ccameraId))
	var dev *C.ACameraDevice

	var ccallbacks C.ACameraDevice_stateCallbacks
	ccallbacks.context = unsafe.Pointer(^reflect.ValueOf(callbacks).Pointer())
	deviceStateKeepLives[ccallbacks.context] = &openCameraContext{
		ccbs: &ccallbacks, cbs: callbacks}
	ccallbacks.onDisconnected = C.ACameraDevice_StateCallback(C.cgoCameraDeviceStateCallbacksOnDisconnected)
	ccallbacks.onError = C.ACameraDevice_ErrorStateCallback(C.cgoCameraDeviceStateCallbacksOnError)

	ret := Status(C.ACameraManager_openCamera(manager.cptr(), ccameraId, &ccallbacks, &dev))
	return (*Device)(dev), ret
}

type openCameraContext struct {
	ccbs *C.ACameraDevice_stateCallbacks
	cbs  DeviceStateCallbacks
}

var deviceStateKeepLives = map[unsafe.Pointer]interface{}{}

//export cgoCameraDeviceStateCallbacksOnDisconnected
func cgoCameraDeviceStateCallbacksOnDisconnected(context unsafe.Pointer, device *C.ACameraDevice) {
	ctx := deviceStateKeepLives[context].(*openCameraContext)
	ctx.cbs.OnDisconnected((*Device)(device))
	delete(deviceStateKeepLives, context)
}

//export cgoCameraDeviceStateCallbacksOnError
func cgoCameraDeviceStateCallbacksOnError(context unsafe.Pointer, device *C.ACameraDevice, error C.int) {
	ctx := deviceStateKeepLives[context].(*openCameraContext)
	ctx.cbs.OnError((*Device)(device), int(error))
	delete(deviceStateKeepLives, context)
}
