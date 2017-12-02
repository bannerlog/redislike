package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	redislike "github.com/bannerlog/redislike/protocol"
)

var flagServerAddress string
var flagCmdlogFilename string

func init() {
	// runtime.GOMAXPROCS(1)
	flag.StringVar(&flagServerAddress, "addr", ":9000", "Start server on host:port")
	flag.StringVar(&flagCmdlogFilename, "cmdlog", "", "Path to command log file")
}

func main() {
	flag.Parse()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v\nShutting down server...\n", sig)
		// Here we can do something usefuly (i.e. gracefuly close all connections
		// and dump strorage to disk).
		os.Exit(0)
	}()

	// storage
	s := newStorage()
	go runExpireMonitor(s)

	// cmdlog
	if flagCmdlogFilename != "" {
		l := newCmdlog(flagCmdlogFilename)
		l.run(s)
	}

	// tcp listner
	li, err := net.Listen("tcp", flagServerAddress)
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()
	log.Printf("Server is running on %s\n", flagServerAddress)
	log.Println("Ready to accept connections")

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handleConnection(s, conn)
	}
}

func handleConnection(s *storage, conn net.Conn) {
	log.Printf("Open connection from %s\n", conn.RemoteAddr())
	defer func() {
		conn.Close()
		log.Printf("Close connection from %s\n", conn.RemoteAddr())
	}()

	for {
		// request part
		req, err := redislike.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		r, err := executeCmd(s, &request{
			cmd:  req.Command,
			argv: req.Args,
			argc: len(req.Args),
		})

		// response part
		var resp *redislike.Response
		if err != nil {
			resp, err = redislike.NewErrResponse(err.Error())
		} else {
			resp, err = redislike.NewOkResponse(r)
		}
		if err != nil {
			log.Println(err)
			return
		}

		resp.Write(conn)
	}
}

func runExpireMonitor(s *storage) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			s.removeExpired()
		}
	}
}
