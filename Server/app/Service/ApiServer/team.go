package ApiServer

import (
	"AITranslatio/app/DAO/ApiDAO"
	DTO "AITranslatio/app/http/validator/validators/Team"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ApiServer struct{}

func CreateApiServer() *ApiServer {
	return &ApiServer{}
}

func (Server *ApiServer) CreateTeam(ctx *gin.Context) error {

	teamName := ctx.GetString("TeamName")
	leaderID := ctx.GetInt64("UserID")
	Introduction := ctx.GetString("Introduction")

	dto := &DTO.TeamDTO{
		teamName,
		leaderID,
		Introduction,
	}

	err := ApiDAO.CreateDAOFactory("mysql").CreateTeam(dto)
	if err != nil {
		return fmt.Errorf("create team err: %v", err)
	}

}
