package channelserver

import (
	"database/sql"
	"erupe-ce/common/stringsupport"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/binpacket"
	"erupe-ce/network/mhfpacket"
	"go.uber.org/zap"
)

type Mail struct {
	ID                   int       `db:"id"`
	SenderID             uint32    `db:"sender_id"`
	RecipientID          uint32    `db:"recipient_id"`
	Subject              string    `db:"subject"`
	Body                 string    `db:"body"`
	Read                 bool      `db:"read"`
	Deleted              bool      `db:"deleted"`
	Locked               bool      `db:"locked"`
	AttachedItemReceived bool      `db:"attached_item_received"`
	AttachedItemID       uint16    `db:"attached_item"`
	AttachedItemAmount   uint16    `db:"attached_item_amount"`
	CreatedAt            time.Time `db:"created_at"`
	IsGuildInvite        bool      `db:"is_guild_invite"`
	SenderName           string    `db:"sender_name"`
}

func (m *Mail) Send(s *Session, transaction *sql.Tx) error {
	query := `
		INSERT INTO mail (sender_id, recipient_id, subject, body, attached_item, attached_item_amount, is_guild_invite)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var err error

	if transaction == nil {
		_, err = s.server.db.Exec(query, m.SenderID, m.RecipientID, m.Subject, m.Body, m.AttachedItemID, m.AttachedItemAmount, m.IsGuildInvite)
	} else {
		_, err = transaction.Exec(query, m.SenderID, m.RecipientID, m.Subject, m.Body, m.AttachedItemID, m.AttachedItemAmount, m.IsGuildInvite)
	}

	if err != nil {
		s.logger.Error(
			"failed to send mail",
			zap.Error(err),
			zap.Uint32("senderID", m.SenderID),
			zap.Uint32("recipientID", m.RecipientID),
			zap.String("subject", m.Subject),
			zap.String("body", m.Body),
			zap.Uint16("itemID", m.AttachedItemID),
			zap.Uint16("itemAmount", m.AttachedItemAmount),
			zap.Bool("isGuildInvite", m.IsGuildInvite),
		)
		return err
	}

	return nil
}

func (m *Mail) MarkRead(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE mail SET read = true WHERE id = $1
	`, m.ID)

	if err != nil {
		s.logger.Error(
			"failed to mark mail as read",
			zap.Error(err),
			zap.Int("mailID", m.ID),
		)
		return err
	}

	return nil
}

func (m *Mail) MarkDeleted(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE mail SET deleted = true WHERE id = $1
	`, m.ID)

	if err != nil {
		s.logger.Error(
			"failed to mark mail as deleted",
			zap.Error(err),
			zap.Int("mailID", m.ID),
		)
		return err
	}

	return nil
}

func (m *Mail) MarkAcquired(s *Session) error {
	_, err := s.server.db.Exec(`
		UPDATE mail SET attached_item_received = true WHERE id = $1
	`, m.ID)

	if err != nil {
		s.logger.Error(
			"failed to mark mail item as claimed",
			zap.Error(err),
			zap.Int("mailID", m.ID),
		)
		return err
	}

	return nil
}

func (m *Mail) MarkLocked(s *Session, locked bool) error {
	_, err := s.server.db.Exec(`
		UPDATE mail SET locked = $1 WHERE id = $2
	`, locked, m.ID)

	if err != nil {
		s.logger.Error(
			"failed to mark mail as locked",
			zap.Error(err),
			zap.Int("mailID", m.ID),
		)
		return err
	}

	return nil
}

func GetMailListForCharacter(s *Session, charID uint32) ([]Mail, error) {
	rows, err := s.server.db.Queryx(`
		SELECT
			m.id,
			m.sender_id,
			m.recipient_id,
			m.subject,
			m.read,
			m.attached_item_received,
			m.attached_item,
			m.attached_item_amount,
			m.created_at,
			m.is_guild_invite,
			m.deleted,
			m.locked,
			c.name as sender_name
		FROM mail m
			JOIN characters c ON c.id = m.sender_id
		WHERE recipient_id = $1 AND m.deleted = false
		ORDER BY m.created_at DESC, id DESC
		LIMIT 32
	`, charID)

	if err != nil {
		s.logger.Error("failed to get mail for character", zap.Error(err), zap.Uint32("charID", charID))
		return nil, err
	}

	defer rows.Close()

	allMail := make([]Mail, 0)

	for rows.Next() {
		mail := Mail{}

		err := rows.StructScan(&mail)

		if err != nil {
			return nil, err
		}

		allMail = append(allMail, mail)
	}

	return allMail, nil
}

func GetMailByID(s *Session, ID int) (*Mail, error) {
	row := s.server.db.QueryRowx(`
		SELECT
			m.id,
			m.sender_id,
			m.recipient_id,
			m.subject,
			m.read,
			m.body,
			m.attached_item_received,
			m.attached_item,
			m.attached_item_amount,
			m.created_at,
			m.is_guild_invite,
			m.deleted,
			m.locked,
			c.name as sender_name
		FROM mail m
			JOIN characters c ON c.id = m.sender_id
		WHERE m.id = $1
		LIMIT 1
	`, ID)

	mail := &Mail{}

	err := row.StructScan(mail)

	if err != nil {
		s.logger.Error(
			"failed to retrieve mail",
			zap.Error(err),
			zap.Int("mailID", ID),
		)
		return nil, err
	}

	return mail, nil
}

func SendMailNotification(s *Session, m *Mail, recipient *Session) {
	senderName, err := getCharacterName(s, m.SenderID)

	if err != nil {
		panic(err)
	}

	bf := byteframe.NewByteFrame()

	notification := &binpacket.MsgBinMailNotify{
		SenderName: senderName,
	}

	notification.Build(bf)

	castedBinary := &mhfpacket.MsgSysCastedBinary{
		CharID:         m.SenderID,
		BroadcastType:  0x00,
		MessageType:    BinaryMessageTypeMailNotify,
		RawDataPayload: bf.Data(),
	}

	castedBinary.Build(bf, s.clientContext)

	recipient.QueueSendMHF(castedBinary)
}

func getCharacterName(s *Session, charID uint32) (string, error) {
	row := s.server.db.QueryRow("SELECT name FROM characters WHERE id = $1", charID)

	charName := ""

	err := row.Scan(&charName)

	if err != nil {
		return "", err
	}

	return charName, nil
}

func handleMsgMhfReadMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadMail)

	mailId := s.mailList[pkt.AccIndex]

	if mailId == 0 {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic("attempting to read mail that doesn't exist in session map")
	}

	mail, err := GetMailByID(s, mailId)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	_ = mail.MarkRead(s)

	bf := byteframe.NewByteFrame()

	body := stringsupport.UTF8ToSJIS(mail.Body)
	bf.WriteNullTerminatedBytes(body)

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfListMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfListMail)

	mail, err := GetMailListForCharacter(s, s.charID)

	if err != nil {
		doAckBufFail(s, pkt.AckHandle, make([]byte, 4))
		panic(err)
	}

	if s.mailList == nil {
		s.mailList = make([]int, 256)
	}

	msg := byteframe.NewByteFrame()

	msg.WriteUint32(uint32(len(mail)))

	startIndex := s.mailAccIndex

	for i, m := range mail {
		accIndex := startIndex + uint8(i)
		s.mailList[accIndex] = m.ID
		s.mailAccIndex++

		itemAttached := m.AttachedItemID != 0

		msg.WriteUint32(m.SenderID)
		msg.WriteUint32(uint32(m.CreatedAt.Unix()))

		msg.WriteUint8(accIndex)
		msg.WriteUint8(uint8(i))

		flags := uint8(0x00)

		if m.Read {
			flags |= 0x01
		}

		if m.Locked {
			flags |= 0x02
		}

		// System message, hides ID
		// flags |= 0x04

		// Workaround until EN mail items are patched
		if s.server.erupeConfig.DevMode && s.server.erupeConfig.DevModeOptions.DisableMailItems {
			if itemAttached {
				flags |= 0x08
			}
		} else {
			if m.AttachedItemReceived {
				flags |= 0x08
			}
		}

		if m.IsGuildInvite {
			flags |= 0x10
		}

		msg.WriteUint8(flags)
		msg.WriteBool(itemAttached)
		msg.WriteUint8(16)
		msg.WriteUint8(21)
		msg.WriteBytes(stringsupport.PaddedString(m.Subject, 16, true))
		msg.WriteBytes(stringsupport.PaddedString(m.SenderName, 21, true))
		if itemAttached {
			msg.WriteUint16(m.AttachedItemAmount)
			msg.WriteUint16(m.AttachedItemID)
		}
	}

	doAckBufSucceed(s, pkt.AckHandle, msg.Data())
}

func handleMsgMhfOprtMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOprtMail)

	mail, err := GetMailByID(s, s.mailList[pkt.AccIndex])
	if err != nil {
		panic(err)
	}

	switch pkt.Operation {
	case mhfpacket.OPERATE_MAIL_DELETE:
		err = mail.MarkDeleted(s)
	case mhfpacket.OPERATE_MAIL_LOCK:
		err = mail.MarkLocked(s, true)
	case mhfpacket.OPERATE_MAIL_UNLOCK:
		err = mail.MarkLocked(s, false)
	case mhfpacket.OPERATE_MAIL_ACQUIRE_ITEM:
		err = mail.MarkAcquired(s)
	}

	if err != nil {
		panic(err)
	}

	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfSendMail(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSendMail)
	query := `
		INSERT INTO mail (sender_id, recipient_id, subject, body, attached_item, attached_item_amount, is_guild_invite)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if pkt.RecipientID == 0 { // Guild mail
		g, err := GetGuildInfoByCharacterId(s, s.charID)
		if err != nil {
			s.logger.Fatal("Failed to get guild info for mail")
		}
		gm, err := GetGuildMembers(s, g.ID, false)
		if err != nil {
			s.logger.Fatal("Failed to get guild members for mail")
		}
		for i := 0; i < len(gm); i++ {
			_, err := s.server.db.Exec(query, s.charID, gm[i].CharID, pkt.Subject, pkt.Body, 0, 0, false)
			if err != nil {
				s.logger.Fatal("Failed to send mail")
			}
		}
	} else {
		_, err := s.server.db.Exec(query, s.charID, pkt.RecipientID, pkt.Subject, pkt.Body, pkt.ItemID, pkt.Quantity, false)
		if err != nil {
			s.logger.Fatal("Failed to send mail")
		}
	}

	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}
