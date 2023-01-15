package channelserver

import (
	"fmt"
	//nolint:staticcheck
	"io/ioutil"
	"os"
	"path/filepath"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
)

func handleMsgSysGetFile(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysGetFile)

	if pkt.IsScenario {
		if s.Server.erupeConfig.DevModeOptions.QuestDebugTools && s.Server.erupeConfig.DevMode {
			s.logger.Debug(
				"Scenario",
				zap.Uint8("CategoryID", pkt.ScenarioIdentifer.CategoryID),
				zap.Uint32("MainID", pkt.ScenarioIdentifer.MainID),
				zap.Uint8("ChapterID", pkt.ScenarioIdentifer.ChapterID),
				zap.Uint8("Flags", pkt.ScenarioIdentifer.Flags),
			)
		}
		filename := fmt.Sprintf("%d_0_0_0_S%d_T%d_C%d", pkt.ScenarioIdentifer.CategoryID, pkt.ScenarioIdentifer.MainID, pkt.ScenarioIdentifer.Flags, pkt.ScenarioIdentifer.ChapterID)
		// Read the scenario file.
		data, err := os.ReadFile(filepath.Join(s.Server.erupeConfig.BinPath, fmt.Sprintf("scenarios/%s.bin", filename)))
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to open file: %s/scenarios/%s.bin", s.Server.erupeConfig.BinPath, filename))
			// This will crash the game.
			doAckBufSucceed(s, pkt.AckHandle, data)
			return
		}
		doAckBufSucceed(s, pkt.AckHandle, data)
	} else {
		if _, err := os.Stat(filepath.Join(s.Server.erupeConfig.BinPath, "quest_override.bin")); err == nil {
			data, err := os.ReadFile(filepath.Join(s.Server.erupeConfig.BinPath, "quest_override.bin"))
			if err != nil {
				panic(err)
			}
			doAckBufSucceed(s, pkt.AckHandle, data)
		} else {
			if s.Server.erupeConfig.DevModeOptions.QuestDebugTools && s.Server.erupeConfig.DevMode {
				s.logger.Debug(
					"Quest",
					zap.String("Filename", pkt.Filename),
				)
			}

			// Get quest file.
			data, err := s.Server.questLoader.QuestBinById(pkt.Filename)

			if err != nil {
				filepath := filepath.Join(s.Server.erupeConfig.BinPath, fmt.Sprintf("quests/%s.bin", pkt.Filename))
				quest, err_bin := ioutil.ReadFile(filepath)

				if err_bin != nil {
					s.logger.Fatal(fmt.Sprintf("Failed to open quest file: quests/%s.bin", pkt.Filename), zap.Error(err_bin))
				}
				doAckBufSucceed(s, pkt.AckHandle, quest)
				return
			}

			doAckBufSucceed(s, pkt.AckHandle, data)
		}
	}
}

func handleMsgMhfLoadFavoriteQuest(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadFavoriteQuest)
	var data []byte
	err := s.Server.db.QueryRow("SELECT savefavoritequest FROM characters WHERE id = $1", s.CharID).Scan(&data)
	if err == nil && len(data) > 0 {
		doAckBufSucceed(s, pkt.AckHandle, data)
	} else {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}
}

func handleMsgMhfSaveFavoriteQuest(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSaveFavoriteQuest)
	// nolint:errcheck
	s.Server.db.Exec("UPDATE characters SET savefavoritequest=$1 WHERE id=$2", pkt.Data, s.CharID)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfEnumerateQuest(s *Session, p mhfpacket.MHFPacket) {
	// local files are easier for now, probably best would be to generate dynamically
	pkt := p.(*mhfpacket.MsgMhfEnumerateQuest)

	data, err := s.Server.questLoader.Quests(pkt.Take, pkt.Skip)

	if err != nil {
		fmt.Printf("questlists/list_%d.bin\n", pkt.Skip)
		fmt.Printf("Failed to load quest: %s\n", err.Error())
		stubEnumerateNoResults(s, pkt.AckHandle)
	} else {
		doAckBufSucceed(s, pkt.AckHandle, data)
	}
}

func handleMsgMhfEnterTournamentQuest(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdBonusQuestInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdBonusQuestInfo)

	udBonusQuestInfos := []struct {
		Unk0      uint8
		Unk1      uint8
		StartTime uint32 // Unix timestamp (seconds)
		EndTime   uint32 // Unix timestamp (seconds)
		Unk4      uint32
		Unk5      uint8
		Unk6      uint8
	}{} // Blank stub array.

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(len(udBonusQuestInfos)))
	for _, q := range udBonusQuestInfos {
		resp.WriteUint8(q.Unk0)
		resp.WriteUint8(q.Unk1)
		resp.WriteUint32(q.StartTime)
		resp.WriteUint32(q.EndTime)
		resp.WriteUint32(q.Unk4)
		resp.WriteUint8(q.Unk5)
		resp.WriteUint8(q.Unk6)
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}
