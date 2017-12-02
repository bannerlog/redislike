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

// ErrEmptyCommand rises when command header is not set
var ErrEmptyCommand = errors.New("Request: Command is not set")

// ErrBadRequest rises on malformed request
var ErrBadRequest = errors.New("Request: Bad request")

// A Request represents an request received by a server
type Request struct {
	Command string
	Args    []string
}

func (r *Request) String() string {
	s := fmt.Sprintf("%d\r\n%d\r\n%s\r\n", len(r.Args)+1, len(r.Command)+2, r.Command)

	for _, v := range r.Args {
		s = fmt.Sprintf("%s%d\r\n%s\r\n", s, len(v)+2, v)
	}

	return s
}

func (r *Request) Write(w io.Writer) error {
	if len(r.Command) < 1 {
		return ErrEmptyCommand
	}

	_, err := w.Write([]byte(r.String()))
	if err != nil {
		return err
	}

	return nil
}

// ReadRequest reads and returns an request from r.
func ReadRequest(r *bufio.Reader) (*Request, error) {
	// read the number of parts
	s, err := r.ReadString('\n')
	if err != nil {
		// if err == io.EOF {
		// 	err = io.ErrUnexpectedEOF
		// }
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

	if len(parts) < 1 {
		return nil, ErrBadRequest
	}

	return &Request{parts[0], parts[1:]}, nil
}

// NewRequest returns a new Request given a comamnd and optional args.
func NewRequest(cmd string, args ...string) (*Request, error) {
	if len(cmd) < 1 {
		return nil, ErrEmptyCommand
	}

	return &Request{cmd, args}, nil
}
