package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var crlf = []byte("\r\n")

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeaders(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed filed line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

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

		name, value, err := parseHeaders(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		if !isToken([]byte(name)) {
			return 0, false, fmt.Errorf("invalid character in header name: %q", name)
		}

		read += idx + len(crlf)

		h[strings.ToLower(name)] = value
	}

	return read, done, nil
}

func isToken(str []byte) bool {
	for _, b := range str {
		isAlphanumeric := (b >= 'A' && b <= 'Z') ||
			(b >= 'a' && b <= 'z') ||
			(b >= '0' && b <= '9')

		if isAlphanumeric {
			continue
		}

		switch b {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}
