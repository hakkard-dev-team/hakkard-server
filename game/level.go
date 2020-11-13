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

func (l *Level) GetLevelExits(c Client, g Game) {
	playerLoc, ok := g.GetLevel(c.Player.Location)
	var exits = ""
	if ok {
		for _, dir := range playerLoc.Exits {
			_, ok := g.GetLevel(dir.Target)
			if ok {
				exits += fmt.Sprintf("When you look %s you see %s. ", dir.Direction, dir.Target)
			}
		}
	}
	c.WriteLineToUser(exits)
}
