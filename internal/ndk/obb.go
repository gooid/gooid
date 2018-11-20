package app

/*
#include <android/obb.h>
*/
import "C"
import "unsafe"

type ObbInfo C.AObbInfo

func (info *ObbInfo) cptr() *C.AObbInfo {
	return (*C.AObbInfo)(info)
}

const AOBBINFO_OVERLAY = 0x0001

/**
 * Scan an OBB and get information about it.
 */
//AObbInfo* AObbScanner_getObbInfo(const char* filename);
func GetObbInfo(filename string) *ObbInfo {
	cfilename := []byte(filename + "\000")
	return (*ObbInfo)(C.AObbScanner_getObbInfo((*C.char)(unsafe.Pointer(&cfilename[0]))))
}

/**
 * Destroy the AObbInfo object. You must call this when finished with the object.
 */
//void AObbInfo_delete(AObbInfo* obbInfo);
func (info *ObbInfo) Delete() {
	C.AObbInfo_delete(info.cptr())
}

/**
 * Get the package name for the OBB.
 */
//const char* AObbInfo_getPackageName(AObbInfo* obbInfo);
func (info *ObbInfo) GetPackageName() string {
	return C.GoString(C.AObbInfo_getPackageName(info.cptr()))
}

/**
 * Get the version of an OBB file.
 */
//int32_t AObbInfo_getVersion(AObbInfo* obbInfo);
func (info *ObbInfo) GetVersion() int {
	return int(C.AObbInfo_getVersion(info.cptr()))
}

/**
 * Get the flags of an OBB file.
 */
//int32_t AObbInfo_getFlags(AObbInfo* obbInfo);
func (info *ObbInfo) GetFlags() int {
	return int(C.AObbInfo_getFlags(info.cptr()))
}
