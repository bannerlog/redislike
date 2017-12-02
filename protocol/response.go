package redislike

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	errType = "ERR"
	okType  = "OK"
)

// ErrWrongResponseType rises when status header is incorrect
var ErrWrongResponseType = errors.New("Response type must be OK or ERR")

// Response represents the response from an request.
type Response struct {
	Type   string
	Values []string
}

func (r *Response) String() string {
	s := fmt.Sprintf("%s\r\n%d\r\n", r.Type, len(r.Values))

	for _, v := range r.Values {
		s = fmt.Sprintf("%s%d\r\n%s\r\n", s, len(v)+2, v)
	}

	return s
}

// IsErr checks if the response is an error status response
func (r *Response) IsErr() bool {
	return r.Type == errType
}

// IsOk checks if the response is an OK status response
func (r *Response) IsOk() bool {
	return r.Type == okType
}

func (r *Response) checkType() error {
	if r.Type == okType || r.Type == errType {
		return nil
	}

	return ErrWrongResponseType
}

func (r *Response) Write(w io.Writer) error {
	if err := r.checkType(); err != nil {
		return err
	}

	_, err := w.Write([]byte(r.String()))
	if err != nil {
		return err
	}

	return nil
}

// ReadResponse reads and returns an response from r.
func ReadResponse(r *bufio.Reader) (*Response, error) {
	resp := new(Response)

	// read response type
	s, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	resp.Type = strings.TrimSpace(s)

	// read the number of parts
	s, err = r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	num, _ := strconv.Atoi(strings.TrimSpace(s))

	// iterate through parts of the request
	parts := []string{}
	for i := 0; i < num; i++ {

		// read next part length
		s, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}
		len, _ := strconv.Atoi(strings.TrimSpace(s))

		// read next part value
		part := bytes.NewBuffer(make([]byte, 0, len))
		if len > 0 {
			buflen := 1024 * 16
			if len < buflen {
				buflen = len
			}
			buf := make([]byte, buflen)
			for {
				remaining := int(len) - part.Len()
				if remaining < 1 {
					break
				}
				n, err := r.Read(buf)
				if n > remaining {
					n = remaining
				}
				part.Write(buf[0:n])

				if err != nil {
					if err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
					return nil, err
				}
			}
		}
		parts = append(parts, strings.TrimSpace(part.String()))
	}

	resp.Values = parts

	return resp, nil
}

// NewOkResponse returns a new OK Response given a body.
func NewOkResponse(args ...string) (*Response, error) {
	return NewResponse(okType, args...)
}

// NewErrResponse returns a new error Response given a status and optional body.
func NewErrResponse(args ...string) (*Response, error) {
	return NewResponse(errType, args...)
}

// NewResponse returns a new Response given a status and optional body.
func NewResponse(rtype string, args ...string) (*Response, error) {
	r := Response{rtype, args}
	if err := r.checkType(); err != nil {
		return nil, err
	}

	return &r, nil
}
