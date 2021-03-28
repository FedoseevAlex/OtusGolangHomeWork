package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"sync"
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

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
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

	t.Run("eof receive check", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient(l.Addr().String(), time.Second, ioutil.NopCloser(in), out)
		require.NoError(t, client.Connect())
		defer func() { require.NoError(t, client.Close()) }()

		go func() {
			defer wg.Done()

			err = client.Receive()
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()
			conn, err := l.Accept()
			require.NoError(t, err)
			defer func() { require.NoError(t, conn.Close()) }()
		}()

		wg.Wait()
	})

	t.Run("send without connection", func(t *testing.T) {
		tc := NewTelnetClient("localhost:4242", time.Second, nil, nil)
		require.ErrorIs(t, tc.Send(), ErrNoConnection)
	})

	t.Run("receive without connection", func(t *testing.T) {
		tc := NewTelnetClient("localhost:4242", time.Second, nil, nil)
		require.ErrorIs(t, tc.Receive(), ErrNoConnection)
	})

	t.Run("close without connection", func(t *testing.T) {
		tc := NewTelnetClient("localhost:4242", time.Second, nil, nil)
		require.ErrorIs(t, tc.Close(), ErrNoConnection)
	})

	t.Run("timeout error", func(t *testing.T) {
		tc := NewTelnetClient("ya.ru:4242", time.Microsecond, nil, nil)

		var netErr net.Error

		rawErr := tc.Connect()
		ok := errors.As(rawErr, &netErr)

		require.True(t, ok, "error doesn't implement net.Error interface")
		require.True(t, netErr.Timeout(), "error is timeout")
	})

	t.Run("connection refused", func(t *testing.T) {
		tc := NewTelnetClient("localhost:9090", time.Second, nil, nil)
		err := tc.Connect()
		require.Contains(t, err.Error(), "connection refused")
	})
}
