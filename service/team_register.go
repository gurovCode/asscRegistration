package service

import (
	api "asscRegsitration"
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (s *service) TeamRegister(ctx context.Context, team api.Team) (int, error) {

	teamDAO, playersDAO := team.DAO()

	teamId, err := s.db.SaveTeam(ctx, teamDAO)
	if err != nil {
		s.logger.Error("failed to save team", zap.Error(err))
		return 0, err
	}

	for k := range playersDAO {
		fmt.Println("PLAYER:", playersDAO[k])
		playersDAO[k].TeamID = teamId
		err = s.db.SavePlayer(ctx, playersDAO[k])
		if err != nil {
			s.logger.Error("failed to save player", zap.Error(err))
			return 0, err
		}
	}

	return teamId, nil
}

func (s *service) validatePlayer(player *api.Player) error {

	//player.
	return nil
}
