package main

import (
	"fmt"
	"github.com/carloseduribeiro/crud-go/configs"
)

func main() {
	config, _ := configs.LoadConfig(".")
	fmt.Println(config.DBDriver)
}
