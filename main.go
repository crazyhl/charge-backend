package main

import (
	"charge/container"
	"charge/initial"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
		os.Exit(2)
	}
	// 存入 container
	c := container.GetContainer()
	c.SetDb(db)
}

func main() {
	//c := container.GetContainer()
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "Cimple-Fiber",
	})

	if !fiber.IsChild() {
		// 只有在主线程的时候才会auto merge 数据结构
		initial.AutoMigrate()
	}

	err := app.Listen(viper.GetString("http-server.host") + ":" + viper.GetString("http-server.port"))
	if err != nil {
		fmt.Println("Start Server Error:", err)
		os.Exit(3)
	}
}
