package domains

import (
	"routers.pub/dbs"
	"routers.pub/infra"
)

var tables []interface{}

func registerAutoMigrate(db interface{}) {
	tables = append(tables, db)
}

func InitDbTables() error {
	db := dbs.GetDb()
	//if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
	//	commons.Log.Error(err)
	//	return err
	//}
	// 自动迁移表结构
	err := db.AutoMigrate(
		tables...,
	)
	if err != nil {
		infra.Log.Error(err)
		return err
	}
	infra.Log.Infof("表初始化成功")
	return nil
}
