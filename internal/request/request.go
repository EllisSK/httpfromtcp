package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type ParserState int

const (
	StateInit = iota
	StateDone
)

const (
	bufferSize = 8
)

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case StateInit:
		rl, bytesConsumed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if bytesConsumed == 0 {
			return 0, nil
		}

		r.RequestLine = rl
		r.State = StateDone
		return bytesConsumed, nil

	case StateDone:
		return 0, fmt.Errorf("trying to read data in a done state")
	}

	return 0, fmt.Errorf("unknown state")
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0

	req := Request{
		State: StateInit,
	}

	for req.State != StateDone {
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		read, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if readToIndex == 0 {
					break
				}
			} else {
				return &req, err
			}
		}

		readToIndex += read

		bc, err := req.parse(buf[:readToIndex])
		if err != nil {
			return &req, err
		}

		if bc > 0 {
			copy(buf, buf[bc:readToIndex])
			readToIndex -= bc
		}

		if err == io.EOF {
			break
		}
	}

	return &req, nil
}

func parseRequestLine(input string) (RequestLine, int, error) {
	seperatedInput := strings.Split(input, "\r\n")

	if len(seperatedInput) == 1 {
		return RequestLine{}, 0, nil
	}

	requestLineString := seperatedInput[0]

	comps := strings.Split(requestLineString, " ")

	if len(comps) != 3 {
		return RequestLine{}, 0, fmt.Errorf("incorrect number of components in request line")
	}

	HttpVersion := strings.Split(comps[2], "/")[1]
	if HttpVersion != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("incorrect http version")
	}

	RequestTarget := comps[1]

	Method := comps[0]
	if !isAllUpper(Method) {
		return RequestLine{}, 0, fmt.Errorf("incorrect method")
	}

	bytesConsumed := len(requestLineString) + 2

	return RequestLine{HttpVersion, RequestTarget, Method}, bytesConsumed, nil
}

func isAllUpper(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
