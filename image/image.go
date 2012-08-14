package image

// #cgo pkg-config: GraphicsMagickWand
// #include <wand/magick_wand.h>
import "C"

import (
	"unsafe"
	"os"
	"io/ioutil"
)

type Image struct {
	wand *C.MagickWand
}

func Open(filename string) (img *Image, err error) {
	img = &Image{wand: C.NewMagickWand()}
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return nil, err
	}
	var data []byte
	if data, err = ioutil.ReadAll(file); err != nil {
		return nil, err
	}

	C.MagickReadImageBlob(img.wand, (*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)))

	return img, nil
}

func (img *Image) Close() {
	C.DestroyMagickWand(img.wand)
}

func (img *Image) Strip() {
	C.MagickStripImage(img.wand)
}

func (img *Image) Resize(width int, height int) {
	wand := img.wand

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

	if r_width > 5 && r_height > 5 {
		C.MagickSampleImage(wand, (C.ulong)(dest_width * 5), (C.ulong)(dest_height * 5))
	}
	C.MagickResizeImage(wand, (C.ulong)(dest_width), (C.ulong)(dest_height), C.LanczosFilter, 1)
	C.MagickCropImage(wand, (C.ulong)(width), (C.ulong)(height), (C.long)(crop_x), (C.long)(crop_y))
}

func (img *Image) Grayscale() {
	C.MagickQuantizeImage(img.wand, 256, C.GRAYColorspace, 1, 0, 0)
}

func (img *Image) GetBytes() []byte {
	size := 0
	buf := C.MagickWriteImageBlob(img.wand, (*C.size_t)(unsafe.Pointer(&size)))
	defer C.MagickRelinquishMemory(unsafe.Pointer(buf))

	return C.GoBytes(unsafe.Pointer(buf), C.int(size))
}

