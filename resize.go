package main

// #cgo pkg-config: GraphicsMagickWand
// #include <wand/magick_wand.h>
import "C"

import (
	"unsafe"
)

func resize(file string, width int, height int) []byte {
	wand := C.NewMagickWand()
	defer C.DestroyMagickWand(wand)

	C.MagickReadImage(wand, C.CString(file))

	cur_width := float64(C.MagickGetImageWidth(wand))
	cur_height := float64(C.MagickGetImageHeight(wand))

	r_width := cur_width / float64(width)
	r_height := cur_height / float64(height)

	ratio := r_height
	if r_width < r_height {
		ratio = r_width
	}

	dest_width := cur_width / ratio
	dest_height := cur_height / ratio

	crop_x := int((dest_width - float64(width)) / 2)
	crop_y := int((dest_height - float64(height)) / 2)

	C.MagickStripImage(wand)
	if r_width > 5 && r_height > 5 {
		C.MagickSampleImage(wand, (C.ulong)(dest_width * 5), (C.ulong)(dest_height * 5))
	}
	C.MagickResizeImage(wand, (C.ulong)(dest_width), (C.ulong)(dest_height), C.LanczosFilter, 1)
	C.MagickCropImage(wand, (C.ulong)(width), (C.ulong)(height), (C.long)(crop_x), (C.long)(crop_y))

	size := 0
	buf := C.MagickWriteImageBlob(wand, (*C.size_t)(unsafe.Pointer(&size)))
	defer C.MagickRelinquishMemory(unsafe.Pointer(buf))

	return C.GoBytes(unsafe.Pointer(buf), C.int(size))
}
