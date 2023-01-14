package channelserver

import (
	"fmt"
	"io"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/common/stringsupport"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
)

func handleMsgMhfPostGuildScout(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfPostGuildScout)

	actorCharGuildData, err := GetCharacterGuildData(s, s.CharID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	if actorCharGuildData == nil || !actorCharGuildData.CanRecruit() {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	guildInfo, err := GetGuildInfoByID(s, actorCharGuildData.GuildID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	hasApplication, err := guildInfo.HasApplicationForCharID(s, pkt.CharID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	if hasApplication {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x04})
		return
	}

	transaction, err := s.Server.db.Begin()

	if err != nil {
		panic(err)
	}

	err = guildInfo.CreateApplication(s, pkt.CharID, GuildApplicationTypeInvited, transaction)

	if err != nil {
		rollbackTransaction(s, transaction)
		doAckBufFail(s, pkt.AckHandle, nil)
		panic(err)
	}

	mail := &Mail{
		SenderID:    s.CharID,
		RecipientID: pkt.CharID,
		Subject:     s.Server.dict["guildInviteName"],
		Body: fmt.Sprintf(
			s.Server.dict["guildInvite"],
			guildInfo.Name,
		),
		IsGuildInvite: true,
	}

	err = mail.Send(s, transaction)

	if err != nil {
		rollbackTransaction(s, transaction)
		doAckBufFail(s, pkt.AckHandle, nil)
		return
	}

	err = transaction.Commit()

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, nil)
		panic(err)
	}

	doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfCancelGuildScout(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCancelGuildScout)

	guildCharData, err := GetCharacterGuildData(s, s.CharID)

	if err != nil {
		panic(err)
	}

	if guildCharData == nil || !guildCharData.CanRecruit() {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	guild, err := GetGuildInfoByID(s, guildCharData.GuildID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	err = guild.CancelInvitation(s, pkt.InvitationID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfAnswerGuildScout(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAnswerGuildScout)
	bf := byteframe.NewByteFrame()
	guild, err := GetGuildInfoByCharacterId(s, pkt.LeaderID)

	if err != nil {
		panic(err)
	}

	app, err := guild.GetApplicationForCharID(s, s.CharID, GuildApplicationTypeInvited)

	if app == nil || err != nil {
		s.logger.Warn(
			"Guild invite missing, deleted?",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", s.CharID),
		)
		bf.WriteUint32(7)
		bf.WriteUint32(guild.ID)
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
		return
	}

	var mail []Mail
	if pkt.Answer {
		err = guild.AcceptApplication(s, s.CharID)
		mail = append(mail, Mail{
			RecipientID:     s.CharID,
			Subject:         s.Server.dict["guildInviteSuccessName"],
			Body:            fmt.Sprintf(s.Server.dict["guildInviteSuccess"], guild.Name),
			IsSystemMessage: true,
		})
		mail = append(mail, Mail{
			SenderID:        s.CharID,
			RecipientID:     pkt.LeaderID,
			Subject:         s.Server.dict["guildInviteAcceptedName"],
			Body:            fmt.Sprintf(s.Server.dict["guildInviteAccepted"], guild.Name),
			IsSystemMessage: true,
		})
	} else {
		err = guild.RejectApplication(s, s.CharID)
		mail = append(mail, Mail{
			RecipientID:     s.CharID,
			Subject:         s.Server.dict["guildInviteRejectName"],
			Body:            fmt.Sprintf(s.Server.dict["guildInviteReject"], guild.Name),
			IsSystemMessage: true,
		})
		mail = append(mail, Mail{
			SenderID:        s.CharID,
			RecipientID:     pkt.LeaderID,
			Subject:         s.Server.dict["guildInviteDeclined"],
			Body:            fmt.Sprintf(s.Server.dict["guildInviteDeclined"], guild.Name),
			IsSystemMessage: true,
		})
	}
	if err != nil {
		bf.WriteUint32(7)
		bf.WriteUint32(guild.ID)
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
	} else {
		bf.WriteUint32(0)
		bf.WriteUint32(guild.ID)
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
		for _, m := range mail {
			// nolint:errcheck
			m.Send(s, nil)
		}
	}
}

func handleMsgMhfGetGuildScoutList(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildScoutList)
	guildInfo, err := GetGuildInfoByCharacterId(s, s.CharID)

	if guildInfo == nil && s.PrevGuildID == 0 {
		doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		return
	} else {
		guildInfo, err = GetGuildInfoByID(s, s.PrevGuildID)
		if guildInfo == nil || err != nil {
			doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
	}

	if guildInfo == nil && s.PrevGuildID == 0 {
		doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		return
	} else {
		guildInfo, err = GetGuildInfoByID(s, s.PrevGuildID)
		if guildInfo == nil || err != nil {
			doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
	}

	rows, err := s.Server.db.Queryx(`
		SELECT c.id, c.name, ga.actor_id
			FROM guild_applications ga 
			JOIN characters c ON c.id = ga.character_id
		WHERE ga.guild_id = $1 AND ga.application_type = 'invited'
	`, guildInfo.ID)

	if err != nil {
		s.logger.Error("failed to retrieve scouted characters", zap.Error(err))
		doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	defer rows.Close()

	bf := byteframe.NewByteFrame()

	bf.SetBE()

	// Result count, we will overwrite this later
	bf.WriteUint32(0x00)

	count := uint32(0)

	for rows.Next() {
		var charName string
		var charID uint32
		var actorID uint32

		err = rows.Scan(&charID, &charName, &actorID)

		if err != nil {
			doAckSimpleFail(s, pkt.AckHandle, nil)
			panic(err)
		}

		// This seems to be used as a unique ID for the invitation sent
		// we can just use the charID and then filter on guild_id+charID when performing operations
		// this might be a problem later with mails sent referencing IDs but we'll see.
		bf.WriteUint32(charID)
		bf.WriteUint32(actorID)
		bf.WriteUint32(charID)
		bf.WriteUint32(uint32(time.Now().Unix()))
		bf.WriteUint16(0x00) // HR?
		bf.WriteUint16(0x00) // GR?

		charNameBytes, _ := stringsupport.ConvertUTF8ToShiftJIS(charName)

		bf.WriteBytes(charNameBytes)
		bf.WriteBytes(make([]byte, 32-len(charNameBytes))) // Fixed length string
		count++
	}

	_, err = bf.Seek(0, io.SeekStart)

	if err != nil {
		panic(err)
	}

	bf.WriteUint32(count)

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetRejectGuildScout(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetRejectGuildScout)

	row := s.Server.db.QueryRow("SELECT restrict_guild_scout FROM characters WHERE id=$1", s.CharID)

	var currentStatus bool

	err := row.Scan(&currentStatus)

	if err != nil {
		s.logger.Error(
			"failed to retrieve character guild scout status",
			zap.Error(err),
			zap.Uint32("charID", s.CharID),
		)
		doAckSimpleFail(s, pkt.AckHandle, nil)
		return
	}

	response := uint8(0x00)

	if currentStatus {
		response = 0x01
	}

	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, response})
}

func handleMsgMhfSetRejectGuildScout(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetRejectGuildScout)

	_, err := s.Server.db.Exec("UPDATE characters SET restrict_guild_scout=$1 WHERE id=$2", pkt.Reject, s.CharID)

	if err != nil {
		s.logger.Error(
			"failed to update character guild scout status",
			zap.Error(err),
			zap.Uint32("charID", s.CharID),
		)
		doAckSimpleFail(s, pkt.AckHandle, nil)
		return
	}

	doAckSimpleSucceed(s, pkt.AckHandle, nil)
}
