package model

import (
	"context"
	"fmt"
	"time"
)

type TeamsDB interface {
	SaveTeam(ctx context.Context, t *TeamDAO) (int, error)
}

type TeamDAO struct {
	ID          int       `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Name        string    `db:"name"`
	Region      string    `db:"region"`
	University  string    `db:"university"`
	FaceitLink  string    `db:"faceit_link"`
	IsConfirmed bool      `db:"is_confirmed"`
}

func (d *db) SaveTeam(ctx context.Context, t *TeamDAO) (int, error) {
	if t == nil {
		return 0, nil
	}

	err := d.conn.QueryRowxContext(ctx, `insert into teams (
	   created_at,
	   updated_at,
	   name,
	   region,
	   university,
	   faceit_link,
	   is_confirmed
	) values (
		current_timestamp,
		current_timestamp,
		$1,
		$2,
	    $3,
	    $4,
	    $5
	) on conflict (name, region, university, faceit_link) do nothing
		returning id`, t.Name, t.Region, t.University, t.FaceitLink, t.IsConfirmed).Scan(&t.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to save team: %s", err)
	}

	return t.ID, nil
}
