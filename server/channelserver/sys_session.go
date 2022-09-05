package channelserver

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"

	"erupe-ce/common/byteframe"
	"erupe-ce/common/stringstack"
	"erupe-ce/common/stringsupport"
	"erupe-ce/network"
	"erupe-ce/network/clientctx"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
)

// Session holds state for the channel server connection.
type Session struct {
	sync.Mutex
	logger        *zap.Logger
	hasLoggerName bool
	Server        *Server `json:"-"`
	rawConn       net.Conn
	cryptConn     *network.CryptConn
	sendPackets   chan []byte
	clientContext *clientctx.ClientContext

	UserEnteredStage bool     `json:"userEnteredStage"` // If the user has entered a stage before
	MySeries         MySeries `json:"mySeries"`
	StageID          string   `json:"stageID"`
	Stage            *Stage   `json:"stage"`
	ReservationStage *Stage   `json:"reservationStage"` // Required for the stateful MsgSysUnreserveStage packet.
	StagePass        string   `json:"stagePass"`        // Temporary storage
	PrevGuildID      uint32   `json:"prevGuildID"`      // Stores the last GuildID used in InfoGuild
	CharID           uint32   `json:"charID"`
	LogKey           []byte   `json:"logKey"`
	SessionStart     int64    `json:"sessionStart"`
	Rights           uint32   `json:"rights"`
	Token            string   `json:"token"`

	Semaphore *Semaphore `json:"Semaphore"` // Required for the stateful MsgSysUnreserveStage packet.

	// A stack containing the stage movement history (push on enter/move, pop on back)
	stageMoveStack *stringstack.StringStack

	// Accumulated index used for identifying mail for a client
	// I'm not certain why this is used, but since the client is sending it
	// I want to rely on it for now as it might be important later.
	MailAccIndex uint8 `json:"mailAccIndex"`
	// Contains the mail list that maps accumulated indexes to mail IDs
	MailList []int `json:"mailList"`

	// For Debuging
	Name string `json:"Name"`
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
		logger:        server.logger.Named(conn.RemoteAddr().String()),
		Server:        server,
		rawConn:       conn,
		hasLoggerName: false,
		cryptConn:     network.NewCryptConn(conn),
		sendPackets:   make(chan []byte, 20),
		clientContext: &clientctx.ClientContext{
			StrConv: &stringsupport.StringConverter{
				Encoding: japanese.ShiftJIS,
			},
		},
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
	bf := byteframe.NewByteFrameFromBytes(data[:2])
	s.logMessage(bf.ReadUint16(), data, "Server", s.Name)
	s.sendPackets <- data
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
	case s.sendPackets <- data:
		// Enqueued properly.
	default:
		// Couldn't enqueue, likely something wrong with the connection.
		s.logger.Warn("Dropped packet for session because of full send buffer, something is probably wrong")
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
		// TODO(Andoryuuta): Test making this into a buffered channel and grouping the packet together before sending.
		rawPacket := <-s.sendPackets

		if rawPacket == nil {
			s.logger.Debug("Got nil from s.SendPackets, exiting send loop")
			return
		}

		// Make a copy of the data.
		terminatedPacket := make([]byte, len(rawPacket))
		copy(terminatedPacket, rawPacket)

		// Append the MSG_SYS_END tailing opcode.
		terminatedPacket = append(terminatedPacket, []byte{0x00, 0x10}...)
		// nolint:errcheck
		s.cryptConn.SendPacket(terminatedPacket)
	}
}

func (s *Session) recvLoop() {
	for {
		pkt, err := s.cryptConn.ReadPacket()

		if err == io.EOF {
			s.logger.Info(fmt.Sprintf("[%s] Disconnected", s.Name))
			logoutPlayer(s)
			return
		}
		if err != nil {
			s.logger.Warn("Error on ReadPacket, exiting recv loop", zap.Error(err))
			logoutPlayer(s)
			return
		}
		s.handlePacketGroup(pkt)
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
		s.rawConn.Close()
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
