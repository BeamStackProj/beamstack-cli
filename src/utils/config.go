package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BeamStackProj/beamstack-cli/src/types"
	"github.com/spf13/viper"
)

var (
	configPath         string
	ValidateConfigOnce sync.Once
)

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

func ValidateCluster() (profile types.Profiles, err error) {
	ValidateConfigOnce.Do(
		func() {
			currentContext, _err := GetCurrentContext()
			if _err != nil {
				err = fmt.Errorf("error getting current context: %v", _err)
				return
			}
			contextsStringMap := viper.GetStringMapString("contexts")

			if _, ok := contextsStringMap[currentContext]; !ok {
				err = fmt.Errorf("cluster not initialized. please run 'beamstack init' to initialize cluster")
				return
			}

			profile, _err = GetProfile(contextsStringMap[currentContext])
			if _err != nil {
				err = fmt.Errorf("error getting current profile: %v", _err)
				return
			}
		},
	)
	return
}
