package redislike

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/bannerlog/redislike/protocol"
)

// ErrCommandResult ...
type ErrCommandResult struct {
	s string
}

func (e *ErrCommandResult) Error() string {
	return e.s
}

// Client represents a wrapepr for server requests and response.
type Client struct {
	conn net.Conn
}

// NewClient returns a new Client given an ip and port of a server.
// If for some reason error occurs, Client will be nil and
// error will be given as a second value. Otherwise second return value is nil.
func NewClient(ip string, port uint16) (*Client, error) {
	c := Client{}
	if err := c.connect(ip, port); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) connect(ip string, port uint16) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	c.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	return nil
}

// Ping pong
func (c *Client) Ping() (string, error) {
	var result string
	return result, c.genericCommand(&result, "PING")
}

// Info returns state of the storage (golang representation of entries and expiries).
func (c *Client) Info() (string, error) {
	var result string
	return result, c.genericCommand(&result, "INFO")
}
func (c *Client) Summary() (string, error) {
	var result string
	return result, c.genericCommand(&result, "INFO", "summary")
}

// Keys returns all keys
func (c *Client) Keys() ([]string, error) {
	var result []string
	return result, c.genericCommand(&result, "KEYS")
}

// Set key to hold the string value. If key already holds a value,
// it is overwritten, regardless of its type. Any previous time
// to live associated with the key is discarded on successful SET operation.
// Return the number of keys that were set.
func (c *Client) Set(key string, value string) (int, error) {
	var result int
	err := c.genericCommand(&result, "SET", key, value)
	return result, err
}

// Get the value of key. If the key does not exist empty string is returned.
// An error is returned if the value stored at key is not a string,
// because GET command only handles string values.
func (c *Client) Get(key string) (string, error) {
	var result string
	return result, c.genericCommand(&result, "GET", key)
}

// Del removes the specified key. A key is ignored if it does not exist.
// Return the number of keys that were removed.
func (c *Client) Del(key string) (int, error) {
	var result int
	return result, c.genericCommand(&result, "DEL", key)
}

// Expire set a timeout on key. After the timeout has expired,
// the key will automatically be deleted. Returns 1 if the timeout was set.
func (c *Client) Expire(key string, sec int) (int, error) {
	var result int
	return result, c.genericCommand(&result, "EXPIRE", key, strconv.Itoa(sec))
}

func (c *Client) LPush(key string, values []string) (int, error) {
	return c.genericPush("LPUSH", key, values)
}

func (c *Client) RPush(key string, values []string) (int, error) {
	return c.genericPush("RPUSH", key, values)
}

func (c *Client) genericPush(cmd string, key string, values []string) (int, error) {
	v := append([]string{key}, values...)
	resp, err := c.request(cmd, v...)
	if err != nil {
		return 0, err
	}

	r, _ := strconv.Atoi(resp.Values[0])
	return r, nil
}

func (c *Client) LSet(key string, idx int, value string) (int, error) {
	var result int
	err := c.genericCommand(&result, "LSET", key, strconv.Itoa(idx), value)
	return result, err
}

func (c *Client) LPop(key string) (string, error) {
	var result string
	return result, c.genericCommand(&result, "LPOP", key)
}

func (c *Client) RPop(key string) (string, error) {
	var result string
	return result, c.genericCommand(&result, "RPOP", key)
}

func (c *Client) LRange(key string, left int, right int) ([]string, error) {
	var result []string
	err := c.genericCommand(&result, "LRANGE", key, strconv.Itoa(left), strconv.Itoa(right))
	return result, err
}

func (c *Client) LIndex(key string, idx int) (string, error) {
	var result string
	err := c.genericCommand(&result, "LINDEX", key, strconv.Itoa(idx))
	return result, err
}

func (c *Client) LLen(key string) (int, error) {
	var result int
	return result, c.genericCommand(&result, "LLEN", key)
}

func (c *Client) HSet(key string, field string, value string) (int, error) {
	var result int
	return result, c.genericCommand(&result, "HSET", key, field, value)
}

func (c *Client) HGet(key string, field string) (string, error) {
	var result string
	return result, c.genericCommand(&result, "HGET", key, field)
}

func (c *Client) HGetAll(key string) (map[string]string, error) {
	var result map[string]string
	return result, c.genericCommand(&result, "HGETALL", key)
}

func (c *Client) HKeys(key string) ([]string, error) {
	var result []string
	return result, c.genericCommand(&result, "HKEYS", key)
}

func (c *Client) HVals(key string) ([]string, error) {
	var result []string
	return result, c.genericCommand(&result, "HVALS", key)
}

func (c *Client) HExists(key string, field string) (int, error) {
	var result int
	return result, c.genericCommand(&result, "HEXISTS", key, field)
}

func (c *Client) HDel(key string, fields []string) (int, error) {
	var result int
	args := append([]string{key}, fields...)
	return result, c.genericCommand(&result, "HDEL", args...)
}

func (c *Client) HLen(key string) (int, error) {
	var result int
	return result, c.genericCommand(&result, "HLEN", key)
}

func (c *Client) genericCommand(result interface{}, cmd string, args ...string) error {
	resp, err := c.request(cmd, args...)
	if err != nil {
		return err
	}

	unmarshalString(resp.Values[0], &result)

	return nil
}

// Send request to the server and return response or error
func (c *Client) request(cmd string, args ...string) (*redislike.Response, error) {
	// send request
	req, err := redislike.NewRequest(cmd, args...)
	if err != nil {
		return nil, err
	}

	if err = req.Write(c.conn); err != nil {
		return nil, err
	}

	// get response
	resp, err := redislike.ReadResponse(bufio.NewReader(c.conn))
	if err != nil {
		return nil, err
	}

	if resp.IsErr() {
		return resp, &ErrCommandResult{resp.Values[0]}
	}

	return resp, nil
}

// Close connection with the server.
func (c *Client) Close() {
	c.conn.Close()
}

func unmarshalString(s string, umval interface{}) error {
	if umval != nil && len(s) > 0 {
		if err := json.Unmarshal([]byte(s), &umval); err != nil {
			return err
		}
	}

	return nil
}
