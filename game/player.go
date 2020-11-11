package game

type Player struct {
	Name       string      `json:"name"`
	Location   string      `json:"location"`
	Attributes []Attribute `json:"attributes"`
	PlayerType string      `json:"playerType"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (p *Player) InitDefaultAttributes() {
	p.Attributes = []Attribute{
		Attribute{
			Name:  "Test",
			Value: 100,
		},
	}
}
