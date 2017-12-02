/*
Cmdlog logs writable commands on disk. It works "almost like" Redis AOF
but simpler and dumber. To run command log you should add -cmdlog flag with path
to the file in which data will be stored.

  server -cmdlog /tmp/cmdlog.log

Every time server starts up, cmdlog restores everything from command log file
into storage. Cmdlog uses same protocol for read and write operations as the server.
*/
package main

import (
	"bufio"
	"io"
	"log"
	"os"

	redislike "github.com/bannerlog/redislike/protocol"
)

type cmdlog struct {
	file    *os.File
	logchan chan *request
}

func newCmdlog(filepath string) *cmdlog {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return &cmdlog{f, make(chan *request)}
}

func (l *cmdlog) listen() {
	for r := range l.logchan {
		<-l.logchan
		l.write(r)
	}
}

func (l *cmdlog) write(r *request) {
	req, err := redislike.NewRequest(r.cmd, r.argv...)
	if err != nil {
		log.Panicln("Could not save command request to disk")
	}

	req.Write(l.file)
}

func (l *cmdlog) restore(s *storage) {
	log.Printf("Restoring storage from %s\n", l.file.Name())
	fr := bufio.NewReader(l.file)
	for {
		req, err := redislike.ReadRequest(fr)
		if err != nil {
			if err != io.EOF {
				log.Print(err)
			}
			break
		}

		r := &request{
			cmd:  req.Command,
			argv: req.Args,
			argc: len(req.Args),
		}

		executeCmd(s, r)
	}

	log.Println("Storage restored successfully")
}

func (l *cmdlog) run(s *storage) {
	l.restore(s)
	go l.listen()
	setCommandLogger(l)
}
