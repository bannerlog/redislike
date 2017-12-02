package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type request struct {
	cmd  string
	argv []string
	argc int
}

var (
	cmdList = map[string]struct {
		fn    func(*storage, *request) (interface{}, error)
		write int
	}{
		"set":     {setCommand, 1},
		"get":     {getCommand, 0},
		"del":     {delCommand, 1},
		"exists":  {existsCommand, 0},
		"expire":  {expireCommand, 1},
		"lpush":   {lpushCommand, 1},
		"rpush":   {rpushCommand, 1},
		"llen":    {llenCommand, 0},
		"lindex":  {lindexCommand, 0},
		"lrange":  {lrangeCommand, 0},
		"lset":    {lsetCommand, 1},
		"lpop":    {lpopCommand, 1},
		"rpop":    {rpopCommand, 1},
		"hset":    {hsetCommand, 1},
		"hget":    {hgetCommand, 0},
		"hgetall": {hgetallCommand, 0},
		"hexists": {hexistsCommand, 0},
		"hvals":   {hvalsCommand, 0},
		"hdel":    {hdelCommand, 1},
		"hkeys":   {hkeysCommand, 0},
		"hlen":    {hlenCommand, 0},
		"keys":    {keysCommand, 0},
		"info":    {infoCommand, 0},
		"ping":    {pingCommand, 0},
	}

	// ErrWrongNumOfArguments ...
	ErrWrongNumOfArguments = errors.New("Wrong number of arguments")
	// ErrBadArguments ...
	ErrBadArguments = errors.New("Invalid arguments of the command")
	// ErrOperationAgainstWrongType ...
	ErrOperationAgainstWrongType = errors.New("Operation against a key holding the wrong kind of value")

	cmdlogger *cmdlog
)

func setCommandLogger(l *cmdlog) {
	cmdlogger = l
}

func executeCmd(s *storage, r *request) (string, error) {
	if c, ok := cmdList[strings.ToLower(r.cmd)]; ok {
		res, err := c.fn(s, r)

		b, e := json.Marshal(res)
		if e != nil {
			panic("Cannot json encode command result")
		}

		if err == nil && c.write == 1 && cmdlogger != nil {
			cmdlogger.logchan <- r
		}

		return string(b), err
	}

	return "", fmt.Errorf("ERR Unknown command '%s'", r.cmd)
}

// EXISTS key
func existsCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	return s.exists(r.argv[0]), nil
}

// DEL key
func delCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 1 {
		return nil, ErrWrongNumOfArguments
	}

	s.del(r.argv[0])

	return 1, nil
}

// EXPIRE key seconds
func expireCommand(s *storage, r *request) (interface{}, error) {
	if r.argc != 2 {
		return nil, ErrWrongNumOfArguments
	}

	sec, err := strconv.ParseInt(r.argv[1], 10, 64)
	if err != nil {
		return nil, ErrBadArguments
	}

	s.setExpire(r.argv[0], time.Now().Unix()+sec)
	return 1, nil

}

// KEYS
func keysCommand(s *storage, r *request) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	sl := []string{}
	for k := range s.entries {
		sl = append(sl, k)
	}

	return sl, nil
}

// INFO [summary]
func infoCommand(s *storage, r *request) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(r.argv) == 1 && r.argv[0] == "summary" {
		kln, xln := s.len()
		return fmt.Sprintf("Number of Keys: %d\nNumber of Expiries %d", kln, xln), nil
	}

	return fmt.Sprintf("%+v", s), nil
}

// PING
func pingCommand(s *storage, r *request) (interface{}, error) {
	return "PONG", nil
}
