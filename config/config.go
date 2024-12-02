package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	// Growtopia Server Configuration
	Host     string `json:"host"`
	Port     string `json:"port"`
	LoginUrl string `json:"loginUrl"`
	ServerCdn string `json:"serverCdn"`

	// Logger Configuration
	Logger bool `json:"isLogging"`

	// Rate Limiter Configuration
	RateLimit         int `json:"rateLimit"`
	RateLimitDuration int `json:"rateLimitDuration"`

	// Geo Location Configuration
	GeoLocation []string `json:"trustedRegions"`
	EnableGeo   bool     `json:"enableGeo"`
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
		Host:              "127.0.0.1",
		Port:              "17777",
		LoginUrl:          "gtsalogin.vercel.app",
		ServerCdn:         "default",
		Logger:            true,
		RateLimit:         150, // 60 requests per minute
		RateLimitDuration: 5,   // 2 minutes of rate limit cooldown
		EnableGeo:         true,
		GeoLocation:       []string{"ID", "SG", "MY"},
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("config.json", data, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
