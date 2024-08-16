package main

import (
	"fmt"

	"github.com/yoruakio/gowebserver/config"
	"github.com/yoruakio/gowebserver/http"
)

func main() {
	config.LoadConfig()

	var config = config.GetConfig()

	fmt.Printf("Config:\n  host: %s\n  port: %s\n", config.Host, config.Port)

	app := http.Initialize()

	http.Start(app)
}
