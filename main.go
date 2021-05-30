package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
)

type Mori struct {
	OsuDir string `json:"osuDir"`
	SourceDir string `json:"sourceDir"`
}

func main() {
	conffile, _ := os.ReadFile(os.Getenv("HOME") + "/.config/mori/mori.json")
	conf := Mori{
		OsuDir: os.Getenv("HOME") + "/.local/share/osu-wine/OSU",
		SourceDir: os.Getenv("HOME") + "/Downloads",
	}
	json.Unmarshal(conffile, &conf)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer watcher.Close()

	conf.Sweep()

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
	fmt.Println("Mori has started up!")
	<-done
}

func (m *Mori) Copy(filename string) {
	dir := ""
	switch filename[len(filename) - 4:] {
	case ".osz":
		dir = filepath.Join(m.OsuDir, "Songs")
	case ".osk":
		dir = filepath.Join(m.OsuDir, "Skins")
	default:
		return
	}

	beatmapname := filepath.Base(filename)
	dest := filepath.Join(dir, beatmapname)
	fmt.Printf("Moving %s to %s\n", filename, dest)
	os.Rename(filename, dest)
}

func handlesig() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for range c {
		fmt.Println("")
		os.Exit(0)
	}
}

func (m *Mori) Sweep() {
	fmt.Println("Beginning sweep of left archives...")
	bmps, _ := filepath.Glob(m.SourceDir + "/*.osz")
	skins, _ := filepath.Glob(m.SourceDir + "/*.osk")

	for _, skin := range skins {
		m.Copy(skin)
	}
	for _, beatmap := range bmps {
		m.Copy(beatmap)
	}
}

