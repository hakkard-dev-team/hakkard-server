package main

import (
	"bufio"
	"fmt"
	"io"

	log "github.com/Matt-Gleich/logoru"

	"encoding/json"
	"io/ioutil"
	"net"
	"os"

	"github.com/hakkard-dev-team/hakkard-server/game"
)


type ServerConfig struct {
	MaxPlayers      int    `json:"maxPlayers"`
	Bind            string `json:"bind"`
	Name            string `json:"name"`
	DefaultLevelKey string `json:"defaultLevel"`
}

func main() {
	log.Info("Parsing Config")
	conf, err := parseConfig()
	if err != nil {
		log.Critical("Could not parse config: " + err.Error())
		os.Exit(-1)
	}
	game := game.InitGame(conf.DefaultLevelKey)

	ln, err := net.Listen("tcp", conf.Bind)
	if err != nil {
		log.Critical(fmt.Sprintf("Error binding network port: %v", err.Error()))
	}
	log.Success(fmt.Sprintf("Listening on %s", ln.Addr()))

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Warning(fmt.Sprintf("Error accepting new connection %v", err))
			continue
		}

		msgchan := make(chan string)
		go handleConnection(conn, game, msgchan)
	}

}

func promptMessage(c net.Conn, bufc *bufio.Reader, message string) string {
	for {
		io.WriteString(c, message+"\n\r")
		answer, _, _ := bufc.ReadLine()
		if string(answer) != "" {
			return string(answer)
		}
	}
}

func handleConnection(c net.Conn, g *game.Game, msgchan chan<- string) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	log.Info(fmt.Sprintf("New Connection from %s", c.RemoteAddr()))

	questions := 0
	var name string
	for {
		if questions >= 3 {
			io.WriteString(c, "See you!\n")
			return
		}

		name = promptMessage(c, bufc, "What name do you wish?")
		ok, plr := g.LoadPlayer(name)
		if ok == false {
			questions++
			io.WriteString(c, fmt.Sprintf("Username %s does not exist\n\r", name))
			answer := promptMessage(c, bufc, "Create it? [y|n]")

			if answer == "y" {
				password := promptMessage(c, bufc, "What is your desired password?")
				io.WriteString(c, "We offer multiple player types. Please choose from this list\n- Mage: Magic User\n-Fighter: Wields Weapons\n\r")
				playerType := promptMessage(c, bufc, "What type do you wish? [mage|fighter]")
				g.CreatePlayer(name, playerType, password)
				break
			}
		} else {
			password := promptMessage(c, bufc, "Password: ")
			pwOk := plr.VerifyPassword(password)
			if !pwOk {
				io.WriteString(c, "Invalid password.")
				continue
			}
			break
		}
	}

	player, ok := g.GetPlayerByName(name)
	if !ok {
		log.Warning("Error getting Player object")
		io.WriteString(c, "Error getting Player object\n\r")
		return
	}

	client := game.NewClient(c, player)

	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\r", client.Player.Name))

	location, ok := g.GetLevel(client.Player.Location)
	if ok {
		location.OnEnterRoom(g, client)
	}

	// Handle I/O
	go client.ReadLinesInto(msgchan, g)
	client.WriteLinesFrom(client.Chan)

}

func parseConfig() (ServerConfig, error) {
	data, err := ioutil.ReadFile("static/server.json")
	if err != nil {
		return ServerConfig{}, err
	}
	bytes := []byte(data)
	var conf ServerConfig
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return ServerConfig{}, err
	}

	return conf, nil
}
