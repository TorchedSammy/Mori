package main

import (
	"log"
	"path/filepath"
	"strings"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
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
				if event.Op&fsnotify.Chmod == fsnotify.Chmod && strings.HasSuffix(event.Name, ".osz") {
					beatmapname := filepath.Base(event.Name)
					dest := os.Getenv("HOME") + "/.local/share/osu-wine/OSU/Songs/" + beatmapname
					log.Printf("Moving %s to osu! song directory to %s\n", event.Name, dest)
					os.Rename(event.Name, dest)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(os.Getenv("HOME") + "/Downloads")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

