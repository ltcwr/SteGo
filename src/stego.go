package main

import (
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
)


func bytesToBits(data []byte) []byte {
	bits := make([]byte, 0, len(data)*8)
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (b>>i)&1)
		}
	}
	return bits
}

func bitsToBytes(bits []byte) []byte {
	bytes := make([]byte, len(bits)/8)
	for i := 0; i < len(bytes); i++ {
		var b byte
		for j := 0; j < 8; j++ {
			b = (b << 1) | bits[i*8+j]
		}
		bytes[i] = b
	}
	return bytes
}

func xorBits(bits, keyBits []byte) {
	for i := range bits {
		bits[i] ^= keyBits[i%len(keyBits)]
	}
}


func encode(input, output, message, key string) error {
	if key == "" {
		return errors.New("empty secret key")
	}

	msgBytes := []byte(message)
	keyBits := bytesToBits([]byte(key))
	msgBits := bytesToBits(msgBytes)

	xorBits(msgBits, keyBits)


	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(msgBytes)))
	lenBits := bytesToBits(lenBuf)

	finalBits := append(lenBits, msgBits...)

	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	capacity := bounds.Dx() * bounds.Dy()
	if len(finalBits) > capacity {
		return errors.New("image is too small")
	}

	outImg := image.NewRGBA(bounds)
	bitIndex := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)

			if bitIndex < len(finalBits) {
				bit := finalBits[bitIndex]
				if bit == 0 && r8%2 == 1 {
					r8--
				} else if bit == 1 && r8%2 == 0 {
					r8++
				}
				bitIndex++
			}

			outImg.Set(x, y, color.RGBA{
				r8,
				uint8(g >> 8),
				uint8(b >> 8),
				uint8(a >> 8),
			})
		}
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, outImg)
}



func decode(input, key string) (string, error) {
	if key == "" {
		return "", errors.New("empty secret key")
	}

	file, err := os.Open(input)
	if err != nil {
		return "", err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width := bounds.Dx()

	getBit := func(i int) byte {
		x := bounds.Min.X + (i % width)
		y := bounds.Min.Y + (i / width)
		r, _, _, _ := img.At(x, y).RGBA()
		return byte((r >> 8) & 1)
	}

	
	lenBits := make([]byte, 32)
	for i := 0; i < 32; i++ {
		lenBits[i] = getBit(i)
	}
	length := binary.BigEndian.Uint32(bitsToBytes(lenBits))

	if length == 0 {
		return "", errors.New("nothing found")
	}

	msgBits := make([]byte, length*8)
	for i := 0; i < int(length*8); i++ {
		msgBits[i] = getBit(32 + i)
	}

	keyBits := bytesToBits([]byte(key))
	xorBits(msgBits, keyBits)

	return string(bitsToBytes(msgBits)), nil
}

