package network

import (
	"bytes"
	"fmt"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	Outgoing bool
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.acceptLoop()

	fmt.Println("TCP Transport Listening To Port: ", t.listenAddr)
	return nil

}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *TCPPeer) readLoop(rpcCh chan RPC) {
	buf := make([]byte, 2048)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			fmt.Printf("read error: %s", err)
			continue
		}

		msg := buf[:n]

		rpcCh <- RPC{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}
	}
}

type TCPTransport struct {
	peerCh     chan *TCPPeer
	listenAddr string
	listener   net.Listener // listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		peerCh:     peerCh,
		listenAddr: addr,
	}
}

func (t *TCPTransport) readLoop(peer *TCPPeer) {
	buf := make([]byte, 2048)
	for {
		n, err := peer.conn.Read(buf)
		if err != nil {
			fmt.Printf("read error: %s", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))

		//handle msg => by server
	}
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Accept error from %+v\n", conn)
			continue
		}

		peer := &TCPPeer{
			conn: conn,
		}

		t.peerCh <- peer

		fmt.Printf("new incoming TCP connection => %+v\n", conn)

		// go t.readLoop(peer)

	}
}

