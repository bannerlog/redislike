# redislike

Simple implementation of Redis-like in-memory cache.

#### How to install
```bash
go get github.com/bannerlog/redislike/server
```

### Documentation
Documentation and Go API for the client package could be read with help of _godoc_. For example:

```
godoc -http 3000

http://localhost:3000/pkg/github.com/bannerlog/redislike/
```

## Usage
### Server
Server is listening on port 9000 by default, though you are free to choose any port you like. Just add the flag -port and port number after the command.

```bash
server -port 9001
```

The easy way to check that server is up and running by sending simple PING command:

```bash
echo "2\r\n8\r\nPING\r\n\r\n" | nc localhost 9000
```

#### Persistence
Cmdlog logs writable commands on disk. It works "almost like" Redis AOF
but simpler and dumber. To run command log you should add -cmdlog flag with path
to the file in which data will be stored.

```bash
server -cmdlog /tmp/cmdlog.log
```

Every time server starts up, cmdlog restores everything from command log file
into storage. Cmdlog uses same protocol for read and write operations as the server.

### Client
Here is simplest client which sends PING command to the server and get response.

```go
package main
import (
  "fmt"
  "log"

  "github.com/bannerlog/redislike/client"
)

func main() {
  c, err := redislike.NewClient("127.0.0.1", 9000)
  if err != nil {
    log.Fatalln(err)
    return nil
  }

  if r, err != c.Ping(); e != nil {
    fmt.Println(r)
  }
}
```

Use _godoc_ to find some more examples and read documentation.
