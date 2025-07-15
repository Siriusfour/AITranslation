package BaseDAO

import (
	"AITranslatio/Src/DTO"
	"AITranslatio/Src/Model"
)

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
