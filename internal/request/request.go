package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"ibrahemassa/http_bootdev/internal/headers"
)

type ParserState int

const (
	initialized    ParserState = 0
	done           ParserState = 1
	parsingHeaders ParserState = 2
	parsingBody    ParserState = 3
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad request-line!")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version!")
var ERROR_READ_DONE = fmt.Errorf("trying to read data in done state!")
var ERROR_UNKOWN_STATE = fmt.Errorf("unkown state!")
var ERROR_UNEXPECTED_EOF = fmt.Errorf("unexpected eof!")
var ERROR_BODY_LENGTH_EXCEEDED = fmt.Errorf("body length is greater than content-length!")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		State:   initialized,
		Headers: headers.NewHeaders(),
	}
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {
	totalParsed := 0
	for r.State != done {
		n, err := r.parseSingle(data[totalParsed:])
		if err != nil {
			return 0, nil
		}
		totalParsed += n
		if n == 0 {
			break
		}
	}
	return totalParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case initialized:
		rl, noBytes, _, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if noBytes == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.State = parsingHeaders
		return noBytes, nil

	case parsingHeaders:
		noBytes, finish, err := r.Headers.Parse(data)
		if err != nil {

			return 0, err
		}

		if finish {
			if isCon := r.Headers.Get("content-length"); isCon != "" {
				r.State = parsingBody
			} else {
				r.State = done
			}
		}
		return noBytes, nil

	case parsingBody:
		conLen, err := strconv.Atoi(r.Headers.Get("content-length"))
		if err != nil {
			return 0, err
		}

		if conLen == 0 {
			r.State = done
			return 0, nil
		}

		if len(data) > conLen {
			return 0, ERROR_BODY_LENGTH_EXCEEDED
		} else if len(data) < conLen {
			return 0, nil
		}

		r.Body = append(r.Body, data...)
		r.State = done

		return 0, nil
	case done:
		return 1, ERROR_READ_DONE

	default:
		fmt.Println("You're such a loser and should stop coding!")
		return 1, ERROR_UNKOWN_STATE
	}
}

func parseRequestLine(request []byte) (*RequestLine, int, bool, error) {
	idx := bytes.Index(request, SEPARATOR)
	done := false
	if idx == -1 {
		return nil, 0, done, nil
	}

	startLine := request[:idx]
	read := idx + len(SEPARATOR)

	tokens := bytes.Split(startLine, []byte(" "))

	if len(tokens) != 3 {
		return nil, read, done, ERROR_BAD_START_LINE
	}

	method := string(tokens[0])
	target := string(tokens[1])
	httpParts := strings.Split(string(tokens[2]), "/")

	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, read, done, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	if !IsUpper(method) {
		return nil, read, done, ERROR_BAD_START_LINE
	}

	rl := &RequestLine{httpParts[1], target, method}

	done = true
	return rl, read, done, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1024)
	idx := 0
	for {
		if request.State == done {
			break
		}
		n, err := reader.Read(buff[idx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				// if request.State != done{
				if request.State == parsingHeaders {
					return nil, fmt.Errorf("incomplete ")
				} else if request.State == parsingBody {
					return nil, ERROR_BODY_LENGTH_EXCEEDED
				}
				break
			}
			return nil, err
		}

		idx += n
		i, err := request.parse(buff[:idx])

		if err != nil {
			return nil, err
		}

		copy(buff, buff[i:idx])
		idx -= i
		if request.State == done {
			if i == 0 && idx == 0 {
				break
			}
		}
	}
	printRequest(*request)
	return request, nil
}

func printRequest(r Request) {
	fmt.Println("Request line:")
	fmt.Println("- Method:", r.RequestLine.Method)
	fmt.Println("- Target:", r.RequestLine.RequestTarget)
	fmt.Println("- Version:", r.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for key, value := range r.Headers {
		if strings.Contains(value, "cur") {
			value = "curl"
		}
		fmt.Printf("- %s: %s\n", key, value)
	}
	fmt.Println("Body:")
	fmt.Println(string(r.Body))
}
