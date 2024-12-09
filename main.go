package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

const defaultListenAddr = ":5001"

// Config holds server configuration
type Config struct {
	ListenAddr string
}

// Message represents a command sent from a peer to the server
type Message struct {
	cmd  Command
	peer *Peer
}

// Server represents the Redis-like server
type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	delPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message

	kv *KV
}

// NewServer initializes a new server with the given configuration
func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		delPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

// Start begins the server's operation, listening for connections and handling messages
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	s.ln = ln
	defer s.ln.Close()

	go s.loop() // Goroutine for main server loop

	slog.Info("Redis-like server running", "listenAddr", s.ListenAddr)
	return s.acceptLoop()
}

// handleMessage processes commands received from peers
func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case ClientCommand:
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil { // Use WriteString to send "OK"
			return err
		}
	case SetCommand:
		if err := s.kv.Set(v.key, v.val); err != nil {
			return err
		}
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString("OK"); err != nil { // Use WriteString to send "OK"
			return err
		}
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		if err := resp.
			NewWriter(msg.peer.conn).
			WriteString(string(val)); err != nil { // Use WriteString to send value
			return err
		}
	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
		}
		_, err := msg.peer.Send(respWriteMap(spec))
		if err != nil {
			return fmt.Errorf("peer send error: %s", err)
		}
	default:
		return fmt.Errorf("unsupported command: %v", msg.cmd)
	}
	return nil
}

// loop handles server events: peer connections, disconnections, and messages
func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("failed to handle message", "error", err)
			}
		case peer := <-s.addPeerCh:
			slog.Info("peer connected", "remoteAddr", peer.conn.RemoteAddr())
			s.peers[peer] = true
		case peer := <-s.delPeerCh:
			slog.Info("peer disconnected", "remoteAddr", peer.conn.RemoteAddr())
			delete(s.peers, peer)
		case <-s.quitCh:
			slog.Info("server shutting down")
			return
		}
	}
}

// acceptLoop handles incoming connections
func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err)
			continue
		}
		go s.handleConn(conn) // Goroutine to handle the new connection
	}
}

// handleConn processes the connection from a peer
func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh, s.delPeerCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "error", err, "remoteAddr", conn.RemoteAddr())
	}
}

// wrapError adds context to an error
func wrapError(context string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", context, err)
	}
	return nil
}

func main() {
	listenAddr := flag.String("listenAddr", defaultListenAddr, "listen address of the Redis-like server")
	flag.Parse()

	server := NewServer(Config{
		ListenAddr: *listenAddr,
	})

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
