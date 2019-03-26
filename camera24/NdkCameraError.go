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
 * @file NdkCameraError.h
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
	"fmt"
)

/*
#include <camera/NdkCameraError.h>
*/
import "C"

type iStatus int

const (
	iSTATUS_OK iStatus = C.ACAMERA_OK

	STATUS_ERROR_BASE iStatus = C.ACAMERA_ERROR_BASE

	/**
	 * Camera operation has failed due to an unspecified cause.
	 */
	STATUS_ERROR_UNKNOWN iStatus = C.ACAMERA_ERROR_UNKNOWN

	/**
	 * Camera operation has failed due to an invalid parameter being passed to the method.
	 */
	STATUS_ERROR_INVALID_PARAMETER iStatus = C.ACAMERA_ERROR_INVALID_PARAMETER

	/**
	 * Camera operation has failed because the camera device has been closed, possibly because a
	 * higher-priority client has taken ownership of the camera device.
	 */
	STATUS_ERROR_CAMERA_DISCONNECTED iStatus = C.ACAMERA_ERROR_CAMERA_DISCONNECTED

	/**
	 * Camera operation has failed due to insufficient memory.
	 */
	STATUS_ERROR_NOT_ENOUGH_MEMORY iStatus = C.ACAMERA_ERROR_NOT_ENOUGH_MEMORY

	/**
	 * Camera operation has failed due to the requested metadata tag cannot be found in input
	 * {@link ACameraMetadata} or {@link ACaptureRequest}.
	 */
	STATUS_ERROR_METADATA_NOT_FOUND iStatus = C.ACAMERA_ERROR_METADATA_NOT_FOUND

	/**
	 * Camera operation has failed and the camera device has encountered a fatal error and needs to
	 * be re-opened before it can be used again.
	 */
	STATUS_ERROR_CAMERA_DEVICE iStatus = C.ACAMERA_ERROR_CAMERA_DEVICE

	/**
	 * Camera operation has failed and the camera service has encountered a fatal error.
	 *
	 * <p>The Android device may need to be shut down and restarted to restore
	 * camera function, or there may be a persistent hardware problem.</p>
	 *
	 * <p>An attempt at recovery may be possible by closing the
	 * ACameraDevice and the ACameraManager, and trying to acquire all resources
	 * again from scratch.</p>
	 */
	STATUS_ERROR_CAMERA_SERVICE iStatus = C.ACAMERA_ERROR_CAMERA_SERVICE

	/**
	 * The {@link ACameraCaptureSession} has been closed and cannnot perform any operation other
	 * than {@link ACameraCaptureSession_close}.
	 */
	STATUS_ERROR_SESSION_CLOSED iStatus = C.ACAMERA_ERROR_SESSION_CLOSED

	/**
	 * Camera operation has failed due to an invalid internal operation. Usually this is due to a
	 * low-level problem that may resolve itself on retry
	 */
	STATUS_ERROR_INVALID_OPERATION iStatus = C.ACAMERA_ERROR_INVALID_OPERATION

	/**
	 * Camera device does not support the stream configuration provided by application in
	 * {@link ACameraDevice_createCaptureSession}.
	 */
	STATUS_ERROR_STREAM_CONFIGURE_FAIL iStatus = C.ACAMERA_ERROR_STREAM_CONFIGURE_FAIL

	/**
	 * Camera device is being used by another higher priority camera API client.
	 */
	STATUS_ERROR_CAMERA_IN_USE iStatus = C.ACAMERA_ERROR_CAMERA_IN_USE

	/**
	 * The system-wide limit for number of open cameras or camera resources has been reached, and
	 * more camera devices cannot be opened until previous instances are closed.
	 */
	STATUS_ERROR_MAX_CAMERA_IN_USE iStatus = C.ACAMERA_ERROR_MAX_CAMERA_IN_USE

	/**
	 * The camera is disabled due to a device policy, and cannot be opened.
	 */
	STATUS_ERROR_CAMERA_DISABLED iStatus = C.ACAMERA_ERROR_CAMERA_DISABLED

	/**
	 * The application does not have permission to open camera.
	 */
	STATUS_ERROR_PERMISSION_DENIED iStatus = C.ACAMERA_ERROR_PERMISSION_DENIED
)

func (s iStatus) Error() string {
	switch s {
	case iSTATUS_OK:
		return "OK"
	case STATUS_ERROR_UNKNOWN:
		return "ERROR_UNKNOWN"
	case STATUS_ERROR_INVALID_PARAMETER:
		return "ERROR_INVALID_PARAMETER"
	case STATUS_ERROR_CAMERA_DISCONNECTED:
		return "ERROR_CAMERA_DISCONNECTED"
	case STATUS_ERROR_NOT_ENOUGH_MEMORY:
		return "ERROR_NOT_ENOUGH_MEMORY"
	case STATUS_ERROR_METADATA_NOT_FOUND:
		return "ERROR_METADATA_NOT_FOUND"
	case STATUS_ERROR_CAMERA_DEVICE:
		return "ERROR_CAMERA_DEVICE"
	case STATUS_ERROR_CAMERA_SERVICE:
		return "ERROR_CAMERA_SERVICE"
	case STATUS_ERROR_SESSION_CLOSED:
		return "ERROR_SESSION_CLOSED"
	case STATUS_ERROR_INVALID_OPERATION:
		return "ERROR_INVALID_OPERATION"
	case STATUS_ERROR_STREAM_CONFIGURE_FAIL:
		return "ERROR_STREAM_CONFIGURE_FAIL"
	case STATUS_ERROR_CAMERA_IN_USE:
		return "ERROR_CAMERA_IN_USE"
	case STATUS_ERROR_MAX_CAMERA_IN_USE:
		return "ERROR_MAX_CAMERA_IN_USE"
	case STATUS_ERROR_CAMERA_DISABLED:
		return "ERROR_CAMERA_DISABLED"
	case STATUS_ERROR_PERMISSION_DENIED:
		return "ERROR_PERMISSION_DENIED"
	default:
		return fmt.Sprintf("ERROR_UNKNOW(%d)", int(s))
	}
}

func Status(i C.camera_status_t) error {
	if i != C.ACAMERA_OK {
		return iStatus(i)
	}
	return nil
}
