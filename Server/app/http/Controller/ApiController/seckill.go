package ApiController

import (
	"AITranslatio/app/http/reposen"
	"AITranslatio/app/types/DTO"
	"github.com/gin-gonic/gin"
)

func (c *ApiController) StartSeckill(ctx *gin.Context) {

	userID := ctx.GetInt64("UserID")

	err := c.Service.StartSeckill(ctx, userID)
	if err != nil {
		reposen.ErrorSystem(ctx, err)
		return
	}
	reposen.OK(ctx, nil)
}

func (c *ApiController) SeckillBuy(ctx *gin.Context) {

	userID := ctx.GetInt64("UserID")

	BuyInfo := &DTO.BuyInfo{}
	err := ctx.Bind(BuyInfo)
	if err != nil {
		reposen.ErrorSystem(ctx, err)
		return
	}

	err = c.Service.SeckillBuy(ctx, userID, BuyInfo.GoodsID, BuyInfo.SeckillID)
	if err != nil {
		reposen.ErrorSystem(ctx, err)
		return
	}
	reposen.OK(ctx, nil)
}
