package main

import (
	"strconv"
	"time"
)

// GET key
func getCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	v := s.get(r.argv[0])
	if v == nil {
		return nil, nil
	}

	if v, ok := v.(string); ok {
		return v, nil
	}

	return nil, ErrOperationAgainstWrongType
}

// SET key value [ttl]
func setCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 2 && r.argc != 3 {
		return 0, ErrWrongNumOfArguments
	}

	if !s.set(r.argv[0], r.argv[1]) {
		return 0, nil
	}

	if r.argc == 3 {
		sec, _ := strconv.ParseInt(r.argv[2], 10, 64)
		s.setExpire(r.argv[0], time.Now().Unix()+sec)
	}

	return 1, nil
}
