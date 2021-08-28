package api

import (
	"asscRegsitration/model"
	"time"

	"github.com/go-openapi/strfmt"
)

type Team struct {
	ID          int       `json:"id,omitempty" example:"1"`
	Name        string    `json:"name"`
	Region      string    `json:"region"`
	University  string    `json:"university"`
	FaceitLink  string    `json:"faceit_link"`
	Players     []*Player `json:"players"`
	IsConfirmed bool      `json:"-"`
}

func (t Team) DAO() (*model.TeamDAO, []*model.PlayerDAO) {

	teamDAO := &model.TeamDAO{
		ID:          t.ID,
		Name:        t.Name,
		Region:      t.Region,
		University:  t.University,
		FaceitLink:  t.FaceitLink,
		IsConfirmed: t.IsConfirmed,
	}

	playersDAO := make([]*model.PlayerDAO, 0, len(t.Players))
	for _, player := range t.Players {
		playersDAO = append(playersDAO, &model.PlayerDAO{
			ID:          player.ID,
			TeamID:      player.TeamID,
			IsCaptain:   player.IsCaptain,
			FirstName:   player.FirstName,
			SecondName:  player.SecondName,
			LastName:    player.LastName,
			Birthdate:   time.Time(*player.Birthdate),
			PhoneNumber: player.PhoneNumber,
			Email:       player.Email,
			DiscordID:   player.DiscordID,
			Document:    player.Document,
		})
	}

	return teamDAO, playersDAO

}

type Player struct {
	ID          int          `json:"id,omitempty" example:"1"`
	TeamID      int          `json:"team_id,omitempty"`
	IsCaptain   bool         `json:"is_captain,omitempty"`
	FirstName   string       `json:"first_name"`
	SecondName  string       `json:"second_name,omitempty"`
	LastName    string       `json:"last_name"`
	Birthdate   *strfmt.Date `json:"birthdate"`
	PhoneNumber string       `json:"phone_number,omitempty"`
	Email       string       `json:"email,omitempty"`
	DiscordID   string       `json:"discord_id,omitempty"`
	Document    string       `json:"document,omitempty"`
}
