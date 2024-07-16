package sync

import (
	"fmt"
	"strings"
	"sync"

	"github.com/tmichov/twingo/internal/config"
)

type FileList struct {
	sync.Mutex
	Files map[string]struct{}
	DeletedFiles map[string]struct{}
}

func NewFileList() *FileList {
	return &FileList{
		Files: make(map[string]struct{}),
		DeletedFiles: make(map[string]struct{}),
	}
}

func (fl *FileList) AddFile(filename string) {
	fl.Lock()
	defer fl.Unlock()

	path := config.AppConfig.SyncFolder;
	homeDir := getHomeDir()

	newName := strings.Replace(filename, strings.Replace(path, "~", homeDir, 1), "", 1);

	fl.Files[newName] = struct{}{}

	if _, ok := fl.DeletedFiles[newName]; ok {
		delete(fl.DeletedFiles, newName)
	}
}

func (fl *FileList) RemoveFile(filename string) {
	fl.Lock()

	defer fl.Unlock()

	path := config.AppConfig.SyncFolder;
	homeDir := getHomeDir()
	newName := strings.Replace(filename, strings.Replace(path, "~", homeDir, 1), "", 1);

	if _, ok := fl.Files[newName]; !ok {
		fl.DeletedFiles[newName] = struct{}{}
	} else {
		delete(fl.Files, newName)
	}
}

func (fl *FileList) PrintFiles() {
	fl.Lock()
	defer fl.Unlock()

	fmt.Println("Files: ")

	for file := range fl.Files {
		fmt.Println(file)
	}

	fmt.Println("Deleted files: ")
	for file := range fl.DeletedFiles {
		fmt.Println(file)
	}
}
