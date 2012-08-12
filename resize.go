package main

// #cgo pkg-config: MagickWand
// #include <wand/magick_wand.h>
import "C"

import (
	"unsafe"
)

var free_wands = make(chan *C.MagickWand, 10)

func init() {
	C.MagickWandGenesis()
}

func resize(file string, width int, height int) []byte {
	var wand *C.MagickWand
	select {
		case wand = <-free_wands:
		default:
			wand = C.NewMagickWand()
	}
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
		C.MagickSampleImage(wand, (C.size_t)(dest_width * 5), (C.size_t)(dest_height * 5))
	}
	C.MagickResizeImage(wand, (C.size_t)(dest_width), (C.size_t)(dest_height), C.LanczosFilter, 1)
	C.MagickCropImage(wand, (C.size_t)(width), (C.size_t)(height), (C.ssize_t)(crop_x), (C.ssize_t)(crop_y))
	C.MagickSetImageCompressionQuality(wand, 65)

	size := 0
	buf := C.MagickGetImageBlob(wand, (*C.size_t)(unsafe.Pointer(&size)))
	defer C.MagickRelinquishMemory(unsafe.Pointer(buf))

	C.ClearMagickWand(wand)
	select {
		case free_wands <- wand:
		default:
			C.DestroyMagickWand(wand)
	}

	return C.GoBytes(unsafe.Pointer(buf), C.int(size))
}
