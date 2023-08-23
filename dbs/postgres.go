package dbs

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"routers.pub/env"
	"routers.pub/infra"
	"routers.pub/utils"
	"time"
)

const (
	POSTGRES = "postgres"
)

// DB 全局数据库变量
var DB *gorm.DB

func InitDatabase() {
	switch env.Conf.System.DatabaseDriver {
	case POSTGRES:
		InitPostgres()
	default:
		InitPostgres()
	}

	db, _ := DB.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
}

// InitPostgres 初始PGSQL数据库
func InitPostgres() {
	//createPostgresDatabase(env.Conf.PostgresSQL.Dbname)
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		env.Conf.PostgresSQL.Host,
		env.Conf.PostgresSQL.Port,
		env.Conf.PostgresSQL.User,
		env.Conf.PostgresSQL.Password,
		env.Conf.PostgresSQL.Dbname,
	)
	// 隐藏密码
	showDsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=****** dbname=%s sslmode=disable",
		env.Conf.PostgresSQL.Host,
		env.Conf.PostgresSQL.Port,
		env.Conf.PostgresSQL.User,
		env.Conf.PostgresSQL.Dbname,
	)
	// SQL日志配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // 彩色打印
		},
	)
	// Log.Info("数据库连接DSN: ", showDsn)
	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{
			// 禁用外键
			DisableForeignKeyConstraintWhenMigrating: true,
			// 打印SQL日志
			Logger: newLogger,
			// // 指定表前缀
			// NamingStrategy: schema.NamingStrategy{
			//	TablePrefix: env.Conf.Postgres.TablePrefix + "_",
			// },
		},
	)
	if err != nil {
		infra.Log.Panicf("初始化Postgres数据库异常: %v", err)
		panic(fmt.Errorf("初始化Postgres数据库异常: %v", err))
	}

	// 开启日志
	if env.Conf.PostgresSQL.LogSwitch {
		db = db.Debug()
	}
	//db.Callback().Create().Before("gorm:create").Register("setupId", setupId)
	// 全局DB赋值
	DB = db
	infra.Log.Infof("初始化Postgres数据库完成!\ndsn: %s", showDsn)
}

//func createPostgresDatabase(dbname string) {
//	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s db= sslmode=disable TimeZone=Asia/Shanghai",
//		env.Conf.PostgresSQL.Host,
//		env.Conf.PostgresSQL.Port,
//		env.Conf.PostgresSQL.User,
//		env.Conf.PostgresSQL.Password,
//	)
//	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//
//	createDatabaseCommand := fmt.Sprintf("CREATE DATABASE %s", dbname)
//	db.Exec(createDatabaseCommand)
//}

func Release(db *gorm.DB, respErr *error) {
	func() {
		if err := recover(); err != nil {
			infra.Log.Errorf("[ERROR-ALERT] api panic: %v\n Stack:%s", err, utils.CallStack(20, 1))
			infra.Log.Errorf("事务异常, 回滚: %v", err)
			db.Rollback()
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			panic(err)
		} else if *respErr != nil {
			infra.Log.Errorf("事务异常, 回滚: %v", respErr)
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
}

func GetDb() *gorm.DB {
	return DB
}

func StartTx(opts ...*sql.TxOptions) *gorm.DB {
	return DB.Begin(opts...)
}
