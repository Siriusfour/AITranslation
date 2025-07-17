package BaseService

import (
	"AITranslatio/Global"
	"AITranslatio/Src/DTO"
)

func (BaseService *BaseService) CreateTeam(CreateTeamDTO *DTO.CreateTeamDTO) error {

	err := BaseService.BaseDAO.CreateTeam(CreateTeamDTO)
	if err != nil {
		return err
	}
	return nil

}

func (BaseService *BaseService) JoinTeam(JoinTeamDTO *DTO.JoinTeamDTO) error {

	//0.申请入库
	err := BaseService.BaseDAO.JoinTeam(JoinTeamDTO)
	if err != nil {
		return err
	}

	//1.向SSE推送信息
	Global.SSEClients.SendNotify(JoinTeamDTO.UserID, JoinTeamDTO.Introduction)

	return nil

}
