package ApiDAO

import (
	"AITranslatio/app/Model/Team"
	"fmt"
)

func (ApiDAO *ApiDAO) CreateTeam(leaderID int64, teamName string, introduction string) error {

	CreateTea := Team.Team{
		LeaderID:     leaderID,
		TeamName:     teamName,
		Introduction: introduction,
	}

	result := ApiDAO.DB_Client.Create(&CreateTea)
	if result.Error != nil {
		return fmt.Errorf("DAO层CreateTeam调用失败:%w", result.Error)
	}

	return nil

}

func (ApiDAO *ApiDAO) JoinTeam(FromUserID int64, TeamID int, Introduction string) error {

	JoinTeam := Team.JoinTeamApplication{
		FromUserID:   FromUserID,
		TeamID:       TeamID,
		Introduction: Introduction,
		Status:       0,
	}

	result := ApiDAO.DB_Client.Create(&JoinTeam)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
