package PublicKey

type WebAthun_PublicKey struct {
	ID  int64  `gorm:"type:int64;not null;column:UserID"`
	Key string `gorm:"type:varchar(128);not null"`
}

func (P *WebAthun_PublicKey) TableName() string {
	return "WebAthun_PublicKey"
}
