package goods

import (
	"gorm.io/gorm"
)

// Goods 商品基本表
type Goods struct {
	gorm.Model
	GoodsName   string  `gorm:"size:128;not null" json:"goods_name"`
	GoodsTitle  string  `gorm:"size:256" json:"goods_title"`
	GoodsImg    string  `gorm:"size:256" json:"goods_img"`
	GoodsDetail string  `gorm:"type:longtext" json:"goods_detail"`
	GoodsPrice  float64 `gorm:"type:decimal(10,2);not null" json:"goods_price"`
	GoodsStock  int     `gorm:"not null" json:"goods_stock"`
}

func (Goods) TableName() string { return "goods" }
