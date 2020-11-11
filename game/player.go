package game

type Player struct {
	Name       string      `json:"name"`
	Location   string      `json:"location"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
