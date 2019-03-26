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
 * @file NdkCameraMetadata.h
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
	"unsafe"
)

/*
#include <stdbool.h>
#include <android/native_window.h>
#include <camera/NdkCameraDevice.h>
#include <camera/NdkCameraMetadata.h>

static void* MetadataConstEntry_getData(ACameraMetadata_const_entry* entry) {
	switch (entry->type) {
	/// Unsigned 8-bit integer (uint8_t)
	case ACAMERA_TYPE_BYTE:
		return (void *)entry->data.u8;

	/// Signed 32-bit integer (int32_t)
	case ACAMERA_TYPE_INT32:
		return (void *)entry->data.i32;

	/// 32-bit float (float)
	case ACAMERA_TYPE_FLOAT:
		return (void *)entry->data.f;

	/// Signed 64-bit integer (int64_t)
	case ACAMERA_TYPE_INT64:
		return (void *)entry->data.i64;

	/// 64-bit float (double)
	case ACAMERA_TYPE_DOUBLE:
		return (void *)entry->data.d;

	/// A 64-bit fraction (ACameraMetadata_rational)
	case ACAMERA_TYPE_RATIONAL:
		return (void *)entry->data.r;
	}
	return (void *)0;
}
*/
import "C"

/**
 * ACameraMetadata is opaque type that provides access to read-only camera metadata like camera
 * characteristics (via {@link ACameraManager_getCameraCharacteristics}) or capture results (via
 * {@link ACameraCaptureSession_captureCallback_result}).
 */
//typedef struct ACameraMetadata ACameraMetadata;
type Metadata C.ACameraMetadata

func (metadata *Metadata) cptr() *C.ACameraMetadata {
	return (*C.ACameraMetadata)(metadata)
}

/**
 * Possible data types of a metadata entry.
 *
 * Keep in sync with system/media/include/system/camera_metadata.h
 */
type Type int

const (
	/// Unsigned 8-bit integer (uint8_t)
	TYPE_BYTE Type = C.ACAMERA_TYPE_BYTE
	/// Signed 32-bit integer (int32_t)
	TYPE_INT32 Type = C.ACAMERA_TYPE_INT32
	/// 32-bit float (float)
	TYPE_FLOAT Type = C.ACAMERA_TYPE_FLOAT
	/// Signed 64-bit integer (int64_t)
	TYPE_INT64 Type = C.ACAMERA_TYPE_INT64
	/// 64-bit float (double)
	TYPE_DOUBLE Type = C.ACAMERA_TYPE_DOUBLE
	/// A 64-bit fraction (ACameraMetadata_rational)
	TYPE_RATIONAL Type = C.ACAMERA_TYPE_RATIONAL
	/// Number of type fields
	NUM_TYPES = C.ACAMERA_NUM_TYPES
)

/**
 * Definition of rational data type in {@link ACameraMetadata}.
 */
//typedef struct ACameraMetadata_rational {
//    int32_t numerator;
//    int32_t denominator;
//} ACameraMetadata_rational;
type MetadataRational C.ACameraMetadata_rational

func (r *MetadataRational) cptr() *C.ACameraMetadata_rational {
	return (*C.ACameraMetadata_rational)(r)
}

func (r *MetadataRational) Numerator() uint32 {
	return uint32(r.cptr().numerator)
}
func (r *MetadataRational) Denominator() uint32 {
	return uint32(r.cptr().denominator)
}

/**
 * A single camera metadata entry.
 *
 * <p>Each entry is an array of values, though many metadata fields may only have 1 entry in the
 * array.</p>
 */
//go typedef struct ACameraMetadata_entry {
/**
 * The tag identifying the entry.
 *
 * <p> It is one of the values defined in {@link NdkCameraMetadataTags.h}, and defines how the
 * entry should be interpreted and which parts of the API provide it.
 * See {@link NdkCameraMetadataTags.h} for more details. </p>
 */
//go     uint32_t tag;

/**
 * The data type of this metadata entry.
 *
 * <p>Must be one of ACAMERA_TYPE_* enum values defined above. A particular tag always has the
 * same type.</p>
 */
//go     uint8_t  type;

/**
 * Count of elements (NOT count of bytes) in this metadata entry.
 */
//go     uint32_t count;

/**
 * Pointer to the data held in this metadata entry.
 *
 * <p>The type field above defines which union member pointer is valid. The count field above
 * defines the length of the data in number of elements.</p>
 */
//go     union {
//go         uint8_t *u8;
//go         int32_t *i32;
//go         float   *f;
//go         int64_t *i64;
//go         double  *d;
//go         ACameraMetadata_rational* r;
//go     } data;
//go } ACameraMetadata_entry;
type MetadataEntry MetadataConstEntry

/**
 * A single read-only camera metadata entry.
 *
 * <p>Each entry is an array of values, though many metadata fields may only have 1 entry in the
 * array.</p>
 */
//go typedef struct ACameraMetadata_const_entry {
/**
 * The tag identifying the entry.
 *
 * <p> It is one of the values defined in {@link NdkCameraMetadataTags.h}, and defines how the
 * entry should be interpreted and which parts of the API provide it.
 * See {@link NdkCameraMetadataTags.h} for more details. </p>
 */
//go     uint32_t tag;

/**
 * The data type of this metadata entry.
 *
 * <p>Must be one of ACAMERA_TYPE_* enum values defined above. A particular tag always has the
 * same type.</p>
 */
//go     uint8_t  type;

/**
 * Count of elements (NOT count of bytes) in this metadata entry.
 */
//go     uint32_t count;

/**
 * Pointer to the data held in this metadata entry.
 *
 * <p>The type field above defines which union member pointer is valid. The count field above
 * defines the length of the data in number of elements.</p>
 */
//go     union {
//go         const uint8_t *u8;
//go         const int32_t *i32;
//go         const float   *f;
//go         const int64_t *i64;
//go         const double  *d;
//go         const ACameraMetadata_rational* r;
//go     } data;
//go } ACameraMetadata_const_entry;
type MetadataConstEntry C.ACameraMetadata_const_entry

func (entry *MetadataConstEntry) cptr() *C.ACameraMetadata_const_entry {
	return (*C.ACameraMetadata_const_entry)(entry)
}

func (entry *MetadataConstEntry) Tag() MetadataTag {
	return MetadataTag(entry.cptr().tag)
}

func (entry *MetadataConstEntry) Type() Type {
	return Type(entry.cptr()._type)
}

func (entry *MetadataConstEntry) Count() int {
	return int(entry.cptr().count)
}

func (entry *MetadataConstEntry) Data() interface{} {
	dataPtr := unsafe.Pointer(C.MetadataConstEntry_getData(entry.cptr()))
	switch entry.Type() {
	/// Unsigned 8-bit integer (uint8_t)
	case TYPE_BYTE:
		return (*[1 << 28]uint8)(dataPtr)[:entry.Count()]

	/// Signed 32-bit integer (int32_t)
	case TYPE_INT32:
		return (*[1 << 28]int32)(dataPtr)[:entry.Count()]

	/// 32-bit float (float)
	case TYPE_FLOAT:
		return (*[1 << 28]float32)(dataPtr)[:entry.Count()]

	/// Signed 64-bit integer (int64_t)
	case TYPE_INT64:
		return (*[1 << 27]int64)(dataPtr)[:entry.Count()]

	/// 64-bit float (double)
	case TYPE_DOUBLE:
		return (*[1 << 27]float64)(dataPtr)[:entry.Count()]

	/// A 64-bit fraction (ACameraMetadata_rational)
	case TYPE_RATIONAL:
		return (*[1 << 27]MetadataRational)(dataPtr)[:entry.Count()]

	default:
	}
	return nil
}

/**
 * Get a metadata entry from an input {@link ACameraMetadata}.
 *
 * <p>The memory of the data field in the returned entry is managed by camera framework. Do not
 * attempt to free it.</p>
 *
 * @param metadata the {@link ACameraMetadata} of interest.
 * @param tag the tag value of the camera metadata entry to be get.
 * @param entry the output {@link ACameraMetadata_const_entry} will be filled here if the method
 *        call succeeeds.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if metadata or entry is NULL.</li>
 *         <li>{@link ACAMERA_ERROR_METADATA_NOT_FOUND} if input metadata does not contain an entry
 *             of input tag value.</li></ul>
 */
//camera_status_t ACameraMetadata_getConstEntry(
//        const ACameraMetadata* metadata, uint32_t tag, /*out*/ACameraMetadata_const_entry* entry);
func (metadata *Metadata) GetConstEntry(tag MetadataTag) (*MetadataConstEntry, error) {
	var v C.ACameraMetadata_const_entry
	ret := Status(C.ACameraMetadata_getConstEntry(metadata.cptr(), C.uint32_t(tag), &v))
	if ret == nil {
		return (*MetadataConstEntry)(&v), ret
	}
	return nil, ret
}

/**
 * List all the entry tags in input {@link ACameraMetadata}.
 *
 * @param metadata the {@link ACameraMetadata} of interest.
 * @param numEntries number of metadata entries in input {@link ACameraMetadata}
 * @param tags the tag values of the metadata entries. Length of tags is returned in numEntries
 *             argument. The memory is managed by ACameraMetadata itself and must NOT be free/delete
 *             by application. Do NOT access tags after calling ACameraMetadata_free.
 *
 * @return <ul>
 *         <li>{@link ACAMERA_OK} if the method call succeeds.</li>
 *         <li>{@link ACAMERA_ERROR_INVALID_PARAMETER} if metadata, numEntries or tags is NULL.</li>
 *         <li>{@link ACAMERA_ERROR_UNKNOWN} if the method fails for some other reasons.</li></ul>
 */
//camera_status_t ACameraMetadata_getAllTags(
//       const ACameraMetadata* metadata, /*out*/int32_t* numEntries, /*out*/const uint32_t** tags);
func (metadata *Metadata) GetAllTags() ([]MetadataTag, error) {
	var numEntries C.int32_t
	var tags *C.uint32_t
	ret := Status(C.ACameraMetadata_getAllTags(metadata.cptr(), &numEntries, &tags))
	if ret == nil {
		return (*[1 << 28]MetadataTag)(unsafe.Pointer(tags))[:numEntries], ret
	} else {
		return nil, ret
	}
}

/**
 * Create a copy of input {@link ACameraMetadata}.
 *
 * <p>The returned ACameraMetadata must be freed by the application by {@link ACameraMetadata_free}
 * after application is done using it.</p>
 *
 * @param src the input {@link ACameraMetadata} to be copied.
 *
 * @return a valid ACameraMetadata pointer or NULL if the input metadata cannot be copied.
 */
//ACameraMetadata* ACameraMetadata_copy(const ACameraMetadata* src);
func (metadata *Metadata) Copy() *Metadata {
	return (*Metadata)(C.ACameraMetadata_copy(metadata.cptr()))
}

/**
 * Free a {@link ACameraMetadata} structure.
 *
 * @param metadata the {@link ACameraMetadata} to be freed.
 */
//void ACameraMetadata_free(ACameraMetadata* metadata);
func (metadata *Metadata) Free() {
	C.ACameraMetadata_free(metadata.cptr())
}

////
func (t Type) String() string {
	switch t {
	case TYPE_BYTE:
		return "BYTE"
	case TYPE_INT32:
		return "INT32"
	case TYPE_FLOAT:
		return "FLOAT"
	case TYPE_INT64:
		return "INT64"
	case TYPE_DOUBLE:
		return "DOUBLE"
	case TYPE_RATIONAL:
		return "RATIONAL"
	default:
		return "UNKNOW"
	}
}
