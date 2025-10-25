package response

import(
	"ibrahemassa/http_bootdev/internal/headers"
	"fmt"
)

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
