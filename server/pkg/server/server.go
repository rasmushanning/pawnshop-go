// Package server implements a TCP server that handles offers from clients.
// The server uses the pawnshop package to handle the offers, and the inventory package to store the inventory.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"pawnshop/server/pkg/inventory"
	"pawnshop/server/pkg/messages"
	"pawnshop/server/pkg/pawnshop"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	addr    = "127.0.0.1:8080"
	bufSize = 128
)

// OfferHandler is an interface that handles offers.
type OfferHandler interface {
	HandleOffer(offer messages.Offer) messages.Answer
}

// PawnShopServer is a TCP server that handles offers from clients and responds to them.
type PawnShopServer struct {
	addr         string
	isRunning    bool
	offerHandler OfferHandler
	listener     net.Listener
	connections  chan net.Conn
	shutdownCtx  context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

/*
Creates a new PawnShopServer with the given inventory size.
If size is less than 1, an error is returned.
*/
func NewPawnShopServer(sz int) (*PawnShopServer, error) {
	if sz < 1 {
		return nil, errors.New("inventory size must be at least 1")
	}

	inv := inventory.NewInventory(sz)

	pawnshop, err := pawnshop.NewPawnShop(inv)
	if err != nil {
		return nil, fmt.Errorf("failed to create new PawnShop, %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	log.Debugf("Created new pawn shop with an inventory of size %d: %s", sz, inv)
	return &PawnShopServer{
		addr:         addr,
		isRunning:    false,
		offerHandler: pawnshop,
		connections:  make(chan net.Conn),
		shutdownCtx:  ctx,
		cancel:       cancel,
		wg:           sync.WaitGroup{},
	}, nil
}

/*
Starts the server and listens for connections.
*/
func (p *PawnShopServer) Start() error {
	l, err := net.Listen("tcp", p.addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	p.listener = l

	// Use a waitgroup to enable graceful shutdown using server.Stop()
	p.wg.Add(2)
	go p.handleConnections()
	go p.acceptConnections()
	log.Infof("Started server, listening at %s", addr)

	// In case of a graceful shutdown, wait for both the
	// acceptConnections and handleConnections goroutines to exit
	p.wg.Wait()

	log.Info("Server has stopped")
	return nil
}

/*
Stops the server and closes the listener.
Returns an error if the listener could not be closed.
*/
func (p *PawnShopServer) Stop() error {
	p.isRunning = false
	p.cancel()

	if p.listener != nil {
		if err := p.listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
	}
	return nil
}

/*
Returns true if the server is running and able to accept new connections, false otherwise.
*/
func (p *PawnShopServer) IsRunning() bool {
	return p.isRunning
}

/*
Accepts new TCP connections and sends any new connections to the connections channel,
which will be handled by the handleConnection function. Supports graceful shutdown.
*/
func (p *PawnShopServer) acceptConnections() {
	defer p.wg.Done()

	p.isRunning = true
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			select {
			case <-p.shutdownCtx.Done():
				log.Debug("acceptConnections received shutdown signal, shutting down...")
				return
			default:
				log.Errorf("Failed to accept connection: %s", err)
			}
		}
		p.connections <- conn
	}
}

/*
Handles TCP connections received from the connections channel.
Supports graceful shutdown.
*/
func (p *PawnShopServer) handleConnections() {
	defer p.wg.Done()

	connsWG := sync.WaitGroup{}

	for {
		select {
		case <-p.shutdownCtx.Done():
			log.Debug("handleConnections received shutdown signal, shutting down...")
			log.Debug("Waiting for all currently handled offers to finish...")
			connsWG.Wait()
			log.Debug("All current offers have finished, shutting down...")
			return
		case conn := <-p.connections:
			connsWG.Add(1)

			go func() {
				p.handleConnection(conn)
				connsWG.Done()
			}()
		}
	}
}

/*
Reads an offer from a connection, handles it and writes the answer back on the connection.
*/
func (p *PawnShopServer) handleConnection(conn net.Conn) {
	buf := make([]byte, bufSize)

	var n int
	n, err := conn.Read(buf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			log.Errorf("Failed to read from connection: %s", err)
		}
		return
	}

	offB := buf[:n]

	var off messages.Offer
	if err = json.Unmarshal(offB, &off); err != nil {
		rejectOffer(conn.Write)
		log.Errorf("Failed to unmarshal offer: %s", err)
		return
	}

	log.Infof("Received offer from client: %s", string(offB))
	ans := p.handleOffer(off)
	ansB, err := json.Marshal(ans)
	if err != nil {
		rejectOffer(conn.Write)
		log.Errorf("Failed to marshal answer: %s", err)
		return
	}

	log.Infof("Sending answer to client: %s", string(ansB))

	if _, err = conn.Write(ansB); err != nil {
		log.Errorf("Failed to write answer: %s", err)
		return
	}
}

/*
Handles an offer and takes appropriate action depending on the Code.
*/
func (p *PawnShopServer) handleOffer(offer messages.Offer) messages.Answer {
	switch offer.Code {
	case messages.PawnCode:
		return p.offerHandler.HandleOffer(offer)
	default:
		return messages.CreateRejectAnswer()
	}
}

/*
Rejects an offer by writing a reject answer on the connection.
*/
func rejectOffer(writeConn func([]byte) (n int, err error)) {
	rejAns := messages.CreateRejectAnswer()
	rejAnsB, err := json.Marshal(rejAns)
	if err != nil {
		log.Errorf("Failed to marshal reject answer: %s", err)
		return
	}

	writeConn(rejAnsB)
}
