package channelserver

import (
	"erupe-ce/common/byteframe"
	"erupe-ce/network/mhfpacket"
	"strings"
)

func handleMsgSysOperateRegister(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysOperateRegister)
	bf := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload)
	s.Server.raviente.Lock()
	switch pkt.SemaphoreID {
	case 4:
		resp := byteframe.NewByteFrame()
		size := 6
		for i := 0; i < len(bf.Data())-1; i += size {
			op := bf.ReadUint8()
			dest := bf.ReadUint8()
			data := bf.ReadUint32()
			resp.WriteUint8(1)
			resp.WriteUint8(dest)
			ref := &s.Server.raviente.state.stateData[dest]
			damageMultiplier := s.Server.raviente.state.damageMultiplier
			switch op {
			case 2:
				resp.WriteUint32(*ref)
				if dest == 28 { // Berserk resurrection tracker
					resp.WriteUint32(*ref + data)
					*ref += data
				} else if dest == 17 { // Berserk poison tracker
					if damageMultiplier == 1 {
						resp.WriteUint32(*ref + data)
						*ref += data
					} else {
						resp.WriteUint32(*ref)
					}
				} else {
					resp.WriteUint32(*ref + data*damageMultiplier)
					*ref += data * damageMultiplier
				}
			case 13:
				fallthrough
			case 14:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				*ref = data
			}
		}
		resp.WriteUint8(0)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	case 5:
		resp := byteframe.NewByteFrame()
		size := 6
		for i := 0; i < len(bf.Data())-1; i += size {
			op := bf.ReadUint8()
			dest := bf.ReadUint8()
			data := bf.ReadUint32()
			resp.WriteUint8(1)
			resp.WriteUint8(dest)
			ref := &s.Server.raviente.support.supportData[dest]
			switch op {
			case 2:
				resp.WriteUint32(*ref)
				resp.WriteUint32(*ref + data)
				*ref += data
			case 13:
				fallthrough
			case 14:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				*ref = data
			}
		}
		resp.WriteUint8(0)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	case 6:
		resp := byteframe.NewByteFrame()
		size := 6
		for i := 0; i < len(bf.Data())-1; i += size {
			op := bf.ReadUint8()
			dest := bf.ReadUint8()
			data := bf.ReadUint32()
			resp.WriteUint8(1)
			resp.WriteUint8(dest)
			switch dest {
			case 0:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.nextTime = data
			case 1:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.startTime = data
			case 2:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.killedTime = data
			case 3:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.postTime = data
			case 4:
				ref := &s.Server.raviente.register.register[0]
				switch op {
				case 2:
					resp.WriteUint32(*ref)
					resp.WriteUint32(*ref + data)
					*ref += data
				case 13:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
					*ref = data
				case 14:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
				}
			case 5:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.carveQuest = data
			case 6:
				ref := &s.Server.raviente.register.register[1]
				switch op {
				case 2:
					resp.WriteUint32(*ref)
					resp.WriteUint32(*ref + data)
					*ref += data
				case 13:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
					*ref = data
				case 14:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
				}
			case 7:
				ref := &s.Server.raviente.register.register[2]
				switch op {
				case 2:
					resp.WriteUint32(*ref)
					resp.WriteUint32(*ref + data)
					*ref += data
				case 13:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
					*ref = data
				case 14:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
				}
			case 8:
				ref := &s.Server.raviente.register.register[3]
				switch op {
				case 2:
					resp.WriteUint32(*ref)
					resp.WriteUint32(*ref + data)
					*ref += data
				case 13:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
					*ref = data
				case 14:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
				}
			case 9:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.maxPlayers = data
			case 10:
				resp.WriteUint32(0)
				resp.WriteUint32(data)
				s.Server.raviente.register.ravienteType = data
			case 11:
				ref := &s.Server.raviente.register.register[4]
				switch op {
				case 2:
					resp.WriteUint32(*ref)
					resp.WriteUint32(*ref + data)
					*ref += data
				case 13:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
					*ref = data
				case 14:
					resp.WriteUint32(0)
					resp.WriteUint32(data)
				}
			default:
				resp.WriteUint32(0)
				resp.WriteUint32(0)
			}
		}
		resp.WriteUint8(0)
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	}
	s.notifyRavi()
	s.Server.raviente.Unlock()
}

func handleMsgSysLoadRegister(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysLoadRegister)
	r := pkt.Unk1
	switch r {
	case 12:
		resp := byteframe.NewByteFrame()
		resp.WriteUint8(0)
		resp.WriteUint8(12)
		resp.WriteUint32(s.Server.raviente.register.nextTime)
		resp.WriteUint32(s.Server.raviente.register.startTime)
		resp.WriteUint32(s.Server.raviente.register.killedTime)
		resp.WriteUint32(s.Server.raviente.register.postTime)
		resp.WriteUint32(s.Server.raviente.register.register[0])
		resp.WriteUint32(s.Server.raviente.register.carveQuest)
		resp.WriteUint32(s.Server.raviente.register.register[1])
		resp.WriteUint32(s.Server.raviente.register.register[2])
		resp.WriteUint32(s.Server.raviente.register.register[3])
		resp.WriteUint32(s.Server.raviente.register.maxPlayers)
		resp.WriteUint32(s.Server.raviente.register.ravienteType)
		resp.WriteUint32(s.Server.raviente.register.register[4])
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	case 29:
		resp := byteframe.NewByteFrame()
		resp.WriteUint8(0)
		resp.WriteUint8(29)
		for _, v := range s.Server.raviente.state.stateData {
			resp.WriteUint32(v)
		}
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	case 25:
		resp := byteframe.NewByteFrame()
		resp.WriteUint8(0)
		resp.WriteUint8(25)
		for _, v := range s.Server.raviente.support.supportData {
			resp.WriteUint32(v)
		}
		doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	}
}

func (s *Session) notifyRavi() {
	var temp mhfpacket.MHFPacket
	raviNotif := byteframe.NewByteFrame()
	temp = &mhfpacket.MsgSysNotifyRegister{RegisterID: 4}
	raviNotif.WriteUint16(uint16(temp.Opcode()))
	// nolint:errcheck // Error return value of `temp.Build` is not checked
	temp.Build(raviNotif, s.clientContext)
	temp = &mhfpacket.MsgSysNotifyRegister{RegisterID: 5}
	raviNotif.WriteUint16(uint16(temp.Opcode()))
	// nolint:errcheck
	temp.Build(raviNotif, s.clientContext)
	temp = &mhfpacket.MsgSysNotifyRegister{RegisterID: 6}
	raviNotif.WriteUint16(uint16(temp.Opcode()))
	// nolint:errcheck
	temp.Build(raviNotif, s.clientContext)
	raviNotif.WriteUint16(0x0010) // End it.
	sema := getRaviSemaphore(s)
	if sema != "" {
		for session := range s.Server.Semaphore[sema].clients {
			session.QueueSend(raviNotif.Data())
		}
	}
}

func getRaviSemaphore(s *Session) string {
	for _, semaphore := range s.Server.Semaphore {
		if strings.HasPrefix(semaphore.StageId, "hs_l0u3B5") && strings.HasSuffix(semaphore.StageId, "4") {
			return semaphore.StageId
		}
	}
	return ""
}

func resetRavi(s *Session) {
	s.Server.raviente.Lock()
	s.Server.raviente.register.nextTime = 0
	s.Server.raviente.register.startTime = 0
	s.Server.raviente.register.killedTime = 0
	s.Server.raviente.register.postTime = 0
	s.Server.raviente.register.ravienteType = 0
	s.Server.raviente.register.maxPlayers = 0
	s.Server.raviente.register.carveQuest = 0
	s.Server.raviente.state.damageMultiplier = 1
	s.Server.raviente.register.register = []uint32{0, 0, 0, 0, 0}
	s.Server.raviente.state.stateData = []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	s.Server.raviente.support.supportData = []uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	s.Server.raviente.Unlock()
}

func handleMsgSysNotifyRegister(s *Session, p mhfpacket.MHFPacket) {}
