package transport

/* The serialized format:
 *    version uint8
 *    service:ubint16
 *    service_specific_error_value:uint8
 *    opcode:ubint16
 *    payload:[]uint8
 */

const RequestVersion byte = 1

var errorVersions = map[Service]byte{
	ClientService: 1, // TODO(dlg): 0 or 1?
	RaftService:   1,
}

type Request struct {
	ServiceId Service
	OpcodeId  Opcode
	Payload   []byte
}

func (r *Request) Marshal() []byte {
	header := make([]byte, 6)
	header[0] = RequestVersion
	PutUbint16(header[1:2], uint16(r.ServiceId))
	header[2] = errorVersions[r.ServiceId]
	PutUbint16(header[4:6], uint16(r.OpcodeId))
	return append(header, r.Payload...)
}
