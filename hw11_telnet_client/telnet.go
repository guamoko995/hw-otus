package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	tc.conn = conn
	return err
}

func (tc *telnetClient) Close() error {
	return tc.conn.Close()
}

func (tc *telnetClient) Send() error {
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *telnetClient) Receive() error {
	_, err := io.Copy(tc.out, tc.conn)
	return err
}
