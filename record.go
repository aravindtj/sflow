package sflow

import (
	"io"
)

type Record interface {
	RecordType() int
	Encode(io.Writer)
}
