package main

import (
	"fmt"

	"github.com/yoruakio/gowebserver/config"
)

func main() {
	config.LoadConfig()

	fmt.Println(config.GetConfig())
}
