package channelserver

import (
	"encoding/hex"
	"erupe-ce/network/mhfpacket"
)

func handleMsgMhfGetUdTacticsPoint(s *Session, p mhfpacket.MHFPacket) {
	// Diva defense interception points
	pkt := p.(*mhfpacket.MsgMhfGetUdTacticsPoint)
	// Temporary canned response
	data, _ := hex.DecodeString("000000A08F0BE2DAE30BE30AE2EAE2E9E2E8E2F5E2F3E2F2E2F1E2BB")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfAddUdTacticsPoint(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfAddUdTacticsPoint)
	stubEnumerateNoResults(s, pkt.AckHandle)
}

func handleMsgMhfGetUdTacticsRewardList(s *Session, p mhfpacket.MHFPacket) {
	// Diva defense interception
	pkt := p.(*mhfpacket.MsgMhfGetUdTacticsRewardList)
	// Temporary canned response
	data, _ := hex.DecodeString("000094000000010732DD00010000000000010732DD00010100000000C8071F2800050100000000C80705C000050000000001901A000001F40000000001901A000001F40100000002580705C00005000000000258071F2800050100000003201A000003E80100000003201A000003E80000000003E81A000004B00100000003E81A000004B00000000004B01A000005DC0100000004B01A000005DC0000000005781A000008FC0100000005781A000008FC0000000006401A000009C40000000006401A000009C40100000007081A00000BB80100000007081A00000BB80000000007D00725FA00010000000007D01A00000CE40000000007D00725FC00010100000007D00725FB00010100000007D00725FA00010100000007D01A00000CE40100000007D00725FC00010000000007D00725FB0001000000000BB80705C00005000000000BB8071F280005010000000FA01A00000DAC000000000FA01A00000DAC0100000013880705C00005000000001388071F2800050100000017700725FE00010100000017700725FD00010100000017700725FF00010100000017700725FD00010000000017700725FE00010000000017700725FF0001000000001B581A00000E74000000001B581A00000E74010000001F400727D00005010000001F400727D000050000000023281A00000FA00000000023281A00000FA00100000027100736EF000100000000271007369600010100000027100736EF00010100000027100736EF0001000000002EE00727D10005010000002EE00727D100050000000036B01D000000010100000036B01D00000001000000003A980737DB0001010000003A980736EF00010000000046500725E600010100000046500725E60001000000004E200738C90001010000004E200736EF00010000000055F01A000010680100000055F01A000010680000000061A80736EF00010000000061A80739A600010100000065900727D200050000000065900727D20005010000007530073A0600010100000075300736EF00010000000075300736EF00010000000075300736EF00010100000084D01D000000020000000084D01D00000002010000009C400727D30005010000009C400727D3000500000000B3B01A0000119400000000B3B01A0000119401000000C3500727D4000500000000C3500727D4000501000000D2F01D0000000300000000D2F01D0000000301000000EA600736EF000100000000EA600736EF000101000000F6181A0000125C00000000F6181A0000125C0100000111700727D500050000000111700727D500050100000119400727D600050100000119400727D600050000000121101D000000040000000121101D000000040100000130B01A000013880000000130B01A000013880100000140500727D700050000000140500727D700050100000148201D000000050000000148201D00000005010000014FF01A000014B4000000014FF01A000014B4010000015F900736EF0001000000015F900736EF00010100000167600729EA00050000000167600729EA0005010000016F301D00000006010000016F301D00000006000000017ED00729EB0005000000017ED00729EB0005010000018E701A0000157C010000018E701A0000157C0000000196401D000000070000000196401D00000007010000019E100729EC0005000000019E100729EC000501000001ADB00727CD000100000001ADB00727CD000101000001BD501D0000000800000001BD501D0000000801000001CCF01A0000164401000001CCF01A0000164400000001E4601D0000000901000001E4601D0000000900000001EC300727CC000101000001EC300727CC0001000000020B701D0000000A000000020B701D0000000A010000023A501A0000170C010000023A501A0000170C0000000249F00736EF00010100000249F00736EF00010000000271001A000017D40100000271001A000017D400000002A7B01A0000189C01000002A7B01A0000189C00000002BF200736EF000100000002BF200736EF000101000002D6901A0000196401000002D6901A00001964000000030D400727CB0001000000030D400727CB00010100000343F01A00001A2C0100000343F01A00001A2C0000000372D0072CB0000F0000000372D0072CB0000F01000003A9801A00001BBC00000003A9801A00001BBC01000003F7A01A000003E800010003F7A01A000003E80101000445C01A000003E80101000445C01A000003E80001005E000000020704020005010000000002070402000500000000000307040200140000000000030704020014010000000005071D200003010000000005071D20000300000000000607040200140100000000060704020014000000000008071D210003010000000008071D21000300000000000A070402001401000000000A070402001400000000000C0722EC000501000000000C0722ED000500000000000C0722F2000500000000000C0722EC000500000000000C0722EF000500000000000C0722ED000501000000000C0722F2000501000000000C0722EF000501000000000D1A000003E801000000000D1A000003E800000000000F07357C000501000000000F07357D000501000000000F07357C000500000000000F07357D00050000000000111A000007D00000000000111A000007D00100000000141C00000001000000000014071D2200030000000000141C00000001010000000014071D22000301000000001607357D000701000000001607357C00070000000000160704020028000000000016070402002801000000001607357C000701000000001607357D0007000000000018071D270003000000000018071D27000301000000001A1A00000BB800000000001A1A00000BB801000000001C07357D000701000000001C070402002801000000001C07357D000700000000001C07357C000700000000001C070402002800000000001C07357C000701000000001E070402003C01000000001E070402003C000000000020071D26000301000000002007357C000700000000002007357D000700000000002007357C000701000000002007357D0007010000000020071D260003000000000023071D280003010000000023071D28000300000000002A070402003C00000000002A070402003C01000000002C0725EE000100000000002C0725EE000101000000002E070402005001000000002E07357D000A01000000002E070402005000000000002E07357C000A00000000002E07357D000A00000000002E07357C000A0100000000300725ED00010000000000300725ED0001010000000032071D200003010000000032071D200003000000000034072C7B0001000000000034072C7B0001010000000037071D210003000000000037071D21000301000000003C0722F1000A00000000003C0722F1000A01000000004107040200500000000000410704020050010000000046071D220003010000000046071D22000300000000004B071D27000301000000004B071D2700030000000000500722F1000F0100000000500722F1000F0000000000550704020050010000000055070402005000000000005A071D26000301000000005A071D26000300000000005F071D28000300000000005F071D2800030100000000641A0000C3500100000000641A0000C3500000002607000E00C8000000010000000307000F0032000000010000000307001000320000000100000003070011003200000001000000030700120032000000010000000307000E0096000000040000000A07000F0028000000040000000A0700100028000000040000000A0700110028000000040000000A0700120028000000040000000A07000E00640000000B0000001907000F001E0000000B00000019070010001E0000000B00000019070011001E0000000B00000019070012001E0000000B0000001907000E00320000001A0000002807000F00140000001A0000002807001000140000001A0000002807001100140000001A0000002807001200140000001A0000002807000E001E000000290000004607000F000A0000002900000046070010000A000000290000004607001100010000002900000046070012000A000000290000004607000E0019000000470000006407000F0008000000470000006407001000080000004700000064070011000100000047000000640700120008000000470000006407000E000F000000650000009607000F0006000000650000009607001000010000006500000096070011000600000065000000960700120006000000650000009607000E000500000097000001F407000F000500000097000001F4070010000500000097000001F4")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfGetUdTacticsFollower(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdTacticsFollower)
	doAckBufSucceed(s, pkt.AckHandle, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func handleMsgMhfGetUdTacticsBonusQuest(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdTacticsBonusQuest)
	// Temporary canned response
	data, _ := hex.DecodeString("14E2F55DCBFE505DCC1A7003E8E2C55DCC6ED05DCC8AF00258E2CE5DCCDF505DCCFB700279E3075DCD4FD05DCD6BF0041AE2F15DCDC0505DCDDC700258E2C45DCE30D05DCE4CF00258E2F55DCEA1505DCEBD7003E8E2C25DCF11D05DCF2DF00258E2CE5DCF82505DCF9E700279E3075DCFF2D05DD00EF0041AE2CE5DD063505DD07F700279E2F35DD0D3D05DD0EFF0028AE2C35DD144505DD160700258E2F05DD1B4D05DD1D0F00258E2CE5DD225505DD241700279E2F55DD295D05DD2B1F003E8E2F25DD306505DD3227002EEE2CA5DD376D05DD392F00258E3075DD3E7505DD40370041AE2F55DD457D05DD473F003E82027313220686F757273273A3A696E74657276616C29202B2027313220686F757273273A3A696E74657276616C2047524F5550204259206D6170204F52444552204259206D61703B2000C7312B000032")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfGetUdTacticsFirstQuestBonus(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdTacticsFirstQuestBonus)
	// Temporary canned response
	data, _ := hex.DecodeString("0500000005DC01000007D002000009C40300000BB80400001194")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfGetUdTacticsRemainingPoint(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsRanking(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfGetUdTacticsRanking)
	// Temporary canned response
	data, _ := hex.DecodeString("00000515000005150000CEB4000003CE000003CE0000CEB44D49444E494748542D414E47454C0000000000000000000000")
	doAckBufSucceed(s, pkt.AckHandle, data)
}

func handleMsgMhfSetUdTacticsFollower(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgMhfGetUdTacticsLog(s *Session, p mhfpacket.MHFPacket) {}
