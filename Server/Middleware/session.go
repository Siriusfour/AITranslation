package Middleware

import (
	"AITranslatio/Global"
	"AITranslatio/app/http/reposen"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func SessionID(app *Global.Infrastructure) gin.HandlerFunc {
	return func(c *gin.Context) {

		SessionID, err := c.Cookie("SessionID")
		if err != nil {
			SessionID = ""
		}

		if SessionID == "" { //cookie里面也没有sessionID , 说明是该设备的第一个请求， 创建一个sessionID
			SessionID = app.SnowflakeManager.GetIDString()
			c.Set("SessionID", SessionID)
		} else { //有的话转化成int64填入ctx
			SessionID_int64, err := strconv.ParseInt(SessionID, 10, 64)
			if err != nil {
				reposen.ErrorSystem(c, fmt.Errorf("生成sessionID失败：%w", err))
			}
			c.Set("SessionID", SessionID_int64)
		}
		c.SetCookie("SessionID", SessionID, 0, "/", "", false, true)
	}
}
