package channelserver

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"erupe-ce/common/byteframe"
	"erupe-ce/common/stringstack"
	"erupe-ce/network"
	"erupe-ce/network/binpacket"
	"erupe-ce/network/clientctx"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
)

type packet struct {
	data        []byte
	nonBlocking bool
}

// Session holds state for the channel server connection.
type Session struct {
	sync.Mutex
	logger        *zap.Logger
	hasLoggerName bool
	Server        *Server `json:"-"`
	rawConn       net.Conn
	cryptConn     *network.CryptConn
	sendPackets   chan packet
	clientContext *clientctx.ClientContext

	UserEnteredStage bool     `json:"userEnteredStage"` // If the user has entered a stage before
	MySeries         MySeries `json:"mySeries"`
	StageID          string   `json:"stageID"`
	Stage            *Stage   `json:"-"`
	ReservationStage *Stage   `json:"-"`           // Required for the stateful MsgSysUnreserveStage packet.
	StagePass        string   `json:"stagePass"`   // Temporary storage
	PrevGuildID      uint32   `json:"prevGuildID"` // Stores the last GuildID used in InfoGuild
	CharID           uint32   `json:"charID"`
	Dev              bool     `json:"dev"`
	LogKey           []byte   `json:"logKey"`
	SessionStart     int64    `json:"sessionStart"`
	Rights           uint32   `json:"rights"`
	courses          []mhfpacket.Course
	Token            string `json:"token"`
	kqf              []byte
	kqfOverride      bool
	Semaphore        *Semaphore `json:"Semaphore"` // Required for the stateful MsgSysUnreserveStage packet.

	// A stack containing the stage movement history (push on enter/move, pop on back)
	stageMoveStack *stringstack.StringStack

	// Accumulated index used for identifying mail for a client
	// I'm not certain why this is used, but since the client is sending it
	// I want to rely on it for now as it might be important later.
	MailAccIndex uint8 `json:"mailAccIndex"`
	// Contains the mail list that maps accumulated indexes to mail IDs
	MailList []int `json:"mailList"`

	// For Debuging
	Name   string `json:"name"`
	closed bool
}

type MySeries struct {
	HouseTier     []byte `json:"houseTier"`
	HouseData     []byte `json:"houseData"`
	BookshelfData []byte `json:"bookshelfData"`
	GalleryData   []byte `json:"galleryData"`
	ToreData      []byte `json:"toreData"`
	GardenData    []byte `json:"gardenData"`
	State         uint8  `json:"state"`
	Password      string `json:"password"`
}

// NewSession creates a new Session type.
func NewSession(server *Server, conn net.Conn) *Session {
	s := &Session{
		logger:           server.logger.Named(conn.RemoteAddr().String()),
		Server:           server,
		rawConn:          conn,
		hasLoggerName:    false,
		Dev:              false,
		cryptConn:        network.NewCryptConn(conn),
		sendPackets:      make(chan packet, 20),
		clientContext:    &clientctx.ClientContext{}, // Unused
		UserEnteredStage: false,
		SessionStart:     Time_Current_Adjusted().Unix(),
		stageMoveStack:   stringstack.New(),
	}
	return s
}

// Start starts the session packet send and recv loop(s).
func (s *Session) Start() {
	go func() {
		s.logger.Info("Channel server got connection!", zap.String("remoteaddr", s.rawConn.RemoteAddr().String()))
		// Unlike the sign and entrance server,
		// the client DOES NOT initalize the channel connection with 8 NULL bytes.
		go s.sendLoop()
		s.recvLoop()
	}()
}

// QueueSend queues a packet (raw []byte) to be sent.
func (s *Session) QueueSend(data []byte) {
	s.logMessage(binary.BigEndian.Uint16(data[0:2]), data, "Server", s.Name)
	select {
	case s.sendPackets <- packet{data, false}:
		// Enqueued data
	default:
		s.logger.Warn("Packet queue too full, flushing!")
		var tempPackets []packet
		for len(s.sendPackets) > 0 {
			tempPacket := <-s.sendPackets
			if !tempPacket.nonBlocking {
				tempPackets = append(tempPackets, tempPacket)
			}
		}
		for _, tempPacket := range tempPackets {
			s.sendPackets <- tempPacket
		}
		s.sendPackets <- packet{data, false}
	}
}

func (s *Session) SetLoggerName(name string) {
	if !s.hasLoggerName {
		s.hasLoggerName = true
		s.logger = s.logger.Named(name)
	}
}

// QueueSendNonBlocking queues a packet (raw []byte) to be sent, dropping the packet entirely if the queue is full.
func (s *Session) QueueSendNonBlocking(data []byte) {
	select {
	case s.sendPackets <- packet{data, true}:
		s.logMessage(binary.BigEndian.Uint16(data[0:2]), data, "Server", s.Name)
	default:
		s.logger.Warn("Packet queue too full, dropping!")
	}
}

// QueueSendMHF queues a MHFPacket to be sent.
func (s *Session) QueueSendMHF(pkt mhfpacket.MHFPacket) {
	// Make the header
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(pkt.Opcode()))

	// Build the packet onto the byteframe.
	err := pkt.Build(bf, s.clientContext)

	if err != nil {
		fmt.Printf("deu ruim2: %s - %s", pkt.Opcode(), err.Error())

	}

	// Queue it.
	s.QueueSend(bf.Data())
}

// QueueAck is a helper function to queue an MSG_SYS_ACK with the given ack handle and data.
func (s *Session) QueueAck(ackHandle uint32, data []byte) {
	bf := byteframe.NewByteFrame()
	bf.WriteUint16(uint16(network.MSG_SYS_ACK))
	bf.WriteUint32(ackHandle)
	bf.WriteBytes(data)
	s.QueueSend(bf.Data())
}

func (s *Session) sendLoop() {
	for {
		if s.closed {
			return
		}
		pkt := <-s.sendPackets
		err := s.cryptConn.SendPacket(append(pkt.data, []byte{0x00, 0x10}...))
		if err != nil {
			s.logger.Warn("Failed to send packet")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *Session) recvLoop() {
	for {
		if s.closed {
			LogoutPlayer(s)
			return
		}
		pkt, err := s.cryptConn.ReadPacket()
		if err == io.EOF {
			s.logger.Info(fmt.Sprintf("[%s] Disconnected", s.Name))
			LogoutPlayer(s)
			return
		}
		if err != nil {
			s.logger.Warn("Error on ReadPacket, exiting recv loop", zap.Error(err))
			LogoutPlayer(s)
			return
		}
		s.handlePacketGroup(pkt)
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *Session) handlePacketGroup(pktGroup []byte) {
	bf := byteframe.NewByteFrameFromBytes(pktGroup)
	opcodeUint16 := bf.ReadUint16()
	opcode := network.PacketID(opcodeUint16)

	// This shouldn't be needed, but it's better to recover and let the connection die than to panic the server.
	// defer func() {
	// 	if r := recover(); r != nil {

	// 		fmt.Printf("[%s] ", s.Name)
	// 		fmt.Println("Recovered from panic", r)

	// 	}
	// }()

	s.logMessage(opcodeUint16, pktGroup, s.Name, "Server")

	if opcode == network.MSG_SYS_LOGOUT {
		s.closed = true
		return
	}
	// Get the packet parser and handler for this opcode.
	mhfPkt := mhfpacket.FromOpcode(opcode)
	if mhfPkt == nil {
		fmt.Println("Got opcode which we don't know how to parse, can't parse anymore for this group")
		return
	}
	// Parse the packet.
	err := mhfPkt.Parse(bf, s.clientContext)
	if err != nil {
		fmt.Printf("\n!!! [%s] %s NOT IMPLEMENTED !!! \n\n\n", s.Name, opcode)
		return
	}
	// Handle the packet.
	handlerTable[opcode](s, mhfPkt)
	// If there is more data on the stream that the .Parse method didn't read, then read another packet off it.
	remainingData := bf.DataFromCurrent()
	if len(remainingData) >= 2 {
		s.handlePacketGroup(remainingData)
	}
}

func ignored(opcode network.PacketID) bool {
	ignoreList := []network.PacketID{
		network.MSG_SYS_END,
		network.MSG_SYS_PING,
		network.MSG_SYS_NOP,
		network.MSG_SYS_TIME,
		network.MSG_SYS_EXTEND_THRESHOLD,
		network.MSG_SYS_POSITION_OBJECT,
		network.MSG_MHF_ENUMERATE_QUEST,
		network.MSG_MHF_SAVEDATA,
	}
	set := make(map[network.PacketID]struct{}, len(ignoreList))
	for _, s := range ignoreList {
		set[s] = struct{}{}
	}
	_, r := set[opcode]
	return r
}

func (s *Session) logMessage(opcode uint16, data []byte, sender string, recipient string) {
	if !s.Server.erupeConfig.DevMode {
		return
	}

	if sender == "Server" && !s.Server.erupeConfig.DevModeOptions.LogOutboundMessages {
		return
	} else if !s.Server.erupeConfig.DevModeOptions.LogInboundMessages {
		return
	}

	opcodePID := network.PacketID(opcode)
	if ignored(opcodePID) {
		return
	}
	fmt.Printf("[%s] -> [%s]\n", sender, recipient)
	fmt.Printf("Opcode: %s\n", opcodePID)
	if len(data) <= s.Server.erupeConfig.DevModeOptions.MaxHexdumpLength {
		fmt.Printf("Data [%d bytes]:\n%s\n", len(data), hex.Dump(data))
	} else {
		fmt.Printf("Data [%d bytes]:\n(Too long!)\n\n", len(data))
	}
}

// BroadcastChatMessage broadcasts a simple chat message to all the sessions.
func (s *Session) SendMessage(SenderName string, message string) {
	// Make the inside of the casted binary
	bf := byteframe.NewByteFrame()
	bf.SetLE()
	msgBinChat := &binpacket.MsgBinChat{
		Unk0:       0,
		Type:       5,
		Flags:      0x80,
		Message:    message,
		SenderName: SenderName,
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

func (s *Session) FindCourse(name string) mhfpacket.Course {
	for _, course := range s.courses {
		for _, alias := range course.Aliases {
			if strings.EqualFold(strings.ToLower(name), strings.ToLower(alias)) {
				return course
			}
		}
	}
	return mhfpacket.Course{}
}
