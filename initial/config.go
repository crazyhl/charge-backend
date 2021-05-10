package initial

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func NewConfig() error {
	// viper
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	currentDirection, currentDirectionErr := os.Getwd()
	configPath := "./config/"
	fmt.Println(currentDirection)
	if currentDirectionErr != nil {
		configPath = currentDirection + "/config/"
	}
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()

	return err
}
