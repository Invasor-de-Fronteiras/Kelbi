package channelserver

import (
	"fmt"
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
	for _, stage := range s.stages {
		count += len(stage.clients)
	}

	message := fmt.Sprintf("Server [%s]: %d players;", s.name, count)

	return message
}

type ListPlayer struct {
	CharName    string
	InQuest     bool
	WeaponEmoji string
	QuestEmoji  string
	StageName   string
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
		":black_circle:",
		":red_circle:",
		":blue_circle:",
		":brown_circle:",
		":green_circle:",
		":purple_circle:",
		":yellow_circle:",
		":orange_circle:",
	}

	bigNameLen := 0

	for _, stage := range s.stages {
		if len(stage.clients) == 0 {
			continue
		}

		questEmoji := questEmojis[len(questEmojis)-1]
		questEmojis = questEmojis[:len(questEmojis)-1]

		isQuest := stage.isQuest()
		for client := range stage.clients {
			char, err := s.getCharacterForUser(int(client.charID))
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

	message := fmt.Sprintf("<:5658sus:969620902385946777> Invasores in Server: [%s ] <:5658sus:969620902385946777> \n========== Total %d ==========\n", s.name, count)
	message += list

	return message
}

func cleanStr(str string) string {
	return strings.ToLower(strings.Trim(str, " "))
}

func getCharInfo(server *Server, charName string) string {
	var s *Stage
	var c *Session

	for _, stage := range server.stages {
		for client := range stage.clients {

			if client.Name == "" {
				continue
			}

			if cleanStr(client.Name) == cleanStr(charName) {
				s = stage
				c = client
			}

		}
	}

	if s == nil {
		return "Character not found"
	}

	objInfo := ""

	obj := server.FindStageObjectByChar(c.charID)
	// server.logger.Info("Found object: %+v", zap.Object("obj", obj))

	if obj != nil {
		objInfo = fmt.Sprintf("X,Y,Z: %f %f %f", obj.x, obj.y, obj.z)
	}

	return fmt.Sprintf("Character: %s\nStage: %s\nStageId: %s\n%s", c.Name, s.GetName(), s.id, objInfo)
}

// onDiscordMessage handles receiving messages from discord and forwarding them ingame.
func (s *Server) onDiscordMessage(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from our bot, or ones that are not in the correct channel.
	if m.Author.ID == ds.State.User.ID || !s.enable {
		return
	}

	args := strings.Split(m.Content, " ")
	commandName := args[0]
	// Move to slash commadns
	if commandName == "!players" {
		ds.ChannelMessageSend(m.ChannelID, PlayerList(s))
		return
	}

	if commandName == "-char" {
		if len(args) < 2 {
			ds.ChannelMessageSend(m.ChannelID, "Usage: !char <char name>")
			return
		}

		charName := strings.Join(args[1:], " ")
		ds.ChannelMessageSend(m.ChannelID, getCharInfo(s, charName))
		return
	}

	if m.ChannelID == s.erupeConfig.Discord.RealtimeChannelID {
		message := fmt.Sprintf("[DISCORD] %s: %s", m.Author.Username, m.Content)
		s.BroadcastChatMessage(s.discordBot.NormalizeDiscordMessage(message))
	}
}
