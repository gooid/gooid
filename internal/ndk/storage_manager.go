package app

/*
#include <stdlib.h>
#include <android/storage_manager.h>
extern void cgoObbCallbackFunc(char* filename, int32_t state, void* data);
*/
import "C"

import "unsafe"

type StorageManager C.AStorageManager

const (
	/*
	 * The OBB container is now mounted and ready for use. Can be returned
	 * as the status for callbacks made during asynchronous OBB actions.
	 */
	OBB_STATE_MOUNTED = C.AOBB_STATE_MOUNTED

	/*
	 * The OBB container is now unmounted and not usable. Can be returned
	 * as the status for callbacks made during asynchronous OBB actions.
	 */
	OBB_STATE_UNMOUNTED = C.AOBB_STATE_UNMOUNTED

	/*
	 * There was an internal system error encountered while trying to
	 * mount the OBB. Can be returned as the status for callbacks made
	 * during asynchronous OBB actions.
	 */
	OBB_STATE_ERROR_INTERNAL = C.AOBB_STATE_ERROR_INTERNAL

	/*
	 * The OBB could not be mounted by the system. Can be returned as the
	 * status for callbacks made during asynchronous OBB actions.
	 */
	OBB_STATE_ERROR_COULD_NOT_MOUNT = C.AOBB_STATE_ERROR_COULD_NOT_MOUNT

	/*
	 * The OBB could not be unmounted. This most likely indicates that a
	 * file is in use on the OBB. Can be returned as the status for
	 * callbacks made during asynchronous OBB actions.
	 */
	OBB_STATE_ERROR_COULD_NOT_UNMOUNT = C.AOBB_STATE_ERROR_COULD_NOT_UNMOUNT

	/*
	 * A call was made to unmount the OBB when it was not mounted. Can be
	 * returned as the status for callbacks made during asynchronous OBB
	 * actions.
	 */
	OBB_STATE_ERROR_NOT_MOUNTED = C.AOBB_STATE_ERROR_NOT_MOUNTED

	/*
	 * The OBB has already been mounted. Can be returned as the status for
	 * callbacks made during asynchronous OBB actions.
	 */
	OBB_STATE_ERROR_ALREADY_MOUNTED = C.AOBB_STATE_ERROR_ALREADY_MOUNTED

	/*
	 * The current application does not have permission to use this OBB.
	 * This could be because the OBB indicates it's owned by a different
	 * package. Can be returned as the status for callbacks made during
	 * asynchronous OBB actions.
	 */
	OBB_STATE_ERROR_PERMISSION_DENIED = C.AOBB_STATE_ERROR_PERMISSION_DENIED
)

/**
 * Obtains a new instance of AStorageManager.
 */
//AStorageManager* AStorageManager_new();
func NewStorageManager() *StorageManager {
	return (*StorageManager)(C.AStorageManager_new())
}
func (mgr *StorageManager) cptr() *C.AStorageManager {
	return (*C.AStorageManager)(mgr)
}

/**
 * Release AStorageManager instance.
 */
//void AStorageManager_delete(AStorageManager* mgr);
func (mgr *StorageManager) Delete() {
	C.AStorageManager_delete(mgr.cptr())
}

/**
 * Callback function for asynchronous calls made on OBB files.
 */
//typedef void (*AStorageManager_obbCallbackFunc)(const char* filename, const int32_t state, void* data);
type ObbCallbackFunc func(filename string, state int)

//export cgoObbCallbackFunc
func cgoObbCallbackFunc(filename *C.char, state C.int, data unsafe.Pointer) {
	(*(*ObbCallbackFunc)(data))(C.GoString(filename), int(state))
}

/**
 * Attempts to mount an OBB file. This is an asynchronous operation.
 */
//void AStorageManager_mountObb(AStorageManager* mgr, const char* filename, const char* key,
//        AStorageManager_obbCallbackFunc cb, void* data);
func (mgr *StorageManager) MountObb(filename, key string, cb ObbCallbackFunc) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	C.AStorageManager_mountObb(mgr.cptr(), cfilename, ckey,
		C.AStorageManager_obbCallbackFunc(C.cgoObbCallbackFunc), unsafe.Pointer(&cb))
}

/**
 * Attempts to unmount an OBB file. This is an asynchronous operation.
 */
//void AStorageManager_unmountObb(AStorageManager* mgr, const char* filename, const int force,
//        AStorageManager_obbCallbackFunc cb, void* data);
func (mgr *StorageManager) UnmountObb(filename string, force bool, cb ObbCallbackFunc) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	cforce := 0
	if force {
		cforce = 1
	}
	C.AStorageManager_unmountObb(mgr.cptr(), cfilename, C.int(cforce),
		C.AStorageManager_obbCallbackFunc(C.cgoObbCallbackFunc), unsafe.Pointer(&cb))
}

/**
 * Check whether an OBB is mounted.
 */
//int AStorageManager_isObbMounted(AStorageManager* mgr, const char* filename);
func (mgr *StorageManager) IsObbMounted(filename string) bool {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return 0 != C.AStorageManager_isObbMounted(mgr.cptr(), cfilename)
}

/**
 * Get the mounted path for an OBB.
 */
//const char* AStorageManager_getMountedObbPath(AStorageManager* mgr, const char* filename);
func (mgr *StorageManager) GetMountedObbPath(filename string) string {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return C.GoString(C.AStorageManager_getMountedObbPath(mgr.cptr(), cfilename))
}
