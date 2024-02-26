package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("test timeout", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		timeout, err := time.ParseDuration("1ns")
		require.NoError(t, err)

		client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(nil), nil)
		err = client.Connect()
		require.NotNil(t, err)

		var tg net.Error
		require.True(t, errors.As(err, &tg))
		require.True(t, tg.Timeout())
	})

	t.Run("connecting to a non-existent server", func(t *testing.T) {
		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("127.0.0.1:", timeout, io.NopCloser(nil), nil)
		err = client.Connect()
		require.NotNil(t, err)

		require.True(t, errors.Is(err, syscall.ECONNREFUSED))
	})
}
