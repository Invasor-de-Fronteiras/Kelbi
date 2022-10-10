package main

import (
	"fmt"

	"erupe-ce/config"
	"erupe-ce/server/channelserver"
	"erupe-ce/server/discordbot"
	"erupe-ce/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	var err error
	zapLogger, _ := utils.NewLogger()
	logger := zapLogger.Named("main")

	logger.Info("Starting Channel-Server")

	// Load the configuration.
	config, err := config.LoadChannelServerConfig()

	if err != nil {
		logger.Fatal("Failed to load config: %s", zap.Error(err))
	}

	// Discord bot
	var discordBot *discordbot.DiscordBot = nil

	if config.Discord.Enabled {
		discordBot, err = discordbot.NewDiscordBot(discordbot.DiscordBotOptions{
			Logger: logger,
			Config: &config.Discord,
		})

		if err != nil {
			logger.Fatal("Failed to create Discord bot: %s", zap.Error(err))
		}

		// Discord bot
		err = discordBot.Start()

		if err != nil {
			logger.Fatal("Failed to start Discord bot: %s", zap.Error(err))
		}

		logger.Info("Discord bot is enabled")
	} else {
		logger.Info("Discord bot is disabled")
	}

	db, err := sqlx.Open("postgres", config.Database.ConnectionString())
	if err != nil {
		logger.Fatal("Failed to open SQL database: %s", zap.Error(err))
	}

	// Clear stale data
	_ = db.MustExec("DELETE FROM sign_sessions")
	_ = db.MustExec("DELETE FROM servers")
	_ = db.MustExec("DELETE FROM cafe_accepted")
	_ = db.MustExec("UPDATE characters SET cafe_time=0")

	si := 0
	ci := 0
	count := 1

	sid := (4096 + si*256) + (16 + ci)
	c := *channelserver.NewServer(&channelserver.ChannelServerOptions{
		ID:         uint16(sid),
		Logger:     logger.Named("channel-" + fmt.Sprint(count)),
		Config:     config,
		DB:         db,
		DiscordBot: discordBot,
	})

	err = c.Start()

	if err != nil {
		logger.Fatal("Failed to start channel server: %s", zap.Error(err))
	}

	logger.Info("Started channel server %d on port %d", count, ce.Port)

}
