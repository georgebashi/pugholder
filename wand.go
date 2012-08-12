package main

// #cgo pkg-config: MagickWand
// #include <wand/magick_wand.h>
import "C"

import (
	"unsafe"
	"fmt"
)
func init() {
	C.MagickWandGenesis()
}

func resize(file string, width int, height int) []byte {
	wand := C.NewMagickWand()
	defer C.DestroyMagickWand(wand)
	C.MagickReadImage(wand, C.CString(file))

	cur_width := float64(int(C.MagickGetImageWidth(wand)))
	cur_height := float64(int(C.MagickGetImageHeight(wand)))

	r_width := float64(cur_width) / float64(width)
	r_height := float64(cur_height) / float64(height)

	ratio := 1.0
	if r_width < r_height {
		ratio = r_width
	} else {
		ratio = r_height
	}

	C.MagickResizeImage(wand, (C.size_t)(cur_width/ratio), (C.size_t)(cur_height/ratio), C.LanczosFilter, 1)
	ex := C.MagickGetExceptionType(wand)
	fmt.Println(ex)

	size := 0
	buf := C.MagickGetImageBlob(wand, (*C.size_t)(unsafe.Pointer(&size)))
	defer C.MagickRelinquishMemory(unsafe.Pointer(buf))

	return C.GoBytes(unsafe.Pointer(buf), C.int(size))
}
