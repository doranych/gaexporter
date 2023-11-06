package export

import "io"

type OutputType int

const (
	OutputTypeFile OutputType = iota + 1
	OutputTypeStdout
)

type Output struct {
	Type OutputType
	Dest io.Writer
}
