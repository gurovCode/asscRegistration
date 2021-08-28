package model

import (
	"context"
	"fmt"
	"time"
)

type PlayersDB interface {
	SavePlayer(ctx context.Context, p *PlayerDAO) error
}

type PlayerDAO struct {
	ID          int       `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	TeamID      int       `db:"team_id"`
	IsCaptain   bool      `db:"is_captain"`
	FirstName   string    `db:"first_name"`
	SecondName  string    `db:"second_name"`
	LastName    string    `db:"last_name"`
	Birthdate   time.Time `db:"birthdate"`
	PhoneNumber string    `db:"phone_number"`
	Email       string    `db:"email"`
	DiscordID   string    `db:"discord_id"`
	Document    string    `db:"document"`
}

func (d *db) SavePlayer(ctx context.Context, p *PlayerDAO) error {
	if p == nil {
		return nil
	}

	_, err := d.conn.NamedExecContext(ctx, `insert into players (
         created_at, 
	     updated_at, 
		 team_id, 
		 is_captain, 
		 first_name, 
		 second_name, 
		 last_name, 
		 birthdate, 
		 phone_number, 
		 email, 
		 discord_id, 
		 document
	) values (
		current_timestamp,
		current_timestamp,
		:team_id,
		:is_captain,
	    :first_name,
	    :second_name,
	    :last_name,
	    :birthdate,
	    :phone_number,
	    :email,
	    :discord_id,
	    :document
	)`, &p)
	if err != nil {
		return fmt.Errorf("failed to save player: %s", err)
	}

	return nil
}
