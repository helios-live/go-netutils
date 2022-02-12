package netutils

import (
	"log"
	"net"

	"github.com/davecgh/go-spew/spew"
)

// PrinterConn prints all bytes that go through it
type PrinterConn struct {
	net.Conn
	Prefix string
}

func (pc PrinterConn) Read(b []byte) (int, error) {
	n, err := pc.Conn.Read(b)
	log.Println("Read:"+pc.Prefix, "\n============================================================================\n", spew.Sdump(err, pc.RemoteAddr().String(), b[0:n]), "\n============================================================================")
	return n, err
}

func (pc PrinterConn) Write(b []byte) (int, error) {
	n, err := pc.Conn.Write(b)
	log.Println("Write:"+pc.Prefix, "\n============================================================================\n", spew.Sdump(err, pc.RemoteAddr().String(), b[0:]), "\n============================================================================")
	return n, err
}
