package channelserver

import (
	"encoding/hex"
	"erupe-ce/common/byteframe"
	"erupe-ce/common/mhfcourse"
	"erupe-ce/common/token"
	"erupe-ce/config"
	"erupe-ce/network/binpacket"
	"erupe-ce/network/mhfpacket"
	"fmt"
	"math"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"go.uber.org/zap"
)

// MSG_SYS_CAST[ED]_BINARY types enum
const (
	BinaryMessageTypeState      = 0
	BinaryMessageTypeChat       = 1
	BinaryMessageTypeQuest      = 2
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

var commands map[string]config.Command

func init() {
	commands = make(map[string]config.Command)
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableCaller = true
	zapLogger, _ := zapConfig.Build()
	// nolint:errcheck
	defer zapLogger.Sync()
	logger := zapLogger.Named("commands")
	cmds := config.ErupeConfig.Commands
	for _, cmd := range cmds {
		commands[cmd.Name] = cmd
		if cmd.Enabled {
			logger.Info(fmt.Sprintf("Command %s: Enabled, prefix: %s", cmd.Name, cmd.Prefix))
		} else {
			logger.Info(fmt.Sprintf("Command %s: Disabled", cmd.Name))
		}
	}
}

func sendDisabledCommandMessage(s *Session, cmd config.Command) {
	sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandDisabled"], cmd.Name))
}

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

func parseChatCommand(s *Session, command string) {
	if strings.HasPrefix(command, commands["PSN"].Prefix) {
		if commands["PSN"].Enabled {
			var id string
			n, err := fmt.Sscanf(command, fmt.Sprintf("%s %%s", commands["PSN"].Prefix), &id)
			if err != nil || n != 1 {
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandPSNError"], commands["PSN"].Prefix))
			} else {
				_, err = s.Server.db.Exec(`UPDATE users u SET psn_id=$1 WHERE u.id=(SELECT c.user_id FROM characters c WHERE c.id=$2)`, id, s.CharID)
				if err == nil {
					sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandPSNSuccess"], id))
				}
			}
		}
	}

	if strings.HasPrefix(command, commands["Reload"].Prefix) {
		// Flush all objects and users and reload
		if commands["Reload"].Enabled {
			sendServerChatMessage(s, s.Server.dict["commandReload"])
			var temp mhfpacket.MHFPacket
			deleteNotif := byteframe.NewByteFrame()
			for _, object := range s.Stage.Objects {
				if object.OwnerCharID == s.CharID {
					continue
				}
				temp = &mhfpacket.MsgSysDeleteObject{ObjID: object.Id}
				deleteNotif.WriteUint16(uint16(temp.Opcode()))
				temp.Build(deleteNotif, s.clientContext)
			}
			for _, session := range s.Server.Sessions {
				if s == session {
					continue
				}
				temp = &mhfpacket.MsgSysDeleteUser{CharID: session.CharID}
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
		} else {
			sendDisabledCommandMessage(s, commands["Reload"])
		}
	}

	if strings.HasPrefix(command, commands["KeyQuest"].Prefix) {
		if commands["KeyQuest"].Enabled {
			if strings.HasPrefix(command, fmt.Sprintf("%s get", commands["KeyQuest"].Prefix)) {
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandKqfGet"], s.kqf))
			} else if strings.HasPrefix(command, fmt.Sprintf("%s set", commands["KeyQuest"].Prefix)) {
				var hexs string
				n, numerr := fmt.Sscanf(command, fmt.Sprintf("%s set %%s", commands["KeyQuest"].Prefix), &hexs)
				if numerr != nil || n != 1 || len(hexs) != 16 {
					sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandKqfSetError"], commands["KeyQuest"].Prefix))
				} else {
					hexd, _ := hex.DecodeString(hexs)
					s.kqf = hexd
					s.kqfOverride = true
					sendServerChatMessage(s, s.Server.dict["commandKqfSetSuccess"])
				}
			}
		} else {
			sendDisabledCommandMessage(s, commands["KeyQuest"])
		}
	}

	if strings.HasPrefix(command, commands["Rights"].Prefix) {
		// Set account rights
		if commands["Rights"].Enabled {
			var v uint32
			n, err := fmt.Sscanf(command, fmt.Sprintf("%s %%d", commands["Rights"].Prefix), &v)
			if err != nil || n != 1 {
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandRightsError"], commands["Rights"].Prefix))
			} else {
				_, err = s.Server.db.Exec("UPDATE users u SET rights=$1 WHERE u.id=(SELECT c.user_id FROM characters c WHERE c.id=$2)", v, s.CharID)
				if err == nil {
					sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandRightsSuccess"], v))
				}
			}
		} else {
			sendDisabledCommandMessage(s, commands["Rights"])
		}
	}

	if strings.HasPrefix(command, commands["Course"].Prefix) {
		if commands["Course"].Enabled {
			var name string
			n, err := fmt.Sscanf(command, fmt.Sprintf("%s %%s", commands["Course"].Prefix), &name)
			if err != nil || n != 1 {
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandCourseError"], commands["Course"].Prefix))
			} else {
				name = strings.ToLower(name)
				for _, course := range mhfcourse.Courses() {
					for _, alias := range course.Aliases() {
						if strings.ToLower(name) == strings.ToLower(alias) {
							if slices.Contains(s.Server.erupeConfig.Courses, config.Course{Name: course.Aliases()[0], Enabled: true}) {
								var delta, rightsInt uint32
								if mhfcourse.CourseExists(course.ID, s.courses) {
									ei := slices.IndexFunc(s.courses, func(c mhfcourse.Course) bool {
										for _, alias := range c.Aliases() {
											if strings.ToLower(name) == strings.ToLower(alias) {
												return true
											}
										}
										return false
									})
									if ei != -1 {
										delta = uint32(-1 * math.Pow(2, float64(course.ID)))
										sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandCourseDisabled"], course.Aliases()[0]))
									}
								} else {
									delta = uint32(math.Pow(2, float64(course.ID)))
									sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandCourseEnabled"], course.Aliases()[0]))
								}
								err = s.Server.db.QueryRow("SELECT rights FROM users u INNER JOIN characters c ON u.id = c.user_id WHERE c.id = $1", s.CharID).Scan(&rightsInt)
								if err == nil {
									s.Server.db.Exec("UPDATE users u SET rights=$1 WHERE u.id=(SELECT c.user_id FROM characters c WHERE c.id=$2)", rightsInt+delta, s.CharID)
								}
								updateRights(s)
							} else {
								sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandCourseLocked"], course.Aliases()[0]))
							}
							return
						}
					}
				}
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandCourseError"], commands["Course"].Prefix))
			}
		} else {
			sendDisabledCommandMessage(s, commands["Course"])
		}
	}

	if strings.HasPrefix(command, commands["Raviente"].Prefix) {
		if commands["Raviente"].Enabled {
			if getRaviSemaphore(s.Server) != nil {
				s.Server.raviente.Lock()
				if !strings.HasPrefix(command, "!ravi ") {
					sendServerChatMessage(s, s.Server.dict["commandRaviNoCommand"])
				} else {
					if strings.HasPrefix(command, "!ravi start") {
						if s.Server.raviente.register.startTime == 0 {
							s.Server.raviente.register.startTime = s.Server.raviente.register.postTime
							sendServerChatMessage(s, s.Server.dict["commandRaviStartSuccess"])
							s.notifyRavi()
						} else {
							sendServerChatMessage(s, s.Server.dict["commandRaviStartError"])
						}
					} else if strings.HasPrefix(command, "!ravi cm") || strings.HasPrefix(command, "!ravi checkmultiplier") {
						sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandRaviMultiplier"], s.Server.raviente.GetRaviMultiplier(s.Server)))
					} else if strings.HasPrefix(command, "!ravi sr") || strings.HasPrefix(command, "!ravi sendres") {
						if s.Server.raviente.state.stateData[28] > 0 {
							sendServerChatMessage(s, s.Server.dict["commandRaviResSuccess"])
							s.Server.raviente.state.stateData[28] = 0
						} else {
							sendServerChatMessage(s, s.Server.dict["commandRaviResError"])
						}
					} else if strings.HasPrefix(command, "!ravi ss") || strings.HasPrefix(command, "!ravi sendsed") {
						sendServerChatMessage(s, s.Server.dict["commandRaviSedSuccess"])
						// Total BerRavi HP
						HP := s.Server.raviente.state.stateData[0] + s.Server.raviente.state.stateData[1] + s.Server.raviente.state.stateData[2] + s.Server.raviente.state.stateData[3] + s.Server.raviente.state.stateData[4]
						s.Server.raviente.support.supportData[1] = HP
					} else if strings.HasPrefix(command, "!ravi rs") || strings.HasPrefix(command, "!ravi reqsed") {
						sendServerChatMessage(s, s.Server.dict["commandRaviRequest"])
						// Total BerRavi HP
						HP := s.Server.raviente.state.stateData[0] + s.Server.raviente.state.stateData[1] + s.Server.raviente.state.stateData[2] + s.Server.raviente.state.stateData[3] + s.Server.raviente.state.stateData[4]
						s.Server.raviente.support.supportData[1] = HP + 12
					} else {
						sendServerChatMessage(s, s.Server.dict["commandRaviError"])
					}
				}
				s.Server.raviente.Unlock()
			} else {
				sendServerChatMessage(s, s.Server.dict["commandRaviNoPlayers"])
			}
		} else {
			sendDisabledCommandMessage(s, commands["Raviente"])
		}
	}

	if strings.HasPrefix(command, commands["Teleport"].Prefix) {
		if commands["Teleport"].Enabled {
			var x, y int16
			n, err := fmt.Sscanf(command, fmt.Sprintf("%s %%d %%d", commands["Teleport"].Prefix), &x, &y)
			if err != nil || n != 2 {
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandTeleportError"], commands["Teleport"].Prefix))
			} else {
				sendServerChatMessage(s, fmt.Sprintf(s.Server.dict["commandTeleportSuccess"], x, y))

				// Make the inside of the casted binary
				payload := byteframe.NewByteFrame()
				payload.SetLE()
				payload.WriteUint8(2) // SetState type(position == 2)
				payload.WriteInt16(x) // X
				payload.WriteInt16(y) // Y
				payloadBytes := payload.Data()

				s.QueueSendMHF(&mhfpacket.MsgSysCastedBinary{
					CharID:         s.CharID,
					MessageType:    BinaryMessageTypeState,
					RawDataPayload: payloadBytes,
				})
			}
		} else {
			sendDisabledCommandMessage(s, commands["Teleport"])
		}
	}
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

			if s.Stage != nil {
				fmt.Printf("Quest Finished: %s\n", s.Stage.QuestFilename)
				SaveQuestRecord(s, frame, "COMPLETED")
			}
		}
	}

	if s.Server.erupeConfig.DevModeOptions.QuestDebugTools && s.Server.erupeConfig.DevMode {
		if pkt.BroadcastType == 0x03 && pkt.MessageType == 0x02 && len(pkt.RawDataPayload) > 32 {
			// This is only correct most of the time
			tmp.ReadBytes(20)
			tmp.SetLE()
			x := tmp.ReadFloat32()
			y := tmp.ReadFloat32()
			z := tmp.ReadFloat32()
			s.logger.Debug("Coord", zap.Float32s("XYZ", []float32{x, y, z}))
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
			dice := fmt.Sprintf("%d", token.RNG().Intn(100)+1)
			roll.WriteUint16(uint16(len(dice) + 1))
			roll.WriteNullTerminatedBytes([]byte(dice))
			roll.WriteNullTerminatedBytes(tmp.ReadNullTerminatedBytes())
			realPayload = roll.Data()
		} else {
			bf := byteframe.NewByteFrameFromBytes(pkt.RawDataPayload)
			bf.SetLE()
			chatMessage := &binpacket.MsgBinChat{}
			chatMessage.Parse(bf)
			if strings.HasPrefix(chatMessage.Message, "!") {
				parseChatCommand(s, chatMessage.Message)
				return
			}
			// Discord integration
			if (pkt.BroadcastType == BroadcastTypeStage && s.Stage.Id == "sl1Ns200p0a0u0") || pkt.BroadcastType == BroadcastTypeWorld {
				s.Server.DiscordChannelSend(chatMessage.SenderName, chatMessage.Message)
			}
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
			if getRaviSemaphore(s.Server) != nil {
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
}

func handleMsgSysCastedBinary(s *Session, p mhfpacket.MHFPacket) {}
