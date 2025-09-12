package BaseDAO

import (
	"AITranslatio/app/Model/team"
	"AITranslatio/app/http/DTO"
)

func (BaseDAO *BaseDAO) CreateTeam(CreateTeamDTO *DTO.CreateTeamDTO) error {

	CreateTea := team.Team{
		LeaderID:     CreateTeamDTO.UserID,
		TeamName:     CreateTeamDTO.TeamName,
		Introduction: CreateTeamDTO.Introduction,
		NoteID:       CreateTeamDTO.NoteID,
	}

	result := BaseDAO.orm.Create(&CreateTea)
	if result != nil {
		return result.Error
	}

	return nil

}

func (BaseDAO *BaseDAO) JoinTeam(JoinTeamDTO *DTO.JoinTeamDTO) error {

	JoinTeam := team.JoinTeamApplication{
		ApplicationID:       JoinTeamDTO.UserID,
		ApplicationNickName: JoinTeamDTO.NickName,
		TeamID:              JoinTeamDTO.JoinTeamID,
		Introduction:        JoinTeamDTO.Introduction,
		Status:              0,
	}

	result := BaseDAO.orm.Create(&JoinTeam)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
