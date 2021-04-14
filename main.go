package main

import (
	"charge/container"
	"charge/controller/account"
	"charge/initial"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// 跨域
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080",
	}))

	if !fiber.IsChild() {
		// 只有在主线程的时候才会auto merge 数据结构
		err := initial.AutoMigrate()
		if err != nil {
			fmt.Println("Migrate db fail:", err)
			os.Exit(4)
		}
	}

	accountGroup := app.Group("/account")
	accountGroup.Post("", account.Add)
	accountGroup.Get("/list", account.List)
	accountGroup.Get("/:id/edit", account.EditDetail)
	accountGroup.Delete("/:id", account.Delete)
	accountGroup.Put("/:id", account.Edit)

	err := app.Listen(viper.GetString("http-server.host") + ":" + viper.GetString("http-server.port"))
	if err != nil {
		fmt.Println("Start Server Error:", err)
		os.Exit(3)
	}

}
