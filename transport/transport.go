// The transport package implements the version 1 RPC-over-TCP
// transport used by LogCabin.
package transport

import (
	"fmt"
	"net"
	"sync/atomic"
)

const (
	magicV1       = uint16(0xdaf4)
	versionV1     = uint16(0x1)
	maxPayloadLen = 16384
)

type transport struct {
	raddr     *net.TCPAddr
	conn      *net.TCPConn
	currentID uint64
}

func New(raddr *net.TCPAddr) (tr *transport, err error) {
	if raddr == nil {
		err = fmt.Errorf("Need a remote addresss to connect to.")
		return
	}
	conn, err := net.DialTCP(raddr.Network(), nil, raddr)
	if err != nil {
		return
	}
	// The correctness of the code depends on no read or write
	// timeouts being set.
	/*
		conn.SetReadBuffer(0)
		conn.SetNoDelay(true)
		conn.SetWriteBuffer(0)
		conn.SetKeepAlive(true)
	*/
	tr = &transport{
		raddr: raddr,
		conn:  conn,
	}
	return
}

// getVersion is NOT thread-safe: it must only be called by New
func (tr *transport) getVersion() (version int, err error) {
	tr.send(VersionMessageID, []byte{})
	messageID, m, err := tr.Receive()
	if err != nil {
		return
	}
	if messageID != VersionMessageID {
		err = fmt.Errorf("getVersion RPC response used the wrong messageID: got %x", messageID)
		return
	}
	if len(m) < 2 {
		err = fmt.Errorf("getVersion RPC response was too short")
		return
	}
	version = int(Ubint16(m))
	return
}

// Close closes the underlying transport.
func (tr *transport) Close() {
	if tr.conn != nil {
		tr.conn.Close()
	}
}

func (tr *transport) Send(m []byte) (messageID uint64, err error) {
	if tr.conn == nil {
		err = fmt.Errorf("connection never opened")
		return
	}

	// Atomically increment currentID
	messageID = atomic.AddUint64(&tr.currentID, 1)
	err = tr.send(messageID, m) // send the message
	return
}

func (tr *transport) send(messageID uint64, m []byte) (err error) {
	if len(m) > 0xffffffff {
		err = fmt.Errorf("message too long")
		return
	}
	payloadLen := uint32(len(m))

	// Serialize the header.
	header := make([]byte, 16)
	PutUbint16(header[0:2], magicV1)
	PutUbint16(header[2:4], versionV1)
	PutUbint32(header[4:8], payloadLen)
	PutUbint64(header[8:16], messageID)

	n, err := tr.conn.Write(append(header, m...))

	if n != len(m) {
		// We've written a partial payload, and this connection's state
		// can't be recovered.
		tr.conn.Close()
		if err == nil {
			// This should never happen, but...
			err = fmt.Errorf("wrote a partial payload (%u out of %u bytes), but without an error", n, len(m))
		}
	}
	return
}

func (tr *transport) ping() error {
	return tr.send(PingMessageID, []byte{})
}

func (tr *transport) Receive() (messageID uint64, m []byte, err error) {
	header := make([]byte, 16)
	n, err := tr.conn.Read(header)
	switch {
	// TODO(dlg): This is stupid.
	case err != nil:
		return
	case n != 16:
		err = fmt.Errorf("couldn't read complete header")
		return
	}

	// Deserialize the header.
	magic := Ubint16(header[0:2])
	version := Ubint16(header[2:4])
	payloadLen := Ubint32(header[4:8])
	messageID = Ubint64(header[8:16])

	// Check the magic and the version.
	if magic != magicV1 || version != versionV1 {
		err = fmt.Errorf("unknown protocol magic or version")
	}
	// Sanity check on payload length.
	if payloadLen > maxPayloadLen {
		err = fmt.Errorf("payload length %d greater than maximum payload length of %d", payloadLen, maxPayloadLen)
		return
	}

	m = make([]byte, payloadLen)
	_, err = tr.conn.Read(m)

	return
}
