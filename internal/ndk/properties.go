package app

/*
#include <stdlib.h>
#include <sys/system_properties.h>
*/
import "C"

import (
	"unsafe"
)

const (
	PROP_NAME_MAX  = C.PROP_NAME_MAX
	PROP_VALUE_MAX = C.PROP_VALUE_MAX
)

func PropGet(k string) string {
	var value [PROP_VALUE_MAX]C.char
	key := C.CString(k)
	defer C.free(unsafe.Pointer(key))
	n := C.__system_property_get((*C.char)(key), &value[0])
	if n > 0 {
		return C.GoString(&value[0])
	}
	return ""
}

func PropVisit(cb func(k, v string)) {
	var name [PROP_NAME_MAX]C.char
	var value [PROP_VALUE_MAX]C.char
	var pi *C.prop_info

	for n := 0; ; n++ {
		pi = C.__system_property_find_nth(C.uint(n))
		if pi == nil {
			break
		}
		C.__system_property_read(pi, &name[0], &value[0])
		cb(C.GoString(&name[0]), C.GoString(&value[0]))
	}
}
