package channelserver

import (
	"erupe-ce/network/mhfpacket"
	"erupe-ce/server/channelserver/compression/deltacomp"
	"erupe-ce/server/channelserver/compression/nullcomp"

	"go.uber.org/zap"
)

func handleMsgMhfLoadPlateData(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadPlateData)
	var data []byte
	err := s.Server.db.QueryRow("SELECT platedata FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if err != nil {
		s.logger.Error("Failed to load platedata", zap.Error(err))
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSavePlateData(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSavePlateData)

	if pkt.IsDataDiff {
		var data []byte

		// Load existing save
		err := s.Server.db.QueryRow("SELECT platedata FROM characters WHERE id = $1", s.CharID).Scan(&data)
		if err != nil {
			s.logger.Error("Failed to load platedata", zap.Error(err))
			doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
			return
		}

		if len(data) > 0 {
			// Decompress
			s.logger.Info("Decompressing...")
			data, err = nullcomp.Decompress(data)
			if err != nil {
				s.logger.Error("Failed to decompress platedata", zap.Error(err))
				doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
				return
			}
		} else {
			// create empty save if absent
			data = make([]byte, 0x1AF20)
		}

		// Perform diff and compress it to write back to db
		s.logger.Info("Diffing...")
		saveOutput, err := nullcomp.Compress(deltacomp.ApplyDataDiff(pkt.RawDataPayload, data))
		if err != nil {
			s.logger.Error("Failed to diff and compress platedata", zap.Error(err))
			doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
			return
		}

		_, err = s.Server.db.Exec("UPDATE characters SET platedata=$1 WHERE id=$2", saveOutput, s.CharID)
		if err != nil {
			s.logger.Error("Failed to save platedata", zap.Error(err))
			doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
			return
		}

		s.logger.Info("Wrote recompressed platedata back to DB")
	} else {
		dumpSaveData(s, pkt.RawDataPayload, "platedata")
		// simply update database, no extra processing
		_, err := s.Server.db.Exec("UPDATE characters SET platedata=$1 WHERE id=$2", pkt.RawDataPayload, s.CharID)
		if err != nil {
			s.logger.Error("Failed to save platedata", zap.Error(err))
		}
	}

	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfLoadPlateBox(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadPlateBox)
	var data []byte
	err := s.Server.db.QueryRow("SELECT platebox FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if err != nil {
		s.logger.Error("Failed to load platebox", zap.Error(err))
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSavePlateBox(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSavePlateBox)

	if pkt.IsDataDiff {
		var data []byte

		// Load existing save
		err := s.Server.db.QueryRow("SELECT platebox FROM characters WHERE id = $1", s.CharID).Scan(&data)
		if err != nil {
			s.logger.Error("Failed to load platebox", zap.Error(err))
			doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
			return
		}

		// Decompress
		if len(data) > 0 {
			// Decompress
			s.logger.Info("Decompressing...")
			data, err = nullcomp.Decompress(data)
			if err != nil {
				s.logger.Error("Failed to decompress platebox", zap.Error(err))
				doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
				return
			}
		} else {
			// create empty save if absent
			data = make([]byte, 0x820)
		}

		// Perform diff and compress it to write back to db
		s.logger.Info("Diffing...")
		saveOutput, err := nullcomp.Compress(deltacomp.ApplyDataDiff(pkt.RawDataPayload, data))
		if err != nil {
			s.logger.Error("Failed to diff and compress platebox", zap.Error(err))
			doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
			return
		}

		_, err = s.Server.db.Exec("UPDATE characters SET platebox=$1 WHERE id=$2", saveOutput, s.CharID)
		if err != nil {
			s.logger.Error("Failed to save platebox", zap.Error(err))
			doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
			return
		}

		s.logger.Info("Wrote recompressed platebox back to DB")
	} else {
		dumpSaveData(s, pkt.RawDataPayload, "platebox")
		// simply update database, no extra processing
		_, err := s.Server.db.Exec("UPDATE characters SET platebox=$1 WHERE id=$2", pkt.RawDataPayload, s.CharID)
		if err != nil {
			s.logger.Error("Failed to save platebox", zap.Error(err))
		}
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfLoadPlateMyset(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadPlateMyset)
	var data []byte
	err := s.Server.db.QueryRow("SELECT platemyset FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if len(data) == 0 {
		s.logger.Error("Failed to load platemyset", zap.Error(err))
		data = make([]byte, 0x780)
	}
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSavePlateMyset(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSavePlateMyset)
	// looks to always return the full thing, simply update database, no extra processing
	dumpSaveData(s, pkt.RawDataPayload, "platemyset")
	_, err := s.Server.db.Exec("UPDATE characters SET platemyset=$1 WHERE id=$2", pkt.RawDataPayload, s.CharID)
	if err != nil {
		s.logger.Error("Failed to save platemyset", zap.Error(err))
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}
