package game

import (
	"fmt"

	log "github.com/Matt-Gleich/logoru"
	"golang.org/x/crypto/bcrypt"
)

type Player struct {
	Name       string      `json:"name"`
	Location   string      `json:"location"`
	Attributes []Attribute `json:"attributes"`
	PlayerType string      `json:"playerType"`
	PwHash     string      `json:"pwHash"`
	Level int `json:"level"`
	Currency int `json:"currency"`
	Inventory []string `json:"inventory"`
	Wearing map[string]string `json:"wearing"`
	Wielding []string `json:"wielding"`
	Experience int `json:"experience"`

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

func (p Player) VerifyPassword(pw string) bool {
	hashedPw := []byte(p.PwHash)
	bytePw := []byte(pw)

	err := bcrypt.CompareHashAndPassword(hashedPw, bytePw)
	if err != nil {
		log.Debug(fmt.Sprintf("Error verifying password: %v", err))
		return false
	}
	return true
}

func (p *Player) SetPassword(pw string) bool {
	bytePw := []byte(pw)

	hash, err := bcrypt.GenerateFromPassword(bytePw, bcrypt.MinCost)
	if err != nil {
		log.Warning(fmt.Sprintf("Error hashing password: %v", err))
		return false
	}
	p.PwHash = string(hash)
	return true
}
