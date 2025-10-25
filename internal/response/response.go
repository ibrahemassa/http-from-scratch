package response

import (
	"io"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

type WriterState int

const(
	WritingStatusLine WriterState = 0
	WritingHeader WriterState = 1
	WritingBody WriterState = 2
)

type Writer struct{
	ioW io.Writer
	writerState WriterState
}

func (w *Writer) Write(data []byte) (int ,error){
	n, err := w.ioW.Write(data)
	return n, err
}

func NewWriter(ioW io.Writer) *Writer{
	return &Writer{
		ioW: ioW,
		writerState: WritingStatusLine,
	}
}



