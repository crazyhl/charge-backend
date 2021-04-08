package main

import (
	"charge/pkg/container"
	"charge/pkg/initial"
	"fmt"
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
	// 存入 container
	c := container.GetContainer()
	c.SetDb(db)
}

func main() {
	c := container.GetContainer()

	fmt.Println(c.GetDb())
}