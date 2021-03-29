package main

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

var ErrNoConnection = errors.New("no active connection")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &SimpleTelnetClient{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

type SimpleTelnetClient struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	conn    net.Conn
}

func (c *SimpleTelnetClient) Connect() (err error) {
	c.conn, err = net.DialTimeout("tcp", c.Address, c.Timeout)
	return
}

func (c *SimpleTelnetClient) Send() (err error) {
	if c.conn == nil {
		return errors.Wrap(ErrNoConnection, "send failed")
	}

	_, err = io.Copy(c.conn, c.In)
	return
}

func (c *SimpleTelnetClient) Receive() (err error) {
	if c.conn == nil {
		return errors.Wrap(ErrNoConnection, "receive failed")
	}

	_, err = io.Copy(c.Out, c.conn)
	return
}

func (c *SimpleTelnetClient) Close() (err error) {
	if c.conn == nil {
		return errors.Wrap(ErrNoConnection, "connection close failed")
	}

	err = c.conn.Close()
	if err != nil {
		return
	}

	err = c.In.Close()
	return
}
