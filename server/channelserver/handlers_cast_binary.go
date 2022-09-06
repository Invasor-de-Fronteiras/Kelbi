package channelserver

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/network/binpacket"
	"erupe-ce/network/mhfpacket"
)

// MSG_SYS_CAST[ED]_BINARY types enum
const (
	BinaryMessageTypeState      = 0
	BinaryMessageTypeChat       = 1
	BinaryMessageTypeData       = 3
	BinaryMessageTypeMailNotify = 4
	BinaryMessageTypeEmote      = 6
)

// MSG_SYS_CAST[ED]_BINARY broadcast types enum
const (
	BroadcastTypeTargeted = 0x01
	BroadcastTypeStage    = 0x03
	BroadcastTypeServer   = 0x06
	BroadcastTypeWorld    = 0x0a
)

func sendServerChatMessage(s *Session, message string) {
	// Make the inside of the casted binary
	bf := byteframe.NewByteFrame()
	bf.SetLE()
	msgBinChat := &binpacket.MsgBinChat{
		Unk0:       0,
		Type:       5,
		Flags:      0x80,
		Message:    message,
		SenderName: "Erupe",
	}
	// nolint:errcheck // Error return value of `.` is not checked
	msgBinChat.Build(bf)

	castedBin := &mhfpacket.MsgSysCastedBinary{
		CharID:         s.CharID,
		MessageType:    BinaryMessageTypeChat,
		RawDataPayload: bf.Data(),
	}

	s.QueueSendMHF(castedBin)
}

func handleMsgSysCastBinary(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysCastBinary)
	tmp := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload)

	if pkt.BroadcastType == 0x03 && pkt.MessageType == 0x03 && len(pkt.RawDataPayload) == 0x10 {
		if tmp.ReadUint16() == 0x0002 && tmp.ReadUint8() == 0x18 {
			_ = tmp.ReadBytes(9)
			tmp.SetLE()
			frame := tmp.ReadUint32()
			sendServerChatMessage(s, fmt.Sprintf("TIME : %d'%d.%03d (%dframe)", frame/30/60, frame/30%60, int(math.Round(float64(frame%30*100)/3)), frame))
		}
	}

	// Parse out the real casted binary payload
	var msgBinTargeted *binpacket.MsgBinTargeted
	var authorLen, msgLen uint16
	var msg []byte

	isDiceCommand := false
	if pkt.MessageType == BinaryMessageTypeChat {
		tmp.SetLE()
		// nolint:errcheck // Error return value of `.` is not checked
		tmp.Seek(int64(0), 0)
		_ = tmp.ReadUint32()
		authorLen = tmp.ReadUint16()
		msgLen = tmp.ReadUint16()
		msg = tmp.ReadNullTerminatedBytes()
	}

	// Customise payload
	realPayload := pkt.RawDataPayload
	if pkt.BroadcastType == BroadcastTypeTargeted {
		tmp.SetBE()
		// nolint:errcheck // Error return value of `.` is not checked
		tmp.Seek(int64(0), 0)
		msgBinTargeted = &binpacket.MsgBinTargeted{}
		err := msgBinTargeted.Parse(tmp)
		if err != nil {
			s.logger.Warn("Failed to parse targeted cast binary")
			return
		}
		realPayload = msgBinTargeted.RawDataPayload
	} else if pkt.MessageType == BinaryMessageTypeChat {
		if msgLen == 6 && string(msg) == "@dice" {
			isDiceCommand = true
			roll := byteframe.NewByteFrame()
			roll.WriteInt16(1) // Unk
			roll.SetLE()
			roll.WriteUint16(4) // Unk
			roll.WriteUint16(authorLen)
			rand.Seed(time.Now().UnixNano())
			dice := fmt.Sprintf("%d", rand.Intn(100)+1)
			roll.WriteUint16(uint16(len(dice) + 1))
			roll.WriteNullTerminatedBytes([]byte(dice))
			roll.WriteNullTerminatedBytes(tmp.ReadNullTerminatedBytes())
			realPayload = roll.Data()
		}
	}

	// Make the response to forward to the other client(s).
	resp := &mhfpacket.MsgSysCastedBinary{
		CharID:         s.CharID,
		BroadcastType:  pkt.BroadcastType, // (The client never uses Type0 upon receiving)
		MessageType:    pkt.MessageType,
		RawDataPayload: realPayload,
	}

	// Send to the proper recipients.
	switch pkt.BroadcastType {
	case BroadcastTypeWorld:
		s.Server.WorldcastMHF(resp, s, nil)
	case BroadcastTypeStage:
		if isDiceCommand {
			s.Stage.BroadcastMHF(resp, nil) // send dice result back to caller
		} else {
			s.Stage.BroadcastMHF(resp, s)
		}
	case BroadcastTypeServer:
		if pkt.MessageType == 1 {
			raviSema := getRaviSemaphore(s)
			if raviSema != "" {
				s.Server.BroadcastMHF(resp, s)
			}
		} else {
			s.Server.BroadcastMHF(resp, s)
		}
	case BroadcastTypeTargeted:
		for _, targetID := range (*msgBinTargeted).TargetCharIDs {
			char := s.Server.FindSessionByCharID(targetID)

			if char != nil {
				char.QueueSendMHF(resp)
			}
		}
	default:
		s.Lock()
		haveStage := s.Stage != nil
		if haveStage {
			s.Stage.BroadcastMHF(resp, s)
		}
		s.Unlock()
	}

	// Handle chat
	if pkt.MessageType == BinaryMessageTypeChat {
		bf := byteframe.NewByteFrameFromBytes(realPayload)

		// IMPORTANT! Casted binary objects are sent _as they are in memory_,
		// this means little endian for LE CPUs, might be different for PS3/PS4/PSP/XBOX.
		bf.SetLE()

		chatMessage := &binpacket.MsgBinChat{}
		// nolint:errcheck // Error return value of `.` is not checked
		chatMessage.Parse(bf)

		fmt.Printf("Got chat message: %+v\n", chatMessage)

		// Flush all objects and users and reload
		if strings.HasPrefix(chatMessage.Message, "!reload") {
			sendServerChatMessage(s, "Reloading players...")
			var temp mhfpacket.MHFPacket
			deleteNotif := byteframe.NewByteFrame()
			for _, object := range s.Stage.Objects {
				if object.OwnerCharID == s.CharID {
					continue
				}
				temp = &mhfpacket.MsgSysDeleteObject{ObjID: object.Id}
				deleteNotif.WriteUint16(uint16(temp.Opcode()))
				// nolint:errcheck
				temp.Build(deleteNotif, s.clientContext)
			}
			for _, session := range s.Server.Sessions {
				if s == session {
					continue
				}
				temp = &mhfpacket.MsgSysDeleteUser{CharID: session.CharID}
				// nolint:nocheck
				deleteNotif.WriteUint16(uint16(temp.Opcode()))
				temp.Build(deleteNotif, s.clientContext)
			}
			deleteNotif.WriteUint16(0x0010)
			s.QueueSend(deleteNotif.Data())
			time.Sleep(500 * time.Millisecond)
			reloadNotif := byteframe.NewByteFrame()
			for _, session := range s.Server.Sessions {
				if s == session {
					continue
				}
				// nolint:errcheck
				temp = &mhfpacket.MsgSysInsertUser{CharID: session.CharID}
				reloadNotif.WriteUint16(uint16(temp.Opcode()))
				temp.Build(reloadNotif, s.clientContext)
				for i := 0; i < 3; i++ {
					temp = &mhfpacket.MsgSysNotifyUserBinary{
						CharID:     session.CharID,
						BinaryType: uint8(i + 1),
					}
					reloadNotif.WriteUint16(uint16(temp.Opcode()))
					temp.Build(reloadNotif, s.clientContext)
				}
			}
			for _, obj := range s.Stage.Objects {
				if obj.OwnerCharID == s.CharID {
					continue
				}
				temp = &mhfpacket.MsgSysDuplicateObject{
					ObjID:       obj.Id,
					X:           obj.X,
					Y:           obj.Y,
					Z:           obj.Z,
					Unk0:        0,
					OwnerCharID: obj.OwnerCharID,
				}
				reloadNotif.WriteUint16(uint16(temp.Opcode()))
				temp.Build(reloadNotif, s.clientContext)
			}
			reloadNotif.WriteUint16(0x0010)
			s.QueueSend(reloadNotif.Data())
		}

		// Set account rights
		// if strings.HasPrefix(chatMessage.Message, "!rights") {
		// 	var v uint32
		// 	n, err := fmt.Sscanf(chatMessage.Message, "!rights %d", &v)
		// 	if err != nil || n != 1 {
		// 		sendServerChatMessage(s, "Error in command. Format: !rights n")
		// 	} else {
		// 		_, err = s.Server.db.Exec("UPDATE users u SET rights=$1 WHERE u.id=(SELECT c.user_id FROM characters c WHERE c.id=$2)", v, s.CharID)
		// 		if err == nil {
		// 			sendServerChatMessage(s, fmt.Sprintf("Set rights integer: %d", v))
		// 		}
		// 	}
		// }

		// Discord integration
		if (pkt.BroadcastType == BroadcastTypeStage && s.Stage.Id == "sl1Ns200p0a0u0") || pkt.BroadcastType == BroadcastTypeWorld {
			s.Server.DiscordChannelSend(chatMessage.SenderName, chatMessage.Message)
		}
		// RAVI COMMANDS V2
		if strings.HasPrefix(chatMessage.Message, "!ravi") {
			if getRaviSemaphore(s) != "" {
				s.Server.raviente.Lock()
				if !strings.HasPrefix(chatMessage.Message, "!ravi ") {
					sendServerChatMessage(s, "No Raviente command specified!")
				} else {
					if strings.HasPrefix(chatMessage.Message, "!ravi start") {
						if s.Server.raviente.register.startTime == 0 {
							s.Server.raviente.register.startTime = s.Server.raviente.register.postTime

							sendServerChatMessage(s, "The Great Slaying will begin in a moment")
							s.notifyRavi()
						} else {
							sendServerChatMessage(s, "The Great Slaying has already begun!")
						}
					} else if strings.HasPrefix(chatMessage.Message, "!ravi sm") || strings.HasPrefix(chatMessage.Message, "!ravi setmultiplier") {
						var num uint16
						n, numerr := fmt.Sscanf(chatMessage.Message, "!ravi sm %d", &num)
						if numerr != nil || n != 1 {
							sendServerChatMessage(s, "Error in command. Format: !ravi sm n")
						} else if s.Server.raviente.state.damageMultiplier == 1 {
							if num > 32 {
								sendServerChatMessage(s, "Raviente multiplier too high, defaulting to 32x")
								s.Server.raviente.state.damageMultiplier = 32
							} else {
								sendServerChatMessage(s, fmt.Sprintf("Raviente multiplier set to %dx", num))
								s.Server.raviente.state.damageMultiplier = uint32(num)
							}
						} else {
							sendServerChatMessage(s, fmt.Sprintf("Raviente multiplier is already set to %dx!", s.Server.raviente.state.damageMultiplier))
						}
					} else if strings.HasPrefix(chatMessage.Message, "!ravi cm") || strings.HasPrefix(chatMessage.Message, "!ravi checkmultiplier") {
						sendServerChatMessage(s, fmt.Sprintf("Raviente multiplier is currently %dx", s.Server.raviente.state.damageMultiplier))
					} else if strings.HasPrefix(chatMessage.Message, "!ravi sr") || strings.HasPrefix(chatMessage.Message, "!ravi sendres") {
						if s.Server.raviente.state.stateData[28] > 0 {
							sendServerChatMessage(s, "Sending resurrection support!")
							s.Server.raviente.state.stateData[28] = 0
						} else {
							sendServerChatMessage(s, "Resurrection support has not been requested!")
						}
					} else if strings.HasPrefix(chatMessage.Message, "!ravi ss") || strings.HasPrefix(chatMessage.Message, "!ravi sendsed") {
						sendServerChatMessage(s, "Sending sedation support if requested!")
						// Total BerRavi HP
						HP := s.Server.raviente.state.stateData[0] + s.Server.raviente.state.stateData[1] + s.Server.raviente.state.stateData[2] + s.Server.raviente.state.stateData[3] + s.Server.raviente.state.stateData[4]
						s.Server.raviente.support.supportData[1] = HP
					} else if strings.HasPrefix(chatMessage.Message, "!ravi rs") || strings.HasPrefix(chatMessage.Message, "!ravi reqsed") {
						sendServerChatMessage(s, "Requesting sedation support!")
						// Total BerRavi HP
						HP := s.Server.raviente.state.stateData[0] + s.Server.raviente.state.stateData[1] + s.Server.raviente.state.stateData[2] + s.Server.raviente.state.stateData[3] + s.Server.raviente.state.stateData[4]
						s.Server.raviente.support.supportData[1] = HP + 12
					} else {
						sendServerChatMessage(s, "Raviente command not recognised!")
					}
				}
				s.Server.raviente.Unlock()
			} else {
				sendServerChatMessage(s, "No one has joined the Great Slaying!")
			}
		}
		// END RAVI COMMANDS V2

		// if strings.HasPrefix(chatMessage.Message, "!tele ") {
		// 	var x, y int16
		// 	n, err := fmt.Sscanf(chatMessage.Message, "!tele %d %d", &x, &y)
		// 	if err != nil || n != 2 {
		// 		sendServerChatMessage(s, "Invalid command. Usage:\"!tele 500 500\"")
		// 	} else {
		// 		sendServerChatMessage(s, fmt.Sprintf("Teleporting to %d %d", x, y))

		// 		// Make the inside of the casted binary
		// 		payload := byteframe.NewByteFrame()
		// 		payload.SetLE()
		// 		payload.WriteUint8(2) // SetState type(position == 2)
		// 		payload.WriteInt16(x) // X
		// 		payload.WriteInt16(y) // Y
		// 		payloadBytes := payload.Data()

		// 		s.QueueSendMHF(&mhfpacket.MsgSysCastedBinary{
		// 			CharID:         s.CharID,
		// 			MessageType:    BinaryMessageTypeState,
		// 			RawDataPayload: payloadBytes,
		// 		})
		// 	}
		// }
	}
}

func handleMsgSysCastedBinary(s *Session, p mhfpacket.MHFPacket) {}
