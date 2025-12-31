package goods

import (
	"time"
)

// SeckillGoods 秒杀商品表
type SeckillGoods struct {
	SeckillID    int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	GoodsID      int64     `gorm:"not null;index" json:"goods_id"`
	SeckillPrice float64   `gorm:"type:decimal(10,2);not null" json:"seckill_price"`
	StockCount   int       `gorm:"not null" json:"stock_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Version      int       `json:"version"`
}

// SeckillOrder 秒杀订单表
type SeckillOrder struct {
	UserID         int64   `gorm:"not null;uniqueIndex:u_uid_gid" json:"user_id"` // 联合唯一索引 part 1
	OrderID        int64   `gorm:"not null" json:"order_id"`
	GoodsID        int64   `gorm:"not null" json:"goods_id"`
	SeckillGoodsID int64   `gorm:"not null;uniqueIndex:u_uid_gid" json:"seckill_goods_id"` // 联合唯一索引 part 2
	Money          float64 `gorm:"type:decimal(10,2)" json:"money"`
	Status         int     `gorm:"default:0" json:"status"` // 0:未支付
}

func (SeckillGoods) TableName() string { return "seckill_goods" }
func (SeckillOrder) TableName() string { return "seckill_order" }
