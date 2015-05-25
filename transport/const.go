package transport

type Service uint16

const (
	ClientService Service = 1
	RaftService           = 2
)

type Opcode uint16

const (
	StateMachineQuery Opcode = iota
	StateMachineCommand
	VerifyRecipient
	GetConfiguration
	SetConfiguration
	GetServerStats
	GetServerInfo
)

// ResponseStatus represents the status of a response from
// the server.
type ResponseStatus uint8

const (
	StatusOkay ResponseStatus = iota
	StatusServiceSpecificError
	StatusInvalidVersion
	StatusInvalidService
	StatusInvalidRequest
)

// RaftPort is the IANA-assigned port for Raft RPC listeners.
const RaftPort = 5254

const (
	PingMessageID    uint64 = 0xffffffffffffffff
	VersionMessageID        = 0xfffffffffffffffe
)
