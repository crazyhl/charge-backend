package main

import (
	"charge/pkg/initial"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	err := initial.NewConfig()
	if err != nil {
		fmt.Println("Load config fail:", err)
		os.Exit(1)
	}
	db, err := initial.NewDb()
	if err != nil {
		fmt.Println("Connect db fail:", err)
	}
	fmt.Println(db)
}

func main() {
	fmt.Println(viper.Get("mysql"))
	fmt.Println(viper.GetString("mysql.charset"))
}