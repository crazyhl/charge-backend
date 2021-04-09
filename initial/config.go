package initial

import "github.com/spf13/viper"

func NewConfig() error {
	// viper
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()

	return err
}
