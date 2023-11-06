package export

import (
	"github.com/pkg/errors"
)

func RunQR(input Input, output Output) error {
	err := processInput(input, output, "qr")
	if err != nil {
		return errors.Wrap(err, "failed to process input")
	}
	return nil
}
