package export

import (
	"bufio"
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
)

func processInput(input Input, output Output, fmt string) error {
	format, err := getInputFormat(input)
	if err != nil {
		return errors.Wrap(err, "failed to get input format")
	}

	switch format {
	case InputFormatImage:
		img, _, _ := image.Decode(input)

		bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
		qrReader := qrcode.NewQRCodeReader()

		result, _ := qrReader.Decode(bmp, nil)

		err = processMigrationUrl(result.String(), output, fmt)
		if err != nil {
			return errors.Wrap(err, "failed to process migration url")
		}
	case InputFormatText:
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			str := scanner.Text()
			err = processMigrationUrl(str, output, fmt)
			if err != nil {
				return errors.Wrap(err, "failed to process migration url")
			}
		}
	}
	return nil
}
