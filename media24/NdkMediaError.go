// Copyright 2018 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
 * Copyright (C) 2014 The Android Open Source Project
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
 * @addtogroup Media
 * @{
 */

/**
 * @file NdkMediaError.h
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

package media

import (
	"fmt"
)

/*
#include "media/NdkMediaError.h"
*/
import "C"

type iStatus int

const (
	MEDIA_OK = C.AMEDIA_OK

	ERROR_BASE              = C.AMEDIA_ERROR_BASE
	ERROR_UNKNOWN           = C.AMEDIA_ERROR_UNKNOWN
	ERROR_MALFORMED         = C.AMEDIA_ERROR_MALFORMED
	ERROR_UNSUPPORTED       = C.AMEDIA_ERROR_UNSUPPORTED
	ERROR_INVALID_OBJECT    = C.AMEDIA_ERROR_INVALID_OBJECT
	ERROR_INVALID_PARAMETER = C.AMEDIA_ERROR_INVALID_PARAMETER
	ERROR_INVALID_OPERATION = C.AMEDIA_ERROR_INVALID_OPERATION

	DRM_ERROR_BASE         = C.AMEDIA_DRM_ERROR_BASE
	DRM_NOT_PROVISIONED    = C.AMEDIA_DRM_NOT_PROVISIONED
	DRM_RESOURCE_BUSY      = C.AMEDIA_DRM_RESOURCE_BUSY
	DRM_DEVICE_REVOKED     = C.AMEDIA_DRM_DEVICE_REVOKED
	DRM_SHORT_BUFFER       = C.AMEDIA_DRM_SHORT_BUFFER
	DRM_SESSION_NOT_OPENED = C.AMEDIA_DRM_SESSION_NOT_OPENED
	DRM_TAMPER_DETECTED    = C.AMEDIA_DRM_TAMPER_DETECTED
	DRM_VERIFY_FAILED      = C.AMEDIA_DRM_VERIFY_FAILED
	DRM_NEED_KEY           = C.AMEDIA_DRM_NEED_KEY
	DRM_LICENSE_EXPIRED    = C.AMEDIA_DRM_LICENSE_EXPIRED

	IMGREADER_ERROR_BASE          = C.AMEDIA_IMGREADER_ERROR_BASE
	IMGREADER_NO_BUFFER_AVAILABLE = C.AMEDIA_IMGREADER_NO_BUFFER_AVAILABLE
	IMGREADER_MAX_IMAGES_ACQUIRED = C.AMEDIA_IMGREADER_MAX_IMAGES_ACQUIRED
	IMGREADER_CANNOT_LOCK_IMAGE   = C.AMEDIA_IMGREADER_CANNOT_LOCK_IMAGE
	IMGREADER_CANNOT_UNLOCK_IMAGE = C.AMEDIA_IMGREADER_CANNOT_UNLOCK_IMAGE
	IMGREADER_IMAGE_NOT_LOCKED    = C.AMEDIA_IMGREADER_IMAGE_NOT_LOCKED
)

func (i iStatus) Error() string {
	switch i {
	case ERROR_UNKNOWN:
		return "ERROR_UNKNOWN"
	case ERROR_MALFORMED:
		return "ERROR_MALFORMED"
	case ERROR_UNSUPPORTED:
		return "ERROR_UNSUPPORTED"
	case ERROR_INVALID_OBJECT:
		return "ERROR_INVALID_OBJECT"
	case ERROR_INVALID_PARAMETER:
		return "ERROR_INVALID_PARAMETER"
	case ERROR_INVALID_OPERATION:
		return "ERROR_INVALID_OPERATION"
	case DRM_NOT_PROVISIONED:
		return "DRM_NOT_PROVISIONED"
	case DRM_RESOURCE_BUSY:
		return "DRM_RESOURCE_BUSY"
	case DRM_DEVICE_REVOKED:
		return "DRM_DEVICE_REVOKED"
	case DRM_SHORT_BUFFER:
		return "DRM_SHORT_BUFFER"
	case DRM_SESSION_NOT_OPENED:
		return "DRM_SESSION_NOT_OPENED"
	case DRM_TAMPER_DETECTED:
		return "DRM_TAMPER_DETECTED"
	case DRM_VERIFY_FAILED:
		return "DRM_VERIFY_FAILED"
	case DRM_NEED_KEY:
		return "DRM_NEED_KEY"
	case DRM_LICENSE_EXPIRED:
		return "DRM_LICENSE_EXPIRED"
	case IMGREADER_NO_BUFFER_AVAILABLE:
		return "IMGREADER_NO_BUFFER_AVAILABLE"
	case IMGREADER_MAX_IMAGES_ACQUIRED:
		return "IMGREADER_MAX_IMAGES_ACQUIRED"
	case IMGREADER_CANNOT_LOCK_IMAGE:
		return "IMGREADER_CANNOT_LOCK_IMAGE"
	case IMGREADER_CANNOT_UNLOCK_IMAGE:
		return "IMGREADER_CANNOT_UNLOCK_IMAGE"
	case IMGREADER_IMAGE_NOT_LOCKED:
		return "IMGREADER_IMAGE_NOT_LOCKED"
	}
	return fmt.Sprintf("ERROR_MEDIA_%d", i)
}

func Status(i C.media_status_t) error {
	if i != MEDIA_OK {
		return iStatus(i)
	}
	return nil
}
