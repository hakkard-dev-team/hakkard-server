package game

type Level struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Exits       []Exit `json:"exits"`
	Description string `json"description"`
}

type Exit struct {
	Target    string `json:"target"`
	IsHidden  bool   `json:"isHidden"`
	Direction string `json:"direction"`
}
