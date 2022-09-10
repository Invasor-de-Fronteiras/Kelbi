package channelserver

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Character struct {
	ID             uint32 `db:"id"`
	IsFemale       bool   `db:"is_female"`
	IsNewCharacter bool   `db:"is_new_character"`
	Name           string `db:"name"`
	UnkDescString  string `db:"unk_desc_string"`
	HRP            uint16 `db:"hrp"`
	GR             uint16 `db:"gr"`
	WeaponType     uint16 `db:"weapon_type"`
	LastLogin      uint32 `db:"last_login"`
}

var weapons = []string{
	"<:gs:970861408227049477>",
	"<:hbg:970861408281563206>",
	"<:hm:970861408239628308>",
	"<:lc:970861408298352660>",
	"<:sns:970861408319315988>",
	"<:lbg:970861408327725137>",
	"<:ds:970861408277368883>",
	"<:ls:970861408319311872>",
	"<:hh:970861408222863400>",
	"<:gl:970861408327720980>",
	"<:bw:970861408294174780>",
	"<:tf:970861408424177744>",
	"<:sw:970861408340283472>",
	"<:ms:970861408411594842>",
}

func (s *Server) getCharacterForUser(uid int) (*Character, error) {
	character := Character{}
	err := s.db.Get(&character, "SELECT id, is_female, is_new_character, name, unk_desc_string, hrp, gr, weapon_type, last_login  FROM characters WHERE id = $1", uid)
	if err != nil {
		return nil, err
	}

	return &character, nil
}

func CountChars(s *Server) string {
	count := 0
	for _, stage := range s.Stages {
		count += len(stage.Clients)
	}

	message := fmt.Sprintf("Server [%s]: %d players;", s.Name, count)

	return message
}

type ListPlayer struct {
	CharName    string
	InQuest     bool
	WeaponEmoji string
	QuestEmoji  string
	StageName   string
}

type CharInfo struct {
	CharID    uint32
	CharName  string
	StageId   string
	StageName string
	IP        string
}

func (p *ListPlayer) toString(length int) string {
	status := ""
	if p.InQuest {
		status = "(in quest " + p.QuestEmoji + ")"
	} else {
		status = p.StageName
	}

	missingSpace := length - len(p.CharName)
	whitespaces := strings.Repeat(" ", missingSpace+5)

	return fmt.Sprintf("%s %s %s %s", p.WeaponEmoji, p.CharName, whitespaces, status)
}

func getPlayerList(s *Server) ([]ListPlayer, int) {
	list := []ListPlayer{}
	questEmojis := []string{
		":white_circle:",
		":red_circle:",
		":blue_circle:",
		":brown_circle:",
		":green_circle:",
		":purple_circle:",
		":yellow_circle:",
		":orange_circle:",
	}

	bigNameLen := 0

	for _, stage := range s.Stages {
		if len(stage.Clients) == 0 {
			continue
		}

		questEmoji := ":black_circle:"

		if len(questEmojis) > 0 {
			questEmoji = questEmojis[len(questEmojis)-1]
			questEmojis = questEmojis[:len(questEmojis)-1]
		}

		isQuest := stage.isQuest()
		for client := range stage.Clients {
			char, err := s.getCharacterForUser(int(client.CharID))
			if err == nil {
				if len(char.Name) > bigNameLen {
					bigNameLen = len(char.Name)
				}

				list = append(list, ListPlayer{
					CharName:    char.Name,
					InQuest:     isQuest,
					QuestEmoji:  questEmoji,
					WeaponEmoji: weapons[char.WeaponType],
					StageName:   stage.GetName(),
				})

			}
		}
	}

	return list, bigNameLen
}

func PlayerList(s *Server) string {
	list := ""
	count := 0
	listPlayers, bigNameLen := getPlayerList(s)

	sort.SliceStable(listPlayers, func(i, j int) bool {
		return listPlayers[i].CharName < listPlayers[j].CharName
	})

	for _, lp := range listPlayers {
		list += lp.toString(bigNameLen) + "\n"
		count++
	}

	message := fmt.Sprintf("<:5658sus:969620902385946777> Invasores in Server: [%s ] <:5658sus:969620902385946777> \n========== Total %d ==========\n", s.Name, count)
	message += list

	return message
}

func debug(s *Server) string {
	list := ""

	for _, stage := range s.Stages {
		if !stage.isQuest() && len(stage.Objects) == 0 {
			continue
		}

		list += fmt.Sprintf("    -> Stage: %s StageId: %s\n", stage.GetName(), stage.Id)
		isQuest := "false"
		hasDeparted := "false"

		if stage.isQuest() {
			isQuest = "true"
		}

		list += fmt.Sprintf("    '-> isQuest: %s\n", isQuest)

		if stage.isQuest() {
			if len(stage.Clients) > 0 {
				hasDeparted = "true"
			}

			list += fmt.Sprintf("    '-> isDeparted: %s\n", hasDeparted)
			list += fmt.Sprintf("    '-> reserveSlots (%d/%d)\n", len(stage.ReservedClientSlots), stage.MaxPlayers)

			for charid := range stage.ReservedClientSlots {
				char, err := s.getCharacterForUser(int(charid))
				if err == nil {
					list += fmt.Sprintf("        '-> %s\n", char.Name)
				}
			}
		}

		list += "    '-> objects: \n"
		for _, obj := range stage.Objects {
			objInfo := fmt.Sprintf("X,Y,Z: %f %f %f", obj.X, obj.Y, obj.Z)
			list += fmt.Sprintf("        '-> ObjectId: %d - %s\n", obj.Id, objInfo)
		}
	}

	message := fmt.Sprintf("Objects in Server: [%s ]\n", s.Name)
	message += list

	return message
}

func questlist(s *Server) string {
	list := ""

	for _, stage := range s.Stages {
		if !stage.isQuest() {
			continue
		}

		hasDeparted := ""
		if len(stage.Clients) > 0 {
			hasDeparted = " - departed"
		}
		list += fmt.Sprintf("    '-> StageId: %s (%d/%d) %s - %s\n", stage.Id, len(stage.ReservedClientSlots), stage.MaxPlayers, hasDeparted, stage.CreatedAt)

		for charid := range stage.ReservedClientSlots {
			char, err := s.getCharacterForUser(int(charid))
			if err == nil {
				list += fmt.Sprintf("        '-> %s\n", char.Name)
			}
		}
	}

	message := fmt.Sprintf("Quests in Server: [%s ]\n", s.Name)
	message += list

	return message
}

func removeStageById(s *Server, stageId string) string {
	if s.Stages[stageId] != nil {
		delete(s.Stages, stageId)
		return "Stage deleted!"
	}

	return "Stage not found!"
}

func cleanStr(str string) string {
	return strings.ToLower(strings.Trim(str, " "))
}

func getCharInfo(server *Server, charName string) string {
	infos := []CharInfo{}

	for _, stage := range server.Stages {
		for client := range stage.Clients {

			if client.Name == "" {
				continue
			}

			if cleanStr(client.Name) == cleanStr(charName) {
				infos = append(infos, CharInfo{
					CharID:    client.CharID,
					CharName:  client.Name,
					IP:        client.rawConn.RemoteAddr().String(),
					StageId:   stage.Id,
					StageName: stage.GetName(),
				})
			}
		}
	}

	if len(infos) == 0 {
		return "Character not found"
	}

	result := ""
	for _, info := range infos {
		objInfo := ""

		result += fmt.Sprintf("Character: %s\nIp: %s\nStage: %s\nStageId: %s\n%s\n\n", info.CharName, info.IP, info.StageName, info.StageId, objInfo)
	}

	return result
}

func disconnectChar(server *Server, charName string) string {
	infos := []CharInfo{}

	for _, stage := range server.Stages {
		for client := range stage.Clients {

			if client.Name == "" {
				continue
			}

			if cleanStr(client.Name) == cleanStr(charName) {
				infos = append(infos, CharInfo{
					CharID:    client.CharID,
					CharName:  client.Name,
					StageId:   stage.Id,
					StageName: stage.GetName(),
				})

				logoutPlayer(client)
			}
		}
	}

	if len(infos) == 0 {
		return "Character not found"
	}

	result := ""
	for _, info := range infos {
		result += fmt.Sprintf("%s Disconnected from stage: %s (%s)\n", info.CharName, info.StageName, info.StageId)
	}

	return result
}

func listSession(server *Server) string {
	infos := []CharInfo{}

	for _, client := range server.Sessions {
		infos = append(infos, CharInfo{
			CharID:    client.CharID,
			CharName:  client.Name,
			StageId:   "",
			StageName: "",
			IP:        client.rawConn.RemoteAddr().String(),
		})
	}

	if len(infos) == 0 {
		return "No Sessions founded"
	}

	sort.SliceStable(infos, func(i, j int) bool {
		return infos[i].CharName < infos[j].CharName
	})

	result := "Sessions:\n"
	for _, info := range infos {
		result += fmt.Sprintf("    CharName: %s\n    Id: %d\n    Ip: %s\n\n", info.CharName, info.CharID, info.IP)
	}

	return result
}

func removeSessionByIp(server *Server, ip string) string {
	infos := []CharInfo{}

	for _, client := range server.Sessions {
		if cleanStr(client.rawConn.RemoteAddr().String()) == cleanStr(ip) {
			infos = append(infos, CharInfo{
				CharID:    client.CharID,
				CharName:  client.Name,
				StageId:   "",
				StageName: "",
				IP:        client.rawConn.RemoteAddr().String(),
			})

			logoutPlayer(client)
		}
	}

	if len(infos) == 0 {
		return "Sessions not found"
	}

	result := "Sessions:\n"
	for _, info := range infos {
		result += fmt.Sprintf("Disconnect Session:\n    Char: %s\n    Id: %d\n    Ip: %s\n\n", info.CharName, info.CharID, info.IP)
	}

	return result
}

func (s *Server) isDiscordAdmin(ds *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, role := range m.Member.Roles {
		for _, id := range s.erupeConfig.Discord.DevRoles {
			if id == role {
				return true
			}
		}
	}

	return false
}

// onDiscordMessage handles receiving messages from discord and forwarding them ingame.
func (s *Server) onDiscordMessage(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from our bot, or ones that are not in the correct channel.
	if m.Author.ID == ds.State.User.ID || !s.Enable {
		return
	}

	// Ignore other channels in devMode
	if s.erupeConfig.Discord.DevMode && m.ChannelID != s.erupeConfig.Discord.RealtimeChannelID {
		return
	}

	args := strings.Split(m.Content, " ")
	commandName := args[0]

	// Move to slash commadns
	if commandName == "!players" {
		// nolint:errcheck // Error return value of `.` is not checked
		ds.ChannelMessageSend(m.ChannelID, PlayerList(s))
		return
	}

	if commandName == "!char" {
		if len(args) < 2 {
			// nolint:errcheck // Error return value of `.` is not checked
			ds.ChannelMessageSend(m.ChannelID, "Usage: !char <char name>")
			return
		}

		charName := strings.Join(args[1:], " ")
		// nolint:errcheck // Error return value of `.` is not checked
		ds.ChannelMessageSend(m.ChannelID, getCharInfo(s, charName))
		return
	}

	if commandName == "!disconnect" && s.isDiscordAdmin(ds, m) {
		if len(args) < 2 {
			// nolint:errcheck
			ds.ChannelMessageSend(m.ChannelID, "Usage: !disconnect <char name>")
			return
		}

		charName := strings.Join(args[1:], " ")
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, disconnectChar(s, charName))
		return
	}

	if commandName == "!debug" && s.isDiscordAdmin(ds, m) {
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, debug(s))
		return
	}

	if commandName == "!kill" && s.isDiscordAdmin(ds, m) {
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, "Dizimando aldeia em 3..2..1")
		os.Exit(1)
		return
	}

	if commandName == "!questlist" && s.isDiscordAdmin(ds, m) {
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, questlist(s))
		return
	}

	if commandName == "!sessions" && s.isDiscordAdmin(ds, m) {
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, listSession(s))
		return
	}

	if commandName == "!remove-session" && s.isDiscordAdmin(ds, m) {
		if len(args) < 2 {
			// nolint:errcheck
			ds.ChannelMessageSend(m.ChannelID, "Usage: !remove-session <session ip>")
			return
		}

		ip := strings.Join(args[1:], " ")
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, removeSessionByIp(s, ip))
		return
	}

	if commandName == "!remove-stage" && s.isDiscordAdmin(ds, m) {
		if len(args) < 2 {
			// nolint:errcheck
			ds.ChannelMessageSend(m.ChannelID, "Usage: !remove-stage <stage id>")
			return
		}

		stageId := strings.Join(args[1:], " ")
		// nolint:errcheck
		ds.ChannelMessageSend(m.ChannelID, removeStageById(s, stageId))
		return
	}

	if m.ChannelID == s.erupeConfig.Discord.RealtimeChannelID {
		message := fmt.Sprintf("[DISCORD] %s: %s", m.Author.Username, m.Content)
		s.BroadcastChatMessage(s.discordBot.NormalizeDiscordMessage(message))
	}
}
