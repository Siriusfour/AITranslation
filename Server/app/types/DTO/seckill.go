package DTO

type BuyInfo struct {
	GoodsID   int64 `json:"goodsId" binding:"required"`
	SeckillID int64 `json:"seckillId" binding:"required"`
}
