package initial

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func NewConfig() error {
	// viper
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	executable, getExecutableErr := os.Executable()
	if getExecutableErr != nil {
		panic(getExecutableErr)
	}
	exPath := filepath.Dir(executable)
	configPath := exPath + "/config/"
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()

	return err
}
