package ApiServer

import (
	"AITranslatio/Global/Consts"
	"fmt"
	"github.com/gin-gonic/gin"
)

func CreateApiServer() *ApiServer {
	return &ApiServer{}
}

func (Server *ApiServer) CreateTeam(ctx *gin.Context) error {

	teamName := ctx.GetString(Consts.ValidatorPrefix + "TeamName")
	leaderID := ctx.GetInt64(Consts.ValidatorPrefix + "UserID")
	Introduction := ctx.GetString(Consts.ValidatorPrefix + "Introduction")

	//TODO 查询该用户已创建多少团体,多于100（配置文件可更改）则拒绝
	err := Server.DAO.CreateTeam(leaderID, teamName, Introduction)
	if err != nil {
		return fmt.Errorf("create team err: %w", err)
	}

	return nil

}

func (Server *ApiServer) JoinTeam(ctx *gin.Context) error {

	FromUserID := ctx.GetInt64(Consts.ValidatorPrefix + "UserID")
	Introduction := ctx.GetString(Consts.ValidatorPrefix + "Introduction")
	//NickName:=ctx.GetString(Consts.ValidatorPrefix + "NickName"),
	TeamID := ctx.GetInt(Consts.ValidatorPrefix + "TeamID")

	//写入数据库
	err := Server.DAO.JoinTeam(FromUserID, TeamID, Introduction)
	if err != nil {
		return fmt.Errorf("DAO层调用JoinTeam失败: %w", err)
	}

	//投放到消息队列，告诉用户“申请成功”，消费者会做“websoket/SSE通知，日志统计"

	return nil
}
