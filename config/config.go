package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	// Growtopia Server Configuration
	Host string `json:"host"`
	Port string `json:"port"`

	// Logger Configuration
	Logger bool `json:"isLogging"`
}

var config Config
var isLoaded bool

func LoadConfig() Config {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		config = CreateConfig()
		isLoaded = true
		return config
	}

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	isLoaded = true
	return config
}

func GetConfig() Config {
	if !isLoaded {
		log.Fatal("LoadConfig() is not called")
	}
	return config
}

func CreateConfig() Config {
	config := Config{
		Host: "127.0.0.1",
		Port: "17091",
		Logger: true,
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("config.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
