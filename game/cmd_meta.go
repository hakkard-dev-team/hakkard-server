package game

import log "github.com/Matt-Gleich/logoru"

func MetaCommands(router *Route) {
	log.Debug("Meta commands...")
	router.On("exit", func(ctx *CmdContext) {
		ctx.Game.OnExit(*ctx.Client)
		ctx.Client.Conn.Close()
	}).Desc("Exits the game")
}
