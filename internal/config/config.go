package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ServerIP   string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	TwinIP     string `json:"twin_ip"`
	TwinPort   string `json:"twin_port"`
	SyncFolder string `json:"sync_folder"`
	IgnoredDirs []string `json:"ignored"`
}

var AppConfig Config

func LoadConfig() {
	data, err := os.ReadFile("config/config.json") 

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}

func SaveConfig(file string) {
	data, err := json.MarshalIndent(AppConfig, "", " 	")
	if err != nil {
		log.Fatalf("Error marshalling config: %v", err)
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}
