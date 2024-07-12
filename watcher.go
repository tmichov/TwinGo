package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()

	if(err != nil) {
		fmt.Println("Error creating watcher", err)
		os.Exit(1)
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				fmt.Println("Event:", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modified file:", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				fmt.Println("Error:", err)
			}
		}
	}()

	homedir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	projectDir := filepath.Join(homedir, "Projects", "Personal", "Notes")

	err = watcher.Add(projectDir)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
