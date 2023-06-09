package channelserver

import (
	"erupe-ce/common/byteframe"
	"erupe-ce/common/stringsupport"
	"erupe-ce/network/mhfpacket"
	"erupe-ce/server/channelserver/compression/deltacomp"
	"erupe-ce/server/channelserver/compression/nullcomp"

	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

func handleMsgMhfLoadPartner(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadPartner)
	var data []byte
	err := s.Server.db.QueryRow("SELECT partner FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		s.logger.Error("Failed to load partner", zap.Error(err))
		data = make([]byte, 9)
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfSavePartner(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSavePartner)
	dumpSaveData(s, pkt.RawDataPayload, "partner")

	_, err := s.Server.db.Exec("UPDATE characters SET partner=$1 WHERE id=$2", pkt.RawDataPayload, s.CharID)
	if err != nil {
		s.logger.Error("Failed to save partner", zap.Error(err))
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfLoadLegendDispatch(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadLegendDispatch)
	bf := byteframe.NewByteFrame()
	legendDispatch := []struct {
		Unk       uint32
		Timestamp uint32
	}{
		{0, uint32(TimeMidnight().Add(-12 * time.Hour).Unix())},
		{0, uint32(TimeMidnight().Add(12 * time.Hour).Unix())},
		{0, uint32(TimeMidnight().Add(36 * time.Hour).Unix())},
	}
	bf.WriteUint8(uint8(len(legendDispatch)))
	for _, dispatch := range legendDispatch {
		bf.WriteUint32(dispatch.Unk)
		bf.WriteUint32(dispatch.Timestamp)
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfLoadHunterNavi(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadHunterNavi)
	var data []byte
	err := s.Server.db.QueryRow("SELECT hunternavi FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		s.logger.Error("Failed to load hunternavi", zap.Error(err))
		data = make([]byte, 0x226)
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSaveHunterNavi(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveHunterNavi)
	if pkt.IsDataDiff {
		var data []byte
		// Load existing save
		err := s.Server.db.QueryRow("SELECT hunternavi FROM characters WHERE id = $1", s.CharID).Scan(&data)
		if err != nil {
			s.logger.Error("Failed to load hunternavi", zap.Error(err))
		}

		// Check if we actually had any hunternavi data, using a blank buffer if not.
		// This is requried as the client will try to send a diff after character creation without a prior MsgMhfSaveHunterNavi packet.
		if len(data) == 0 {
			data = make([]byte, 0x226)
		}

		// Perform diff and compress it to write back to db
		s.logger.Info("Diffing...")
		saveOutput := deltacomp.ApplyDataDiff(pkt.RawDataPayload, data)

		_, err = s.Server.db.Exec("UPDATE characters SET hunternavi=$1 WHERE id=$2", saveOutput, s.CharID)
		if err != nil {
			s.logger.Error("Failed to save hunternavi", zap.Error(err))
		}
		s.logger.Info("Wrote recompressed hunternavi back to DB")
	} else {
		dumpSaveData(s, pkt.RawDataPayload, "hunternavi")
		// simply update database, no extra processing
		_, err := s.Server.db.Exec("UPDATE characters SET hunternavi=$1 WHERE id=$2", pkt.RawDataPayload, s.CharID)
		if err != nil {
			s.logger.Error("Failed to save hunternavi", zap.Error(err))
		}
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfMercenaryHuntdata(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfMercenaryHuntdata)
	if pkt.Unk0 == 1 {
		// Format:
		// uint8 Hunts
		// struct Hunt
		//   uint32 HuntID
		//   uint32 MonID
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 1))
	} else {
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 0))
	}
}

func handleMsgMhfEnumerateMercenaryLog(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateMercenaryLog)
	bf := byteframe.NewByteFrame()
	bf.WriteUint32(0)
	// Format:
	// struct Log
	//   uint32 Timestamp
	//   []byte Name (len 18)
	//   uint8 Unk
	//   uint8 Unk
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfCreateMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCreateMercenary)
	bf := byteframe.NewByteFrame()
	var nextID uint32
	// nolint:errcheck
	_ = s.Server.db.QueryRow("SELECT nextval('rasta_id_seq')").Scan(&nextID)
	s.Server.db.Exec("UPDATE characters SET rasta_id=$1 WHERE id=$2", nextID, s.CharID)
	bf.WriteUint32(nextID)
	doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfSaveMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveMercenary)
	dumpSaveData(s, pkt.MercData, "mercenary")
	if len(pkt.MercData) > 0 {
		temp := byteframe.NewByteFrameFromBytes(pkt.MercData)
		// nolint:errcheck
		s.Server.db.Exec("UPDATE characters SET savemercenary=$1, rasta_id=$2 WHERE id=$3", pkt.MercData, temp.ReadUint32(), s.CharID)
	}
	// nolint:errcheck
	s.Server.db.Exec("UPDATE characters SET gcp=$1, pact_id=$2 WHERE id=$3", pkt.GCP, pkt.PactMercID, s.CharID)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfReadMercenaryW(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMercenaryW)
	if pkt.Op > 0 {
		bf := byteframe.NewByteFrame()
		var pactID uint32
		var name string
		var cid uint32

		s.Server.db.QueryRow("SELECT pact_id FROM characters WHERE id=$1", s.CharID).Scan(&pactID)
		if pactID > 0 {
			s.Server.db.QueryRow("SELECT name, id FROM characters WHERE rasta_id = $1", pactID).Scan(&name, &cid)
			bf.WriteUint8(1) // numLends
			bf.WriteUint32(pactID)
			bf.WriteUint32(cid)
			bf.WriteBool(false) // ?
			bf.WriteUint32(uint32(TimeAdjusted().Add(time.Hour * 24 * -8).Unix()))
			bf.WriteUint32(uint32(TimeAdjusted().Add(time.Hour * 24 * -1).Unix()))
			bf.WriteBytes(stringsupport.PaddedString(name, 18, true))
		} else {
			bf.WriteUint8(0)
		}

		if pkt.Op < 2 {
			var loans uint8
			temp := byteframe.NewByteFrame()
			rows, _ := s.Server.db.Query("SELECT name, id, pact_id FROM characters WHERE pact_id=(SELECT rasta_id FROM characters WHERE id=$1)", s.CharID)
			for rows.Next() {
				loans++
				rows.Scan(&name, &cid, &pactID)
				temp.WriteUint32(pactID)
				temp.WriteUint32(cid)
				temp.WriteUint32(uint32(TimeAdjusted().Add(time.Hour * 24 * -8).Unix()))
				temp.WriteUint32(uint32(TimeAdjusted().Add(time.Hour * 24 * -1).Unix()))
				temp.WriteBytes(stringsupport.PaddedString(name, 18, true))
			}
			bf.WriteUint8(loans)
			bf.WriteBytes(temp.Data())
		}
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
		return
	}
	var data []byte
	var gcp uint32
	// nolint:errcheck
	s.Server.db.QueryRow("SELECT savemercenary FROM characters WHERE id = $1", s.CharID).Scan(&data)
	// nolint:errcheck
	s.Server.db.QueryRow("SELECT COALESCE(gcp, 0) FROM characters WHERE id = $1", s.CharID).Scan(&gcp)

	resp := byteframe.NewByteFrame()
	resp.WriteUint16(0)
	if len(data) == 0 {
		resp.WriteBool(false)
	} else {
		resp.WriteBool(true)
		resp.WriteBytes(data)
	}
	resp.WriteUint32(gcp)
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfReadMercenaryM(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMercenaryM)
	var data []byte
	// nolint:errcheck
	s.Server.db.QueryRow("SELECT savemercenary FROM characters WHERE id = $1", pkt.CharID).Scan(&data)
	resp := byteframe.NewByteFrame()
	if len(data) == 0 {
		resp.WriteBool(false)
	} else {
		resp.WriteBytes(data)
	}
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfContractMercenary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfContractMercenary)
	switch pkt.Op {
	case 0:
		s.Server.db.Exec("UPDATE characters SET pact_id=$1 WHERE id=$2", pkt.PactMercID, s.CharID)
	case 1: // Cancel lend
		s.Server.db.Exec("UPDATE characters SET pact_id=0 WHERE id=$1", s.CharID)
	case 2: // Cancel loan
		s.Server.db.Exec("UPDATE characters SET pact_id=0 WHERE id=$1", pkt.CID)
	}
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfLoadOtomoAirou(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadOtomoAirou)
	var data []byte
	err := s.Server.db.QueryRow("SELECT otomoairou FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		s.logger.Error("Failed to load otomoairou", zap.Error(err))
		data = make([]byte, 10)
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSaveOtomoAirou(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveOtomoAirou)
	dumpSaveData(s, pkt.RawDataPayload, "otomoairou")
	decomp, err := nullcomp.Decompress(pkt.RawDataPayload[1:])
	if err != nil {
		s.logger.Error("Failed to decompress airou", zap.Error(err))
		doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		return
	}
	bf := byteframe.NewByteFrameFromBytes(decomp)
	save := byteframe.NewByteFrame()
	var catsExist uint8
	save.WriteUint8(0)

	cats := bf.ReadUint8()
	for i := 0; i < int(cats); i++ {
		dataLen := bf.ReadUint32()
		catID := bf.ReadUint32()
		if catID == 0 {
			_ = s.Server.db.QueryRow("SELECT nextval('airou_id_seq')").Scan(&catID)
		}
		exists := bf.ReadBool()
		data := bf.ReadBytes(uint(dataLen) - 5)
		if exists {
			catsExist++
			save.WriteUint32(dataLen)
			save.WriteUint32(catID)
			save.WriteBool(exists)
			save.WriteBytes(data)
		}
	}
	save.WriteBytes(bf.DataFromCurrent())
	// nolint:errcheck
	save.Seek(0, 0)
	save.WriteUint8(catsExist)
	comp, err := nullcomp.Compress(save.Data())
	if err != nil {
		s.logger.Error("Failed to compress airou", zap.Error(err))
	} else {
		comp = append([]byte{0x01}, comp...)
		// nolint:errcheck
		s.Server.db.Exec("UPDATE characters SET otomoairou=$1 WHERE id=$2", comp, s.CharID)
	}
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfEnumerateAiroulist(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateAiroulist)
	resp := byteframe.NewByteFrame()
	if _, err := os.Stat(filepath.Join(s.Server.erupeConfig.BinPath, "airoulist.bin")); err == nil {
		data, _ := os.ReadFile(filepath.Join(s.Server.erupeConfig.BinPath, "airoulist.bin"))
		resp.WriteBytes(data)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
		return
	}
	airouList := getGuildAirouList(s)
	resp.WriteUint16(uint16(len(airouList)))
	resp.WriteUint16(uint16(len(airouList)))
	for _, cat := range airouList {
		resp.WriteUint32(cat.CatID)
		resp.WriteBytes(cat.CatName)
		resp.WriteUint32(cat.Experience)
		resp.WriteUint8(cat.Personality)
		resp.WriteUint8(cat.Class)
		resp.WriteUint8(cat.WeaponType)
		resp.WriteUint16(cat.WeaponID)
		resp.WriteUint32(0) // 32 bit unix timestamp, either time at which the cat stops being fatigued or the time at which it started
	}
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

// CatDefinition holds values needed to populate the guild cat list
type CatDefinition struct {
	CatID       uint32
	CatName     []byte
	CurrentTask uint8
	Personality uint8
	Class       uint8
	Experience  uint32
	WeaponType  uint8
	WeaponID    uint16
}

func getGuildAirouList(s *Session) []CatDefinition {
	var guild *Guild
	var err error
	var guildCats []CatDefinition

	// returning 0 cats on any guild issues
	// can probably optimise all of the guild queries pretty heavily
	guild, err = GetGuildInfoByCharacterId(s, s.CharID)
	if err != nil {
		return guildCats
	}

	// Get cats used recently
	// Retail reset at midday, 12 hours is a midpoint
	tempBanDuration := 43200 - (1800) // Minus hunt time
	bannedCats := make(map[uint32]int)
	var csvTemp string
	rows, err := s.Server.db.Query(`SELECT cats_used
	FROM guild_hunts gh
	INNER JOIN characters c
	ON gh.host_id = c.id
	WHERE c.id=$1 AND gh.return+$2>$3`, s.CharID, tempBanDuration, TimeAdjusted().Unix())
	if err != nil {
		s.logger.Warn("Failed to get recently used airous", zap.Error(err))
	}
	for rows.Next() {
		// nolint:errcheck // Error return value of `rows.Scan` is not checked
		rows.Scan(&csvTemp)
		for i, j := range stringsupport.CSVElems(csvTemp) {
			bannedCats[uint32(j)] = i
		}
	}

	// ellie's GetGuildMembers didn't seem to pull leader?
	rows, err = s.Server.db.Query(`SELECT c.otomoairou
	FROM characters c
	INNER JOIN guild_characters gc
	ON gc.character_id = c.id
	WHERE gc.guild_id = $1 AND c.otomoairou IS NOT NULL
	ORDER BY c.id ASC
	LIMIT 60;`, guild.ID)
	if err != nil {
		s.logger.Warn("Selecting otomoairou based on guild failed", zap.Error(err))
		return guildCats
	}

	for rows.Next() {
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			s.logger.Warn("select failure", zap.Error(err))
			continue
		} else if len(data) == 0 {
			// non extant cats that aren't null in DB
			continue
		}
		// first byte has cat existence in general, can skip if 0
		if data[0] == 1 {
			decomp, err := nullcomp.Decompress(data[1:])
			if err != nil {
				s.logger.Warn("decomp failure", zap.Error(err))
				continue
			}
			bf := byteframe.NewByteFrameFromBytes(decomp)
			cats := GetCatDetails(bf)
			for _, cat := range cats {
				_, exists := bannedCats[cat.CatID]
				if cat.CurrentTask == 4 && !exists {
					guildCats = append(guildCats, cat)
				}
			}
		}
	}
	return guildCats
}

func GetCatDetails(bf *byteframe.ByteFrame) []CatDefinition {
	catCount := bf.ReadUint8()
	cats := make([]CatDefinition, catCount)
	for x := 0; x < int(catCount); x++ {
		var catDef CatDefinition
		// cat sometimes has additional bytes for whatever reason, gift items? timestamp?
		// until actual variance is known we can just seek to end based on start
		catDefLen := bf.ReadUint32()
		catStart, _ := bf.Seek(0, io.SeekCurrent)

		catDef.CatID = bf.ReadUint32()
		// nolint:errcheck
		bf.Seek(1, io.SeekCurrent)        // unknown value, probably a bool
		catDef.CatName = bf.ReadBytes(18) // always 18 len, reads first null terminated string out of section and discards rest
		catDef.CurrentTask = bf.ReadUint8()
		// nolint:errcheck
		bf.Seek(16, io.SeekCurrent) // appearance data and what is seemingly null bytes
		catDef.Personality = bf.ReadUint8()
		catDef.Class = bf.ReadUint8()
		// nolint:errcheck
		bf.Seek(5, io.SeekCurrent)          // affection and colour sliders
		catDef.Experience = bf.ReadUint32() // raw cat rank points, doesn't have a rank
		// nolint:errcheck
		bf.Seek(1, io.SeekCurrent)         // bool for weapon being equipped
		catDef.WeaponType = bf.ReadUint8() // weapon type, presumably always 6 for melee?
		catDef.WeaponID = bf.ReadUint16()  // weapon id
		// nolint:errcheck
		bf.Seek(catStart+int64(catDefLen), io.SeekStart)
		cats[x] = catDef
	}
	return cats
}
