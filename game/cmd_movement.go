package game

import (
	"strings"

	"github.com/Matt-Gleich/logoru"
)

func MovementCommands(router *Route) {
	logoru.Info("Movement commands...")

	router.On("go", func(ctx *CmdContext) {
		playerLoc, ok := ctx.Game.GetLevel(ctx.Client.Player.Location)
		if ok {
			for _, oneDir := range playerLoc.Exits {
				if strings.ToLower(oneDir.Direction) == strings.ToLower(ctx.Args) {
					target, ok := ctx.Game.GetLevel(oneDir.Target)
					if ok {
						ctx.Client.Player.Location = string(target.Key)
						ctx.Game.SavePlayer(ctx.Client.Player)
						target.OnEnterRoom(ctx.Game, *ctx.Client)
					}
				}
			}
		}
	})
}
