package pnm

import (
	"bufio"
	"image"
	"io"
)

func Decode(r io.Reader) (image.Image, error) {
	return nil, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	return image.Config{}, nil
}

func decodeHeader(r io.Reader) (header, error) {
	return header{}, nil
}

type header struct {
	typ  int
	w, h int
}

type decoder struct {
	r io.Reader
	b *bufio.Reader
	header
}

func init() {
	image.RegisterFormat("pbm", "P1", Decode, DecodeConfig)
	image.RegisterFormat("pgm", "P2", Decode, DecodeConfig)
	image.RegisterFormat("ppm", "P3", Decode, DecodeConfig)
	image.RegisterFormat("pbm", "P4", Decode, DecodeConfig)
	image.RegisterFormat("pgm", "P5", Decode, DecodeConfig)
	image.RegisterFormat("ppm", "P6", Decode, DecodeConfig)
}
