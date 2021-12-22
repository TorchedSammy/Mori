package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pborman/getopt"
)

type Mori struct {
	OsuDir string `json:"osuDir"`
	SourceDir string `json:"sourceDir"`
	SweepTime string `json:"sweepTime"`
	AutoExtract bool `json:"autoExtract"`
	FileEvents []string `json:"fileEvents"`
}

func main() {
	userconf, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config dir: ", err)
		os.Exit(1)
	}
	defaultConf := filepath.Join(userconf, "Mori", "mori.json")
	defaultOsuDir := ""
	switch runtime.GOOS {
	case "windows":
		defaultOsuDir = filepath.Join(userconf, "osu!")
	case "darwin":
		defaultOsuDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "osu!")
	default:
		defaultOsuDir = "~/.local/share/osu-wine/OSU"
	}
	helpflag := getopt.BoolLong("help", 'h', "Prints Mori flags (this message)")
	confPath := getopt.StringLong("config", 'C', defaultConf, "")
	getopt.Parse()

	if *helpflag {
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	homedir, _ := os.UserHomeDir()
	conffile, _ := os.ReadFile(*confPath)
	conf := Mori{
		OsuDir: defaultOsuDir,
		SourceDir: "~/Downloads",
		SweepTime: "5m",
		AutoExtract: true,
		FileEvents: []string{"chmod", "rename"},
	}
	json.Unmarshal(conffile, &conf)
	conf.OsuDir = strings.Replace(conf.OsuDir, "~", homedir, 1)
	conf.SourceDir = strings.Replace(conf.SourceDir, "~", homedir, 1)
	sweepDuration, err := time.ParseDuration(conf.SweepTime)
	if err != nil {
		fmt.Println("Could not parse sweep time, `" + conf.SweepTime + "` is invalid")
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer watcher.Close()

	fmt.Println("Sourcing archives from", conf.SourceDir)

	conf.Sweep()
	interval := sweepDuration
	ticker := time.NewTicker(interval)
	done := make(chan bool)

	go func() {
		fmt.Println("Sweeping every", sweepDuration)
		for {
			select {
			case <-ticker.C:
				conf.Sweep()
			}
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				for _, eventType := range conf.FileEvents {
					if event.Op.String() == eventType {
						conf.Copy(event.Name)
					}
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
		os.Exit(1)
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
		if m.AutoExtract {
			err := extract(filename)
			if err != nil {
				fmt.Println("error trying to extract skin: ", err)
				return
			}
			os.Remove(filename) // delete skin file after extraction
			filename = strings.TrimSuffix(filename, ".osk")
		}
	default:
		return
	}

	beatmapname := filepath.Base(filename)
	dest := filepath.Join(dir, beatmapname)
	fmt.Printf("Moving %s to %s\n", beatmapname, dest)
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

func extract(src string) error {
	dest := strings.TrimSuffix(src, ".osk")

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
    defer r.Close()

	for _, file := range r.File {
		fPath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		output, err := os.OpenFile(fPath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		cont, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(output, cont)
		cont.Close()
		output.Close()

		if err != nil {
			return err
		}
    }

    return nil
}
