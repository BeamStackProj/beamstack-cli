package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var configPath string

func InitConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configPath = filepath.Join(homeDir, ".beamstack", "config")

	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetDefault("PROGRESS_BAR_WIDTH", 30)
}
