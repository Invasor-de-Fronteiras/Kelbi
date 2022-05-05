package channelserver

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// onDiscordMessage handles receiving messages from discord and forwarding them ingame.
func (s *Server) onDiscordMessage(ds *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from our bot, or ones that are not in the correct channel.
	if m.Author.ID == ds.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")
	commandName := args[0]
	// Move to slash commadns
	if commandName == "!players" && s.enable {
		ds.ChannelMessageSend(m.ChannelID, PlayerList(s))
		return
	}

	if m.ChannelID != s.erupeConfig.Discord.RealtimeChannelID {
		message := fmt.Sprintf("[DISCORD] %s: %s", m.Author.Username, m.Content)
		s.BroadcastChatMessage(message)
	}
}

func PlayerList(s *Server) string {
	list := ""
	count := 0
	for _, stage := range s.stages {
		for client := range stage.clients {
			list = list + "- " + client.Name + "\n"
			count += 1
		}
	}

	message := fmt.Sprintf("List of players: Server [%s]\n ######### Total %d ######### \n", s.name, count)
	message += list

	return message
}
