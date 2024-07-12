package sync

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/tmichov/twingo/internal/config"
)

func Watcher() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("Watcher event channel closed")
					return
				}

				fmt.Println("Event: ", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modified file: ", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("Watcher error channel closed")
					return
				}
				
				fmt.Println("Error: ", err)
			}
		}
	}()

	err = filepath.WalkDir(getFilePath(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}


		if d.IsDir() {
			for _, dir := range config.AppConfig.IgnoredDirs {
				if d.IsDir() && (d.Name() == string(dir)) {
					fmt.Println("Ignoring directory: ", path)
					return filepath.SkipDir
				}
			}
			fmt.Println("Watching directory: ", path)
			return watcher.Add(path)
		}

		return nil
	});
	
	if err != nil {
		log.Fatalf("Error walking path: %v", err)
	}

	<-make(chan bool)
}

func getFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	filePath := config.AppConfig.SyncFolder;

	fmt.Println("Syncing folder: ", filePath)

	// replace ~ with home directory
	if filePath[0] == '~' {
		filePath = homeDir + filePath[1:]
	}

	return filePath
}
