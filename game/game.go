package game

import (
	log "github.com/Matt-Gleich/logoru"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Game struct {
	players      map[string]Player
	levels       map[string]Level
	DefaultLevel Level
	Route        *Route
}

func InitGame(defaultLevelKey string) *Game {
	log.Info("Initializing Game...")
	game := &Game{
		players: make(map[string]Player),
		levels:  make(map[string]Level),
	}

	game.initLevels()

	router := NewRouter()
	game.initCommands(router)
	game.Route = router

	// Set Default Level
	defaultLevel, err := game.GetLevel(defaultLevelKey)
	if err != true {
		log.Critical("Invalid default level! Check static/server.json")
		os.Exit(-1)
	}
	game.DefaultLevel = defaultLevel

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
			log.Warning(fmt.Sprintf("File %s could not be loader", path))
			log.Warning(fileIoErr.Error())
			return fileIoErr
		}

		level := Level{}
		if jsonErr := json.Unmarshal(fileContent, &level); jsonErr != nil {
			log.Warning(fmt.Sprintf("File %s could not be parsed", path))
			log.Warning(jsonErr.Error())
			return jsonErr
		}
		log.Debug(fmt.Sprintf("Loaded level %s", info.Name()))
		g.addLevel(level)
		return nil
	}
	return filepath.Walk("./content/levels/", levelWalker)
}

// Initializes commands
func (g *Game) initCommands(router *Route) error {
	log.Info("Initializing Commands...")
	MetaCommands(router)
	MovementCommands(router)

	return nil
}

func (g *Game) GetLevel(key string) (Level, bool) {
	level, ok := g.levels[key]
	return level, ok
}

func (g Game) addLevel(lvl Level) error {
	g.levels[lvl.Key] = lvl
	return nil
}
func (g Game) addPlayer(plr Player) error {
	g.players[plr.Name] = plr
	return nil
}

func (g *Game) getPlayerFileName(playerName string) string {
	return "./static/data/" + playerName + ".player"
}

func (g *Game) LoadPlayer(playerName string) (bool, Player) {
	playerFileLocation := g.getPlayerFileName(playerName)

	log.Debug(fmt.Sprintf("Loading player %s", playerName))
	fileContent, fileIoErr := ioutil.ReadFile(playerFileLocation)
	if fileIoErr != nil {
		log.Warning(fmt.Sprintf("Could not load player file %s: %v", playerFileLocation, fileIoErr))
		return false, Player{}
	}
	player := Player{}
	if jsonErr := json.Unmarshal(fileContent, &player); jsonErr != nil {
		log.Warning(fmt.Sprintf("Invalid player file %s: %v", playerFileLocation, jsonErr))
		return false, Player{}
	}
	log.Success(fmt.Sprintf("Loaded player %s", playerName))
	g.addPlayer(player)

	return true, player
}

func (g *Game) CreatePlayer(playerName string, playerType string, pw string) {
	playerFileLocation := g.getPlayerFileName(playerName)

	log.Debug(fmt.Sprintf("Creating player %s", playerName))
	if _, err := os.Stat(playerFileLocation); err == nil {
		g.LoadPlayer(playerName)
		log.Debug(fmt.Sprintf("Player %s already exists, loading...", playerName))
		return
	}
	player := Player{
		Name:       playerName,
		PlayerType: playerType,
		Location:   g.DefaultLevel.Key,
	}
	g.addPlayer(player)
	player.InitDefaultAttributes()
	player.SetPassword(pw)
	g.SavePlayer(player)
}

func (g *Game) GetPlayerByName(playerName string) (Player, bool) {
	player, ok := g.players[playerName]
	return player, ok
}

func (g *Game) SavePlayer(player Player) bool {
	data, err := json.MarshalIndent(player, "", "	")
	if err == nil {
		playerFileLocation := g.getPlayerFileName(player.Name)

		if ioerror := ioutil.WriteFile(playerFileLocation, data, 0600); ioerror != nil {
			log.Warning(fmt.Sprintf("Could not save player %s: %v", player.Name, ioerror))
			return false
		}
	} else {
		log.Warning(fmt.Sprintf("Could not Marshal player data %s: %v", player.Name, err))
		return false
	}
	return true
}

func (g *Game) OnExit(client Client) {
	g.SavePlayer(client.Player)
	client.WriteLineToUser("Bye!")
}
