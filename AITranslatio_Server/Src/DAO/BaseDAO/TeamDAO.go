package BaseDAO

import (
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/Model"
)

func (BaseDAO *BaseDAO) CreateTeam(CreateTeamDTO *DTO.CreateTeamDTO) error {

	CreateTea := Model.Team{
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

	JoinTeam := Model.JoinTeamApplication{
		ApplicationID:       JoinTeamDTO.UserID,
		ApplicationNickName: JoinTeamDTO.NickName,
		TeamID:              JoinTeamDTO.JoinTeamID,
		Introduction:        JoinTeamDTO.Introduction,
		Status:              0,
	}

	result := BaseDAO.orm.Create(&JoinTeam)
	if result != nil {
		return result.Error
	}

	return nil
}
