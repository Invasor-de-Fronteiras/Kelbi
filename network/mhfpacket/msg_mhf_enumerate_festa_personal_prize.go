package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfEnumerateFestaPersonalPrize represents the MSG_MHF_ENUMERATE_FESTA_PERSONAL_PRIZE
type MsgMhfEnumerateFestaPersonalPrize struct {
	AckHandle uint32
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfEnumerateFestaPersonalPrize) Opcode() network.PacketID {
	return network.MSG_MHF_ENUMERATE_FESTA_PERSONAL_PRIZE
}

// Parse parses the packet from binary
func (m *MsgMhfEnumerateFestaPersonalPrize) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfEnumerateFestaPersonalPrize) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
