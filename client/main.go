package client

import (
	"bytes"
	"context"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	address string
}

func New(address string) *Client {
	return &Client{
		address: address,
	}
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	conn, err := net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	wr := resp.NewWriter(buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})

	_, err = io.Copy(conn, buf)
	return err
}
