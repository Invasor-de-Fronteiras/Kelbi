package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds the global server-wide config.
type Config struct {
	Host             string `mapstructure:"Host"`
	BinPath          string `mapstructure:"BinPath"`
	DisableSoftCrash bool   // Disables the 'Press Return to exit' dialog allowing scripts to reboot the server automatically
	DevMode          bool

	DevModeOptions DevModeOptions
	Discord        Discord
	Database       Database
	ServerHttp     ServerHttp
	Launcher       Launcher
	Sign           Sign

	Language               string
	FeaturedWeapons        int    // Number of Active Feature weapons to generate daily
	HideLoginNotice        bool   // Hide the Erupe notice on login
	LoginNotice            string // MHFML string of the login notice displayed
	PatchServerManifest    string // Manifest patch server override
	PatchServerFile        string // File patch server override
	ScreenshotAPIURL       string // Destination for screenshots uploaded to BBS
	DeleteOnSaveCorruption bool   // Attempts to save corrupted data will flag the save for deletion

	Commands []Command
	Courses  []Course
	SignV2   SignV2
	Channel  Channel
	Entrance Entrance
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

// Discord holds the discord integration config.
type Discord struct {
	Enabled           bool
	BotToken          string
	ServerID          string
	RealtimeChannelID string
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
	Port  int
	Token string
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

type Channel struct {
	Enabled bool
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

// LoadConfig loads the given config toml file.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("DevModeOptions.SaveDumps", SaveDumpOptions{
		Enabled:   false,
		OutputDir: "savedata",
	})

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

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
