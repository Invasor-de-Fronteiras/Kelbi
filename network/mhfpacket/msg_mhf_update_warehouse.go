package mhfpacket

import (
	"errors"
	"erupe-ce/common/byteframe"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
)

type WarehouseStack struct {
	ID        uint32
	Index     uint16
	EquipType uint16
	ItemID    uint16
	Quantity  uint16
	Data      []byte
}

// MsgMhfUpdateWarehouse represents the MSG_MHF_UPDATE_WAREHOUSE
type MsgMhfUpdateWarehouse struct {
	AckHandle uint32
	BoxType   string
	BoxIndex  uint8
	Updates   []WarehouseStack
}

// Opcode returns the ID associated with this packet type.
func (m *MsgMhfUpdateWarehouse) Opcode() network.PacketID {
	return network.MSG_MHF_UPDATE_WAREHOUSE
}

// Parse parses the packet from binary
func (m *MsgMhfUpdateWarehouse) Parse(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	m.AckHandle = bf.ReadUint32()
	boxType := bf.ReadUint8()
	switch boxType {
	case 0:
		m.BoxType = "item"
	case 1:
		m.BoxType = "equip"
	}
	m.BoxIndex = bf.ReadUint8()
	changes := int(bf.ReadUint16())
	var stackUpdate WarehouseStack
	for i := 0; i < changes; i++ {
		switch boxType {
		case 0:
			stackUpdate.ID = bf.ReadUint32()
			stackUpdate.Index = bf.ReadUint16()
			stackUpdate.ItemID = bf.ReadUint16()
			stackUpdate.Quantity = bf.ReadUint16()
			_ = bf.ReadUint16() // Unk
			m.Updates = append(m.Updates, stackUpdate)
		case 1:
			stackUpdate.ID = bf.ReadUint32()
			stackUpdate.Index = bf.ReadUint16()
			stackUpdate.EquipType = bf.ReadUint16()
			stackUpdate.ItemID = bf.ReadUint16()
			stackUpdate.Data = bf.ReadBytes(56)
			m.Updates = append(m.Updates, stackUpdate)
		}
	}
	_ = bf.ReadUint16()
	return nil
}

// Build builds a binary packet from the current data.
func (m *MsgMhfUpdateWarehouse) Build(bf *byteframe.ByteFrame, ctx *clientctx.ClientContext) error {
	return errors.New("NOT IMPLEMENTED")
}
