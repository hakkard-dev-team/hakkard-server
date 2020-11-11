package game

import (
	"bufio"
	"io"
	"net"
	"strings"
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
