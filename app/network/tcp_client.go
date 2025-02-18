package network

import (
	"fmt"
	"io"
	"net"
)

type TcpClient struct {
	conn net.Conn
}

func NewTcpClient(address string) (*TcpClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("tcp client connect error: %v", err)
	}
	return &TcpClient{conn: conn}, nil
}

func (c *TcpClient) Send(request []byte) ([]byte, error) {
	_, err := c.conn.Write(request)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	count, err := c.conn.Read(response)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return response[:count], nil
}

func (c *TcpClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
