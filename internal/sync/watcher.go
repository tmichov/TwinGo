package sync

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/tmichov/twingo/internal/config"
)

func Watcher(filelist *FileList) {
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

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("Deleted in watcher: ", event.Name)
					
					// remove all watchers for all subdirectories
					list := watcher.WatchList()
					for _, item := range list {
						if strings.HasPrefix(item, event.Name) {
							watcher.Remove(item)
						}
					}

					filelist.DeletedItem(event.Name)

					continue;
				}

				info, err := os.Stat(event.Name)
				if err != nil {
					fmt.Println("Error getting file info: ", err)
					continue
				}

				fmt.Println("Event: ", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Modified file: ", event.Name)
					filelist.AddWatchedFile(event.Name, false)
					continue
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					if info.IsDir() {
						fmt.Println("Created folder: ", event.Name)
						filelist.AddWatchedFile(event.Name, true)

						watcher.Add(event.Name)
					} else {
						fmt.Println("Created file in watched: ", event.Name)
						filelist.AddWatchedFile(event.Name, false)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("Watcher errorj29k channel closed")
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

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	return homeDir
}

func getFilePath() string {
	homeDir := getHomeDir()

	filePath := config.AppConfig.SyncFolder;

	fmt.Println("Syncing folder: ", filePath)

	if filePath[0] == '~' {
		filePath = homeDir + filePath[1:]
	}

	return filePath
}
