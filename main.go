package main

import (
	"bytes"
	"fmt"
	"os"

	"image"
	_ "image/png"
)

type Converter interface {
	Convert(in string, out string) error
}

func rgb332(r, g, b uint32, buf *bytes.Buffer) {
	// Convert to RGB332 format
	r = r >> 5
	g = g >> 5
	b = b >> 6
	pixel := byte((r << 5) | (g << 2) | b)
	buf.WriteByte(pixel)
}
func rgb565(r, g, b uint32, buf *bytes.Buffer) {
	// Convert to RGB565 format
	r = r >> 3
	g = g >> 2
	b = b >> 3
	value := (r << 11) | (g << 5) | b
	buf.WriteByte(byte(value >> 8))
	buf.WriteByte(byte(value & 0xFF))
}

func converterImpl(in string, out string, conv func(uint32, uint32, uint32, *bytes.Buffer)) error {
	imageFile, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("error opening image file: %w", err)
	}
	defer imageFile.Close()
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var buf bytes.Buffer
	buf.WriteByte(byte(width & 0xFF))
	buf.WriteByte(byte((width >> 8) & 0xFF))
	buf.WriteByte(byte(height & 0xFF))
	buf.WriteByte(byte((height >> 8) & 0xFF))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// RGBA values are 16-bit, we need to convert them to 8-bit
			r = r >> 8
			g = g >> 8
			b = b >> 8

			conv(r, g, b, &buf)
		}
	}
	outputFile, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()
	_, err = outputFile.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}

	return nil
}

type RGB332Converter struct{}

func (c *RGB332Converter) Convert(in string, out string) error {
	return converterImpl(in, out, rgb332)
}

type RGB565Converter struct{}

func (c *RGB565Converter) Convert(in string, out string) error {
	return converterImpl(in, out, rgb565)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("usage: %v <rbg332 or rgb565> <:inputFileName> <:outputFilename>", os.Args[0])
		return
	}

	var converter Converter
	switch os.Args[1] {
	case "rgb332":
		converter = &RGB332Converter{}
	case "rgb565":
		converter = &RGB565Converter{}
	default:
		fmt.Printf("error: %v is not a valid converter(should be rgb332 or rgb565)", os.Args[1])
		return
	}

	err := converter.Convert(os.Args[2], os.Args[3])
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
}
