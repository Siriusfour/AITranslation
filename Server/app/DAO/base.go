package DAO

//func ChooseDB_Conn(logger *zap.Logger, sqlType string) *gorm.DB {
//
//	var DB_Client *gorm.DB
//	sqlType = strings.Trim(sqlType, " ")
//	switch strings.ToLower(sqlType) {
//
//	case "mysql":
//
//		if Global.MySQL_Client == nil {
//			logger.Error(MyErrors.ErrorsGormNotInitGlobalPointer, zap.String("sqlType", sqlType))
//		}
//
//		DB_Client = Global.MySQL_Client
//
//	case "postgressql":
//		if Global.PostgreSQL_Client == nil {
//			logger.Error(MyErrors.ErrorsGormNotInitGlobalPointer, zap.String("sqlType", sqlType))
//		}
//
//		DB_Client = Global.PostgreSQL_Client
//
//	}
//	return DB_Client
//}
