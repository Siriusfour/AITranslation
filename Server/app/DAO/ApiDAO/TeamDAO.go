package ApiDAO

import (
	"AITranslatio/Global"
	"AITranslatio/app/DAO"
	"AITranslatio/app/Model/Team"
	"fmt"
	"gorm.io/gorm"
)

type ApiDAO struct {
	DB_Client *gorm.DB
}

func CreateDAOFactory(sqlType string) *ApiDAO {
	return &ApiDAO{
		DB_Client: DAO.ChooseDB_Conn(sqlType),
	}
}

func (ApiDAO *ApiDAO) CreateTeam(leaderID int64, teamName string, introduction string) error {

	CreateTea := Team.Team{
		LeaderID:     leaderID,
		TeamName:     teamName,
		Introduction: introduction,
	}

	result := Global.MySQL_Client.Create(&CreateTea)
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

	result := Global.MySQL_Client.Create(&JoinTeam)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
