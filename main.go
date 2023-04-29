package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erupe-ce/config"
	"erupe-ce/server/channelserver"
	"erupe-ce/server/discordbot"
	"erupe-ce/server/entranceserver"
	httpserver "erupe-ce/server/http-server"
	"erupe-ce/server/launcherserver"
	"erupe-ce/server/signserver"
	"erupe-ce/server/signv2server"
	"erupe-ce/utils"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var erupeConfig *config.Config

// Temporary DB auto clean on startup for quick development & testing.
func cleanDB(db *sqlx.DB) {
	_ = db.MustExec("DELETE FROM guild_characters")
	_ = db.MustExec("DELETE FROM guilds")
	_ = db.MustExec("DELETE FROM characters")
	_ = db.MustExec("DELETE FROM sign_sessions")
	_ = db.MustExec("DELETE FROM users")
}

func modifyToSentryLogger(log *zap.Logger, client *sentry.Client) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(client))

	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}

	log = zapsentry.AttachCoreToLogger(core, log)

	// to use breadcrumbs feature - create new scope explicitly
	// and attach after attaching the core
	return log.With(zapsentry.NewScope())
}

func main() {
	var err error
	zapLogger, _ := utils.NewLogger()
	logger := zapLogger.Named("main")

	logger.Info("Starting Erupe")

	// Load the configuration.
	erupeConfig, err = config.LoadConfig()
	if err != nil {
		preventClose(fmt.Sprintf("Failed to load config: %s", err.Error()))
	}

	if erupeConfig.Sentry.Enabled {
		sentryClient, err := sentry.NewClient(sentry.ClientOptions{
			Dsn: erupeConfig.Sentry.DSN,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: 1.0,
		})

		if err != nil {
			logger.Fatal("sentry.Init: %s", zap.Error(err))
		}

		logger = modifyToSentryLogger(logger, sentryClient)
		defer sentry.Flush(2 * time.Second)
	}

	if erupeConfig.Database.Password == "" {
		preventClose("Database password is blank")
	}

	if net.ParseIP(erupeConfig.Host) == nil {
		ips, _ := net.LookupIP(erupeConfig.Host)
		for _, ip := range ips {
			if ip != nil {
				erupeConfig.Host = ip.String()
				break
			}
		}
		if net.ParseIP(erupeConfig.Host) == nil {
			preventClose("Invalid host address")
		}
	}

	// Discord bot
	var discordBot *discordbot.DiscordBot = nil

	if erupeConfig.Discord.Enabled {
		bot, err := discordbot.NewDiscordBot(discordbot.DiscordBotOptions{
			Logger: logger,
			Config: erupeConfig,
		})

		if err != nil {
			preventClose(fmt.Sprintf("Failed to create Discord bot: %s", err.Error()))
		}

		// Discord bot
		err = bot.Start()

		if err != nil {
			preventClose(fmt.Sprintf("Failed to start Discord bot: %s", err.Error()))
		}

		discordBot = bot
		logger.Info("Discord bot is enabled")
	} else {
		logger.Info("Discord bot is disabled")
	}

	// Create the postgres DB pool.
	connectString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname= %s sslmode=disable",
		erupeConfig.Database.Host,
		erupeConfig.Database.Port,
		erupeConfig.Database.User,
		erupeConfig.Database.Password,
		erupeConfig.Database.Database,
	)

	db, err := sqlx.Open("postgres", connectString)
	if err != nil {
		preventClose(fmt.Sprintf("Failed to open SQL database: %s", err.Error()))
	}

	// Test the DB connection.
	err = db.Ping()
	if err != nil {
		preventClose(fmt.Sprintf("Failed to ping database: %s", err.Error()))
	}

	db.SetMaxOpenConns(50)
	logger.Info("Connected to database")

	// Clear stale data
	_ = db.MustExec("DELETE FROM sign_sessions")
	_ = db.MustExec("DELETE FROM servers")

	// Clean the DB if the option is on.
	if erupeConfig.DevMode && erupeConfig.DevModeOptions.CleanDB {
		logger.Info("Cleaning DB")
		cleanDB(db)
		logger.Info("Done cleaning DB")
	}

	// Now start our server(s).

	// Launcher HTTP server.
	var launcherServer *launcherserver.Server
	if erupeConfig.Launcher.Enabled {
		launcherServer = launcherserver.NewServer(
			&launcherserver.Config{
				Logger:                   logger.Named("launcher"),
				ErupeConfig:              erupeConfig,
				DB:                       db,
				UseOriginalLauncherFiles: erupeConfig.Launcher.UseOriginalLauncherFiles,
			})
		err = launcherServer.Start()
		if err != nil {
			preventClose(fmt.Sprintf("Failed to start launcher server: %s", err.Error()))
		}
		logger.Info("Started launcher server")
	}

	// Entrance server.

	var entranceServer *entranceserver.Server
	if config.ErupeConfig.Entrance.Enabled {
		entranceServer = entranceserver.NewServer(
			&entranceserver.Config{
				Logger:      logger.Named("entrance"),
				ErupeConfig: config.ErupeConfig,
				DB:          db,
			})
		err = entranceServer.Start()
		if err != nil {
			preventClose(fmt.Sprintf("Failed to start entrance server: %s", err.Error()))
		}
		logger.Info("Started entrance server")
	}

	// Sign server.

	var signServer *signserver.Server
	if config.ErupeConfig.Sign.Enabled {
		signServer = signserver.NewServer(
			&signserver.Config{
				Logger:      logger.Named("sign"),
				ErupeConfig: config.ErupeConfig,
				DB:          db,
			})
		err = signServer.Start()
		if err != nil {
			preventClose(fmt.Sprintf("Failed to start sign server: %s", err.Error()))
		}
		logger.Info("Started sign server")
	}

	// New Sign server
	var newSignServer *signv2server.Server
	if config.ErupeConfig.SignV2.Enabled {
		newSignServer = signv2server.NewServer(
			&signv2server.Config{
				Logger:      logger.Named("sign"),
				ErupeConfig: config.ErupeConfig,
				DB:          db,
			})
		err = newSignServer.Start()
		if err != nil {
			preventClose(fmt.Sprintf("Failed to start sign-v2 server: %s", err.Error()))
		}
		logger.Info("Started new sign server")
	}

	var channels []*channelserver.Server

	if config.ErupeConfig.Channel.Enabled {
		channelQuery := ""
		si := 0
		ci := 0
		count := 1
		for _, ee := range config.ErupeConfig.Entrance.Entries {
			for i, ce := range ee.Channels {
				sid := (4096 + si*256) + (16 + ci)
				c := *channelserver.NewServer(&channelserver.Config{
					ID:          uint16(sid),
					Logger:      logger.Named("channel-" + fmt.Sprint(count)),
					ErupeConfig: erupeConfig,
					DB:          db,
					Name:        ee.Name,
					Enable:      ce.MaxPlayers > 0,
					DiscordBot:  discordBot,
				})
				if ee.IP == "" {
					c.IP = config.ErupeConfig.Host
				} else {
					c.IP = ee.IP
				}
				c.Port = ce.Port
				err = c.Start()
				if err != nil {
					preventClose(fmt.Sprintf("Failed to start channel server: %s", err.Error()))
				} else {
					channelQuery += fmt.Sprintf(`INSERT INTO servers (server_id, season, current_players, world_name, world_description, land) VALUES (%d, %d, 0, '%s', '%s', %d);`, sid, si%3, ee.Name, ee.Description, i+1)
					channels = append(channels, &c)
					logger.Info(fmt.Sprintf("Started channel server %d on port %d", count, ce.Port))
					ci++
					count++
				}
			}
			ci = 0
			si++
		}

		// Register all servers in DB
		_ = db.MustExec(channelQuery)

		for _, c := range channels {
			c.Channels = channels
		}

		if erupeConfig.ServerHttp.Enabled {
			httpContext := httpserver.HttpServerContext{
				Servers:     channels,
				ErupeConfig: erupeConfig,
				Address:     fmt.Sprintf("0.0.0.0:%d", erupeConfig.ServerHttp.Port),
				Token:       erupeConfig.ServerHttp.Token,
			}

			go httpserver.RunHttpServer(&httpContext)
		}
	}

	// Wait for exit or interrupt with ctrl+C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info("Trying to shutdown gracefully")

	if erupeConfig.Channel.Enabled {
		for _, c := range channels {
			c.Shutdown()
		}
	}

	if erupeConfig.Sign.Enabled {
		signServer.Shutdown()
	}

	if erupeConfig.SignV2.Enabled {
		newSignServer.Shutdown()
	}

	if erupeConfig.Entrance.Enabled {
		entranceServer.Shutdown()
	}

	if erupeConfig.Launcher.Enabled {
		launcherServer.Shutdown()
	}

	time.Sleep(1 * time.Second)
}

func wait() {
	for {
		time.Sleep(time.Millisecond * 100)
	}
}

func preventClose(text string) {
	if erupeConfig.DisableSoftCrash {
		os.Exit(0)
	}
	fmt.Println("\nFailed to start Erupe:\n" + text)
	go wait()
	fmt.Println("\nPress Enter/Return to exit...")
	fmt.Scanln()
	os.Exit(0)
}
