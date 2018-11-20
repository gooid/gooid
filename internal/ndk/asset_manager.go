package app

/*
#include <stdlib.h>
#include <android/asset_manager.h>
#if __ANDROID_API__ < 13
static off64_t AAsset_seek64(AAsset* asset, off64_t offset, int whence) {
	return AAsset_seek(asset, offset, whence);
}
static off64_t AAsset_getLength64(AAsset* asset) {
	return AAsset_getLength(asset);
}
static off64_t AAsset_getRemainingLength64(AAsset* asset) {
	return AAsset_getRemainingLength(asset);
}
#endif
*/
import "C"

import (
	"fmt"
	"io"
	"unsafe"
)

type AssetManager C.AAssetManager

func (mgr *AssetManager) cptr() *C.AAssetManager {
	return (*C.AAssetManager)(mgr)
}

type AssetDir C.AAssetDir

func (assetDir *AssetDir) cptr() *C.AAssetDir {
	return (*C.AAssetDir)(assetDir)
}

type Asset C.AAsset

func (asset *Asset) cptr() *C.AAsset {
	return (*C.AAsset)(asset)
}

/* Available modes for opening assets */
const (
	ASSET_MODE_UNKNOWN   = C.AASSET_MODE_UNKNOWN
	ASSET_MODE_RANDOM    = C.AASSET_MODE_RANDOM
	ASSET_MODE_STREAMING = C.AASSET_MODE_STREAMING
	ASSET_MODE_BUFFER    = C.AASSET_MODE_BUFFER
)

/**
 * Open the named directory within the asset hierarchy.  The directory can then
 * be inspected with the AAssetDir functions.  To open the top-level directory,
 * pass in "" as the dirName.
 *
 * The object returned here should be freed by calling AAssetDir_close().
 */
//AAssetDir* AAssetManager_openDir(AAssetManager* mgr, const char* dirName);
func (mgr *AssetManager) OpenDir(dirName string) *AssetDir {
	cdirName := C.CString(dirName)
	defer C.free(unsafe.Pointer(cdirName))
	return (*AssetDir)(C.AAssetManager_openDir(mgr.cptr(), cdirName))
}

/**
 * Open an asset.
 *
 * The object returned here should be freed by calling AAsset_close().
 */
//AAsset* AAssetManager_open(AAssetManager* mgr, const char* filename, int mode);
func (mgr *AssetManager) Open(filename string, mode int) *Asset {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	return (*Asset)(C.AAssetManager_open(mgr.cptr(), cfilename, C.int(mode)))
}

/**
 * Iterate over the files in an asset directory.  A NULL string is returned
 * when all the file names have been returned.
 *
 * The returned file name is suitable for passing to AAssetManager_open().
 *
 * The string returned here is owned by the AssetDir implementation and is not
 * guaranteed to remain valid if any other calls are made on this AAssetDir
 * instance.
 */
//const char* AAssetDir_getNextFileName(AAssetDir* assetDir);
func (assetDir *AssetDir) GetNextFileName() string {
	return C.GoString(C.AAssetDir_getNextFileName(assetDir.cptr()))
}

/**
 * Reset the iteration state of AAssetDir_getNextFileName() to the beginning.
 */
//void AAssetDir_rewind(AAssetDir* assetDir);
func (assetDir *AssetDir) Rewind() {
	C.AAssetDir_rewind(assetDir.cptr())
}

/**
 * Close an opened AAssetDir, freeing any related resources.
 */
//void AAssetDir_close(AAssetDir* assetDir);
func (assetDir *AssetDir) Close() {
	C.AAssetDir_close(assetDir.cptr())
}

/**
 * Attempt to read 'count' bytes of data from the current offset.
 *
 * Returns the number of bytes read, zero on EOF, or < 0 on error.
 */
//int AAsset_read(AAsset* asset, void* buf, size_t count);
func (asset *Asset) Read(buf []byte) (int, error) {
	ret := int(C.AAsset_read(asset.cptr(), unsafe.Pointer(&buf[0]), C.size_t(len(buf))))
	if ret == 0 {
		return ret, io.EOF
	} else if ret < 0 {
		return 0, fmt.Errorf("ASSET: Error code (%d)", ret)
	}
	return ret, nil
}

/**
 * Seek to the specified offset within the asset data.  'whence' uses the
 * same constants as lseek()/fseek().
 *
 * Returns the new position on success, or (off_t) -1 on error.
 */
//off_t AAsset_seek(AAsset* asset, off_t offset, int whence);

/**
 * Seek to the specified offset within the asset data.  'whence' uses the
 * same constants as lseek()/fseek().
 *
 * Uses 64-bit data type for large files as opposed to the 32-bit type used
 * by AAsset_seek.
 *
 * Returns the new position on success, or (off64_t) -1 on error.
 */
//off64_t AAsset_seek64(AAsset* asset, off64_t offset, int whence);
func (asset *Asset) Seek(offset int64, whence int) (int64, error) {
	ret := int64(C.AAsset_seek64(asset.cptr(), C.off64_t(offset), C.int(whence)))
	if ret < 0 {
		return 0, fmt.Errorf("ASSET: Seek fail.")
	}
	return ret, nil
}

/**
 * Close the asset, freeing all associated resources.
 */
//void AAsset_close(AAsset* asset);
func (asset *Asset) Close() error {
	C.AAsset_close(asset.cptr())
	return nil
}

/**
 * Get a pointer to a buffer holding the entire contents of the assset.
 *
 * Returns NULL on failure.
 */
//const void* AAsset_getBuffer(AAsset* asset);
func (asset *Asset) GetBuffer() []byte {
	cptr := C.AAsset_getBuffer(asset.cptr())
	if cptr == nil {
		return nil
	}
	size := asset.Length()
	return (*[1 << 26]byte)(unsafe.Pointer(cptr))[:size]
}

/**
 * Report the total size of the asset data.
 */
//off_t AAsset_getLength(AAsset* asset);

/**
 * Report the total size of the asset data. Reports the size using a 64-bit
 * number insted of 32-bit as AAsset_getLength.
 */
//off64_t AAsset_getLength64(AAsset* asset);
func (asset *Asset) Length() int64 {
	return int64(C.AAsset_getLength64(asset.cptr()))
}

/**
 * Report the total amount of asset data that can be read from the current position.
 */
//off_t AAsset_getRemainingLength(AAsset* asset);

/**
 * Report the total amount of asset data that can be read from the current position.
 *
 * Uses a 64-bit number instead of a 32-bit number as AAsset_getRemainingLength does.
 */
//off64_t AAsset_getRemainingLength64(AAsset* asset);
func (asset *Asset) GetRemainingLength() int64 {
	return int64(C.AAsset_getRemainingLength64(asset.cptr()))
}

/**
 * Open a new file descriptor that can be used to read the asset data. If the
 * start or length cannot be represented by a 32-bit number, it will be
 * truncated. If the file is large, use AAsset_openFileDescriptor64 instead.
 *
 * Returns < 0 if direct fd access is not possible (for example, if the asset is
 * compressed).
 */
//int AAsset_openFileDescriptor(AAsset* asset, off_t* outStart, off_t* outLength);

/**
 * Open a new file descriptor that can be used to read the asset data.
 *
 * Uses a 64-bit number for the offset and length instead of 32-bit instead of
 * as AAsset_openFileDescriptor does.
 *
 * Returns < 0 if direct fd access is not possible (for example, if the asset is
 * compressed).
 */
//int AAsset_openFileDescriptor64(AAsset* asset, off64_t* outStart, off64_t* outLength);
//func (asset *Asset) OpenFileDescriptor() int {}

/**
 * Returns whether this asset's internal buffer is allocated in ordinary RAM (i.e. not
 * mmapped).
 */
//int AAsset_isAllocated(AAsset* asset);
func (asset *Asset) IsAllocated() bool {
	return 0 != C.AAsset_isAllocated(asset.cptr())
}
