package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tmichov/twingo/internal/config"
	"github.com/tmichov/twingo/internal/sync"
)

func main() {
	config.LoadConfig()


	//displayUserInfo()
	getUserInput()

	filelist := sync.NewFileList()

	go sync.Watcher(filelist)
	go sync.Send(filelist)

	select {}
}

func displayUserInfo() {
	response, err := http.Get("https://ipinfo.io/ip")
	if(err != nil) {
		log.Fatal(err)
	}

	defer response.Body.Close()

	ip, err := io.ReadAll(response.Body)
	if(err != nil) {
		log.Fatal(err)
	}

	fmt.Println("Your IP is: ", string(ip))
}

func getUserInput() {
	var twinIP, twinPort, syncFolder string

	if config.AppConfig.TwinIP != "" {
		fmt.Println("Current Twin IP: ", config.AppConfig.TwinIP)
	} else {
		fmt.Print("Enter Twin IP: ")
		fmt.Scanln(&twinIP)
		config.AppConfig.TwinIP = twinIP
	}

	if config.AppConfig.TwinPort != "" {
		fmt.Println("Current Twin Port: ", config.AppConfig.TwinPort)
	} else {
		fmt.Print("Enter Twin Port: ")
		fmt.Scanln(&twinPort)
		config.AppConfig.TwinPort = twinPort
	}

	if config.AppConfig.SyncFolder != "" {
		fmt.Println("Current Sync Folder: ", config.AppConfig.SyncFolder)
	} else {
		fmt.Print("Enter Sync Folder: ")
		fmt.Scanln(&syncFolder)
		config.AppConfig.SyncFolder = syncFolder
	}

	config.SaveConfig("config/config.json")
}
