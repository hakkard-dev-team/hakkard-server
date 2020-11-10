package main

import (
	log "github.com/Matt-Gleich/logoru"

	"encoding/json"
	"io/ioutil"
	"os"
)

type Server struct {
	Config ServerConfig
}

type ServerConfig struct {
	MaxPlayers int    `json:"maxPlayers"`
	Bind       string `json:"bind"`
	Name       string `json:"name"`
}

func main() {
	log.Info("Parsing Config")
	conf, err := parseConfig()
	if err != nil {
		log.Critical("Could not parse config: " + err.Error())
		os.Exit(-1)
	}
	log.Debug(conf.MaxPlayers)
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
