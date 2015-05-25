package transport

import "encoding/binary"

var (
	PutUbint64 = binary.BigEndian.PutUint64
	PutUbint32 = binary.BigEndian.PutUint32
	PutUbint16 = binary.BigEndian.PutUint16
	Ubint64    = binary.BigEndian.Uint64
	Ubint32    = binary.BigEndian.Uint32
	Ubint16    = binary.BigEndian.Uint16
)
