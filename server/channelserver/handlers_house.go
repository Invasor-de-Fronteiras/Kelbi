package channelserver

import (
	"erupe-ce/common/byteframe"
	ps "erupe-ce/common/pascalstring"
	"erupe-ce/common/stringsupport"
	"erupe-ce/network/mhfpacket"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
)

const warehouseNamesQuery = `
SELECT
COALESCE(item0name, ''),
COALESCE(item1name, ''),
COALESCE(item2name, ''),
COALESCE(item3name, ''),
COALESCE(item4name, ''),
COALESCE(item5name, ''),
COALESCE(item6name, ''),
COALESCE(item7name, ''),
COALESCE(item8name, ''),
COALESCE(item9name, ''),
COALESCE(equip0name, ''),
COALESCE(equip1name, ''),
COALESCE(equip2name, ''),
COALESCE(equip3name, ''),
COALESCE(equip4name, ''),
COALESCE(equip5name, ''),
COALESCE(equip6name, ''),
COALESCE(equip7name, ''),
COALESCE(equip8name, ''),
COALESCE(equip9name, '')
FROM warehouse
`

func handleMsgMhfUpdateInterior(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateInterior)
	_, err := s.Server.db.Exec("UPDATE characters SET house=$1 WHERE id=$2", pkt.InteriorData, s.CharID)
	if err != nil {
		panic(err)
	}
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

type HouseData struct {
	CharID uint32 `db:"id"`
	HRP    uint16 `db:"hrp"`
	GR     uint16 `db:"gr"`
	Name   string `db:"name"`
}

func handleMsgMhfEnumerateHouse(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateHouse)
	bf := byteframe.NewByteFrame()
	var houses []HouseData
	switch pkt.Method {
	case 1:
		var friendsList string
		// nolint:errcheck // Error return value of `` is not checked
		s.Server.db.QueryRow("SELECT friends FROM characters WHERE id=$1", s.CharID).Scan(&friendsList)
		cids := stringsupport.CSVElems(friendsList)
		for _, cid := range cids {
			house := HouseData{}
			row := s.Server.db.QueryRowx("SELECT id, hrp, gr, name FROM characters WHERE id=$1", cid)
			err := row.StructScan(&house)
			if err != nil {
				panic(err)
			} else {
				houses = append(houses, house)
			}
		}
	case 2:
		guild, err := GetGuildInfoByCharacterId(s, s.CharID)
		if err != nil {
			break
		}
		guildMembers, err := GetGuildMembers(s, guild.ID, false)
		if err != nil {
			break
		}
		for _, member := range guildMembers {
			house := HouseData{}
			row := s.Server.db.QueryRowx("SELECT id, hrp, gr, name FROM characters WHERE id=$1", member.CharID)
			err := row.StructScan(&house)
			if err != nil {
				panic(err)
			} else {
				houses = append(houses, house)
			}
		}
	case 3:
		house := HouseData{}
		row := s.Server.db.QueryRowx("SELECT id, hrp, gr, name FROM characters WHERE name ILIKE $1", fmt.Sprintf(`%%%s%%`, pkt.Name))
		err := row.StructScan(&house)
		if err != nil {
			panic(err)
		} else {
			houses = append(houses, house)
		}
	case 4:
		house := HouseData{}
		row := s.Server.db.QueryRowx("SELECT id, hrp, gr, name FROM characters WHERE id=$1", pkt.CharID)
		err := row.StructScan(&house)
		if err != nil {
			panic(err)
		} else {
			houses = append(houses, house)
		}
	case 5: // Recent visitors
		break
	}
	var exists int
	for _, house := range houses {
		for _, session := range s.Server.Sessions {
			if session.CharID == house.CharID {
				exists++
				bf.WriteUint32(house.CharID)
				bf.WriteUint8(session.MySeries.State)
				if len(session.MySeries.Password) > 0 {
					bf.WriteUint8(3)
				} else {
					bf.WriteUint8(0)
				}
				bf.WriteUint16(house.HRP)
				bf.WriteUint16(house.GR)
				ps.Uint8(bf, house.Name, true)
				break
			}
		}
	}
	resp := byteframe.NewByteFrame()
	resp.WriteUint16(uint16(exists))
	resp.WriteBytes(bf.Data())
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfUpdateHouse(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateHouse)
	// 01 = closed
	// 02 = open anyone
	// 03 = open friends
	// 04 = open guild
	// 05 = open friends+guild
	s.MySeries.State = pkt.State
	s.MySeries.Password = pkt.Password
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfLoadHouse(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadHouse)
	bf := byteframe.NewByteFrame()
	if pkt.Destination != 9 && len(pkt.Password) > 0 && pkt.CheckPass {
		for _, session := range s.Server.Sessions {
			if session.CharID == pkt.CharID && pkt.Password != session.MySeries.Password {
				doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
				return
			}
		}
	}

	var furniture []byte
	err := s.Server.db.QueryRow("SELECT house FROM characters WHERE id=$1", pkt.CharID).Scan(&furniture)
	if err != nil {
		panic(err)
	}
	if furniture == nil {
		furniture = make([]byte, 20)
	}

	switch pkt.Destination {
	case 3: // Others house
		for _, session := range s.Server.Sessions {
			if session.CharID == pkt.CharID {
				bf.WriteBytes(session.MySeries.HouseTier)
				bf.WriteBytes(session.MySeries.HouseData)
				bf.WriteBytes(make([]byte, 19)) // Padding?
				bf.WriteBytes(furniture)
			}
		}
	case 4: // Bookshelf
		for _, session := range s.Server.Sessions {
			if session.CharID == pkt.CharID {
				bf.WriteBytes(session.MySeries.BookshelfData)
			}
		}
	case 5: // Gallery
		for _, session := range s.Server.Sessions {
			if session.CharID == pkt.CharID {
				bf.WriteBytes(session.MySeries.GalleryData)
			}
		}
	case 8: // Tore
		for _, session := range s.Server.Sessions {
			if session.CharID == pkt.CharID {
				bf.WriteBytes(session.MySeries.ToreData)
			}
		}
	case 9: // Own house
		bf.WriteBytes(furniture)
	case 10: // Garden
		for _, session := range s.Server.Sessions {
			if session.CharID == pkt.CharID {
				bf.WriteBytes(session.MySeries.GardenData)
				c, d := getGookData(s, pkt.CharID)
				bf.WriteUint16(c)
				bf.WriteUint16(0)
				bf.WriteBytes(d)
			}
		}
	}
	if len(bf.Data()) == 0 {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
	} else {
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
	}
}

func handleMsgMhfGetMyhouseInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetMyhouseInfo)

	var data []byte
	err := s.Server.db.QueryRow("SELECT trophy FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if err != nil {
		panic(err)
	}
	if len(data) > 0 {
		doAckBufSucceed(s, pkt.AckHandle, data)
	} else {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}
}

func handleMsgMhfUpdateMyhouseInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateMyhouseInfo)

	_, err := s.Server.db.Exec("UPDATE characters SET trophy=$1 WHERE id=$2", pkt.Unk0, s.CharID)
	if err != nil {
		panic(err)
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfLoadDecoMyset(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadDecoMyset)
	var data []byte
	err := s.Server.db.QueryRow("SELECT decomyset FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if err != nil {
		s.logger.Fatal("Failed to get preset decorations savedata from db", zap.Error(err))
	}

	if len(data) > 0 {
		doAckBufSucceed(s, pkt.AckHandle, data)
		//doAckBufSucceed(s, pkt.AckHandle, data)
	} else {
		// set first byte to 1 to avoid pop up every time without save
		body := make([]byte, 0x226)
		body[0] = 1
		doAckBufSucceed(s, pkt.AckHandle, body)
	}
}

func handleMsgMhfSaveDecoMyset(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveDecoMyset)
	// https://gist.github.com/Andoryuuta/9c524da7285e4b5ca7e52e0fc1ca1daf
	var loadData []byte
	bf := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload[1:]) // skip first unk byte
	err := s.Server.db.QueryRow("SELECT decomyset FROM characters WHERE id = $1", s.CharID).Scan(&loadData)
	if err != nil {
		s.logger.Fatal("Failed to get preset decorations savedata from db", zap.Error(err))
	} else {
		numSets := bf.ReadUint8() // sets being written
		// empty save
		if len(loadData) == 0 {
			loadData = []byte{0x01, 0x00}
		}

		savedSets := loadData[1] // existing saved sets
		// no sets, new slice with just first 2 bytes for appends later
		if savedSets == 0 {
			loadData = []byte{0x01, 0x00}
		}
		for i := 0; i < int(numSets); i++ {
			writeSet := bf.ReadUint16()
			dataChunk := bf.ReadBytes(76)
			setBytes := append([]byte{uint8(writeSet >> 8), uint8(writeSet & 0xff)}, dataChunk...)
			for x := 0; true; x++ {
				if x == int(savedSets) {
					// appending set
					if loadData[len(loadData)-1] == 0x10 {
						// sanity check for if there was a messy manual import
						loadData = append(loadData[:len(loadData)-2], setBytes...)
					} else {
						loadData = append(loadData, setBytes...)
					}
					savedSets++
					break
				}
				currentSet := loadData[3+(x*78)]
				if int(currentSet) == int(writeSet) {
					// replacing a set
					loadData = append(loadData[:2+(x*78)], append(setBytes, loadData[2+((x+1)*78):]...)...)
					break
				} else if int(currentSet) > int(writeSet) {
					// inserting before current set
					loadData = append(loadData[:2+((x)*78)], append(setBytes, loadData[2+((x)*78):]...)...)
					savedSets++
					break
				}
			}
			loadData[1] = savedSets // update set count
		}
		_, err := s.Server.db.Exec("UPDATE characters SET decomyset=$1 WHERE id=$2", loadData, s.CharID)
		if err != nil {
			s.logger.Fatal("Failed to update decomyset savedata in db", zap.Error(err))
		}
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

type Title struct {
	ID       uint16    `db:"id"`
	Acquired time.Time `db:"unlocked_at"`
	Updated  time.Time `db:"updated_at"`
}

func handleMsgMhfEnumerateTitle(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateTitle)
	var count uint16
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(0)
	bf.WriteUint16(0) // Unk
	rows, err := s.Server.db.Queryx("SELECT id, unlocked_at, updated_at FROM titles WHERE char_id=$1", s.CharID)
	if err != nil {
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
		return
	}
	for rows.Next() {
		title := &Title{}
		err = rows.StructScan(&title)
		if err != nil {
			continue
		}
		count++
		bf.WriteUint16(title.ID)
		bf.WriteUint16(0) // Unk
		bf.WriteUint32(uint32(title.Acquired.Unix()))
		bf.WriteUint32(uint32(title.Updated.Unix()))
	}
	bf.Seek(0, io.SeekStart)
	bf.WriteUint16(count)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfAcquireTitle(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAcquireTitle)
	var exists int
	err := s.Server.db.QueryRow("SELECT count(*) FROM titles WHERE id=$1 AND char_id=$2", pkt.TitleID, s.CharID).Scan(&exists)
	if err != nil || exists == 0 {
		s.Server.db.Exec("INSERT INTO titles VALUES ($1, $2, now(), now())", pkt.TitleID, s.CharID)
	} else {
		s.Server.db.Exec("UPDATE titles SET updated_at=now()")
	}
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfResetTitle(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfOperateWarehouse(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOperateWarehouse)
	var t int
	err := s.Server.db.QueryRow("SELECT character_id FROM warehouse WHERE character_id=$1", s.CharID).Scan(&t)
	if err != nil {
		s.Server.db.Exec("INSERT INTO warehouse (character_id) VALUES ($1)", s.CharID)
	}
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(pkt.Operation)
	switch pkt.Operation {
	case 0:
		var count uint8
		itemNames := make([]string, 10)
		equipNames := make([]string, 10)
		s.Server.db.QueryRow(fmt.Sprintf("%s WHERE character_id=$1", warehouseNamesQuery), s.CharID).Scan(&itemNames[0],
			&itemNames[1], &itemNames[2], &itemNames[3], &itemNames[4], &itemNames[5], &itemNames[6], &itemNames[7], &itemNames[8], &itemNames[9], &equipNames[0],
			&equipNames[1], &equipNames[2], &equipNames[3], &equipNames[4], &equipNames[5], &equipNames[6], &equipNames[7], &equipNames[8], &equipNames[9])
		bf.WriteUint32(0)
		bf.WriteUint16(10000) // Usages
		temp := byteframe.NewByteFrame()
		for i, name := range itemNames {
			if len(name) > 0 {
				count++
				temp.WriteUint8(0)
				temp.WriteUint8(uint8(i))
				ps.Uint8(temp, name, true)
			}
		}
		for i, name := range equipNames {
			if len(name) > 0 {
				count++
				temp.WriteUint8(1)
				temp.WriteUint8(uint8(i))
				ps.Uint8(temp, name, true)
			}
		}
		bf.WriteUint8(count)
		bf.WriteBytes(temp.Data())
	case 1:
		bf.WriteUint8(0)
	case 2:
		s.Server.db.Exec(fmt.Sprintf("UPDATE warehouse SET %s%dname=$1 WHERE character_id=$2", pkt.BoxType, pkt.BoxIndex), pkt.Name, s.CharID)
	case 3:
		bf.WriteUint32(0)     // Usage renewal time, >1 = disabled
		bf.WriteUint16(10000) // Usages
	case 4:
		bf.WriteUint32(0)
		bf.WriteUint16(10000) // Usages
		bf.WriteUint8(0)
	}
	// Opcodes
	// 0 = Get box names
	// 1 = Commit usage
	// 2 = Rename
	// 3 = Get usage limit
	// 4 = Get gift box names (doesn't do anything?)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func addWarehouseGift(s *Session, boxType string, giftStack mhfpacket.WarehouseStack) {
	giftBox := getWarehouseBox(s, boxType, 10)
	if boxType == "item" {
		exists := false
		for i, stack := range giftBox {
			if stack.ItemID == giftStack.ItemID {
				exists = true
				giftBox[i].Quantity += giftStack.Quantity
				break
			}
		}
		if exists == false {
			giftBox = append(giftBox, giftStack)
		}
	} else {
		giftBox = append(giftBox, giftStack)
	}
	s.Server.db.Exec(fmt.Sprintf("UPDATE warehouse SET %s10=$1 WHERE character_id=$2", boxType), boxToBytes(giftBox, boxType), s.CharID)
}

func getWarehouseBox(s *Session, boxType string, boxIndex uint8) []mhfpacket.WarehouseStack {
	var data []byte
	s.Server.db.QueryRow(fmt.Sprintf("SELECT %s%d FROM warehouse WHERE character_id=$1", boxType, boxIndex), s.CharID).Scan(&data)
	if len(data) > 0 {
		box := byteframe.NewByteFrameFromBytes(data)
		numStacks := box.ReadUint16()
		stacks := make([]mhfpacket.WarehouseStack, numStacks)
		for i := 0; i < int(numStacks); i++ {
			if boxType == "item" {
				stacks[i].ID = box.ReadUint32()
				stacks[i].Index = box.ReadUint16()
				stacks[i].ItemID = box.ReadUint16()
				stacks[i].Quantity = box.ReadUint16()
				box.ReadUint16()
			} else {
				stacks[i].ID = box.ReadUint32()
				stacks[i].Index = box.ReadUint16()
				stacks[i].EquipType = box.ReadUint16()
				stacks[i].ItemID = box.ReadUint16()
				stacks[i].Data = box.ReadBytes(56)
			}
		}
		return stacks
	} else {
		return make([]mhfpacket.WarehouseStack, 0)
	}
}

func boxToBytes(stacks []mhfpacket.WarehouseStack, boxType string) []byte {
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(len(stacks)))
	for i, stack := range stacks {
		if boxType == "item" {
			bf.WriteUint32(stack.ID)
			bf.WriteUint16(uint16(i + 1))
			bf.WriteUint16(stack.ItemID)
			bf.WriteUint16(stack.Quantity)
			bf.WriteUint16(0)
		} else {
			bf.WriteUint32(stack.ID)
			bf.WriteUint16(uint16(i + 1))
			bf.WriteUint16(stack.EquipType)
			bf.WriteUint16(stack.ItemID)
			bf.WriteBytes(stack.Data)
		}
	}
	bf.WriteUint16(0)
	return bf.Data()
}

func handleMsgMhfEnumerateWarehouse(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateWarehouse)
	box := getWarehouseBox(s, pkt.BoxType, pkt.BoxIndex)
	if len(box) > 0 {
		doAckBufSucceed(s, pkt.AckHandle, boxToBytes(box, pkt.BoxType))
	} else {
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
	}
}

func handleMsgMhfUpdateWarehouse(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateWarehouse)
	box := getWarehouseBox(s, pkt.BoxType, pkt.BoxIndex)
	// Update existing stacks
	var newStacks []mhfpacket.WarehouseStack
	for _, update := range pkt.Updates {
		exists := false
		if pkt.BoxType == "item" {
			for i, stack := range box {
				if stack.Index == update.Index {
					exists = true
					box[i].Quantity = update.Quantity
					break
				}
			}
		} else {
			for i, stack := range box {
				if stack.Index == update.Index {
					exists = true
					box[i].ItemID = update.ItemID
					break
				}
			}
		}
		if exists == false {
			newStacks = append(newStacks, update)
		}
	}
	// Append new stacks
	for _, stack := range newStacks {
		box = append(box, stack)
	}
	// Slice empty stacks
	var cleanedBox []mhfpacket.WarehouseStack
	for _, stack := range box {
		if pkt.BoxType == "item" {
			if stack.Quantity > 0 {
				cleanedBox = append(cleanedBox, stack)
			}
		} else {
			if stack.ItemID != 0 {
				cleanedBox = append(cleanedBox, stack)
			}
		}
	}
	s.Server.db.Exec(fmt.Sprintf("UPDATE warehouse SET %s%d=$1 WHERE character_id=$2", pkt.BoxType, pkt.BoxIndex), boxToBytes(cleanedBox, pkt.BoxType), s.CharID)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}
