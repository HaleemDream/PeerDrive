package network

// HeaderLength size
const HeaderLength = 9

// ChunkSize size (bytes)
const ChunkSize = 1024

// MaxPayloadSizeMB limit payload size request
const MaxPayloadSizeMB = 1

// MaxPieceRequest limits number of pieces peers to request at a time
const MaxPieceRequest = MaxPayloadSizeMB * 1024 * 1024 / ChunkSize

// MessageType represents first significant byte
type MessageType uint8

// MessageType constants
const (
	RequestPieceInformation MessageType = 0
	SendingPieceInformation MessageType = 1
	RequestPieces           MessageType = 2
	SendingPieces           MessageType = 3
	TerminateConnection     MessageType = 255
)

// Header represents header data in buffer
type Header struct {
	MessageType MessageType
	PieceCount  uint32
	FileIndex   uint32
}
