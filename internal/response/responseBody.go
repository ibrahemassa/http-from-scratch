package response

import(
	"fmt"
)

func (w *Writer) WriteBody(p []byte) (int, error){
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}

	n, err := w.Write(p)
	if err != nil{
		return 0, err
	}

	return n, nil
}

func (w *Writer) writeChunk(p []byte) (int, error){
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("cannot write chunked body in state %d", w.writerState)
	}
	content := fmt.Sprintf("%X\r\n%s\r\n", len(p), p)
	_, err := w.WriteBody([]byte(content))
	if err != nil{
		return 0, err
	}

	return len(p), nil
}

func (w *Writer) WriteChunckedBody(p []byte) (int, error){
	if w.writerState != WritingBody {
		return 0, fmt.Errorf("cannot write chunked body done in state %d", w.writerState)
	}
	n := len(p)

	_, err := w.WriteBody([]byte(fmt.Sprintf("%X\r\n", n)))
	if err != nil{
			return 0, err
	}
	
	_, err = w.WriteBody(p[:n])
	if err != nil{
			return 0, err
	}

	_, err = w.WriteBody([]byte("\r\n"))
	if err != nil{
			return 0, err
	}

	return n, nil
}

func (w *Writer) WriteChunckedBodyDone() (int, error){
	_, err := w.WriteBody([]byte("0\r\n"))
	return 0, err
}



