package DTO

type RegisterDTO struct {
	UserID    int64
	UserName  string
	Password  string
	Salt      string
	Email     string
	EamilCode string

	//验证webAuthn随机数的数据
	Verify    string
	Timestamp int64
	Domain    string
}
