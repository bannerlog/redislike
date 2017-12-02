package main

import (
	"errors"
	"strconv"
)

const (
	listHead = 0
	listTale = 1
)

// ErrListValueOutOfRange ...
var ErrListValueOutOfRange = errors.New("ERR Value is out of range")

// ErrListEmpty ...
var ErrListEmpty = errors.New("List is empty")

func findListEntry(s *storage, k string) ([]string, error) {
	v := s.get(k)
	if v == nil {
		return make([]string, 0), nil
	}

	if v, ok := v.([]string); ok {
		return v, nil
	}

	return nil, ErrOperationAgainstWrongType
}

// LPUSH key value [value ...]
// Return value is the length of the list after the push operations.
func lpushCommand(s *storage, r *request) (interface{}, error) {
	if r.argc < 2 {
		return nil, ErrWrongNumOfArguments
	}

	return pushGenericCommand(s, r, listTale)
}

// RPUSH key value [value ...]
// Return value is the length of the list after the push operations.
func rpushCommand(s *storage, r *request) (interface{}, error) {
	if r.argc < 2 {
		return nil, ErrWrongNumOfArguments
	}

	return pushGenericCommand(s, r, listHead)
}

// Return value is the length of the list after the push operations.
func pushGenericCommand(s *storage, r *request, where int) (int, error) {
	l, err := findListEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}

	if where == listHead {
		l = append(r.argv[1:], l...)
	} else {
		l = append(l, r.argv[1:]...)
	}

	s.set(r.argv[0], l)

	return len(l), nil
}

// LLEN key
func llenCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	l, err := findListEntry(s, r.argv[0])
	if err != nil {
		return nil, err
	}

	return len(l), nil
}

// LINDEX key index
func lindexCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 2 {
		return nil, ErrWrongNumOfArguments
	}

	idx, err := strconv.Atoi(r.argv[1])
	if err != nil {
		return nil, ErrBadArguments
	}

	l, err := findListEntry(s, r.argv[0])
	if err != nil {
		return nil, err
	}
	if idx <= len(l)-1 {
		return l[idx], nil
	}

	return nil, nil
}

// LRANGE key start stop
func lrangeCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 3 {
		return nil, ErrWrongNumOfArguments
	}

	start, strerr := strconv.Atoi(r.argv[1])
	stop, stperr := strconv.Atoi(r.argv[2])
	if strerr != nil || stperr != nil {
		return nil, ErrListValueOutOfRange
	}

	list, err := findListEntry(s, r.argv[0])
	if err != nil {
		return nil, err
	}

	len := len(list)
	if len < 1 {
		return nil, ErrListEmpty
	}

	if start >= len {
		return nil, ErrListValueOutOfRange
	}
	if start < 0 {
		if start = len + start; start < 0 {
			start = 0
		}
	}
	if stop < 0 {
		if stop = len + stop + 1; stop < 0 {
			stop = 0
		}
	}
	if stop > len {
		stop = len
	}
	if start > stop {
		return nil, ErrListValueOutOfRange
	}

	return list[start:stop], nil
}

// LSET key index value
func lsetCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 3 {
		return nil, ErrWrongNumOfArguments
	}

	list, err := findListEntry(s, r.argv[0])
	if err != nil {
		return nil, err
	}

	idx, _ := strconv.Atoi(r.argv[1])
	len := len(list)
	if idx < 0 || idx > len {
		return nil, ErrListValueOutOfRange
	}

	list[idx] = r.argv[2]
	s.set(r.argv[0], list)

	return 1, nil
}

// LPOP key
func lpopCommand(s *storage, r *request) (interface{}, error) {
	return popGenericCommand(s, r, listHead)
}

// RPOP key
func rpopCommand(s *storage, r *request) (interface{}, error) {
	return popGenericCommand(s, r, listTale)
}

func popGenericCommand(s *storage, r *request, where int) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	list, err := findListEntry(s, r.argv[0])
	if err != nil {
		return nil, err
	}
	if len(list) < 1 {
		return nil, nil
	}

	var x string
	if where == listHead {
		x, list = list[0], list[1:]
	} else {
		x, list = list[len(list)-1], list[:len(list)-1]
	}

	if len(list) > 0 {
		s.set(r.argv[0], list)
	} else {
		s.del(r.argv[0])
	}

	return x, nil
}
