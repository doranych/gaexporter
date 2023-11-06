package export

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"io"
)

type Input interface {
	io.Reader
	io.Seeker
}

type InputFormat int

const (
	InputFormatImage = iota + 1
	InputFormatText
)

func getInputFormat(input Input) (InputFormat, error) {
	_, _, err := image.Decode(input)
	defer input.Seek(0, io.SeekStart)
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return InputFormatText, nil
		}
		return 0, err
	}
	return InputFormatImage, nil
}
