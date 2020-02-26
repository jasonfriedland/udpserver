/*
An echo UDP server; useful for mocking a backend required for tests. You could
use it for sending StatsD metrics like this:

	s := udpserver.New(8125)
	s.Serve()

	m := metrics.New("127.0.0.1:8125", "hello")
	m.Count("foo", 73)
	m.Count("bar", 79)

	s.Close()

Which would result in the following similar output:

	2020/02/26 12:18:01 server listening on: 0.0.0.0:8125
	2020/02/26 12:18:01 received: hello.foo:73|c|#env:dev
	2020/02/26 12:18:01 received: hello.bar:79|c|#env:dev
	2020/02/26 12:18:01 shutting down...

*/

package udpserver

import (
	"log"
	"net"
)

// bufSize is the size of the incoming data buffer.
const bufSize = 1024

// UDPServer wraps a UDP connection.
type UDPServer struct {
	port int          // port to listen on
	conn *net.UDPConn // underlying connection
	done chan bool    // signals shutdown
}

// New builds a new UDPServer, listening on the specified port.
func New(port int) *UDPServer {
	return &UDPServer{
		port: port,
		done: make(chan bool),
	}
}

// Serve listens on the specified port and prints received data to stdout.
func (u *UDPServer) Serve() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: u.port, Zone: ""})
	if err != nil {
		log.Fatal(err)
	}
	u.conn = conn
	log.Printf("server listening on: 0.0.0.0:%d\n", u.port)

	buf := make([]byte, bufSize)
	go func() {
		for {
			n, _, _ := u.conn.ReadFromUDP(buf)
			log.Printf("received: %s\n", string(buf[0:n]))

			select {
			case <-u.done:
				log.Println("shutting down...")
				return
			default:
			}
		}
	}()
}

// Close closes the connection.
func (u *UDPServer) Close() {
	u.done <- true
	u.conn.Close()
}
