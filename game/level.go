package game

import (
	"fmt"
)

type Level struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Exits       []Exit `json:"exits"`
	Description string `json:"description"`
}

type Exit struct {
	Target    string `json:"target"`
	IsHidden  bool   `json:"isHidden"`
	Direction string `json:"direction"`
}

func (l *Level) OnEnterRoom(g *Game, c Client) {
	c.WriteLineToUser(fmt.Sprintf("You are at %s", l.Name))

	if l.Description != "" {
		c.WriteLineToUser(fmt.Sprintf("> %s", l.Description))
	}

}
