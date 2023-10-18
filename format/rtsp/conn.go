package rtsp

import (
	"net"
	"time"
)

type ConnWithTimeout struct {
	Timeout time.Duration
	net.Conn
}

func (c *ConnWithTimeout) Read(p []byte) (n int, err error) {
	if c.Timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(c.Timeout))
	}
	return c.Conn.Read(p)
}

func (c *ConnWithTimeout) Write(p []byte) (n int, err error) {
	if c.Timeout > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	}
	return c.Conn.Write(p)
}
