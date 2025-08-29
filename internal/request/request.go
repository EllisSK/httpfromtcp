package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLineString := strings.Split(string(data), "\r\n")[0]

	requestLine, err := parseRequestLine(requestLineString)
	if err != nil {
		return nil, err
	}

	request := Request{requestLine}

	return &request, nil
}

func parseRequestLine(input string) (RequestLine, error) {
	comps := strings.Split(input, " ")

	if len(comps) != 3 {
		return RequestLine{}, fmt.Errorf("incorrect number of components in request line")
	}

	HttpVersion := strings.Split(comps[2], "/")[1]
	if HttpVersion != "1.1" {
		return RequestLine{}, fmt.Errorf("incorrect http version")
	}

	RequestTarget := comps[1]

	Method := comps[0]
	if !isAllUpper(Method) {
		return RequestLine{}, fmt.Errorf("incorrect method")
	}

	return RequestLine{HttpVersion, RequestTarget, Method}, nil
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
