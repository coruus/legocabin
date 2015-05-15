// The transport package implements the version 1 RPC-over-TCP
// transport used by LogCabin.
package transport

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	magicV1       = uint16(0xdaf4)
	versionV1     = uint16(0x1)
	maxPayloadLen = 16384
)

var (
	PutUbint64 = binary.BigEndian.PutUint64
	PutUbint32 = binary.BigEndian.PutUint32
	PutUbint16 = binary.BigEndian.PutUint16
	Ubint64    = binary.BigEndian.Uint64
	Ubint32    = binary.BigEndian.Uint32
	Ubint16    = binary.BigEndian.Uint16
)

type transport struct {
	raddr *net.TCPAddr
	conn  *net.TCPConn
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

func (tr *transport) Close() {
	if tr.conn != nil {
		tr.conn.Close()
		tr.conn = nil
	}
}

func (tr *transport) Send(m []byte) (messageID uint64, err error) {
	if tr.conn == nil {
		err = fmt.Errorf("connection closed")
		return
	}
	if len(m) > 0xffffffff {
		err = fmt.Errorf("message too long")
		return
	}

	payloadLen := uint32(len(m))

	// Generate a random message ID.
	rawid := make([]byte, 4)
	_, err = rand.Read(rawid)
	messageID = Ubint64(rawid)

	// Serialize the header.
	header := make([]byte, 16)
	PutUbint16(header[0:2], magicV1)
	PutUbint16(header[2:4], versionV1)
	PutUbint32(header[4:8], payloadLen)
	copy(header[8:16], rawid)

	tr.conn.Write(append(header, m...))

	return
}

func (tr *transport) Recv() (messageID uint64, m []byte, err error) {
	header := make([]byte, 16)
	n, err := tr.conn.Read(header)
	switch {
	// TODO(dlg): This is stupid.
	case err != nil:
		return
	case n != 16:
		err = fmt.Errorf("couldn't read complete header")
		// TODO(dlg): Is this right? Need to check what syscalls are used.
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
