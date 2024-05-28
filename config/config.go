package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the global server-wide config.
type Config struct {
	Host                        string `mapstructure:"Host"`
	Teste                       string `mapstructure:"TESTE"`
	BinPath                     string `mapstructure:"BinPath"`
	DisableSoftCrash            bool   // Disables the 'Press Return to exit' dialog allowing scripts to reboot the server automatically
	ClearAllServersFromDatabase bool   //Clear all servers and sessions from database
	DevMode                     bool
	Sentry                      Sentry
	DevModeOptions              DevModeOptions
	Discord                     Discord
	Database                    Database
	ServerHttp                  ServerHttp
	Launcher                    Launcher
	Sign                        Sign

	Language               string
	HideLoginNotice        bool     // Hide the Erupe notice on login
	LoginNotices           []string // MHFML string of the login notices displayed
	PatchServers           PatchServers
	ScreenshotAPIURL       string // Destination for screenshots uploaded to BBS
	DeleteOnSaveCorruption bool   // Attempts to save corrupted data will flag the save for deletion
	QuestLoader            string // QuestLoader type (db, questlist, events_folder)
	GameplayOptions        GameplayOptions

	Commands []Command
	Courses  []Course
	SignV2   SignV2
	Channel  Channel
	Entrance Entrance
	Entry    Entry
}

type PatchServer struct {
	PatchServerManifest string // Manifest patch server override
	PatchServerFile     string // File patch server override
}

type PatchServers struct {
	En PatchServer
	Jp PatchServer
}

// DevModeOptions holds various debug/temporary options for use while developing Erupe.
type DevModeOptions struct {
	ServerName          string // To get specific instance server about (Current Players/Event Week)
	HideLoginNotice     bool   // Hide the Erupe notice on login
	LoginNotice         string // MHFML string of the login notice displayed
	CleanDB             bool   // Automatically wipes the DB on server reset.
	MaxLauncherHR       bool   // Sets the HR returned in the launcher to HR9 so that you can join non-beginner worlds.
	FixedStageID        bool   // Causes all move_stage to use the ID sl1Ns200p0a0u0 to get you into all stages
	LogInboundMessages  bool   // Log all messages sent to the server
	LogOutboundMessages bool   // Log all messages sent to the clients
	MaxHexdumpLength    int    // Maximum number of bytes printed when logs are enabled
	DivaEvent           int    // Diva Defense event status
	FestaEvent          int    // Hunter's Festa event status
	TournamentEvent     int    // VS Tournament event status
	MezFesEvent         bool   // MezFes status
	MezFesAlt           bool   // Swaps out Volpakkun for Tokotoko
	DisableTokenCheck   bool   // Disables checking login token exists in the DB (security risk!)
	DisableMailItems    bool   // Hack to prevent english versions of MHF from crashing
	DisableReturnBoost  bool
	SaveDumps           SaveDumpOptions

	AutoCreateAccount bool // Automatically create accounts if they don't exist
	QuestDebugTools   bool // Enable various quest debug logs
}

type SaveDumpOptions struct {
	Enabled   bool
	OutputDir string
}

// GameplayOptions has various gameplay modifiers
type GameplayOptions struct {
	FeaturedWeapons     int    // Number of Active Feature weapons to generate daily
	MaximumNP           int    // Maximum number of NP held by a player
	MaximumRP           uint16 // Maximum number of RP held by a player
	DisableLoginBoost   bool   // Disables the Login Boost system
	DisableBoostTime    bool   // Disables the daily NetCafe Boost Time
	BoostTimeDuration   int    // The number of minutes NetCafe Boost Time lasts for
	GuildMealDuration   int    // The number of minutes a Guild Meal can be activated for after cooking
	BonusQuestAllowance uint32 // Number of Bonus Point Quests to allow daily
	DailyQuestAllowance uint32 // Number of Daily Quests to allow daily
	MezfesSoloTickets   uint32 // Number of solo tickets given weekly
	MezfesGroupTickets  uint32 // Number of group tickets given weekly
}

// Discord holds the discord integration config.
type Discord struct {
	Enabled           bool
	BotToken          string
	ServerID          string
	RealtimeChannelID string
	LogsChannelID     string
	DevRoles          []string
	DevMode           bool
}

// Command is a channelserver chat command
type Command struct {
	Name    string
	Enabled bool
	Prefix  string
}

// Course represents a course within MHF
type Course struct {
	Name    string
	Enabled bool
}

// Database holds the postgres database config.
type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Server Http.
type ServerHttp struct {
	Port    int
	Token   string
	Enabled bool
}

// Launcher holds the launcher server config.
type Launcher struct {
	Enabled                  bool // Enables the launcher server to be served on port 80
	Path                     string
	Port                     int
	UseOriginalLauncherFiles bool
}

// Sign holds the sign server config.
type Sign struct {
	Enabled bool
	Port    int
}

// SignV2 holds the new sign server config
type SignV2 struct {
	Enabled bool
	Port    int
}

type Sentry struct {
	Enabled bool
	DSN     string
}

type Channel struct {
	Enabled bool
}

// Entrance holds the entrance server config.
type Entry struct {
	Enabled     bool
	IP          string
	Port        uint16
	Server      uint16
	Land        uint16
	Name        string // Server name, 66 byte null terminated Shift-JIS(JP) or Big5(TW).
	Description string // Server description
}

// Entrance holds the entrance server config.
type Entrance struct {
	Enabled bool
	Port    uint16
	Entries []EntranceServerInfo
}

// EntranceServerInfo represents an entry in the serverlist.
type EntranceServerInfo struct {
	IP          string
	Type        uint8  // Server type. 0=?, 1=open, 2=cities, 3=newbie, 4=bar
	Season      uint8  // Server activity. 0 = green, 1 = orange, 2 = blue
	Recommended uint8  // Something to do with server recommendation on 0, 3, and 5.
	Name        string // Server name, 66 byte null terminated Shift-JIS(JP) or Big5(TW).
	Description string // Server description
	// 4096(PC, PS3/PS4)?, 8258(PC, PS3/PS4)?, 8192 == nothing?
	// THIS ONLY EXISTS IF Binary8Header.type == "SV2", NOT "SVR"!
	AllowedClientFlags uint32

	Channels []EntranceChannelInfo
}

// EntranceChannelInfo represents an entry in a server's channel list.
type EntranceChannelInfo struct {
	Port           uint16
	MaxPlayers     uint16
	CurrentPlayers uint16
}

var ErupeConfig *Config

func init() {
	var err error
	ErupeConfig, err = LoadConfig()
	if err != nil {
		preventClose(fmt.Sprintf("Failed to load config: %s", err.Error()))
	}

}

// getOutboundIP4 gets the preferred outbound ip4 of this machine
// From https://stackoverflow.com/a/37382208
func getOutboundIP4() net.IP {
	conn, err := net.Dial("udp4", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.To4()
}

func BoolToString(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func LoadConfigFromEnv(config *Config) {
	config.Launcher.Enabled = viper.GetBool("Launcher.Enabled")
	config.Sign.Enabled = viper.GetBool("Sign.Enabled")
	config.ServerHttp.Enabled = viper.GetBool("ServerHttp.Enabled")
	config.Discord.Enabled = viper.GetBool("Discord.Enabled")
	config.Entrance.Enabled = viper.GetBool("Entrance.Enabled")
	config.Channel.Enabled = viper.GetBool("Channel.Enabled")

	println("Launcher: " + BoolToString(config.Launcher.Enabled))
	println("Discord: " + BoolToString(config.Discord.Enabled))
	println("Entrance: " + BoolToString(config.Entrance.Enabled))
	println("Channel: " + BoolToString(config.Channel.Enabled))
	println("Launcher: " + BoolToString(config.Launcher.Enabled))

	config.Entry.Enabled = viper.GetBool("Entry.Enabled")
	config.Entry.IP = viper.GetString("Entry.IP")
	config.Entry.Port = viper.GetUint16("Entry.Port")
	config.Entry.Server = viper.GetUint16("Entry.Server")
	config.Entry.Land = viper.GetUint16("Entry.Land")
	config.Entry.Name = viper.GetString("Entry.Name")
	config.Entry.Description = viper.GetString("Entry.Description")
}

// LoadConfig loads the given config toml file.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/config")
	viper.AddConfigPath(".")

	viper.SetDefault("DevModeOptions.SaveDumps", SaveDumpOptions{
		Enabled:   false,
		OutputDir: "savedata",
	})

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("ERUPE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	c := &Config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	LoadConfigFromEnv(c)

	if c.Host == "" {
		c.Host = getOutboundIP4().To4().String()
	}

	return c, nil
}

func preventClose(text string) {
	if ErupeConfig.DisableSoftCrash {
		os.Exit(0)
	}
	fmt.Println("\nFailed to start Erupe:\n" + text)
	go wait()
	fmt.Println("\nPress Enter/Return to exit...")
	fmt.Scanln()
	os.Exit(0)
}

func wait() {
	for {
		time.Sleep(time.Millisecond * 100)
	}
}
