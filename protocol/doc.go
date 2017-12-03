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

Rsponse message structure looks like this. As in the request part above,
the response must have CRLF at the end of each line.

  Response type (could be OK or ERR)
  Number of parts
  Part length
  Response part
  ...

Example of OK response.

  OK\r\n
  1\r\n
  18\r\n
  A hash field 1\r\n
  32\r\n
  Second value of a hash field\r\n

Example of ERR response.

  ERR\r\n
  1\r\n
  29\r\n
  Wrong number of arguments\r\n

*/
package redislike
