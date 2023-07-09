package svg

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"

	"github.com/nfnt/resize"
)

var (
	icnsHeader  = []byte{0x69, 0x63, 0x6e, 0x73}
	sizes       = []int{16, 32, 64, 128, 256, 512, 1024}
	sizeToTypes = map[int][]string{
		16:   {"icp4"},
		32:   {"icp5", "ic11"},
		64:   {"icp6", "ic12"},
		128:  {"ic07"},
		256:  {"ic08", "ic13"},
		512:  {"ic09", "ic14"},
		1024: {"ic10"},
	}
)

func iconset(img image.Image) ([]byte, error) {
	iconBuffers := new(bytes.Buffer)
	for _, s := range sizes {
		imgBuf := new(bytes.Buffer)
		resized := resize.Resize(uint(s), uint(s), img, resize.MitchellNetravali)
		if err := png.Encode(imgBuf, resized); err != nil {
			return nil, err
		}

		lenByt := make([]byte, 4, 4)
		binary.BigEndian.PutUint32(lenByt, uint32(imgBuf.Len()+8))

		// Iterate through every OSType and append the icon to iconBuffers
		for _, ostype := range sizeToTypes[s] {
			if _, err := iconBuffers.Write([]byte(ostype)); err != nil {
				return nil, err
			}
			if _, err := iconBuffers.Write(lenByt); err != nil {
				return nil, err
			}
			if _, err := iconBuffers.Write(imgBuf.Bytes()); err != nil {
				return nil, err
			}
		}
	}

	lenByt := make([]byte, 4, 4)
	binary.BigEndian.PutUint32(lenByt, uint32(iconBuffers.Len()+8))

	data := icnsHeader
	data = append(data, lenByt...)
	data = append(data, iconBuffers.Bytes()...)

	return data, nil
}
