package channelserver

import (
	"math"
	"math/rand"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
	timeServerFix "erupe-ce/server/channelserver/timeserver"
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
	StartTime      time.Time
	ActiveFeatures uint32
	Unk1           uint16
}

func handleMsgMhfGetWeeklySchedule(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetWeeklySchedule)
	persistentEventSchedule := make([]activeFeature, 8) // generate day after weekly restart
	for x := -1; x < 7; x++ {
		feat := generateActiveWeapons(14) // number of active weapons
		// TODO: only generate this once per restart (server should be restarted weekly)
		// then load data from db instead of regenerating
		persistentEventSchedule[x+1] = activeFeature{
			StartTime:      Time_Current_Midnight().Add(time.Duration(24*x) * time.Hour),
			ActiveFeatures: uint32(feat),
			Unk1:           0,
		}
	}

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(len(persistentEventSchedule)))                           // Entry count, client only parses the first 7 or 8.
	resp.WriteUint32(uint32(Time_Current_Adjusted().Add(-5 * time.Minute).Unix())) // 5 minutes ago server time

	for _, es := range persistentEventSchedule {
		resp.WriteUint32(uint32(es.StartTime.Unix()))
		resp.WriteUint32(es.ActiveFeatures)
		resp.WriteUint16(es.Unk1)
	}
	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func generateActiveWeapons(count int) int {
	nums := make([]int, 0)
	var result int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		num := r.Intn(14)
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	for _, num := range nums {
		result += int(math.Pow(2, float64(num)))
	}
	return result
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

func handleMsgMhfGetUdSchedule(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdSchedule)
	var t = timeServerFix.Tstatic_midnight()
	var event int = s.Server.erupeConfig.DevModeOptions.DivaEvent

	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	// Events with time limits are Festival with Sign up, Soul Week and Winners Weeks
	// Diva Defense with Prayer, Interception and Song weeks
	// Mezeporta Festival with simply 'available' being a weekend thing
	resp := byteframe.NewByteFrame()
	resp.WriteUint32(0x1d5fda5c) // Unk (1d5fda5c, 0b5397df)

	if event == 1 {
		resp.WriteUint32(uint32(midnight.Add(24 * 21 * time.Hour).Unix())) // Week 1 Timestamp, Festi start?
	} else {
		resp.WriteUint32(uint32(midnight.Add(-24 * 21 * time.Hour).Unix())) // Week 1 Timestamp, Festi start?
	}

	if event == 2 {
		resp.WriteUint32(uint32(midnight.Add(24 * 14 * time.Hour).Unix())) // Week 2 Timestamp
		resp.WriteUint32(uint32(midnight.Add(24 * 14 * time.Hour).Unix())) // Week 2 Timestamp
	} else {
		resp.WriteUint32(uint32(midnight.Add(-24 * 14 * time.Hour).Unix())) // Week 2 Timestamp
		resp.WriteUint32(uint32(midnight.Add(-24 * 14 * time.Hour).Unix())) // Week 2 Timestamp
	}

	if event == 3 {
		resp.WriteUint32(uint32(midnight.Add((24) * 7 * time.Hour).Unix()))  // Diva Defense Interception
		resp.WriteUint32(uint32(midnight.Add((24) * 14 * time.Hour).Unix())) // Diva Defense Greeting Song
	} else {
		resp.WriteUint32(uint32(midnight.Add((-24) * 7 * time.Hour).Unix()))  // Diva Defense Interception
		resp.WriteUint32(uint32(midnight.Add((-24) * 14 * time.Hour).Unix())) // Diva Defense Greeting Song
	}

	resp.WriteUint16(0x19) // Unk 00011001
	resp.WriteUint16(0x2d) // Unk 00101101
	resp.WriteUint16(0x02) // Unk 00000010
	resp.WriteUint16(0x02) // Unk 00000010

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetUdInfo(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdInfo)
	// Message that appears on the Diva Defense NPC and triggers the green exclamation mark
	udInfos := []struct {
		Text      string
		StartTime time.Time
		EndTime   time.Time
	}{
		/*{
			Text:      " ~C17【Erupe】 is dead event!\n\n■Features\n~C18 Dont bother walking around!\n~C17 Take down your DB by doing \n~C17 nearly anything!",
			StartTime: Time_static().Add(time.Duration(-5) * time.Minute), // Event started 5 minutes ago,
			EndTime:   Time_static().Add(time.Duration(24) * time.Hour),   // Event ends in 5 minutes,
		}, */
	}

	resp := byteframe.NewByteFrame()
	resp.WriteUint8(uint8(len(udInfos)))
	for _, udInfo := range udInfos {
		resp.WriteBytes(fixedSizeShiftJIS(udInfo.Text, 1024))
		resp.WriteUint32(uint32(udInfo.StartTime.Unix()))
		resp.WriteUint32(uint32(udInfo.EndTime.Unix()))
	}

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfGetBoostTime(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetBoostTime)

	doAckBufSucceed(s, pkt.AckHandle, []byte{})
	updateRights(s)
}

func handleMsgMhfGetBoostTimeLimit(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetBoostTimeLimit)
	doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetBoostRight(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetBoostRight)
	doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfPostBoostTimeQuestReturn(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfPostBoostTimeQuestReturn)
	doAckSimpleSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfStartBoostTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostBoostTime(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfPostBoostTimeLimit(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfSetRestrictionEvent(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfSetRestrictionEvent)
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}
