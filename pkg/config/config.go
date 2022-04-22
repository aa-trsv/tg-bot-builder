package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken  string
	DirPath        string
	FilePathDamper string
	DBPath         string `mapstructure:"db_file"`
	Debug          bool   `mapstructure:"debug"`

	Messages   Messages
	SSHCommand SSHCommand
}

type Messages struct {
	Responses
}

type Responses struct {
	Start          string `mapstructure:"start"`
	Rebuild        string `mapstructure:"rebuild"`
	UnknownCommand string `mapstructure:"unknown_command"`
}

type SSHCommand struct {
	Build   string `mapstructure:"build"`
	Rebuild string `mapstructure:"rebuild"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	if err := viper.BindEnv("debug"); err != nil {
		return err
	}
	if err := viper.BindEnv("dir_path"); err != nil {
		return err
	}
	if err := viper.BindEnv("file_path_dumper"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.DirPath = viper.GetString("dir_path")
	cfg.FilePathDamper = viper.GetString("file_path_dumper")
	cfg.Debug = viper.GetBool("debug")

	return nil
}
