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
 * @file NdkCaptureRequest.h
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

/*
#include <camera/NdkCaptureRequest.h>
*/
import "C"

import (
	"unsafe"

	"github.com/gooid/gooid/internal/ndk"
)

// Container for output targets
type OutputTargets C.ACameraOutputTargets

// Container for a single output target
type OutputTarget C.ACameraOutputTarget

func (target *OutputTarget) cptr() *C.ACameraOutputTarget {
	return (*C.ACameraOutputTarget)(target)
}

/**
 * ACaptureRequest is an opaque type that contains settings and output targets needed to capture
 * a single image from camera device.
 *
 * <p>ACaptureRequest contains the configuration for the capture hardware (sensor, lens, flash),
 * the processing pipeline, the control algorithms, and the output buffers. Also
 * contains the list of target {@link ANativeWindow}s to send image data to for this
 * capture.</p>
 *
 * <p>ACaptureRequest is created by {@link ACameraDevice_createCaptureRequest}.
 *
 * <p>ACaptureRequest is given to {@link ACameraCaptureSession_capture} or
 * {@link ACameraCaptureSession_setRepeatingRequest} to capture images from a camera.</p>
 *
 * <p>Each request can specify a different subset of target {@link ANativeWindow}s for the
 * camera to send the captured data to. All the {@link ANativeWindow}s used in a request must
 * be part of the {@link ANativeWindow} list given to the last call to
 * {@link ACameraDevice_createCaptureSession}, when the request is submitted to the
 * session.</p>
 *
 * <p>For example, a request meant for repeating preview might only include the
 * {@link ANativeWindow} for the preview SurfaceView or SurfaceTexture, while a
 * high-resolution still capture would also include a {@link ANativeWindow} from a
 * {@link AImageReader} configured for high-resolution JPEG images.</p>
 *
 * @see ACameraDevice_createCaptureRequest
 * @see ACameraCaptureSession_capture
 * @see ACameraCaptureSession_setRepeatingRequest
 */
type CaptureRequest C.ACaptureRequest

func (request *CaptureRequest) cptr() *C.ACaptureRequest {
	return (*C.ACaptureRequest)(request)
}

/**
 * Create a ACameraOutputTarget object.
 *
 * <p>The ACameraOutputTarget is used in {@link ACaptureRequest_addTarget} method to add an output
 * {@link ANativeWindow} to ACaptureRequest. Use {@link ACameraOutputTarget_free} to free the object
 * and its memory after application no longer needs the {@link ACameraOutputTarget}.</p>
 *
 * @param window the {@link ANativeWindow} to be associated with the {@link ACameraOutputTarget}
 * @param output the output {@link ACameraOutputTarget} will be stored here if the
 *                  method call succeeds.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds. The created ACameraOutputTarget will
 *                                be filled in the output argument.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if window or output is NULL.</li></ul>
 *
 * @see ACaptureRequest_addTarget
 */
//camera_status_t ACameraOutputTarget_create(ANativeWindow* window, ACameraOutputTarget** output);
func CameraOutputTargetCreate(window *app.Window) (*OutputTarget, error) {
	var output *C.ACameraOutputTarget
	ret := Status(C.ACameraOutputTarget_create((*C.ANativeWindow)(window.Pointer()), &output))
	return (*OutputTarget)(output), ret
}

/**
 * Free a ACameraOutputTarget object.
 *
 * @param output the {@link ACameraOutputTarget} to be freed.
 *
 * @see ACameraOutputTarget_create
 */
//void ACameraOutputTarget_free(ACameraOutputTarget* output);
func (output *OutputTarget) Free() {
	C.ACameraOutputTarget_free(output.cptr())
}

/**
 * Add an {@link ACameraOutputTarget} object to {@link ACaptureRequest}.
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param output the output {@link ACameraOutputTarget} to be added to capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request or output is NULL.</li></ul>
 */
//camera_status_t ACaptureRequest_addTarget(ACaptureRequest* request,
//        const ACameraOutputTarget* output);
func (request *CaptureRequest) AddTarget(output *OutputTarget) error {
	return Status(C.ACaptureRequest_addTarget(request.cptr(), output.cptr()))
}

/**
 * Remove an {@link ACameraOutputTarget} object from {@link ACaptureRequest}.
 *
 * <p>This method has no effect if the ACameraOutputTarget does not exist in ACaptureRequest.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param output the output {@link ACameraOutputTarget} to be removed from capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request or output is NULL.</li></ul>
 */
//camera_status_t ACaptureRequest_removeTarget(ACaptureRequest* request,
//        const ACameraOutputTarget* output);
func (request *CaptureRequest) RemoveTarget(output *OutputTarget) error {
	return Status(C.ACaptureRequest_removeTarget(request.cptr(), output.cptr()))
}

/**
 * Get a metadata entry from input {@link ACaptureRequest}.
 *
 * <p>The memory of the data field in returned entry is managed by camera framework. Do not
 * attempt to free it.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be get.
 * @param entry the output {@link ACameraMetadata_const_entry} will be filled here if the method
 *        call succeeeds.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if metadata or entry is NULL.</li>
 *         <li>{@link ACAMERA_ERROR_METADATA_NOT_FOUND} if the capture request does not contain an
 *             entry of input tag value.</li></ul>
 */
//camera_status_t ACaptureRequest_getConstEntry(
//        const ACaptureRequest* request, uint32_t tag, ACameraMetadata_const_entry* entry);
func (request *CaptureRequest) GetConstEntry(tag MetadataTag) (*MetadataConstEntry, error) {
	var ent C.ACameraMetadata_const_entry
	ret := Status(C.ACaptureRequest_getConstEntry(request.cptr(), C.uint32_t(tag), &ent))
	return (*MetadataConstEntry)(&ent), ret
}

/*
 * List all the entry tags in input {@link ACaptureRequest}.
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param numEntries number of metadata entries in input {@link ACaptureRequest}
 * @param tags the tag values of the metadata entries. Length of tags is returned in numEntries
 *             argument. The memory is managed by ACaptureRequest itself and must NOT be free/delete
 *             by application. Calling ACaptureRequest_setEntry_* methods will invalidate previous
 *             output of ACaptureRequest_getAllTags. Do not access tags after calling
 *             ACaptureRequest_setEntry_*. To get new list of tags after updating capture request,
 *             application must call ACaptureRequest_getAllTags again. Do NOT access tags after
 *             calling ACaptureRequest_free.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request, numEntries or tags is NULL.</li>
 *         <li>{@link ACAMERA_ERROR_UNKNOWN} if the method fails for some other reasons.</li></ul>
 */
//camera_status_t ACaptureRequest_getAllTags(
//        const ACaptureRequest* request, /*out*/int32_t* numTags, /*out*/const uint32_t** tags);
func (request *CaptureRequest) GetAllTags() ([]uint32, error) {
	var numEntries C.int32_t
	var tags *C.uint32_t
	ret := Status(C.ACaptureRequest_getAllTags(request.cptr(), &numEntries, &tags))
	return (*[1 << 28]uint32)(unsafe.Pointer(&tags))[:numEntries], ret
}

/**
 * Set/change a camera capture control entry with unsigned 8 bits data type.
 *
 * <p>Set count to 0 and data to NULL to remove a tag from the capture request.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be set.
 * @param count number of elements to be set in data argument
 * @param data the entries to be set/change in the capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request is NULL, count is larger than
 *             zero while data is NULL, the data type of the tag is not unsigned 8 bits, or
 *             the tag is not controllable by application.</li></ul>
 */
//camera_status_t ACaptureRequest_setEntry_u8(
//        ACaptureRequest* request, uint32_t tag, uint32_t count, const uint8_t* data);
func (request *CaptureRequest) SetEntryU8(tag MetadataTag, data []uint8) error {
	return Status(C.ACaptureRequest_setEntry_u8(request.cptr(), C.uint32_t(tag), C.uint32_t(len(data)), (*C.uint8_t)(&data[0])))
}

/**
 * Set/change a camera capture control entry with signed 32 bits data type.
 *
 * <p>Set count to 0 and data to NULL to remove a tag from the capture request.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be set.
 * @param count number of elements to be set in data argument
 * @param data the entries to be set/change in the capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request is NULL, count is larger than
 *             zero while data is NULL, the data type of the tag is not signed 32 bits, or
 *             the tag is not controllable by application.</li></ul>
 */
//camera_status_t ACaptureRequest_setEntry_i32(
//        ACaptureRequest* request, uint32_t tag, uint32_t count, const int32_t* data);
func (request *CaptureRequest) SetEntryI32(tag MetadataTag, data []int32) error {
	return Status(C.ACaptureRequest_setEntry_i32(request.cptr(), C.uint32_t(tag), C.uint32_t(len(data)), (*C.int32_t)(&data[0])))
}

/**
 * Set/change a camera capture control entry with float data type.
 *
 * <p>Set count to 0 and data to NULL to remove a tag from the capture request.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be set.
 * @param count number of elements to be set in data argument
 * @param data the entries to be set/change in the capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request is NULL, count is larger than
 *             zero while data is NULL, the data type of the tag is not float, or
 *             the tag is not controllable by application.</li></ul>
 */
//camera_status_t ACaptureRequest_setEntry_float(
//        ACaptureRequest* request, uint32_t tag, uint32_t count, const float* data);
func (request *CaptureRequest) SetEntryF32(tag MetadataTag, data []float32) error {
	return Status(C.ACaptureRequest_setEntry_float(request.cptr(), C.uint32_t(tag), C.uint32_t(len(data)), (*C.float)(&data[0])))
}

/**
 * Set/change a camera capture control entry with signed 64 bits data type.
 *
 * <p>Set count to 0 and data to NULL to remove a tag from the capture request.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be set.
 * @param count number of elements to be set in data argument
 * @param data the entries to be set/change in the capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request is NULL, count is larger than
 *             zero while data is NULL, the data type of the tag is not signed 64 bits, or
 *             the tag is not controllable by application.</li></ul>
 */
//camera_status_t ACaptureRequest_setEntry_i64(
//        ACaptureRequest* request, uint32_t tag, uint32_t count, const int64_t* data);
func (request *CaptureRequest) SetEntryI64(tag MetadataTag, data []int64) error {
	return Status(C.ACaptureRequest_setEntry_i64(request.cptr(), C.uint32_t(tag), C.uint32_t(len(data)), (*C.int64_t)(&data[0])))
}

/**
 * Set/change a camera capture control entry with double data type.
 *
 * <p>Set count to 0 and data to NULL to remove a tag from the capture request.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be set.
 * @param count number of elements to be set in data argument
 * @param data the entries to be set/change in the capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request is NULL, count is larger than
 *             zero while data is NULL, the data type of the tag is not double, or
 *             the tag is not controllable by application.</li></ul>
 */
//camera_status_t ACaptureRequest_setEntry_double(
//        ACaptureRequest* request, uint32_t tag, uint32_t count, const double* data);
func (request *CaptureRequest) SetEntryF64(tag MetadataTag, data []float64) error {
	return Status(C.ACaptureRequest_setEntry_double(request.cptr(), C.uint32_t(tag), C.uint32_t(len(data)), (*C.double)(&data[0])))
}

/**
 * Set/change a camera capture control entry with rational data type.
 *
 * <p>Set count to 0 and data to NULL to remove a tag from the capture request.</p>
 *
 * @param request the {@link ACaptureRequest} of interest.
 * @param tag the tag value of the camera metadata entry to be set.
 * @param count number of elements to be set in data argument
 * @param data the entries to be set/change in the capture request.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if request is NULL, count is larger than
 *             zero while data is NULL, the data type of the tag is not rational, or
 *             the tag is not controllable by application.</li></ul>
 */
//camera_status_t ACaptureRequest_setEntry_rational(
//        ACaptureRequest* request, uint32_t tag, uint32_t count,
//        const ACameraMetadata_rational* data);
func (request *CaptureRequest) SetEntryRational(tag MetadataTag, data []MetadataRational) error {
	cdata := make([]C.ACameraMetadata_rational, len(data))
	for i, d := range data {
		cdata[i].numerator = C.int32_t(d.numerator)
		cdata[i].denominator = C.int32_t(d.denominator)
	}
	return Status(C.ACaptureRequest_setEntry_rational(request.cptr(), C.uint32_t(tag), C.uint32_t(len(data)), &cdata[0]))
}

/**
 * Free a {@link ACaptureRequest} structure.
 *
 * @param request the {@link ACaptureRequest} to be freed.
 */
//void ACaptureRequest_free(ACaptureRequest* request);
func (request *CaptureRequest) Free() {
	C.ACaptureRequest_free(request.cptr())
}
