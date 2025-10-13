package response

import (
	"fmt"
	"ibrahemassa/http_bootdev/internal/headers"
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


func (w *Writer) WriteStatusLine(statusCode StatusCode) error{
	if w.writerState != WritingStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = WritingHeader }()

	var err error = nil
	reasonPhrase := ""
	switch statusCode {
	case OK:
		reasonPhrase = "OK"
	case BadRequest:
		reasonPhrase = "Bad Request"
	case InternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	_, err = w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase)))
	return err
}

func GetDefaultHeader(contentLen int, mimeType string) headers.Headers {
	h := headers.NewHeaders()
	h.Replace("content-length", fmt.Sprintf("%d", contentLen))
	h.Replace("connection", "close")
	if mimeType != ""{
		h.Replace("content-type", mimeType)
	} else {
		h.Replace("content-type", "text/plain")
	}

	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error{
	if w.writerState != WritingHeader {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = WritingBody }()

	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error){
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("cannot write status line in state %d", w.writerState)
	}

	n, err := w.Write(p)
	if err != nil{
		return 0, err
	}

	return n, nil
}
