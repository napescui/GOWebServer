package main

import (
	"github.com/yoruakio/gowebserver/config"
	"github.com/yoruakio/gowebserver/http"
	"github.com/yoruakio/gowebserver/logger"
)

func main() {
	config.LoadConfig()

	var config = config.GetConfig()

	logger.Infof("Configuration: %+v ", config)

	app := http.Initialize()
	http.Start(app)
}
