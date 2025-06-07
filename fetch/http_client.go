package fetch

import (
	"context"
	"net"
	"net/http"
	"time"
)

type timeoutConn struct {
	net.Conn
	idle time.Duration
}

func (c *timeoutConn) Read(b []byte) (int, error) {
	// bump the read deadline before each Read
	_ = c.Conn.SetReadDeadline(time.Now().Add(c.idle))
	return c.Conn.Read(b)
}

func (c *timeoutConn) Write(b []byte) (int, error) {
	// also enforce write deadlines
	_ = c.Conn.SetWriteDeadline(time.Now().Add(c.idle))
	return c.Conn.Write(b)
}

func newIdleTimeoutClient(idleTimeout time.Duration) *http.Client {
	// clone the default Transport (to inherit all defaults)
	base := http.DefaultTransport.(*http.Transport).Clone()

	// wrap DialContext to install our timeoutConn
	base.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		d := &net.Dialer{
			Timeout:   30 * time.Second, // connect timeout (optional)
			KeepAlive: 30 * time.Second, // keep-alive period
		}
		rawConn, err := d.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		return &timeoutConn{Conn: rawConn, idle: idleTimeout}, nil
	}

	return &http.Client{Transport: base}
}
