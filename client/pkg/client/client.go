/*
Package client provides a client for the pawnshop server.
It is only used for acceptance tests and therefore does not have any rigorous testing or much supporting documentation.
*/
package client

import (
	"encoding/json"
	"fmt"
	"net"
	"pawnshop/server/pkg/messages"
)

const (
	addr       = "127.0.0.1:8080"
	bufferSize = 128
)

/*
Client is a lightweight client for the pawn shop server.
*/
type Client struct {
}

/*
Runs the client with the given offer.
It only supports sending a single offer to the server, and does not care about the response.
A real client would want to be able to send multiple offers to the server as well as
handle any response from the server appropriately.
*/
func (c *Client) Run(offer messages.Offer) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer func() {
		conn.Close()
		fmt.Println("Client: Closed connection to server")
	}()

	b, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("failed to marshal offer: %w", err)
	}

	if _, err = conn.Write(b); err != nil {
		return fmt.Errorf("failed to write offer: %w", err)
	}

	buf := make([]byte, bufferSize)

	if _, err = conn.Read(buf); err != nil {
		return fmt.Errorf("failed to read answer: %w", err)
	}

	return nil
}
