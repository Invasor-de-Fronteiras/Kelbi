package mhfpacket

import (
 "errors"

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgSysReserve1A7 represents the MSG_SYS_reserve1A7
type MsgSysReserve1A7 struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgSysReserve1A7) Opcode() network.PacketID {
	return network.MSG_SYS_reserve1A7
}

// Parse parses the packet from binary
func (m *MsgSysReserve1A7) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgSysReserve1A7) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
