package mhfpacket

import ( 
 "errors" 

 	"erupe-ce/network/clientctx"
	"erupe-ce/network"
	"erupe-ce/common/byteframe"
)

// MsgMhfInfoJoint represents the MSG_MHF_INFO_JOINT
type MsgMhfInfoJoint struct{}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfInfoJoint) Opcode() network.PacketID {
	return network.MSG_MHF_INFO_JOINT
}

// Parse parses the packet from binary
func (m *MsgMhfInfoJoint) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}

// Build builds a binary packet from the current data.
func (m *MsgMhfInfoJoint) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
