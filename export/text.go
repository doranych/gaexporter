package export

import (
	"github.com/pkg/errors"
)

func RunText(input Input, output Output) error {
	err := processInput(input, output, "txt")
	if err != nil {
		return errors.Wrap(err, "failed to process input")
	}
	return nil
}
