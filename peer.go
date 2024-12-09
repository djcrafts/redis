package main

import (
	"io"
	"log/slog"
	"net"

	"github.com/tidwall/resp"
)

// Peer represents a connected client
type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

// Send sends a message to the peer
func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

// NewPeer initializes a new peer
func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

// readLoop continuously reads and processes commands from the peer
func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	defer p.conn.Close() // Ensure connection is closed on exit

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}
		if err != nil {
			slog.Error("error reading from peer", "error", err, "remoteAddr", p.conn.RemoteAddr())
			p.delCh <- p
			break
		}

		var cmd Command
		if v.Type() == resp.Array {
			rawCMD := v.Array()[0].String()
			switch rawCMD {
			case CommandClient:
				cmd = ClientCommand{value: v.Array()[1].String()}
			case CommandGET:
				cmd = GetCommand{key: v.Array()[1].Bytes()}
			case CommandSET:
				cmd = SetCommand{key: v.Array()[1].Bytes(), val: v.Array()[2].Bytes()}
			case CommandHELLO:
				cmd = HelloCommand{value: v.Array()[1].String()}
			default:
				slog.Warn("unhandled command received", "command", rawCMD)
			}

			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}
		}
	}
	return nil
}
