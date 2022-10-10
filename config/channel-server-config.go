package config

import (
	"errors"
	"net"

	"github.com/spf13/viper"
)

// ChannelServerConfig holds the global server-wide config.
type ChannelServerConfig struct {
	EntranceChannelInfo

	Name    string
	Host    string `mapstructure:"Host"`
	BinPath string `mapstructure:"BinPath"`

	DevMode bool

	DevModeOptions DevModeOptions
	Discord        DiscordConfig
	Database       DatabaseConfig
}

func (conf ChannelServerConfig) Enabled() bool {
	return conf.MaxPlayers > 0
}

// LoadConfig loads the given config toml file.
func LoadChannelServerConfig() (erupeConfig *ChannelServerConfig, err error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("DevModeOptions.SaveDumps", SaveDumpOptions{
		Enabled:   false,
		OutputDir: "savedata",
	})

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	erupeConfig = &ChannelServerConfig{}
	err = viper.Unmarshal(erupeConfig)
	if err != nil {
		return
	}

	if erupeConfig.Host == "" {
		erupeConfig.Host = getOutboundIP4().To4().String()
	}

	if erupeConfig.Database.Password == "" {
		err = errors.New("DATABASE PASSWORD IS BLACK")
		return
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
			err = errors.New("INVALID HOST ADDRESS")
			return
		}
	}

	return
}
