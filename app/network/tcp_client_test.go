package network

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTcpClient(t *testing.T) {
	const address = "localhost:8080"
	response := []byte("ok")
	request := make([]byte, 1024)
	listener, err := net.Listen("tcp", address)
	require.NoError(t, err)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		_, err = conn.Read(request)
		require.NoError(t, err)

		_, err = conn.Write([]byte(response))
		require.NoError(t, err)

	}()

	testData := map[string]struct {
		request []byte
		address string
		client  func(string) (*TcpClient, error)

		expected []byte
		err      error
	}{
		"invalid address": {
			request: request,
			address: "localhost:8081",

			client: func(address string) (*TcpClient, error) {
				return nil, fmt.Errorf("tcp client connect error")
			},

			expected: response,
			err:      errors.New("tcp client connect error"),
		},
		"send ok": {
			request: request,
			address: address,

			client: func(address string) (*TcpClient, error) {
				conn, err := net.Dial("tcp", address)
				if err != nil {
					return nil, err
				}
				return &TcpClient{conn: conn}, nil
			},
			expected: response,
			err:      nil,
		},
	}

	for name, test := range testData {
		t.Run(name, func(t *testing.T) {
			client, err := test.client(test.address)
			require.Equal(t, test.err, err)
			if client == nil {
				return
			}

			response, err := client.Send(test.request)
			require.Equal(t, test.expected, response)
			if test.err != nil {
				assert.Error(t, err, test.err)
				assert.EqualError(t, err, test.err.Error())
			}

		})
	}
}
