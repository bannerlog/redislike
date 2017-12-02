package main

import "errors"

// ErrHashEmpty ...
var ErrHashEmpty = errors.New("Hash is empty")

func findHashEntry(s *storage, k string) (map[string]string, error) {
	v := s.get(k)
	if v == nil {
		return make(map[string]string, 0), nil
	}

	if v, ok := v.(map[string]string); ok {
		return v, nil
	}

	return nil, ErrOperationAgainstWrongType
}

// HSET key field value
func hsetCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 3 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}

	h[r.argv[1]] = r.argv[2]
	s.set(r.argv[0], h)

	return 1, nil
}

// HGET key field
func hgetCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 2 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}

	if v, ok := h[r.argv[1]]; ok {
		return v, nil
	}

	return nil, nil
}

// HGETALL key
func hgetallCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}
	if len(h) < 1 {
		return nil, ErrHashEmpty
	}

	return h, nil
}

// HEXISTS key field
func hexistsCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 2 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}

	if _, ok := h[r.argv[1]]; ok {
		return 1, nil
	}

	return 0, nil
}

// HVALS key
func hvalsCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}
	if len(h) < 1 {
		return nil, ErrHashEmpty
	}

	vs := make([]string, 0, len(h))
	for _, i := range h {
		vs = append(vs, i)
	}
	return vs, nil

}

// HDEL key field [field...]
func hdelCommand(s *storage, r *request) (interface{}, error) {
	if r.argc < 2 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}
	if len(h) < 1 {
		return 0, nil
	}

	var deleted int
	sln := len(h)
	for i := 1; i < r.argc; i++ {
		delete(h, r.argv[i])
	}
	fln := len(h)
	deleted = sln - fln

	if fln < 1 {
		s.del(r.argv[0])
	} else {
		s.set(r.argv[0], h)
	}

	return deleted, nil
}

// HKEYS key
func hkeysCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}
	if len(h) < 1 {
		return nil, ErrHashEmpty
	}

	ks := make([]string, 0, len(h))
	for k := range h {
		ks = append(ks, k)
	}
	return ks, nil

}

// HLEN key
func hlenCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	h, err := findHashEntry(s, r.argv[0])
	if err != nil {
		return 0, err
	}

	return len(h), nil
}
