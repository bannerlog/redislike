package redislike

import (
	"bufio"
	"fmt"
	"log"

	"github.com/bannerlog/redislike/client"
)

func ExampleNewRequest() {
	conn, err := redislike.NewClient("127.0.0.1", 9000)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	req, err := redislike.NewRequest(cmd, "INFO", "summary")
	if err != nil {
		panic(err)
	}

	if err = req.Write(c.conn); err != nil {
		return nil, err
	}
}

func ExampleNewResponse() {
	conn, err := redislike.NewClient("127.0.0.1", 9000)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	resp, err := redislike.ReadResponse(bufio.NewReader(conn))
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v", resp)
}
