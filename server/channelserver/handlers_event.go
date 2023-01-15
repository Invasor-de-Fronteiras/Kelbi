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
	rows, _ := s.Server.db.Queryx(`SELECT start_time, featured FROM feature_weapon WHERE start_time=$1 OR start_time=$2`, Time_Current_Midnight().Add(-24*time.Hour), Time_Current_Midnight())
	for rows.Next() {
		var feature activeFeature
		// nolint:errcheck
		rows.StructScan(&feature)
		features = append(features, feature)
	}

	if len(features) < 2 {
		if len(features) == 0 {
			feature := generateFeatureWeapons(s.Server.erupeConfig.FeaturedWeapons, Time_Current_Midnight().Add(-24*time.Hour))
			// feature.StartTime =
			features = append(features, feature)
			// nolint:errcheck
			s.Server.db.Exec(`INSERT INTO feature_weapon VALUES ($1, $2)`, feature.StartTime, feature.ActiveFeatures)
		}
		feature := generateFeatureWeapons(s.Server.erupeConfig.FeaturedWeapons, Time_Current_Midnight())
		// feature.StartTime =
		features = append(features, feature)
		// nolint:errcheck
		s.Server.db.Exec(`INSERT INTO feature_weapon VALUES ($1, $2)`, feature.StartTime, feature.ActiveFeatures)
	}

	bf := byteframe.NewByteFrame()
	bf.WriteUint8(2)
	bf.WriteUint32(uint32(Time_Current_Adjusted().Add(-5 * time.Minute).Unix()))
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
	{name: "bw", id: 1024},
	{name: "lbg", id: 32},
	{name: "hbg", id: 2},
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
	WeekReq, WeekCount, LastWeek uint8
	Available                    bool
	Expiration                   uint32
}

func handleMsgMhfGetKeepLoginBoostStatus(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetKeepLoginBoostStatus)

	var loginBoostStatus []loginBoost
	insert := false
	boostState, err := s.Server.db.Query("SELECT week_req, week_count, available, end_time, last_week FROM login_boost_state WHERE char_id=$1 ORDER BY week_req ASC", s.CharID)
	if err != nil {
		panic(err)
	}
	for boostState.Next() {
		var boost loginBoost
		err = boostState.Scan(&boost.WeekReq, &boost.WeekCount, &boost.Available, &boost.Expiration, &boost.LastWeek)
		if err != nil {
			panic(err)
		}
		loginBoostStatus = append(loginBoostStatus, boost)
	}
	if len(loginBoostStatus) == 0 {
		// create default Entries (should only been week 1 with )
		insert = true
		loginBoostStatus = []loginBoost{
			{
				WeekReq:    1,    // weeks needed
				WeekCount:  1,    // weeks passed
				Available:  true, // available
				Expiration: 0,    //uint32(t.Add(120 * time.Minute).Unix()), // uncomment to enable permanently
				LastWeek:   0,
			},
			{
				WeekReq:    2,
				WeekCount:  2,
				Available:  true,
				Expiration: 0,
				LastWeek:   0,
			},
			{
				WeekReq:    3,
				WeekCount:  3,
				Available:  true,
				Expiration: 0,
				LastWeek:   0,
			},
			{
				WeekReq:    4,
				WeekCount:  4,
				Available:  true,
				Expiration: 0,
				LastWeek:   0,
			},
			{
				WeekReq:    5,
				WeekCount:  5,
				Available:  true,
				Expiration: 0,
				LastWeek:   0,
			},
		}
	}
	resp := byteframe.NewByteFrame()
	_, week := time.Now().ISOWeek()
	CurrentWeek := uint8((week % 5) + 1)

	for d := range loginBoostStatus {
		if loginBoostStatus[d].LastWeek != CurrentWeek &&
			loginBoostStatus[d].WeekCount < loginBoostStatus[d].WeekReq {
			if (loginBoostStatus[d].LastWeek == 5 && CurrentWeek == 1) || (loginBoostStatus[d].LastWeek+1 == CurrentWeek) {
				loginBoostStatus[d].WeekCount += 1
			} else {
				loginBoostStatus[d].WeekCount = 0
			}
		}

		loginBoostStatus[d].LastWeek = CurrentWeek

		if loginBoostStatus[d].Expiration > 0 && uint32(time.Now().In(time.FixedZone("UTC+1", 1*60*60)).Unix()) >= loginBoostStatus[d].Expiration {
			loginBoostStatus[d].Expiration = 0
			loginBoostStatus[d].WeekCount = 0
		}
		if !insert {
			_, err := s.Server.db.Exec(`UPDATE login_boost_state SET week_count=$1, end_time=$2, available=$3, last_week=$4 WHERE char_id=$5 AND week_req=$6`, loginBoostStatus[d].WeekCount, loginBoostStatus[d].Expiration, loginBoostStatus[d].Available, loginBoostStatus[d].LastWeek, s.CharID, loginBoostStatus[d].WeekReq)
			if err != nil {
				panic(err)
			}
		}
	}
	for _, v := range loginBoostStatus {
		if insert {
			_, err := s.Server.db.Exec(`INSERT INTO login_boost_state (char_id, week_req, week_count, available, last_week, end_time) VALUES ($1,$2,$3,$4,$5, $6)`, s.CharID, v.WeekReq, v.WeekCount, v.Available, v.LastWeek, v.Expiration)
			if err != nil {
				panic(err)
			}
		}
		resp.WriteUint8(v.WeekReq)
		resp.WriteBool(v.WeekCount >= v.WeekReq)
		resp.WriteUint8(v.WeekCount)
		resp.WriteUint32(v.Expiration)
	}
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfUseKeepLoginBoost(s *Session, p mhfpacket.MHFPacket) {
	// Directly interacts with MhfGetKeepLoginBoostStatus
	// TODO: make these states persistent on a per character basis
	pkt := p.(*mhfpacket.MsgMhfUseKeepLoginBoost)
	var t = time.Now().In(time.FixedZone("UTC+1", 1*60*60))
	resp := byteframe.NewByteFrame()
	resp.WriteUint8(0)

	// response is end timestamp based on input
	switch pkt.BoostWeekUsed {
	case 1:
		t = t.Add(120 * time.Minute)
		resp.WriteUint32(uint32(t.Unix()))
	case 2:
		t = t.Add(240 * time.Minute)
		resp.WriteUint32(uint32(t.Unix()))
	case 3:
		t = t.Add(120 * time.Minute)
		resp.WriteUint32(uint32(t.Unix()))
	case 4:
		t = t.Add(180 * time.Minute)
		resp.WriteUint32(uint32(t.Unix()))
	case 5:
		t = t.Add(240 * time.Minute)
		resp.WriteUint32(uint32(t.Unix()))
	}
	_, err := s.Server.db.Exec(`UPDATE login_boost_state SET end_time=$1 WHERE char_id=$2 AND week_req=$3`, uint32(t.Unix()), s.CharID, pkt.BoostWeekUsed)

	if err != nil {
		panic(err)
	}
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetRestrictionEvent)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}
