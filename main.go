package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
)

type TenshoConfig struct {
	OsuDir string `json:"osuDir"`
	SourceDir string `json:"sourceDir"`
}

func main() {
	conffile, _ := os.ReadFile(os.Getenv("HOME") + "/.config/tensho/tensho.json")
	conf := TenshoConfig{
		OsuDir: os.Getenv("HOME") + "/.local/share/osu-wine/OSU",
		SourceDir: os.Getenv("HOME") + "/Downloads",
	}
	json.Unmarshal(conffile, &conf)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
						conf.Copy(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(conf.SourceDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	go handlesig()
	<-done
}

func (c *TenshoConfig) Copy(filename string) {
	dir := ""
	switch filename[len(filename) - 4:] {
	case ".osz":
		dir = filepath.Join(c.OsuDir, "Songs")
	case ".osk":
		dir = filepath.Join(c.OsuDir, "Skins")
	default:
		return
	}

	beatmapname := filepath.Base(filename)
	dest := filepath.Join(dir, beatmapname)
	fmt.Printf("Moving %s to osu! song directory to %s\n", filename, dest)
	os.Rename(filename, dest)
}

func handlesig() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for range c {
		os.Exit(0)
	}
}
