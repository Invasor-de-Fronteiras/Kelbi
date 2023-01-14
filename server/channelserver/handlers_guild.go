package channelserver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"erupe-ce/common/byteframe"
	ps "erupe-ce/common/pascalstring"
	"erupe-ce/common/stringsupport"
	"erupe-ce/network/mhfpacket"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type FestivalColour string

const (
	FestivalColourNone FestivalColour = "none"
	FestivalColourRed  FestivalColour = "red"
	FestivalColourBlue FestivalColour = "blue"
)

var FestivalColourCodes = map[FestivalColour]uint8{
	FestivalColourBlue: 0x00,
	FestivalColourRed:  0x01,
	FestivalColourNone: 0xFF,
}

type GuildApplicationType string

const (
	GuildApplicationTypeApplied GuildApplicationType = "applied"
	GuildApplicationTypeInvited GuildApplicationType = "invited"
)

type Guild struct {
	ID             uint32         `db:"id"`
	Name           string         `db:"name"`
	MainMotto      uint8          `db:"main_motto"`
	SubMotto       uint8          `db:"sub_motto"`
	CreatedAt      time.Time      `db:"created_at"`
	MemberCount    uint16         `db:"member_count"`
	RankRP         uint32         `db:"rank_rp"`
	EventRP        uint32         `db:"event_rp"`
	Comment        string         `db:"comment"`
	PugiName1      string         `db:"pugi_name_1"`
	PugiName2      string         `db:"pugi_name_2"`
	PugiName3      string         `db:"pugi_name_3"`
	PugiOutfit1    uint8          `db:"pugi_outfit_1"`
	PugiOutfit2    uint8          `db:"pugi_outfit_2"`
	PugiOutfit3    uint8          `db:"pugi_outfit_3"`
	PugiOutfits    uint32         `db:"pugi_outfits"`
	Recruiting     bool           `db:"recruiting"`
	FestivalColour FestivalColour `db:"festival_colour"`
	Souls          uint32         `db:"souls"`
	Rank           uint16         `db:"rank"`
	AllianceID     uint32         `db:"alliance_id"`
	Icon           *GuildIcon     `db:"icon"`

	GuildLeader
}

type GuildLeader struct {
	LeaderCharID uint32 `db:"leader_id"`
	LeaderName   string `db:"leader_name"`
}

type GuildIconPart struct {
	Index    uint16
	ID       uint16
	Page     uint8
	Size     uint8
	Rotation uint8
	Red      uint8
	Green    uint8
	Blue     uint8
	PosX     uint16
	PosY     uint16
}

type GuildApplication struct {
	ID              int                  `db:"id"`
	GuildID         uint32               `db:"guild_id"`
	CharID          uint32               `db:"character_id"`
	ActorID         uint32               `db:"actor_id"`
	ApplicationType GuildApplicationType `db:"application_type"`
	CreatedAt       time.Time            `db:"created_at"`
}

type GuildIcon struct {
	Parts []GuildIconPart
}

func (gi *GuildIcon) Scan(val interface{}) (err error) {
	switch v := val.(type) {
	case []byte:
		err = json.Unmarshal(v, &gi)
	case string:
		err = json.Unmarshal([]byte(v), &gi)
	}

	return
}

func (gi *GuildIcon) Value() (valuer driver.Value, err error) {
	return json.Marshal(gi)
}

const guildInfoSelectQuery = `
SELECT
	g.id,
	g.name,
	g.rank_rp,
	g.event_rp,
	g.main_motto,
	g.sub_motto,
	g.created_at,
	g.leader_id,
	lc.name as leader_name,
	g.comment,
	COALESCE(g.pugi_name_1, '') AS pugi_name_1,
	COALESCE(g.pugi_name_2, '') AS pugi_name_2,
	COALESCE(g.pugi_name_3, '') AS pugi_name_3,
	g.pugi_outfit_1,
	g.pugi_outfit_2,
	g.pugi_outfit_3,
	g.pugi_outfits,
	recruiting,
	COALESCE((SELECT team FROM festa_registrations fr WHERE fr.guild_id = g.id), 'none') AS festival_colour,
	(SELECT SUM(souls) FROM guild_characters gc WHERE gc.guild_id = g.id) AS souls,
	CASE
		WHEN rank_rp <= 48 THEN rank_rp/24
		WHEN rank_rp <= 288 THEN rank_rp/48+1
		WHEN rank_rp <= 504 THEN rank_rp/72+3
		WHEN rank_rp <= 1080 THEN (rank_rp-24)/96+5
		WHEN rank_rp < 1200 THEN 16
		ELSE 17
	END rank,
	COALESCE((
		SELECT id FROM guild_alliances ga WHERE
	 	ga.parent_id = g.id OR
	 	ga.sub1_id = g.id OR
	 	ga.sub2_id = g.id
	), 0) AS alliance_id,
	icon,
	(SELECT count(1) FROM guild_characters gc WHERE gc.guild_id = g.id) AS member_count
	FROM guilds g
	JOIN guild_characters lgc ON lgc.character_id = leader_id
	JOIN characters lc on leader_id = lc.id
`

func (guild *Guild) Save(s *Session) error {
	_, err := s.Server.db.Exec(`
		UPDATE guilds SET main_motto=$2, sub_motto=$3, comment=$4, pugi_name_1=$5, pugi_name_2=$6, pugi_name_3=$7,
		pugi_outfit_1=$8, pugi_outfit_2=$9, pugi_outfit_3=$10, pugi_outfits=$11, icon=$12, leader_id=$13 WHERE id=$1
	`, guild.ID, guild.MainMotto, guild.SubMotto, guild.Comment, guild.PugiName1, guild.PugiName2, guild.PugiName3,
		guild.PugiOutfit1, guild.PugiOutfit2, guild.PugiOutfit3, guild.PugiOutfits, guild.Icon, guild.GuildLeader.LeaderCharID)

	if err != nil {
		s.logger.Error("failed to update guild data", zap.Error(err), zap.Uint32("guildID", guild.ID))
		return err
	}

	return nil
}

func (guild *Guild) CreateApplication(s *Session, charID uint32, applicationType GuildApplicationType, transaction *sql.Tx) error {

	query := `
		INSERT INTO guild_applications (guild_id, character_id, actor_id, application_type)
		VALUES ($1, $2, $3, $4)
	`

	var err error

	if transaction == nil {
		_, err = s.Server.db.Exec(query, guild.ID, charID, s.CharID, applicationType)
	} else {
		_, err = transaction.Exec(query, guild.ID, charID, s.CharID, applicationType)
	}

	if err != nil {
		s.logger.Error(
			"failed to add guild application",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		return err
	}

	return nil
}

func (guild *Guild) Disband(s *Session) error {
	transaction, err := s.Server.db.Begin()

	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return err
	}

	_, err = transaction.Exec("DELETE FROM guild_characters WHERE guild_id = $1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild characters", zap.Error(err), zap.Uint32("guildId", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec("DELETE FROM guilds WHERE id = $1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild", zap.Error(err), zap.Uint32("guildID", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec("DELETE FROM guild_alliances WHERE parent_id=$1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild alliance", zap.Error(err), zap.Uint32("guildID", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec("UPDATE guild_alliances SET sub1_id=sub2_id, sub2_id=NULL WHERE sub1_id=$1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild from alliance", zap.Error(err), zap.Uint32("guildID", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec("UPDATE guild_alliances SET sub2_id=NULL WHERE sub2_id=$1", guild.ID)

	if err != nil {
		s.logger.Error("failed to remove guild from alliance", zap.Error(err), zap.Uint32("guildID", guild.ID))
		rollbackTransaction(s, transaction)
		return err
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	s.logger.Info("Character disbanded guild", zap.Uint32("charID", s.CharID), zap.Uint32("guildID", guild.ID))

	return nil
}

func (guild *Guild) RemoveCharacter(s *Session, charID uint32) error {
	_, err := s.Server.db.Exec("DELETE FROM guild_characters WHERE character_id=$1", charID)

	if err != nil {
		s.logger.Error(
			"failed to remove character from guild",
			zap.Error(err),
			zap.Uint32("charID", charID),
			zap.Uint32("guildID", guild.ID),
		)

		return err
	}

	return nil
}

func (guild *Guild) AcceptApplication(s *Session, charID uint32) error {
	transaction, err := s.Server.db.Begin()

	if err != nil {
		s.logger.Error("failed to start db transaction", zap.Error(err))
		return err
	}

	_, err = transaction.Exec(`DELETE FROM guild_applications WHERE character_id = $1`, charID)

	if err != nil {
		s.logger.Error("failed to accept character's guild application", zap.Error(err))
		rollbackTransaction(s, transaction)
		return err
	}

	_, err = transaction.Exec(`
		INSERT INTO guild_characters (guild_id, character_id, order_index)
		VALUES ($1, $2, (SELECT MAX(order_index) + 1 FROM guild_characters WHERE guild_id = $1))
	`, guild.ID, charID)

	if err != nil {
		s.logger.Error(
			"failed to add applicant to guild",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		rollbackTransaction(s, transaction)
		return err
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit db transaction", zap.Error(err))
		rollbackTransaction(s, transaction)
		return err
	}

	return nil
}

// This is relying on the fact that invitation ID is also character ID right now
// if invitation ID changes, this will break.
func (guild *Guild) CancelInvitation(s *Session, charID uint32) error {
	_, err := s.Server.db.Exec(
		`DELETE FROM guild_applications WHERE character_id = $1 AND guild_id = $2 AND application_type = 'invited'`,
		charID, guild.ID,
	)

	if err != nil {
		s.logger.Error(
			"failed to cancel guild invitation",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		return err
	}

	return nil
}

func (guild *Guild) RejectApplication(s *Session, charID uint32) error {
	_, err := s.Server.db.Exec(
		`DELETE FROM guild_applications WHERE character_id = $1 AND guild_id = $2 AND application_type = 'applied'`,
		charID, guild.ID,
	)

	if err != nil {
		s.logger.Error(
			"failed to reject guild application",
			zap.Error(err),
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", charID),
		)
		return err
	}

	return nil
}

func (guild *Guild) ArrangeCharacters(s *Session, charIDs []uint32) error {
	transaction, err := s.Server.db.Begin()

	if err != nil {
		s.logger.Error("failed to start db transaction", zap.Error(err))
		return err
	}

	for i, id := range charIDs {
		_, err := transaction.Exec("UPDATE guild_characters SET order_index = $1 WHERE character_id = $2", 2+i, id)

		if err != nil {
			err = transaction.Rollback()

			if err != nil {
				s.logger.Error("failed to rollback db transaction", zap.Error(err))
			}

			return err
		}
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit db transaction", zap.Error(err))
		return err
	}

	return nil
}

func (guild *Guild) GetApplicationForCharID(s *Session, charID uint32, applicationType GuildApplicationType) (*GuildApplication, error) {
	application := &GuildApplication{}
	query := `
        SELECT * from guild_applications
        WHERE character_id = $1 AND guild_id = $2 AND application_type = $3
    `

	err := s.Server.db.Get(&application, query, charID, guild.ID, applicationType)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		s.logger.Error(
			"failed to retrieve guild application for character",
			zap.Error(err),
			zap.Uint32("charID", charID),
			zap.Uint32("guildID", guild.ID),
		)
		return nil, err
	}

	return application, nil
}

func (guild *Guild) HasApplicationForCharID(s *Session, charID uint32) (bool, error) {
	row := s.Server.db.QueryRowx(`
		SELECT 1 from guild_applications WHERE character_id = $1 AND guild_id = $2
	`, charID, guild.ID)

	num := 0

	err := row.Scan(&num)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		s.logger.Error(
			"failed to retrieve guild applications for character",
			zap.Error(err),
			zap.Uint32("charID", charID),
			zap.Uint32("guildID", guild.ID),
		)
		return false, err
	}

	return true, nil
}

func CreateGuild(s *Session, guildName string) (int32, error) {
	transaction, err := s.Server.db.Begin()

	if err != nil {
		s.logger.Error("failed to start db transaction", zap.Error(err))
		return 0, err
	}

	if err != nil {
		panic(err)
	}

	guildResult, err := transaction.Query(
		"INSERT INTO guilds (name, leader_id) VALUES ($1, $2) RETURNING id",
		guildName, s.CharID,
	)

	if err != nil {
		s.logger.Error("failed to create guild", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	var guildId int32

	guildResult.Next()

	err = guildResult.Scan(&guildId)

	if err != nil {
		s.logger.Error("failed to retrieve guild ID", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	err = guildResult.Close()

	if err != nil {
		s.logger.Error("failed to finalise query", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	_, err = transaction.Exec(`
		INSERT INTO guild_characters (guild_id, character_id)
		VALUES ($1, $2)
	`, guildId, s.CharID)

	if err != nil {
		s.logger.Error("failed to add character to guild", zap.Error(err))
		rollbackTransaction(s, transaction)
		return 0, err
	}

	err = transaction.Commit()

	if err != nil {
		s.logger.Error("failed to commit guild creation", zap.Error(err))
		return 0, err
	}

	return guildId, nil
}

func rollbackTransaction(s *Session, transaction *sql.Tx) {
	err := transaction.Rollback()

	if err != nil {
		s.logger.Error("failed to rollback transaction", zap.Error(err))
	}
}

func GetGuildInfoByID(s *Session, guildID uint32) (*Guild, error) {
	rows, err := s.Server.db.Queryx(fmt.Sprintf(`
		%s
		WHERE g.id = $1
		LIMIT 1
	`, guildInfoSelectQuery), guildID)

	if err != nil {
		s.logger.Error("failed to retrieve guild", zap.Error(err), zap.Uint32("guildID", guildID))
		return nil, err
	}

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow {
		return nil, nil
	}

	return buildGuildObjectFromDbResult(rows, err, s)
}

func GetGuildInfoByCharacterId(s *Session, charID uint32) (*Guild, error) {
	rows, err := s.Server.db.Queryx(fmt.Sprintf(`
		%s
		WHERE EXISTS(
				SELECT 1
				FROM guild_characters gc1
				WHERE gc1.character_id = $1
				  AND gc1.guild_id = g.id
			)
		   OR EXISTS(
				SELECT 1
				FROM guild_applications ga
				WHERE ga.character_id = $1
				  AND ga.guild_id = g.id
				  AND ga.application_type = 'applied'
			)
		LIMIT 1
	`, guildInfoSelectQuery), charID)

	if err != nil {
		s.logger.Error("failed to retrieve guild for character", zap.Error(err), zap.Uint32("charID", charID))
		return nil, err
	}

	defer rows.Close()

	hasRow := rows.Next()

	if !hasRow {
		return nil, nil
	}

	return buildGuildObjectFromDbResult(rows, err, s)
}

//nolint:staticcheck
func buildGuildObjectFromDbResult(result *sqlx.Rows, err error, s *Session) (*Guild, error) {
	guild := &Guild{}
	//nolint:staticcheck
	err = result.StructScan(guild)

	if err != nil {
		s.logger.Error("failed to retrieve guild data from database", zap.Error(err))
		return nil, err
	}

	return guild, nil
}

func handleMsgMhfCreateGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCreateGuild)

	guildId, err := CreateGuild(s, pkt.Name)

	if err != nil {
		bf := byteframe.NewByteFrame()

		// No reasoning behind these values other than they cause a 'failed to create'
		// style message, it's better than nothing for now.
		bf.WriteUint32(0x01010101)

		doAckSimpleFail(s, pkt.AckHandle, bf.Data())
		return
	}

	bf := byteframe.NewByteFrame()

	bf.WriteUint32(uint32(guildId))

	doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfOperateGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOperateGuild)

	guild, err := GetGuildInfoByID(s, pkt.GuildID)

	if err != nil {
		return
	}

	characterGuildInfo, err := GetCharacterGuildData(s, s.CharID)

	if err != nil {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	bf := byteframe.NewByteFrame()

	switch pkt.Action {
	case mhfpacket.OPERATE_GUILD_DISBAND:
		if guild.LeaderCharID != s.CharID {
			s.logger.Warn(fmt.Sprintf("character '%d' is attempting to manage guild '%d' without permission", s.CharID, guild.ID))
			return
		}

		err = guild.Disband(s)
		response := 0x01

		if err != nil {
			// All successful acks return 0x01, assuming 0x00 is failure
			response = 0x00
		}

		bf.WriteUint32(uint32(response))
	case mhfpacket.OPERATE_GUILD_RESIGN:
		guildMembers, err := GetGuildMembers(s, guild.ID, false)
		if err == nil {
			sort.Slice(guildMembers[:], func(i, j int) bool {
				return guildMembers[i].OrderIndex < guildMembers[j].OrderIndex
			})
			for i := 1; i < len(guildMembers); i++ {
				if !guildMembers[i].AvoidLeadership {
					guild.LeaderCharID = guildMembers[i].CharID
					guildMembers[0].OrderIndex = guildMembers[i].OrderIndex
					guildMembers[i].OrderIndex = 1
					// nolint:errcheck
					guildMembers[0].Save(s)
					// nolint:errcheck
					guildMembers[i].Save(s)
					bf.WriteUint32(guildMembers[i].CharID)
					break
				}
			}
			// nolint:errcheck
			guild.Save(s)
		}
	case mhfpacket.OPERATE_GUILD_APPLY:
		err = guild.CreateApplication(s, s.CharID, GuildApplicationTypeApplied, nil)

		if err == nil {
			bf.WriteUint32(guild.LeaderCharID)
		}
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
		return
	case mhfpacket.OPERATE_GUILD_LEAVE:
		var err error

		if characterGuildInfo.IsApplicant {
			err = guild.RejectApplication(s, s.CharID)
		} else {
			err = guild.RemoveCharacter(s, s.CharID)
		}

		response := 0x01
		if err != nil {
			// All successful acks return 0x01, assuming 0x00 is failure
			response = 0x00
		} else {
			mail := Mail{
				RecipientID:     s.CharID,
				Subject:         "Withdrawal",
				Body:            fmt.Sprintf("You have withdrawn from 「%s」.", guild.Name),
				IsSystemMessage: true,
			}
			// nolint:errcheck
			mail.Send(s, nil)
		}

		bf.WriteUint32(uint32(response))
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
		return
	case mhfpacket.OPERATE_GUILD_DONATE_RANK:
		bf.WriteBytes(handleDonateRP(s, uint16(pkt.Data1.ReadUint32()), guild, false))
	case mhfpacket.OPERATE_GUILD_SET_APPLICATION_DENY:
		s.Server.db.Exec("UPDATE guilds SET recruiting=false WHERE id=$1", guild.ID)
	case mhfpacket.OPERATE_GUILD_SET_APPLICATION_ALLOW:
		s.Server.db.Exec("UPDATE guilds SET recruiting=true WHERE id=$1", guild.ID)
	case mhfpacket.OPERATE_GUILD_SET_AVOID_LEADERSHIP_TRUE:
		handleAvoidLeadershipUpdate(s, pkt, true)
	case mhfpacket.OPERATE_GUILD_SET_AVOID_LEADERSHIP_FALSE:
		handleAvoidLeadershipUpdate(s, pkt, false)
	case mhfpacket.OPERATE_GUILD_UPDATE_COMMENT:
		if !characterGuildInfo.IsLeader && !characterGuildInfo.IsSubLeader() {
			doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
			return
		}
		guild.Comment = stringsupport.SJISToUTF8(pkt.Data2.ReadNullTerminatedBytes())
		// nolint:errcheck
		guild.Save(s)
	case mhfpacket.OPERATE_GUILD_UPDATE_MOTTO:
		if !characterGuildInfo.IsLeader && !characterGuildInfo.IsSubLeader() {
			doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
			return
		}
		_ = pkt.Data1.ReadUint16()
		guild.SubMotto = pkt.Data1.ReadUint8()
		guild.MainMotto = pkt.Data1.ReadUint8()
		guild.Save(s)
	case mhfpacket.OPERATE_GUILD_RENAME_PUGI_1:
		handleRenamePugi(s, pkt.Data2, guild, 1)
	case mhfpacket.OPERATE_GUILD_RENAME_PUGI_2:
		handleRenamePugi(s, pkt.Data2, guild, 2)
	case mhfpacket.OPERATE_GUILD_RENAME_PUGI_3:
		handleRenamePugi(s, pkt.Data2, guild, 3)
	case mhfpacket.OPERATE_GUILD_CHANGE_PUGI_1:
		handleChangePugi(s, uint8(pkt.Data1.ReadUint32()), guild, 1)
	case mhfpacket.OPERATE_GUILD_CHANGE_PUGI_2:
		handleChangePugi(s, uint8(pkt.Data1.ReadUint32()), guild, 2)
	case mhfpacket.OPERATE_GUILD_CHANGE_PUGI_3:
		handleChangePugi(s, uint8(pkt.Data1.ReadUint32()), guild, 3)
	case mhfpacket.OPERATE_GUILD_UNLOCK_OUTFIT:
		// TODO: This doesn't implement blocking, if someone unlocked the same outfit at the same time
		s.Server.db.Exec(`UPDATE guilds SET pugi_outfits=pugi_outfits+$1 WHERE id=$2`, int(math.Pow(float64(pkt.Data1.ReadUint32()), 2)), guild.ID)
	case mhfpacket.OPERATE_GUILD_DONATE_EVENT:
		bf.WriteBytes(handleDonateRP(s, uint16(pkt.Data1.ReadUint32()), guild, true))
	case mhfpacket.OPERATE_GUILD_EVENT_EXCHANGE:
		rp := uint16(pkt.Data1.ReadUint32())
		var balance uint32
		s.Server.db.QueryRow(`UPDATE guilds SET event_rp=event_rp-$1 WHERE id=$2 RETURNING event_rp`, rp, guild.ID).Scan(&balance)
		bf.WriteUint32(balance)
	default:
		panic(fmt.Sprintf("unhandled operate guild action '%d'", pkt.Action))
	}

	if len(bf.Data()) > 0 {
		doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
	} else {
		doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
	}
}

func handleRenamePugi(s *Session, bf *byteframe.ByteFrame, guild *Guild, num int) {
	name := stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
	switch num {
	case 1:
		guild.PugiName1 = name
	case 2:
		guild.PugiName2 = name
	default:
		guild.PugiName3 = name
	}
	// nolint:errcheck
	guild.Save(s)
}

func handleChangePugi(s *Session, outfit uint8, guild *Guild, num int) {
	switch num {
	case 1:
		guild.PugiOutfit1 = outfit
	case 2:
		guild.PugiOutfit2 = outfit
	case 3:
		guild.PugiOutfit3 = outfit
	}
	guild.Save(s)
}

func handleDonateRP(s *Session, amount uint16, guild *Guild, isEvent bool) []byte {
	bf := byteframe.NewByteFrame()
	bf.WriteUint32(0)
	saveData, err := GetCharacterSaveData(s, s.CharID)
	if err != nil {
		return bf.Data()
	}
	saveData.RP -= amount
	saveData.Save(s)
	updateSQL := "UPDATE guilds SET rank_rp = rank_rp + $1 WHERE id = $2"
	if isEvent {
		updateSQL = "UPDATE guilds SET event_rp = event_rp + $1 WHERE id = $2"
	}
	s.Server.db.Exec(updateSQL, amount, guild.ID)
	bf.Seek(0, 0)
	bf.WriteUint32(uint32(saveData.RP))
	return bf.Data()
}

func handleAvoidLeadershipUpdate(s *Session, pkt *mhfpacket.MsgMhfOperateGuild, avoidLeadership bool) {
	characterGuildData, err := GetCharacterGuildData(s, s.CharID)

	if err != nil {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	characterGuildData.AvoidLeadership = avoidLeadership

	err = characterGuildData.Save(s)

	if err != nil {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfOperateGuildMember(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOperateGuildMember)

	guild, err := GetGuildInfoByCharacterId(s, pkt.CharID)

	if err != nil || guild == nil {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	actorCharacter, err := GetCharacterGuildData(s, s.CharID)

	if err != nil || (!actorCharacter.IsSubLeader() && guild.LeaderCharID != s.CharID) {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	var mail Mail
	switch pkt.Action {
	case mhfpacket.OPERATE_GUILD_MEMBER_ACTION_ACCEPT:
		err = guild.AcceptApplication(s, pkt.CharID)
		mail = Mail{
			RecipientID:     pkt.CharID,
			Subject:         "Accepted!",
			Body:            fmt.Sprintf("Your application to join 「%s」 was accepted.", guild.Name),
			IsSystemMessage: true,
		}
	case mhfpacket.OPERATE_GUILD_MEMBER_ACTION_REJECT:
		err = guild.RejectApplication(s, pkt.CharID)
		mail = Mail{
			RecipientID:     pkt.CharID,
			Subject:         "Rejected",
			Body:            fmt.Sprintf("Your application to join 「%s」 was rejected.", guild.Name),
			IsSystemMessage: true,
		}
	case mhfpacket.OPERATE_GUILD_MEMBER_ACTION_KICK:
		err = guild.RemoveCharacter(s, pkt.CharID)
		mail = Mail{
			RecipientID:     pkt.CharID,
			Subject:         "Kicked",
			Body:            fmt.Sprintf("You were kicked from 「%s」.", guild.Name),
			IsSystemMessage: true,
		}
	default:
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		s.logger.Warn(fmt.Sprintf("unhandled operateGuildMember action '%d'", pkt.Action))
	}

	if err != nil {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
	} else {

		// nolint:errcheck // Error return value of `.` is not checked
		mail.Send(s, nil)

		session := s.Server.FindSessionByCharID(pkt.CharID)
		if session != nil {
			SendMailNotification(s, &mail, session)
		}

		doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
	}
}

func handleMsgMhfInfoGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfInfoGuild)

	var guild *Guild
	var err error

	if pkt.GuildID > 0 {
		guild, err = GetGuildInfoByID(s, pkt.GuildID)
	} else {
		guild, err = GetGuildInfoByCharacterId(s, s.CharID)
	}

	if err == nil && guild != nil {
		s.PrevGuildID = guild.ID

		guildName := stringsupport.UTF8ToSJIS(guild.Name)
		guildComment := stringsupport.UTF8ToSJIS(guild.Comment)
		guildLeaderName := stringsupport.UTF8ToSJIS(guild.LeaderName)

		characterGuildData, err := GetCharacterGuildData(s, s.CharID)
		characterJoinedAt := uint32(0xFFFFFFFF)

		if characterGuildData != nil && characterGuildData.JoinedAt != nil {
			characterJoinedAt = uint32(characterGuildData.JoinedAt.Unix())
		}

		if err != nil {
			resp := byteframe.NewByteFrame()
			resp.WriteUint32(0) // Count
			resp.WriteUint8(0)  // Unk, read if count == 0.

			doAckBufSucceed(s, pkt.AckHandle, resp.Data())
			return
		}

		bf := byteframe.NewByteFrame()

		bf.WriteUint32(guild.ID)
		bf.WriteUint32(guild.LeaderCharID)
		bf.WriteUint16(guild.Rank)
		bf.WriteUint16(guild.MemberCount)

		bf.WriteUint8(guild.MainMotto)
		bf.WriteUint8(guild.SubMotto)

		// Unk appears to be static
		bf.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

		bf.WriteBool(!guild.Recruiting)

		if characterGuildData == nil || characterGuildData.IsApplicant {
			bf.WriteUint16(0x00)
		} else if guild.LeaderCharID == s.CharID {
			bf.WriteUint16(0x01)
		} else {
			bf.WriteUint16(0x02)
		}

		bf.WriteUint32(uint32(guild.CreatedAt.Unix()))
		bf.WriteUint32(characterJoinedAt)
		bf.WriteUint8(uint8(len(guildName)))
		bf.WriteUint8(uint8(len(guildComment)))
		bf.WriteUint8(uint8(5)) // Length of unknown string below
		bf.WriteUint8(uint8(len(guildLeaderName)))
		bf.WriteBytes(guildName)
		bf.WriteBytes(guildComment)
		bf.WriteUint8(FestivalColourCodes[guild.FestivalColour])
		bf.WriteUint32(guild.RankRP)
		bf.WriteBytes(guildLeaderName)
		bf.WriteBytes([]byte{0x00, 0x00, 0x00, 0x00}) // Unk
		bf.WriteBool(false)                           // isReturnGuild
		bf.WriteBool(false)                           // earnedSpecialHall
		bf.WriteBytes([]byte{0x02, 0x02})             // Unk
		bf.WriteUint32(guild.EventRP)
		ps.Uint8(bf, guild.PugiName1, true)
		ps.Uint8(bf, guild.PugiName2, true)
		ps.Uint8(bf, guild.PugiName3, true)
		bf.WriteUint8(guild.PugiOutfit1)
		bf.WriteUint8(guild.PugiOutfit2)
		bf.WriteUint8(guild.PugiOutfit3)
		bf.WriteUint8(guild.PugiOutfit1)
		bf.WriteUint8(guild.PugiOutfit2)
		bf.WriteUint8(guild.PugiOutfit3)
		bf.WriteUint32(guild.PugiOutfits)

		// Unk flags
		bf.WriteUint8(0x3C) // also seen as 0x32 on JP and 0x64 on TW

		bf.WriteBytes([]byte{
			0x00, 0x00, 0xD6, 0xD8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		})

		if guild.AllianceID > 0 {
			alliance, err := GetAllianceData(s, guild.AllianceID)
			if err != nil {
				bf.WriteUint32(0) // Error, no alliance
			} else {
				bf.WriteUint32(alliance.ID)
				bf.WriteUint32(uint32(alliance.CreatedAt.Unix()))
				bf.WriteUint16(alliance.TotalMembers)
				bf.WriteUint16(0) // Unk0
				ps.Uint16(bf, alliance.Name, true)
				if alliance.SubGuild1ID > 0 {
					if alliance.SubGuild2ID > 0 {
						bf.WriteUint8(3)
					} else {
						bf.WriteUint8(2)
					}
				} else {
					bf.WriteUint8(1)
				}
				bf.WriteUint32(alliance.ParentGuildID)
				bf.WriteUint32(0) // Unk1
				if alliance.ParentGuildID == guild.ID {
					bf.WriteUint16(1)
				} else {
					bf.WriteUint16(0)
				}
				bf.WriteUint16(alliance.ParentGuild.Rank)
				bf.WriteUint16(alliance.ParentGuild.MemberCount)
				ps.Uint16(bf, alliance.ParentGuild.Name, true)
				ps.Uint16(bf, alliance.ParentGuild.LeaderName, true)
				if alliance.SubGuild1ID > 0 {
					bf.WriteUint32(alliance.SubGuild1ID)
					bf.WriteUint32(0) // Unk1
					if alliance.SubGuild1ID == guild.ID {
						bf.WriteUint16(1)
					} else {
						bf.WriteUint16(0)
					}
					bf.WriteUint16(alliance.SubGuild1.Rank)
					bf.WriteUint16(alliance.SubGuild1.MemberCount)
					ps.Uint16(bf, alliance.SubGuild1.Name, true)
					ps.Uint16(bf, alliance.SubGuild1.LeaderName, true)
				}
				if alliance.SubGuild2ID > 0 {
					bf.WriteUint32(alliance.SubGuild2ID)
					bf.WriteUint32(0) // Unk1
					if alliance.SubGuild2ID == guild.ID {
						bf.WriteUint16(1)
					} else {
						bf.WriteUint16(0)
					}
					bf.WriteUint16(alliance.SubGuild2.Rank)
					bf.WriteUint16(alliance.SubGuild2.MemberCount)
					ps.Uint16(bf, alliance.SubGuild2.Name, true)
					ps.Uint16(bf, alliance.SubGuild2.LeaderName, true)
				}
			}
		} else {
			bf.WriteUint32(0) // No alliance
		}

		applicants, err := GetGuildMembers(s, guild.ID, true)
		if err != nil || (characterGuildData != nil && !characterGuildData.CanRecruit()) {
			bf.WriteUint16(0)
		} else {
			bf.WriteUint16(uint16(len(applicants)))
			for _, applicant := range applicants {
				bf.WriteUint32(applicant.CharID)
				bf.WriteUint16(0)
				bf.WriteUint16(0)
				bf.WriteUint16(applicant.HRP)
				bf.WriteUint16(applicant.GR)
				ps.Uint8(bf, applicant.Name, true)
			}
		}

		bf.WriteUint16(0x0000) // lenAllianceApplications

		/*
			alliance application format
			uint16 numapplicants (above)

			uint32 guild id
			uint32 guild leader id (for mail)
			uint32 unk (always null in pcap)
			uint16 unk (always 0001 in pcap)
			uint16 len guild name
			string nullterm guild name
			uint16 len guild leader name
			string nullterm guild leader name
		*/

		if guild.Icon != nil {
			bf.WriteUint8(uint8(len(guild.Icon.Parts)))

			for _, p := range guild.Icon.Parts {
				bf.WriteUint16(p.Index)
				bf.WriteUint16(p.ID)
				bf.WriteUint8(p.Page)
				bf.WriteUint8(p.Size)
				bf.WriteUint8(p.Rotation)
				bf.WriteUint8(p.Red)
				bf.WriteUint8(p.Green)
				bf.WriteUint8(p.Blue)
				bf.WriteUint16(p.PosX)
				bf.WriteUint16(p.PosY)
			}
		} else {
			bf.WriteUint8(0x00)
		}
		bf.WriteUint8(0) // Unk

		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
	} else {
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 5))
	}
}

func handleMsgMhfEnumerateGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuild)

	var guilds []*Guild
	var alliances []*GuildAlliance
	var rows *sqlx.Rows
	var err error
	bf := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload)

	switch pkt.Type {
	case mhfpacket.ENUMERATE_GUILD_TYPE_GUILD_NAME:
		bf.ReadBytes(8)
		searchTerm := fmt.Sprintf(`%%%s%%`, stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes()))
		rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s WHERE g.name ILIKE $1 OFFSET $2 LIMIT 11`, guildInfoSelectQuery), searchTerm, pkt.Page*10)
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_LEADER_NAME:
		bf.ReadBytes(8)
		searchTerm := fmt.Sprintf(`%%%s%%`, stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes()))
		rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s WHERE lc.name ILIKE $1 OFFSET $2 LIMIT 11`, guildInfoSelectQuery), searchTerm, pkt.Page*10)
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_LEADER_ID:
		ID := bf.ReadUint32()
		rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s WHERE leader_id = $1`, guildInfoSelectQuery), ID)
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_ORDER_MEMBERS:
		if pkt.Sorting {
			rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s ORDER BY member_count DESC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		} else {
			rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s ORDER BY member_count ASC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		}
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_ORDER_REGISTRATION:
		if pkt.Sorting {
			rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s ORDER BY id ASC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		} else {
			rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s ORDER BY id DESC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		}
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_ORDER_RANK:
		if pkt.Sorting {
			rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s ORDER BY rank_rp DESC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		} else {
			rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s ORDER BY rank_rp ASC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		}
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_MOTTO:
		mainMotto := bf.ReadUint16()
		subMotto := bf.ReadUint16()
		rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s WHERE main_motto = $1 AND sub_motto = $2 OFFSET $3 LIMIT 11`, guildInfoSelectQuery), mainMotto, subMotto, pkt.Page*10)
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	case mhfpacket.ENUMERATE_GUILD_TYPE_RECRUITING:
		// Assume the player wants the newest guilds with open recruitment
		rows, err = s.Server.db.Queryx(fmt.Sprintf(`%s WHERE recruiting=true ORDER BY id DESC OFFSET $1 LIMIT 11`, guildInfoSelectQuery), pkt.Page*10)
		if err == nil {
			for rows.Next() {
				guild, _ := buildGuildObjectFromDbResult(rows, err, s)
				guilds = append(guilds, guild)
			}
		}
	}

	if pkt.Type > 8 {
		var tempAlliances []*GuildAlliance
		rows, err = s.Server.db.Queryx(allianceInfoSelectQuery)
		if err == nil {
			for rows.Next() {
				alliance, _ := buildAllianceObjectFromDbResult(rows, err, s)
				tempAlliances = append(tempAlliances, alliance)
			}
		}
		switch pkt.Type {
		case mhfpacket.ENUMERATE_ALLIANCE_TYPE_ALLIANCE_NAME:
			bf.ReadBytes(8)
			searchTerm := stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
			for _, alliance := range tempAlliances {
				if strings.Contains(alliance.Name, searchTerm) {
					alliances = append(alliances, alliance)
				}
			}
		case mhfpacket.ENUMERATE_ALLIANCE_TYPE_LEADER_NAME:
			bf.ReadBytes(8)
			searchTerm := stringsupport.SJISToUTF8(bf.ReadNullTerminatedBytes())
			for _, alliance := range tempAlliances {
				if strings.Contains(alliance.ParentGuild.LeaderName, searchTerm) {
					alliances = append(alliances, alliance)
				}
			}
		case mhfpacket.ENUMERATE_ALLIANCE_TYPE_LEADER_ID:
			ID := bf.ReadUint32()
			for _, alliance := range tempAlliances {
				if alliance.ParentGuild.LeaderCharID == ID {
					alliances = append(alliances, alliance)
				}
			}
		case mhfpacket.ENUMERATE_ALLIANCE_TYPE_ORDER_MEMBERS:
			if pkt.Sorting {
				sort.Slice(tempAlliances, func(i, j int) bool {
					return tempAlliances[i].TotalMembers > tempAlliances[j].TotalMembers
				})
			} else {
				sort.Slice(tempAlliances, func(i, j int) bool {
					return tempAlliances[i].TotalMembers < tempAlliances[j].TotalMembers
				})
			}
			alliances = tempAlliances
		case mhfpacket.ENUMERATE_ALLIANCE_TYPE_ORDER_REGISTRATION:
			if pkt.Sorting {
				sort.Slice(tempAlliances, func(i, j int) bool {
					return tempAlliances[i].CreatedAt.Unix() > tempAlliances[j].CreatedAt.Unix()
				})
			} else {
				sort.Slice(tempAlliances, func(i, j int) bool {
					return tempAlliances[i].CreatedAt.Unix() < tempAlliances[j].CreatedAt.Unix()
				})
			}
			alliances = tempAlliances
		}
	}

	if err != nil || (guilds == nil && alliances == nil) {
		stubEnumerateNoResults(s, pkt.AckHandle)
		return
	}

	bf = byteframe.NewByteFrame()

	if pkt.Type > 8 {
		hasNextPage := false
		if len(alliances) > 10 {
			hasNextPage = true
			alliances = alliances[:10]
		}
		bf.WriteUint16(uint16(len(alliances)))
		bf.WriteBool(hasNextPage)
		for _, alliance := range alliances {
			bf.WriteUint32(alliance.ID)
			bf.WriteUint32(alliance.ParentGuild.LeaderCharID)
			bf.WriteUint16(alliance.TotalMembers)
			bf.WriteUint16(0x0000)
			if alliance.SubGuild1ID == 0 && alliance.SubGuild2ID == 0 {
				bf.WriteUint16(1)
			} else if alliance.SubGuild1ID > 0 && alliance.SubGuild2ID == 0 || alliance.SubGuild1ID == 0 && alliance.SubGuild2ID > 0 {
				bf.WriteUint16(2)
			} else {
				bf.WriteUint16(3)
			}
			bf.WriteUint32(uint32(alliance.CreatedAt.Unix()))
			ps.Uint8(bf, alliance.Name, true)
			ps.Uint8(bf, alliance.ParentGuild.LeaderName, true)
			bf.WriteUint8(0x01) // Unk
			bf.WriteBool(true)  // TODO: Enable GuildAlliance applications
		}
	} else {
		hasNextPage := false
		if len(guilds) > 10 {
			hasNextPage = true
			guilds = guilds[:10]
		}
		bf.WriteUint16(uint16(len(guilds)))
		bf.WriteBool(hasNextPage)
		for _, guild := range guilds {
			bf.WriteUint32(guild.ID)
			bf.WriteUint32(guild.LeaderCharID)
			bf.WriteUint16(guild.MemberCount)
			bf.WriteUint16(0x0000)     // Unk
			bf.WriteUint16(guild.Rank) // OR guilds in alliance
			bf.WriteUint32(uint32(guild.CreatedAt.Unix()))
			ps.Uint8(bf, guild.Name, true)
			ps.Uint8(bf, guild.LeaderName, true)
			bf.WriteUint8(0x01) // Unk
			bf.WriteBool(!guild.Recruiting)
		}
	}

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfArrangeGuildMember(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfArrangeGuildMember)

	guild, err := GetGuildInfoByID(s, pkt.GuildID)

	if err != nil {
		s.logger.Error(
			"failed to respond to ArrangeGuildMember message",
			zap.Uint32("charID", s.CharID),
		)
		return
	}

	if guild.LeaderCharID != s.CharID {
		s.logger.Error("non leader attempting to rearrange guild members!",
			zap.Uint32("charID", s.CharID),
			zap.Uint32("guildID", guild.ID),
		)
		return
	}

	err = guild.ArrangeCharacters(s, pkt.CharIDs)

	if err != nil {
		s.logger.Error(
			"failed to respond to ArrangeGuildMember message",
			zap.Uint32("charID", s.CharID),
			zap.Uint32("guildID", guild.ID),
		)
		return
	}

	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfEnumerateGuildMember(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuildMember)

	var guild *Guild
	var err error

	if pkt.GuildID > 0 {
		guild, err = GetGuildInfoByID(s, pkt.GuildID)
	} else {
		guild, err = GetGuildInfoByCharacterId(s, s.CharID)
	}

	if guild == nil && s.PrevGuildID > 0 {
		guild, err = GetGuildInfoByID(s, s.PrevGuildID)
	}

	if guild != nil {
		isApplicant, _ := guild.HasApplicationForCharID(s, s.CharID)
		if isApplicant {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
	}

	if guild == nil && s.PrevGuildID > 0 {
		guild, err = GetGuildInfoByID(s, s.PrevGuildID)
	}

	if guild != nil {
		isApplicant, _ := guild.HasApplicationForCharID(s, s.CharID)
		if isApplicant {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
	}

	if guild == nil && s.PrevGuildID > 0 {
		guild, err = GetGuildInfoByID(s, s.PrevGuildID)
	}

	if err != nil {
		s.logger.Warn("failed to retrieve guild sending no result message")
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 2))
		return
	} else if guild == nil {
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 2))
		return
	}

	guildMembers, err := GetGuildMembers(s, guild.ID, false)

	if err != nil {
		s.logger.Error("failed to retrieve guild")
		return
	}

	alliance, err := GetAllianceData(s, guild.AllianceID)
	if err != nil {
		s.logger.Error("Failed to get alliance data")
		return
	}

	bf := byteframe.NewByteFrame()

	bf.WriteUint16(guild.MemberCount)

	sort.Slice(guildMembers[:], func(i, j int) bool {
		return guildMembers[i].OrderIndex < guildMembers[j].OrderIndex
	})

	for _, member := range guildMembers {
		bf.WriteUint32(member.CharID)
		bf.WriteUint16(member.HRP)
		bf.WriteUint16(member.GR)
		bf.WriteUint16(member.WeaponID)
		if member.WeaponType == 1 || member.WeaponType == 5 || member.WeaponType == 10 { // If weapon is ranged
			bf.WriteUint16(0x0700)
		} else {
			bf.WriteUint16(0x0600)
		}
		bf.WriteUint8(member.OrderIndex)
		bf.WriteBool(member.AvoidLeadership)
		ps.Uint8(bf, member.Name, true)
	}

	for _, member := range guildMembers {
		bf.WriteUint32(member.LastLogin)
	}

	if guild.AllianceID > 0 {
		bf.WriteUint16(alliance.TotalMembers - guild.MemberCount)
		if guild.ID != alliance.ParentGuildID {
			mems, err := GetGuildMembers(s, alliance.ParentGuildID, false)
			if err != nil {
				panic(err)
			}
			for _, m := range mems {
				bf.WriteUint32(m.CharID)
			}
		}
		if guild.ID != alliance.SubGuild1ID {
			mems, err := GetGuildMembers(s, alliance.SubGuild1ID, false)
			if err != nil {
				panic(err)
			}
			for _, m := range mems {
				bf.WriteUint32(m.CharID)
			}
		}
		if guild.ID != alliance.SubGuild2ID {
			mems, err := GetGuildMembers(s, alliance.SubGuild2ID, false)
			if err != nil {
				panic(err)
			}
			for _, m := range mems {
				bf.WriteUint32(m.CharID)
			}
		}
	} else {
		bf.WriteUint16(0)
	}

	for range guildMembers {
		bf.WriteUint32(0x00) // Unk
	}

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetGuildManageRight(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildManageRight)

	guild, err := GetGuildInfoByCharacterId(s, s.CharID)

	if guild == nil && s.PrevGuildID != 0 {
		guild, err = GetGuildInfoByID(s, s.PrevGuildID)
		s.PrevGuildID = 0
		if guild == nil || err != nil {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
	}

	if guild == nil && s.PrevGuildID != 0 {
		guild, err = GetGuildInfoByID(s, s.PrevGuildID)
		s.PrevGuildID = 0
		if guild == nil || err != nil {
			doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
	}

	if err != nil {
		s.logger.Warn("failed to respond to manage rights message")
		return
	} else if guild == nil {
		bf := byteframe.NewByteFrame()
		bf.WriteUint16(0x00) // Unk
		bf.WriteUint16(0x00) // Member count

		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
		return
	}

	bf := byteframe.NewByteFrame()

	bf.WriteUint16(0x00) // Unk
	bf.WriteUint16(guild.MemberCount)

	members, _ := GetGuildMembers(s, guild.ID, false)

	for _, member := range members {
		bf.WriteUint32(member.CharID)
		bf.WriteBool(member.Recruiter)
		bf.WriteBytes(make([]byte, 3))
	}

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetUdGuildMapInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdGuildMapInfo)

	data, _ := hex.DecodeString("00050000013600000137013500000000E2DF000000000204000000640100000001019901350000E2DF0000000001000000044C000000000001FE01FF00000000000000000D0000001036000000000001FC01FD00000000000000000B0000000F0A000000000001FB01FC00000000000000000A0000000E740000000000019B013700000000000000000F00000011620000000000019601FB0000000000000000090000000DDE0000000000013700D400000000000000001000000011F80000000000013201960000000000000000080000000D48000000000000D40070000000000000000011000000128E000000000000CE01320000000000000000070000000CB200000000000070006F00000000000000001200000013240000000000006F006E00000000000000001300000013BA0000000000006E006D00000000000000001400000014500000000000006D0000000000000000000015020000157C0000000000006A00CE0000000000000000060000000C1C00000000000069006A0000000000000000050000000B860000000000006800690000000000000000040000000AF00000000000006700680000000000000000030000000A5A00000000000066006700000000000000000200000009C4000000000001FD01FE01990000000000000C0300000FA00000000000006500660000000000000000010100000000000000000001FF019B00000000000000000E00000010CC0000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013700000138013200000000E2E0000000000204000000640100000002019701320000E2E00000000001000000044C00000000000193012E00000000000000000E000000319C0000000000013701360000000000000000060000001EDC00000000000136019900000000000000000700000021340000000000012E012D00000000000000000F00000033F40000000000019801FB01970000000000000903000025E4000000000001F9019400000000000000000C0000002CEC0000000000012D00C9000000000000000010000000364C000000000000D401370000000000000000050000001C84000000000000C9006600000000000000001100000036B00000000000007000D40000000000000000040000001A2C000000000001FA01F900000000000000000B0000002A940000000000006F007000000000000000000300000017D40000000000006E006F000000000000000002000000157C0000000000006D006E00000000000000000101000000000000000000006900000000000000000000150200004362000000000001FB01FA00000000000000000A000000283C0000000000006800690000000000000000140000003B6000000000000067006800000000000000001300000039D0000000000001990198000000000000000008000000238C000000000000660067000000000000000012000000390800000000000194019300000000000000000D0000002F44000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001380000013901FF00000000E2E10000000001040000044C0000000003019B013601FF0000000000000D0300003BC4000000000001FA000000000000000000001502000055F0000000000001F901FA00000000000000001400000051A400000000000199013400000000000000000F00000042040000000000019501F90000000000000000130000004E8400000000000138019B00000000000000000C00000038A400000000000136019900000000000000000E0000003EE4000000000001340133000000000000000010000000452400000000000133013200000000000000001100000048440000000000013201950000000000000000120000004B64000000000000D4013800000000000000000B0000003584000000000000D1006E0000000000000000070000002904000000000000CD006A0000000000000000030000001C840000000000007000D400000000000000000A00000032640000000000006F00700000000000000000090000002F440000000000006E006F0000000000000000080000002C240000000000006C00D100000000000000000600000025E40000000000006B006C00000000000000000500000022C40000000000006A006B0000000000000000040000001FA40000000000006800CD000000000000000002000000196400000000000067006800000000000000000101000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001390000013A006500000000E2E200000000020400000064010000000500C900650000E2E20000000001000000044C00000000000133019700000000000000000C040000445C0000000004013700D20000000000000000130000005FB4000000000000CA00CB00C9000000000000060300002CEC000000000001FF019B00000000000000001100000057E4000000000001FE01FF00000000000000001000000053FC000000000001FD01FE00000000000000000F0000005014000000000001FA01F90000000000000000010100000000000000000001F901940000000000000000020000001D4C0000000000019B01370000000000000000120000005BCC0000000000019801FD00000000000000000E0000004C2C00000000000197019800000000000000000D0000004844000000000001940193000000000000000003000000213400000000000193012E000000000000000004000000251C0000000000012E00CA0000000000000000050000002904000000000000D2006E000000000000000014000000639C000000000000CF013300000000000000000B0000004074000000000000CB006800000000000000000700000030D40000000000006E000000000000000000001502000075300000000000006A00CF00000000000000000A0000003C8C00000000000069006A00000000000000000900000038A400000000000068006900000000000000000800000034BC0000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000013A0000013701F800000000E2E30000000001040000044C000000000601FD01FC00000000000000000700000034BC000000000001FC01FB00000000000000000800000038A4000000000001FB01FA0000000000000000090000003C8C000000000001FA01F900000000000000000A00000040740000000000019B019A0000000000000000050000002CEC0000000000019A01FD00000000000000000600000030D400000000000195013200000000000000000C000000484400000000000138019B00000000000000000400000029040000000000013300CF00000000000000000E000000501400000000000132013300000000000000000D0000004C2C000000000000D40138000000000000000003000000251C000000000000D300D40000000000000000020000002134000000000000CF006A00000000000000000F00000053FC000000000000CD00CC0000000000000000110000005BCC000000000000CC00CB0000000000000000120000005FB4000000000000CB00CA000000000000000013000000639C000000000000CA00C90000000000000000140000006784000000000000C90000000000000000000015020000FDE80000000000006E00D30000000000000000010100000000000000000001F9019501F80000000000000B030000445C0000000000006A00CD00000000000000001000000057E400000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009000000010738AD00010001000102000000011A000007D00002006302000000020738AD00020001000102000000021A000007D00002006302000000031A000003E80001006301000000041A000003E80001006301000000050738AD00020001000102000000051A000007D00002006302000000061A000003E800010063010100000136000117000000000000044C019901350000E2DF000000000100000000000000000064013500000000E2DF00000000020400000000000000000000650066000000000000000001010100000000000009C40066006700000000000000000200000000000000000A5A0067006800000000000000000300000000000000000AF00068006900000000000000000400000000000000000B860069006A00000000000000000500000000000000000C1C006A00CE00000000000000000600000000000000000CB200CE013200000000000000000700000000000000000D480132019600000000000000000800000000000000000DDE019601FB00000000000000000900000000000000000E7401FB01FC00000000000000000A00000000000000000F0A01FC01FD00000000000000000B00000000000000000FA001FD01FE01990000000000000C0300000000000000103601FE01FF00000000000000000D000000000000000010CC01FF019B00000000000000000E00000000000000001162019B013700000000000000000F000000000000000011F8013700D40000000000000000100000000000000000128E00D40070000000000000000011000000000000000013240070006F000000000000000012000000000000000013BA006F006E00000000000000001300000000000000001450006E006D0000000000000000140000000000000000157C006D000000000000000000001502000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfGetGuildTargetMemberNum(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildTargetMemberNum)

	var guild *Guild
	var err error

	if pkt.GuildID == 0x0 {
		guild, err = GetGuildInfoByCharacterId(s, s.CharID)
	} else {
		guild, err = GetGuildInfoByID(s, pkt.GuildID)
	}

	if err != nil || guild == nil {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x02})
		return
	}

	bf := byteframe.NewByteFrame()

	bf.WriteUint16(0x0)
	bf.WriteUint16(guild.MemberCount - 1)

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfEnumerateGuildItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuildItem)
	var boxContents []byte
	bf := byteframe.NewByteFrame()
	err := s.Server.db.QueryRow("SELECT item_box FROM guilds WHERE id = $1", int(pkt.GuildId)).Scan(&boxContents)
	if err != nil {
		s.logger.Fatal("Failed to get guild item box contents from db", zap.Error(err))
	} else {
		if len(boxContents) == 0 {
			bf.WriteUint32(0x00)
		} else {
			amount := len(boxContents) / 4
			bf.WriteUint16(uint16(amount))
			bf.WriteUint32(0x00)
			bf.WriteUint16(0x00)
			for i := 0; i < amount; i++ {
				bf.WriteUint32(binary.BigEndian.Uint32(boxContents[i*4 : i*4+4]))
				if i+1 != amount {
					bf.WriteUint64(0x00)
				}
			}
		}
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

type Item struct {
	ItemId uint16
	Amount uint16
}

func handleMsgMhfUpdateGuildItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateGuildItem)

	// Get item cache from DB
	var boxContents []byte
	var oldItems []Item
	err := s.Server.db.QueryRow("SELECT item_box FROM guilds WHERE id = $1", int(pkt.GuildId)).Scan(&boxContents)
	if err != nil {
		s.logger.Fatal("Failed to get guild item box contents from db", zap.Error(err))
	} else {
		amount := len(boxContents) / 4
		oldItems = make([]Item, amount)
		for i := 0; i < amount; i++ {
			oldItems[i].ItemId = binary.BigEndian.Uint16(boxContents[i*4 : i*4+2])
			oldItems[i].Amount = binary.BigEndian.Uint16(boxContents[i*4+2 : i*4+4])
		}
	}

	// Update item stacks
	newItems := make([]Item, len(oldItems))
	copy(newItems, oldItems)
	for i := 0; i < int(pkt.Amount); i++ {
		for j := 0; j <= len(oldItems); j++ {
			if j == len(oldItems) {
				var newItem Item
				newItem.ItemId = pkt.Items[i].ItemId
				newItem.Amount = pkt.Items[i].Amount
				newItems = append(newItems, newItem)
				break
			}
			if pkt.Items[i].ItemId == oldItems[j].ItemId {
				newItems[j].Amount = pkt.Items[i].Amount
				break
			}
		}
	}

	// Delete empty item stacks
	for i := len(newItems) - 1; i >= 0; i-- {
		if int(newItems[i].Amount) == 0 {
			copy(newItems[i:], newItems[i+1:])
			newItems[len(newItems)-1] = make([]Item, 1)[0]
			newItems = newItems[:len(newItems)-1]
		}
	}

	// Create new item cache
	bf := byteframe.NewByteFrame()
	for i := 0; i < len(newItems); i++ {
		bf.WriteUint16(newItems[i].ItemId)
		bf.WriteUint16(newItems[i].Amount)
	}

	// Upload new item cache
	_, err = s.Server.db.Exec("UPDATE guilds SET item_box = $1 WHERE id = $2", bf.Data(), int(pkt.GuildId))
	if err != nil {
		s.logger.Fatal("Failed to update guild item box contents in db", zap.Error(err))
	}

	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfUpdateGuildIcon(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateGuildIcon)

	guild, err := GetGuildInfoByID(s, pkt.GuildID)

	if err != nil {
		panic(err)
	}

	characterInfo, err := GetCharacterGuildData(s, s.CharID)

	if err != nil {
		panic(err)
	}

	if !characterInfo.IsSubLeader() && !characterInfo.IsLeader {
		s.logger.Warn(
			"character without leadership attempting to update guild icon",
			zap.Uint32("guildID", guild.ID),
			zap.Uint32("charID", s.CharID),
		)
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	icon := &GuildIcon{}

	icon.Parts = make([]GuildIconPart, pkt.PartCount)

	for i, p := range pkt.IconParts {
		icon.Parts[i] = GuildIconPart{
			Index:    p.Index,
			ID:       p.ID,
			Page:     p.Page,
			Size:     p.Size,
			Rotation: p.Rotation,
			Red:      p.Red,
			Green:    p.Green,
			Blue:     p.Blue,
			PosX:     p.PosX,
			PosY:     p.PosY,
		}
	}

	guild.Icon = icon

	err = guild.Save(s)

	if err != nil {
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}

	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfReadGuildcard(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReadGuildcard)

	resp := byteframe.NewByteFrame()
	resp.WriteUint32(0)
	resp.WriteUint32(0)
	resp.WriteUint32(0)
	resp.WriteUint32(0)
	resp.WriteUint32(0)
	resp.WriteUint32(0)
	resp.WriteUint32(0)
	resp.WriteUint32(0)

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetGuildMissionList(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildMissionList)

	decoded, err := hex.DecodeString("000694610000023E000112990023000100000200015DDD232100069462000002F30000005F000C000200000300025DDD232100069463000002EA0000005F0006000100000100015DDD23210006946400000245000000530010000200000400025DDD232100069465000002B60001129B0019000100000200015DDD232100069466000003DC0000001B0010000100000600015DDD232100069467000002DA000112A00019000100000400015DDD232100069468000002A800010DEF0032000200000200025DDD2321000694690000045500000022003C000200000600025DDD23210006946A00000080000122D90046000200000300025DDD23210006946B000001960000003B000A000100000100015DDD23210006946C0000049200000046005A000300000600035DDD23210006946D000000A4000000260018000200000600025DDD23210006946E0000017A00010DE40096000300000100035DDD23210006946F000001BE0000005E0014000200000400025DDD2355000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	if err != nil {
		panic(err)
	}

	doAckBufSucceed(s, pkt.AckHandle, decoded)
}

func handleMsgMhfGetGuildMissionRecord(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildMissionRecord)

	// No guild mission records = 0x190 empty bytes
	doAckBufSucceed(s, pkt.AckHandle, make([]byte, 0x190))
}

func handleMsgMhfAddGuildMissionCount(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAddGuildMissionCount)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfSetGuildMissionTarget(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetGuildMissionTarget)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfCancelGuildMissionTarget(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCancelGuildMissionTarget)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

type GuildMeal struct {
	ID      uint32 `db:"id"`
	MealID  uint32 `db:"meal_id"`
	Level   uint32 `db:"level"`
	Expires uint32 `db:"expires"`
}

func handleMsgMhfLoadGuildCooking(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfLoadGuildCooking)

	guild, _ := GetGuildInfoByCharacterId(s, s.CharID)
	data, err := s.Server.db.Queryx("SELECT id, meal_id, level, expires FROM guild_meals WHERE guild_id = $1", guild.ID)
	if err != nil {
		s.logger.Fatal("Failed to get guild meals from db", zap.Error(err))
	}
	temp := byteframe.NewByteFrame()
	count := 0
	for data.Next() {
		mealData := &GuildMeal{}
		err = data.StructScan(&mealData)
		if err != nil {
			s.logger.Fatal("Failed to scan meal data", zap.Error(err))
		}
		if mealData.Expires > uint32(Time_Current_Adjusted().Add(-60*time.Minute).Unix()) {
			count++
			temp.WriteUint32(mealData.ID)
			temp.WriteUint32(mealData.MealID)
			temp.WriteUint32(mealData.Level)
			temp.WriteUint32(mealData.Expires)
		}
	}
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(count))
	bf.WriteBytes(temp.Data())
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfRegistGuildCooking(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfRegistGuildCooking)
	guild, _ := GetGuildInfoByCharacterId(s, s.CharID)
	if pkt.OverwriteID != 0 {
		_, err := s.Server.db.Exec("DELETE FROM guild_meals WHERE id = $1", pkt.OverwriteID)
		if err != nil {
			s.logger.Fatal("Failed to delete meal in db", zap.Error(err))
		}
	}
	_, err := s.Server.db.Exec("INSERT INTO guild_meals (guild_id, meal_id, level, expires) VALUES ($1, $2, $3, $4)", guild.ID, pkt.MealID, pkt.Success, Time_Current_Adjusted().Add(30*time.Minute).Unix())
	if err != nil {
		s.logger.Fatal("Failed to register meal in db", zap.Error(err))
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x01, 0x00})
}

func handleMsgMhfGetGuildWeeklyBonusMaster(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildWeeklyBonusMaster)

	// Values taken from brand new guild capture
	doAckBufSucceed(s, pkt.AckHandle, make([]byte, 0x28))
}
func handleMsgMhfGetGuildWeeklyBonusActiveCount(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetGuildWeeklyBonusActiveCount)
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0x3C) // Active count
	bf.WriteUint8(0x3C) // Current active count
	bf.WriteUint8(0x00) // New active count
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGuildHuntdata(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGuildHuntdata)
	bf := byteframe.NewByteFrame()
	switch pkt.Operation {
	case 0: // Unk
		doAckBufSucceed(s, pkt.AckHandle, []byte{})
	case 1: // Get Huntdata
		bf.WriteUint8(0) // Entries
		/* Entry format
		uint32 UnkID
		uint32 MonID
		*/
		doAckBufSucceed(s, pkt.AckHandle, bf.Data())
	case 2: // Unk, controls glow
		doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00})
	}
}

type MessageBoardPost struct {
	Type      uint32 `db:"post_type"`
	StampID   uint32 `db:"stamp_id"`
	Title     string `db:"title"`
	Body      string `db:"body"`
	AuthorID  uint32 `db:"author_id"`
	Timestamp uint64 `db:"created_at"`
	LikedBy   string `db:"liked_by"`
}

func handleMsgMhfEnumerateGuildMessageBoard(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateGuildMessageBoard)
	guild, _ := GetGuildInfoByCharacterId(s, s.CharID)

	msgs, err := s.Server.db.Queryx("SELECT post_type, stamp_id, title, body, author_id, (EXTRACT(epoch FROM created_at)::int) as created_at, liked_by FROM guild_posts WHERE guild_id = $1 AND post_type = $2 ORDER BY created_at DESC", guild.ID, int(pkt.BoardType))
	if err != nil {
		s.logger.Fatal("Failed to get guild messages from db", zap.Error(err))
	}

	bf := byteframe.NewByteFrame()
	noMsgs := true
	postCount := 0
	for msgs.Next() {
		noMsgs = false
		postCount++
		postData := &MessageBoardPost{}
		err = msgs.StructScan(&postData)
		if err != nil {
			s.logger.Fatal("Failed to get guild messages from db", zap.Error(err))
		}

		bf.WriteUint32(postData.Type)
		bf.WriteUint32(postData.AuthorID)
		bf.WriteUint64(postData.Timestamp)
		likedBySlice := strings.Split(postData.LikedBy, ",")
		if likedBySlice[0] == "" {
			bf.WriteUint32(0)
		} else {
			bf.WriteUint32(uint32(len(likedBySlice)))
		}
		bf.WriteBool(stringsupport.CSVContains(postData.LikedBy, int(s.CharID)))
		bf.WriteUint32(postData.StampID)
		ps.Uint32(bf, postData.Title, true)
		ps.Uint32(bf, postData.Body, true)
	}
	if noMsgs {
		doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
	} else {
		data := byteframe.NewByteFrame()
		data.WriteUint32(uint32(postCount))
		data.WriteBytes(bf.Data())
		doAckBufSucceed(s, pkt.AckHandle, data.Data())
	}
}

func handleMsgMhfUpdateGuildMessageBoard(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUpdateGuildMessageBoard)
	bf := byteframe.NewByteFrameFromBytes(pkt.Request)
	guild, _ := GetGuildInfoByCharacterId(s, s.CharID)
	if guild == nil {
		if pkt.MessageOp == 5 {
			doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
			return
		}
		doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
		return
	}
	var titleConv, bodyConv string
	switch pkt.MessageOp {
	case 0: // Create message
		postType := bf.ReadUint32() // 0 = message, 1 = news
		stampId := bf.ReadUint32()
		titleLength := bf.ReadUint32()
		bodyLength := bf.ReadUint32()
		title := bf.ReadBytes(uint(titleLength))
		body := bf.ReadBytes(uint(bodyLength))
		titleConv = stringsupport.SJISToUTF8(title)
		bodyConv = stringsupport.SJISToUTF8(body)
		_, err := s.Server.db.Exec("INSERT INTO guild_posts (guild_id, author_id, stamp_id, post_type, title, body) VALUES ($1, $2, $3, $4, $5, $6)", guild.ID, s.CharID, int(stampId), int(postType), titleConv, bodyConv)
		if err != nil {
			s.logger.Fatal("Failed to add new guild message to db", zap.Error(err))
		}
		// TODO: if there are too many messages, purge excess
		_, err = s.Server.db.Exec("")
		if err != nil {
			s.logger.Fatal("Failed to remove excess guild messages from db", zap.Error(err))
		}
	case 1: // Delete message
		postType := bf.ReadUint32()
		timestamp := bf.ReadUint64()
		_, err := s.Server.db.Exec("DELETE FROM guild_posts WHERE post_type = $1 AND (EXTRACT(epoch FROM created_at)::int) = $2 AND guild_id = $3", int(postType), int(timestamp), guild.ID)
		if err != nil {
			s.logger.Fatal("Failed to delete guild message from db", zap.Error(err))
		}
	case 2: // Update message
		postType := bf.ReadUint32()
		timestamp := bf.ReadUint64()
		titleLength := bf.ReadUint32()
		bodyLength := bf.ReadUint32()
		title := bf.ReadBytes(uint(titleLength))
		body := bf.ReadBytes(uint(bodyLength))
		titleConv = stringsupport.SJISToUTF8(title)
		bodyConv = stringsupport.SJISToUTF8(body)
		_, err := s.Server.db.Exec("UPDATE guild_posts SET title = $1, body = $2 WHERE post_type = $3 AND (EXTRACT(epoch FROM created_at)::int) = $4 AND guild_id = $5", titleConv, bodyConv, int(postType), int(timestamp), guild.ID)
		if err != nil {
			s.logger.Fatal("Failed to update guild message in db", zap.Error(err))
		}
	case 3: // Update stamp
		postType := bf.ReadUint32()
		timestamp := bf.ReadUint64()
		stampId := bf.ReadUint32()
		_, err := s.Server.db.Exec("UPDATE guild_posts SET stamp_id = $1 WHERE post_type = $2 AND (EXTRACT(epoch FROM created_at)::int) = $3 AND guild_id = $4", int(stampId), int(postType), int(timestamp), guild.ID)
		if err != nil {
			s.logger.Fatal("Failed to update guild message stamp in db", zap.Error(err))
		}
	case 4: // Like message
		postType := bf.ReadUint32()
		timestamp := bf.ReadUint64()
		likeState := bf.ReadBool()
		var likedBy string
		err := s.Server.db.QueryRow("SELECT liked_by FROM guild_posts WHERE post_type = $1 AND (EXTRACT(epoch FROM created_at)::int) = $2 AND guild_id = $3", int(postType), int(timestamp), guild.ID).Scan(&likedBy)
		if err != nil {
			s.logger.Fatal("Failed to get guild message like data from db", zap.Error(err))
		} else {
			if likeState {
				likedBy = stringsupport.CSVAdd(likedBy, int(s.CharID))
				_, err := s.Server.db.Exec("UPDATE guild_posts SET liked_by = $1 WHERE post_type = $2 AND (EXTRACT(epoch FROM created_at)::int) = $3 AND guild_id = $4", likedBy, int(postType), int(timestamp), guild.ID)
				if err != nil {
					s.logger.Fatal("Failed to like guild message in db", zap.Error(err))
				}
			} else {
				likedBy = stringsupport.CSVRemove(likedBy, int(s.CharID))
				_, err := s.Server.db.Exec("UPDATE guild_posts SET liked_by = $1 WHERE post_type = $2 AND (EXTRACT(epoch FROM created_at)::int) = $3 AND guild_id = $4", likedBy, int(postType), int(timestamp), guild.ID)
				if err != nil {
					s.logger.Fatal("Failed to unlike guild message in db", zap.Error(err))
				}
			}
		}
	case 5: // Check for new messages
		var timeChecked int
		var newPosts int
		err := s.Server.db.QueryRow("SELECT (EXTRACT(epoch FROM guild_post_checked)::int) FROM characters WHERE id = $1", s.CharID).Scan(&timeChecked)
		if err != nil {
			s.logger.Fatal("Failed to get last guild post check timestamp from db", zap.Error(err))
		} else {
			_, err = s.Server.db.Exec("UPDATE characters SET guild_post_checked = $1 WHERE id = $2", time.Now(), s.CharID)
			if err != nil {
				s.logger.Fatal("Failed to update guild post check timestamp in db", zap.Error(err))
			} else {
				err = s.Server.db.QueryRow("SELECT COUNT(*) FROM guild_posts WHERE guild_id = $1 AND (EXTRACT(epoch FROM created_at)::int) > $2 AND author_id != $3", guild.ID, timeChecked, s.CharID).Scan(&newPosts)
				if err != nil {
					s.logger.Fatal("Failed to check for new guild posts in db", zap.Error(err))
				} else {
					if newPosts > 0 {
						doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x01})
						return
					}
				}
			}
		}
	}
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfEntryRookieGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEntryRookieGuild)
	doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfUpdateForceGuildRank(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfAddGuildWeeklyBonusExceptionalUser(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAddGuildWeeklyBonusExceptionalUser)
	// TODO: record pkt.NumUsers to DB
	// must use addition
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGenerateUdGuildMap(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfUpdateGuild(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetGuildManageRight(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetGuildManageRight)
	// nolint:errcheck
	s.Server.db.Exec("UPDATE guild_characters SET recruiter=$1 WHERE character_id=$2", pkt.Allowed, pkt.CharID)
	doAckBufSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfCheckMonthlyItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfCheckMonthlyItem)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x01})
	// TODO: Implement month-by-month tracker, 0 = Not claimed, 1 = Claimed
	// Also handles HLC and EXC items, IDs = 064D, 076B
}

func handleMsgMhfAcquireMonthlyItem(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAcquireMonthlyItem)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfEnumerateInvGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateInvGuild)
	stubEnumerateNoResults(s, pkt.AckHandle)
}

func handleMsgMhfOperationInvGuild(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfOperationInvGuild)
	doAckSimpleFail(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfUpdateGuildcard(s *Session, p mhfpacket.MHFPacket) {}
