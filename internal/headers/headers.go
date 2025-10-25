package headers

import (
	"bytes"
	"fmt"
	// "io"
	"strings"
)

type Headers map[string]string

var crlf = []byte("\r\n")
var specialChars = map[byte]bool{
	'!': true, '#': true, '$': true, '%': true, '&': true,
	'\'': true, '*': true, '+': true, '-': true, '.': true,
	'^': true, '_': true, '`': true, '|': true, '~': true,
}
var ERROR_FIELD_NAME_WHITESPACE = fmt.Errorf("whitespace in the field name!")
var ERROR_BAD_FIELD_LINE = fmt.Errorf("bad field line!")
var ERROR_KEY_NOT_FOUND = fmt.Errorf("key not found in the headers map!")
var ERROR_KEY_INVALID = fmt.Errorf("key contains invalid character!")

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Replace(key string, value string) {
	h[strings.ToLower(key)] = value
}

func (h Headers) Set(key string, value string) {
	key = strings.ToLower(key)
	old, ok := h[key]
	if ok {
		h[strings.ToLower(key)] = old + ", " + value
	} else {
		h[strings.ToLower(key)] = value
	}
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

func (h Headers) Delete(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func keyIsValid(key string) bool {
	if len(key) < 1 {
		return false
	}

	for i := 0; i < len(key); i++ {
		c := key[i]

		if (c >= 'A' && c <= 'Z') ||
			(c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			specialChars[c] {
			continue
		}

		return false
	}
	return true
}

func parseFieldLine(line []byte) (string, string, error) {
	tokens := bytes.SplitN(line, []byte(":"), 2)
	if len(tokens) != 2 {
		return "", "", ERROR_BAD_FIELD_LINE
	}

	if bytes.Contains(tokens[0], []byte(" ")) {
		return "", "", ERROR_FIELD_NAME_WHITESPACE
	}

	key := string(tokens[0])
	field := string(bytes.TrimSpace(tokens[1]))

	return key, field, nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	var e error = nil
	// idx := bytes.Index(data[read:], crlf)

	// if idx  == -1 {
	// 	return 0, done, nil
	// }
	//
	// if idx == 0{
	// 	return 2, true, nil
	// }

	for {
		idx := bytes.Index(data[read:], crlf)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			read += len(crlf)
			break
		}

		line := data[read : read+idx]

		key, field, err := parseFieldLine(line)
		if err != nil {
			e = err
			break
		}

		if keyIsValid(key) {
		} else {
			e = ERROR_KEY_INVALID
			break
		}

		h.Set(key, field)
		read += idx + len(crlf)
	}
	//
	// if !done && e == nil{
	// 	return read, done, io.EOF
	// }
	return read, done, e
}
