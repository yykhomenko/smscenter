package server

import (
	"bufio"
	"bytes"
	"io"
	"net"

	"github.com/fiorix/go-smpp/smpp/pdu"
)

type conn struct {
	rwc net.Conn
	r   *bufio.Reader
	w   *bufio.Writer
}

func newConn(c net.Conn) *conn {
	return &conn{
		rwc: c,
		r:   bufio.NewReader(c),
		w:   bufio.NewWriter(c),
	}
}

// RemoteAddr implements the Conn interface.
func (c *conn) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}

// Read reads PDU off the wire.
func (c *conn) Read() (pdu.Body, error) {
	return pdu.Decode(c.r)
}

// Write implements the Conn interface.
func (c *conn) Write(p pdu.Body) error {
	var b bytes.Buffer
	err := p.SerializeTo(&b)
	if err != nil {
		return err
	}
	_, err = io.Copy(c.w, &b)
	if err != nil {
		return err
	}
	return c.w.Flush()
}

// Close implements the Conn interface.
func (c *conn) Close() error {
	return c.rwc.Close()
}
