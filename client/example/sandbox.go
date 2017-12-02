package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/bannerlog/redislike/client"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	if len(os.Args) < 1 {
		os.Exit(1)
	}

	fnmap := map[string]func(){
		"generate": generate,
		"summary":  summary,
		"info":     info,
		"test":     test,
		"expire":   expire,
	}

	if fn, ok := fnmap[os.Args[1]]; ok {
		fn()
	}
}

func generate() {
	c := conn()
	defer c.Close()

	var n int
	if len(os.Args) > 1 {
		n, _ = strconv.Atoi(os.Args[2])
	}
	if n < 1 {
		n = 1
	}

	for i := 0; i < n; i++ {

		k := randString(randNumber(5, 25))
		_, err := c.Set(k, randString(randNumber(100, 1000)))
		if err != nil {
			fmt.Printf("\n\n%s\n", err)
			break
		}
		exp := randNumber(0, 1000)
		if l := exp % 2; l != 0 {
			_, err = c.Expire(k, exp)
			if err != nil {
				fmt.Printf("\n\n%s\n", err)
				break
			}
		}

		_, err = c.LPush(randString(randNumber(5, 25)), []string{"1", "2", "3", "4", "5"})
		if err != nil {
			fmt.Printf("\n\n%s\n", err)
			break
		}
	}
}

func summary() {
	c := conn()
	defer c.Close()

	print(c.Summary())
}

func info() {
	c := conn()
	defer c.Close()

	print(c.Info())
}

func expire() {
	c := conn()
	defer c.Close()

	print(c.Ping())
	print(c.Set("key-1", "value-1"))
	print(c.Set("key-2", "value-2"))
	print(c.Set("key-3", "value-3"))
	print(c.Expire("key-1", 3))
	print(c.Expire("key-2", 7))
}

func test() {
	c := conn()
	defer c.Close()

	print(c.Ping())

	// overview key base
	print(c.Info())

	{ // STRING
		fmt.Println("-----  STRING -----")
		print(c.Set("mykey", "my string value"))
		print(c.Expire("mykey", 100))
		print(c.Get("mykey"))
		print(c.Info())
		// print(c.Del("mykey"))
		print(c.Info())
	}

	{ // LIST
		fmt.Println("-----  LIST -----")
		print(c.LPush("mylist", []string{"1", "2", "3"}))
		print(c.RPush("mylist", []string{"40", "50", "60"}))
		print(c.LRange("mylist", 1, -2))
		print(c.LRange("mylist", 2, 1))
		print(c.LSet("mylist", 4, "79"))
		print(c.LIndex("mylist", 4))
		print(c.LPop("mylist"))
		print(c.RPop("mylist"))
		print(c.LLen("mylist"))
		print(c.Info())
		// c.Del("mylist")
	}

	{ // HASH
		fmt.Println("-----  HASH -----")
		print(c.HSet("myhash", "field", "value"))
		print(c.HSet("myhash", "test me", "Hello, 世界"))
		print(c.HSet("myhash", "123", "98457 HG"))
		print(c.HLen("myhash"))
		print(c.HGetAll("myhash"))
		print(c.HGet("myhash", "test me"))
		print(c.HKeys("myhash"))
		print(c.HVals("myhash"))
		print(c.HExists("myhash", "123"))
		print(c.HExists("myhash", "noname"))
		print(c.HDel("myhash", []string{"field", "nokey", "123"}))
		print(c.Info())
		c.Del("myhash")
	}

	{ // finaly
		fmt.Println("----- finaly -----")
		c.LPush("mylist", []string{"1", "2", "3"})
		c.HSet("myhash", "field", "value")
		c.Set("mykey", "my string value")
		print(c.Keys())
		c.Del("mylist")
		// c.Del("myhash")
		// c.Del("mykey")
	}

}

func conn() *redislike.Client {
	c, err := redislike.NewClient("127.0.0.1", 9000)
	if err != nil {
		log.Fatalln(err)
	}

	return c
}

func print(r interface{}, err interface{}) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func randNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
