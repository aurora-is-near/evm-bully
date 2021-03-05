package nearapi

import (
	"errors"
	"fmt"

	"github.com/ybbus/jsonrpc/v2"
)

type Connection struct {
	c jsonrpc.RPCClient
}

func NewConnection(nodeURL string) *Connection {
	var c Connection
	c.c = jsonrpc.NewClient(nodeURL)
	return &c
}

func (c *Connection) call(method string, params ...interface{}) (interface{}, error) {
	res, err := c.c.Call(method, params...)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		if res.Error.Data != nil {
			return nil, fmt.Errorf("nearapi: jsonrpc: %d: %s: %v",
				res.Error.Code, res.Error.Message, res.Error.Data)
		}
		return nil, fmt.Errorf("nearapi: jsonrpc: %d: %s",
			res.Error.Code, res.Error.Message)
	}
	if res.Result == nil {
		return nil, errors.New("nearapi: JSON-RPC result is nil")
	}
	return res.Result, nil
}

func (c *Connection) State(accountID string) (interface{}, error) {
	res, err := c.call("query", map[string]string{
		"request_type": "view_account",
		"finality":     "final",
		"account_id":   accountID,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
