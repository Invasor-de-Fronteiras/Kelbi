package channelserver

import (
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
)

func handleMsgMhfRegisterEvent(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfRegisterEvent)
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(pkt.Unk2)
	bf.WriteUint8(pkt.Unk4)
	bf.WriteUint16(0x1142)
	doAckSimpleSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfReleaseEvent(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfReleaseEvent)

	// Do this ack manually because it uses a non-(0|1) error code
	/*
		_ACK_SUCCESS = 0
		_ACK_ERROR = 1

		_ACK_EINPROGRESS = 16
		_ACK_ENOENT = 17
		_ACK_ENOSPC = 18
		_ACK_ETIMEOUT = 19

		_ACK_EINVALID = 64
		_ACK_EFAILED = 65
		_ACK_ENOMEM = 66
		_ACK_ENOTEXIT = 67
		_ACK_ENOTREADY = 68
		_ACK_EALREADY = 69
		_ACK_DISABLE_WORK = 71
	*/
	s.QueueSendMHF(&mhfpacket.MsgSysAck{
		AckHandle:        pkt.AckHandle,
		IsBufferResponse: false,
		ErrorCode:        0x41,
		AckData:          []byte{0x00, 0x00, 0x00, 0x00},
	})
}

func handleMsgMhfEnumerateEvent(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfEnumerateEvent)
	stubEnumerateNoResults(s, pkt.AckHandle)
}

type activeFeature struct {
	StartTime      time.Time `db:"start_time"`
	ActiveFeatures uint32    `db:"featured"`
}

func handleMsgMhfGetWeeklySchedule(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetWeeklySchedule)

	var features []activeFeature
	rows, _ := s.Server.db.Queryx(`SELECT start_time, featured FROM feature_weapon WHERE start_time=$1 OR start_time=$2`, TimeMidnight().Add(-24*time.Hour), TimeMidnight())
	for rows.Next() {
		var feature activeFeature
		// nolint:errcheck
		rows.StructScan(&feature)
		features = append(features, feature)
	}

	if len(features) < 2 {
		if len(features) == 0 {
			feature := generateFeatureWeapons(s.Server.erupeConfig.GameplayOptions.FeaturedWeapons, TimeMidnight().Add(-24*time.Hour))
			// feature.StartTime =
			features = append(features, feature)
			// nolint:errcheck
			s.Server.db.Exec(`INSERT INTO feature_weapon VALUES ($1, $2)`, feature.StartTime, feature.ActiveFeatures)
		}
		feature := generateFeatureWeapons(s.Server.erupeConfig.GameplayOptions.FeaturedWeapons, TimeMidnight())
		// feature.StartTime =
		features = append(features, feature)
		// nolint:errcheck
		s.Server.db.Exec(`INSERT INTO feature_weapon VALUES ($1, $2)`, feature.StartTime, feature.ActiveFeatures)
	}

	bf := byteframe.NewByteFrame()
	bf.WriteUint8(uint8(len(features)))
	bf.WriteUint32(uint32(TimeAdjusted().Add(-5 * time.Minute).Unix()))
	for _, feature := range features {
		bf.WriteUint32(uint32(feature.StartTime.Unix()))
		bf.WriteUint32(feature.ActiveFeatures)
		bf.WriteUint16(0)
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

type Weapon struct {
	name string
	id   uint32
}

var weaponIds = []Weapon{
	{name: "sns", id: 16},
	{name: "ds", id: 64},
	{name: "gs", id: 1},
	{name: "ls", id: 128},
	{name: "hm", id: 4},
	{name: "hh", id: 256},
	{name: "lc", id: 8},
	{name: "gl", id: 512},
	{name: "sw", id: 4096},
	{name: "tf", id: 2048},
	{name: "ms", id: 8192},
	{name: "lbg", id: 32},
	{name: "hbg", id: 2},
	{name: "bw", id: 1024},
}

func generateFeatureWeapons(count int, startTime time.Time) activeFeature {
	if count > 14 {
		count = 14
	}
	nums := make([]uint32, 0)

	dayOfYear := startTime.YearDay()
	weaponIndex := dayOfYear % 14

	var result uint32 = 0

	for len(nums) < count {
		if weaponIndex >= 14 {
			weaponIndex = 0
		}

		weapon := weaponIds[weaponIndex]
		exist := false
		for _, v := range nums {
			if v == weapon.id {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, weapon.id)
		}

		weaponIndex++
	}

	for _, num := range nums {
		result += num
	}

	return activeFeature{ActiveFeatures: result, StartTime: startTime}
}

type loginBoost struct {
	WeekReq    uint8 `db:"week_req"`
	WeekCount  uint8
	Active     bool
	Expiration time.Time `db:"expiration"`
	Reset      time.Time `db:"reset"`
}

func handleMsgMhfGetKeepLoginBoostStatus(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetKeepLoginBoostStatus)

	bf := byteframe.NewByteFrame()

	var loginBoosts []loginBoost
	rows, err := s.Server.db.Queryx("SELECT week_req, expiration, reset FROM login_boost WHERE char_id=$1 ORDER BY week_req", s.CharID)
	if err != nil || s.Server.erupeConfig.GameplayOptions.DisableLoginBoost {
		rows.Close()
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 35))
		return
	}
	for rows.Next() {
		var temp loginBoost
		rows.StructScan(&temp)
		loginBoosts = append(loginBoosts, temp)
	}
	if len(loginBoosts) == 0 {
		temp := TimeWeekStart()
		loginBoosts = []loginBoost{
			{WeekReq: 1, Expiration: temp},
			{WeekReq: 2, Expiration: temp},
			{WeekReq: 3, Expiration: temp},
			{WeekReq: 4, Expiration: temp},
			{WeekReq: 5, Expiration: temp},
		}
		for _, boost := range loginBoosts {
			s.Server.db.Exec(`INSERT INTO login_boost VALUES ($1, $2, $3, $4)`, s.CharID, boost.WeekReq, boost.Expiration, time.Time{})
		}
	}

	for _, boost := range loginBoosts {
		// Reset if next week
		if !boost.Reset.IsZero() && boost.Reset.Before(TimeAdjusted()) {
			boost.Expiration = TimeWeekStart()
			boost.Reset = time.Time{}
			s.Server.db.Exec(`UPDATE login_boost SET expiration=$1, reset=$2 WHERE char_id=$3 AND week_req=$4`, boost.Expiration, boost.Reset, s.CharID, boost.WeekReq)
		}

		boost.WeekCount = uint8((TimeAdjusted().Unix()-boost.Expiration.Unix())/604800 + 1)

		if boost.WeekCount >= boost.WeekReq {
			boost.Active = true
			boost.WeekCount = boost.WeekReq
		}

		// Show reset timer on expired boosts
		if boost.Reset.After(TimeAdjusted()) {
			boost.Active = true
			boost.WeekCount = 0
		}

		bf.WriteUint8(boost.WeekReq)
		bf.WriteBool(boost.Active)
		bf.WriteUint8(boost.WeekCount)
		if !boost.Reset.IsZero() {
			bf.WriteUint32(uint32(boost.Expiration.Unix()))
		} else {
			bf.WriteUint32(0)
		}
	}

	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfUseKeepLoginBoost(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfUseKeepLoginBoost)
	var expiration time.Time
	bf := byteframe.NewByteFrame()
	bf.WriteUint8(0)
	switch pkt.BoostWeekUsed {
	case 1:
		fallthrough
	case 3:
		expiration = TimeAdjusted().Add(120 * time.Minute)
	case 4:
		expiration = TimeAdjusted().Add(180 * time.Minute)
	case 2:
		fallthrough
	case 5:
		expiration = TimeAdjusted().Add(240 * time.Minute)
	}
	bf.WriteUint32(uint32(expiration.Unix()))
	s.Server.db.Exec(`UPDATE login_boost SET expiration=$1, reset=$2 WHERE char_id=$3 AND week_req=$4`, expiration, TimeWeekNext(), s.CharID, pkt.BoostWeekUsed)
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}

func handleMsgMhfGetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetRestrictionEvent)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}
