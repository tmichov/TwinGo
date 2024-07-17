package sync

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tmichov/twingo/internal/config"
)

type WatchedFile struct {
	path string	
	isDir bool
}

type DeletedFile struct {
	path string
}

type FileList struct {
	sync.Mutex
	Files map[string]WatchedFile
	DeletedFiles map[string]DeletedFile
}

func NewFileList() *FileList {
	return &FileList{
		Files: make(map[string]WatchedFile),
		DeletedFiles: make(map[string]DeletedFile),
	}
}

func (fl *FileList) AddWatchedFile(filename string, isDir bool) {
	fl.Lock()
	defer fl.Unlock()

	path := config.AppConfig.SyncFolder;
	homeDir := getHomeDir()

	newName := strings.Replace(filename , strings.Replace(path, "~", homeDir, 1), "", 1);

	fl.Files[filename] = WatchedFile{path: newName, isDir: isDir}

	if _, ok := fl.DeletedFiles[filename]; ok {
		delete(fl.DeletedFiles, filename)
	}
}

func (fl *FileList) DeletedItem(filename string) {
	fl.Lock()
	defer fl.Unlock()

	path := config.AppConfig.SyncFolder;
	homeDir := getHomeDir()

	newName := strings.Replace(filename , strings.Replace(path, "~", homeDir, 1), "", 1);

	if _, ok := fl.Files[filename]; ok {
		delete(fl.Files, filename)
	} else {
		fl.DeletedFiles[filename] = DeletedFile{path: newName}
	}
}

func (fl *FileList) PrintFiles() {
	fl.Lock()
	defer fl.Unlock()

	for k, v := range fl.Files {
		fmt.Println("File: ", k, v)
	}

	for k, v := range fl.DeletedFiles {
		fmt.Println("Deleted: ", k, v)
	}
}

