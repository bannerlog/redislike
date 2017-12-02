/*
A client communicates with a server by creating a TCP connection to a specific port
(by default it's 9000). Protocol by it self looks simliar to the RESP (REdis Serialization
Protocol).

Request message to the server consists of parts shown below.

  Number of parts
  Part length
  Command name
  Part length
  Command argument
  ...

There must be at least a command name part and each line should be ended with CRLF.
Here is an example of PING request.

  1\r\n
  8\r\n
  PING\r\n

Server doesn't close a connection after response was sent. It keeps connection
open until the client won't close it. So you can send as many request as you want.
That is just a test application, so no limits except hardware.

*/
package redislike
