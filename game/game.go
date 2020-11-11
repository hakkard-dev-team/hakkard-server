package game

import (
	log "github.com/Matt-Gleich/logoru"

	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Game struct {
	players      map[string]Player
	levels       map[string]Level
	DefaultLevel Level
}

func InitGame() *Game {
	log.Info("Initializing Game...")
	game := &Game{
		players: make(map[string]Player),
		levels:  make(map[string]Level),
	}

	game.initLevels()

	return game
}

func (g Game) initLevels() error {
	log.Info("Initializing Levels...")

	levelWalker := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileContent, fileIoErr := ioutil.ReadFile(path)
		if fileIoErr != nil {
			log.Warning("File %s could not be loader", path)
			log.Warning(fileIoErr.Error())
			return fileIoErr
		}

		level := Level{}
		if jsonErr := json.Unmarshal(fileContent, &level); jsonErr != nil {
			log.Warning("File %s could not be parsed", path)
			log.Warning(jsonErr.Error())
			return jsonErr
		}
		log.Debug("Loaded level %s", info.Name())
		return nil
	}
	return filepath.Walk("./content/levels/", levelWalker)
}
