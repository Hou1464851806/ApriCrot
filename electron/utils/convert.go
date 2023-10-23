package utils

import (
	"syscall"
	"unsafe"
)

func Str2Ptr(s string) uintptr {
	p, _ := syscall.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(p))
}

func Ptr2Str(p uintptr) string {
	return syscall.UTF16ToString(*(*[]uint16)(unsafe.Pointer(&p)))
}

func Bool2Ptr(b bool) uintptr {
	if b {
		return uintptr(1)
	}
	return uintptr(0)
}

func Ptr2Bool(p uintptr) bool {
	if p == 0 {
		return false
	}
	return true
}

func Float32Ptr(f float32) uintptr {
	return uintptr(*(*uint32)(unsafe.Pointer(&f)))
}

func PtrToFloat32(p uintptr) float32 {
	u := uint32(p)
	return *(*float32)(unsafe.Pointer(&u))
}

func Int2Ptr(i int) uintptr {
	return uintptr(i)
}

func Ptr2Int(u uintptr) int {
	return int(u)
}
