package game

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/Matt-Gleich/logoru"
)

type Client struct {
	Conn   net.Conn
	Name   string
	Player Player
	Chan   chan string
}

func NewClient(c net.Conn, player Player) Client {
	return Client{
		Conn:   c,
		Name:   player.Name,
		Player: player,
		Chan:   make(chan string),
	}
}

func (c Client) WriteLineToUser(msg string) {
	io.WriteString(c.Conn, msg+"\n\r")
}

func (c Client) ReadLines(ch chan<- string) {
	bufc := bufio.NewReader(c.Conn)

	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}

		userLine := strings.TrimSpace(line)

		if userLine == "" {
			continue
		}
		c.WriteLineToUser("You wrote: " + userLine)
	}
}

func (c Client) ReadLinesInto(ch chan<- string, g *Game) {
	bufc := bufio.NewReader(c.Conn)

	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}
		userLine := strings.TrimSpace(line)

		if userLine == "" {
			continue
		}

		lineParts := strings.SplitN(userLine, " ", 2)
		var command, commandText string
		if len(lineParts) > 0 {
			command = lineParts[0]
		}
		if len(lineParts) > 1 {
			commandText = lineParts[1]
		}

		log.Debug(fmt.Sprintf("Command by %s: %s %s", c.Player.Name, command, commandText))

		if ok := g.Route.FindAndExecute(g, &c, command, commandText); !ok {
			c.WriteLineToUser("Huh?")
		}

		/*		switch command {
		case "look":
			playerLoc, ok := g.GetLevel(c.Player.Location)
			if ok {
				for _, dir := range playerLoc.Exits {
					place, ok := g.GetLevel(dir.Target)
					if ok {
						c.WriteLineToUser(fmt.Sprintf("When you look %s you see %s", dir.Direction, place.Name))
					}
				}
			}
			c.WriteLineToUser(playerLoc.Description)
		default:
			continue
		}*/
	}
}

func (c Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.Conn, msg)
		if err != nil {
			return
		}
	}
}
