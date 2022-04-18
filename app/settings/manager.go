package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/karrick/godirwalk"
	"github.com/wieku/danser-go/framework/env"
	"github.com/wieku/danser-go/framework/files"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var fileStorage *fileformat
var filePath string
var watcher *fsnotify.Watcher

func initStorage() {
	fileStorage = &fileformat{
		General:     General,
		ApiV2:       ApiV2,
		Graphics:    Graphics,
		Audio:       Audio,
		Input:       Input,
		Gameplay:    Gameplay,
		Skin:        Skin,
		Cursor:      Cursor,
		Objects:     Objects,
		Playfield:   Playfield,
		CursorDance: CursorDance,
		Knockout:    Knockout,
		Recording:   Recording,
	}
}

func LoadSettings(version string) bool {
	if err := os.MkdirAll(env.ConfigDir(), 0755); err != nil {
		panic(err)
	}

	migrateSettings()

	initStorage()

	fileName := "default"
	if version != "" {
		fileName = version
	}

	fileName += ".json"

	filePath = filepath.Join(env.ConfigDir(), fileName)

	file, err := os.Open(filePath)

	if os.IsNotExist(err) {
		saveSettings(filePath, fileStorage, true)
		return true
	} else if err != nil {
		panic(err)
	} else {
		load(file, fileStorage)

		file.Close()

		saveSettings(filePath, fileStorage, false) // this is done to save additions from the current format
	}

	if !RECORD {
		setupWatcher(filePath)
	}

	return false
}

func setupWatcher(file string) {
	var err error

	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("SettingsManager: Detected", file, "modification, reloading...")

					time.Sleep(time.Millisecond * 200)

					sFile, _ := os.Open(event.Name)

					load(sFile, fileStorage)

					sFile.Close()

					saveSettings(filePath, fileStorage, false) // re-save the file if it contained backslash fixes
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("error:", err)
			}
		}
	}()

	abs, _ := filepath.Abs(file)

	err = watcher.Add(abs)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseWatcher() {
	if watcher != nil {
		err := watcher.Close()
		if err != nil {
			log.Println(err)
		}

		watcher = nil
	}
}

func load(file *os.File, target interface{}) {
	log.Println(fmt.Sprintf(`SettingsManager: Loading "%s"`, file.Name()))

	// I hope it won't backfire, replacing \ or \\\\\\\ with \\ so JSON can parse it as \

	data, err := io.ReadAll(files.NewUnicodeReader(file))
	if err != nil {
		panic(err)
	}

	tG := target.(*fileformat)
	tG.srcData = data

	str := string(data)
	str = regexp.MustCompile(`\\+`).ReplaceAllString(str, `\`)
	str = strings.ReplaceAll(str, `\`, `\\`)

	if err = json.Unmarshal([]byte(str), target); err != nil {
		panic(fmt.Sprintf("SettingsManager: Failed to parse %s! Please re-check the file for mistakes. Error: %s", file.Name(), err))
	}

	migrateCursorDance(target)
}

func Save() {
	saveSettings(filePath, fileStorage, false)
}

func saveSettings(path string, source interface{}, forceSave bool) {
	tG := source.(*fileformat)

	data, err := json.MarshalIndent(source, "", "\t")
	if err != nil {
		panic(err)
	}

	if forceSave || !bytes.Equal(data, tG.srcData) { // Don't rewrite the file unless necessary
		log.Println(fmt.Sprintf(`SettingsManager: Saving current settings to "%s"`, path))

		tG.srcData = data

		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(err)
		}

		if err = os.WriteFile(path, data, 0644); err != nil {
			panic(err)
		}
	}
}

func GetFormat() *fileformat {
	return fileStorage
}

func migrateSettings() {
	_ = godirwalk.Walk(env.DataDir(), &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if osPathname != env.DataDir() && de.IsDir() {
				return godirwalk.SkipThis
			}

			if !strings.HasSuffix(osPathname, ".json") || !strings.HasPrefix(osPathname, "settings") {
				return nil
			}

			var destName string

			if osPathname == "settings.json" {
				destName = "default.json"
			} else {
				destName = strings.TrimPrefix(osPathname, "settings-")
			}

			err := os.Rename(osPathname, filepath.Join(env.ConfigDir(), destName))
			if err != nil {
				panic(err)
			}

			return nil
		},
		Unsorted:            true,
		FollowSymbolicLinks: true,
	})
}

func migrateCursorDance(target interface{}) {
	tG := target.(*fileformat)

	if tG.Dance == nil {
		return
	}

	movers := make([]*mover, 0, len(tG.Dance.Movers))
	spinners := make([]*spinner, 0, len(tG.Dance.Spinners))

	for _, m := range tG.Dance.Movers {
		movers = append(movers, &mover{
			Mover:             m,
			SliderDance:       tG.Dance.SliderDance,
			RandomSliderDance: tG.Dance.RandomSliderDance,
		})
	}

	for _, m := range tG.Dance.Spinners {
		spinners = append(spinners, &spinner{
			Mover:  m,
			Radius: tG.Dance.SpinnerRadius,
		})
	}

	tG.CursorDance.Movers = movers
	tG.CursorDance.Spinners = spinners

	tG.CursorDance.Battle = tG.Dance.Battle
	tG.CursorDance.DoSpinnersTogether = tG.Dance.DoSpinnersTogether
	tG.CursorDance.TAGSliderDance = tG.Dance.TAGSliderDance

	tG.CursorDance.MoverSettings.Bezier = []*bezier{
		tG.Dance.Bezier,
	}

	tG.CursorDance.MoverSettings.Flower = []*flower{
		tG.Dance.Flower,
	}

	tG.CursorDance.MoverSettings.HalfCircle = []*circular{
		tG.Dance.HalfCircle,
	}

	tG.CursorDance.MoverSettings.Spline = []*spline{
		tG.Dance.Spline,
	}

	tG.CursorDance.MoverSettings.Momentum = []*momentum{
		tG.Dance.Momentum,
	}

	tG.CursorDance.MoverSettings.ExGon = []*exgon{
		tG.Dance.ExGon,
	}

	tG.Dance = nil
}
