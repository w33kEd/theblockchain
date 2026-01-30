package network

import (
	"crypto"
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
	BlockTime time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	blockTime time.Duration
	memPool *TxPool
	isValidator bool
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		blockTime: opts.BlockTime,
		memPool: NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator{
				s.CreateNewBlock()
			}
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) CreateNewBlock() error {
	fmt.Println("creating new block")
	return nil
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// handle messages
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
