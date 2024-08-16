package main

import (
	"github.com/yoruakio/gowebserver/config"
	"github.com/yoruakio/gowebserver/logger"
)

func main() {
	config.LoadConfig()

	logger.Infof("Host: %s", config.GetConfig().Host)
}
