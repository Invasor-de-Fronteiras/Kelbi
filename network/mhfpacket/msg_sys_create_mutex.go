package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysCreateMutex represents the MSG_SYS_CREATE_MUTEX
type MsgSysCreateMutex struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysCreateMutex) Opcode() network.PacketID {
	return network.MSG_SYS_CREATE_MUTEX
}

// Parse parses the packet from binary
func (m *MsgSysCreateMutex) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysCreateMutex) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
