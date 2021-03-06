// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package rpc implements the Ethereum JSON-RPC API.
package rpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/pnomarev/go-ethereum/rpc/comms"
	"github.com/pnomarev/go-ethereum/rpc/shared"
)

// Xeth is a native API interface to a remote node.
type Xeth struct {
	client comms.EthereumClient
	reqId  uint32
}

// NewXeth constructs a new native API interface to a remote node.
func NewXeth(client comms.EthereumClient) *Xeth {
	return &Xeth{
		client: client,
	}
}

// Call invokes a method with the given parameters are the remote node.
func (self *Xeth) Call(method string, params []interface{}) (map[string]interface{}, error) {
	// Assemble the json RPC request
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req := &shared.Request{
		Id:      atomic.AddUint32(&self.reqId, 1),
		Jsonrpc: "2.0",
		Method:  method,
		Params:  data,
	}
	// Send the request over and retrieve the response
	if err := self.client.Send(req); err != nil {
		return nil, err
	}
	res, err := self.client.Recv()
	if err != nil {
		return nil, err
	}
	// Ensure the response is valid, and extract the results
	success, isSuccessResponse := res.(*shared.SuccessResponse)
	failure, isFailureResponse := res.(*shared.ErrorResponse)
	switch {
	case isFailureResponse:
		return nil, fmt.Errorf("Method invocation failed: %v", failure.Error)

	case isSuccessResponse:
		return success.Result.(map[string]interface{}), nil

	default:
		return nil, fmt.Errorf("Invalid response type: %v", reflect.TypeOf(res))
	}
}
