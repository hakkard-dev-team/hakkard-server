package game

import (
	"fmt"

	log "github.com/Matt-Gleich/logoru"
)

func MetaCommands(router *Route) {
	log.Debug("Meta commands...")
	router.On("exit", func(ctx *CmdContext) {
		ctx.Game.OnExit(*ctx.Client)
		ctx.Client.Conn.Close()
	}).Desc("Exits the game")

	router.On("look", func(ctx *CmdContext) {
		playerLoc, ok := ctx.Game.GetLevel(ctx.Client.Player.Location)
		c := ctx.Client
		if ok {
			c.WriteLineToUser(playerLoc.Description)
			c.WriteLineToUser("")
			var exits = ""
			for _, v := range playerLoc.Exits {
				place := v.Target
				if ok {
					exits += fmt.Sprintf("When you look %s, you see %s", v.Direction, place)
				}
			}
			c.WriteLineToUser(exits)
		}
	}).Desc("Looks around")

}
