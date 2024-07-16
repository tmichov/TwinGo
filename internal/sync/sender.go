package sync

import (
	"fmt"
	"time"
)

func Send(filelist *FileList) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			fmt.Println("Sending files...")

			filelist.PrintFiles()
		}
	}()
}
