package env

import (
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

var Conf *conf = &conf{}

type (
	conf struct {
		System      *System
		Server      *Server
		PostgresSQL *PostgresSQL
		Redis       *Redis
		Logs        *LogsConfig
	}

	Mode string

	System struct {
		DatabaseDriver string `env:"DATABASE_DRIVER" default:"postgres"`
		Env            string `env:"ENV" default:"prod"`
	}

	Server struct {
		ServerUrl string `env:"SERVER_URL" default:"https://chat.aak1247.hive-intel.com"` // 外部访问地址
		Port      string `env:"SERVER_PORT" default:"8080"`                               // 服务兼容端口
		Ip        string `env:"SERVER_IP" default:"0.0.0.0"`                              // 服务监听ip
		Mode      Mode   `env:"mode" default:"dev"`
	}

	PostgresSQL struct {
		Host            string `env:"POSTGRESQL_HOST" default:"127.0.0.1"`
		Port            int    `env:"POSTGRESQL_PORT" default:"5432"`
		User            string `env:"POSTGRESQL_USER" default:"pgsu"`
		Password        string `env:"POSTGRESQL_PASSWORD" default:"password"`
		Dbname          string `env:"POSTGRESQL_DBNAME" default:"routers_pub"`
		Sslmode         string `env:"POSTGRESQL_SSLMODE" default:"disable"`
		MaxIdleConns    int    `env:"POSTGRESQL_MAXIDLECONNS" default:"10"`
		MaxOpenConns    int    `env:"POSTGRESQL_MAXOPENCONNS" default:"30"`
		ConnMaxLifeTime string `env:"POSTGRESQL_CONNMAXLIFETIME" default:"60s"`
		LogSwitch       bool   `env:"LOGSWITCH" default:"true"`
	}

	Redis struct {
		Host       string        `env:"REDIS_HOST" default:"127.0.0.1"`
		Port       string        `env:"REDIS_PORT" default:"6379"`
		Password   string        `env:"REDIS_PASSWORD" default:""`
		ExpireTime time.Duration `env:"REDIS_CACHE_EXPIRETIME" default:"300s"`
		MaxIdle    int           `env:"REDIS_MAXIDLE" default:"1024"`
		MaxActive  int           `env:"REDIS_MAXACTIVE " default:"1024"`
	}

	LogsConfig struct {
		Level      zapcore.Level `env:"LOGS_LEVEL" default:"0"`
		Path       string        `env:"LOGS_PATH" default:"./log"`
		MaxSize    int           `env:"MAX_SIZE" default:"1024"`
		MaxBackups int           `env:"MAX_BACKUPS" default:"10"`
		MaxAge     int           `env:"MAX_AGE" default:"7"`
		Compress   bool          `env:"COMPRESS"  default:"true"`
	}
)

const (
	RELEASE Mode = "release"
	DEV     Mode = "dev"
)

func GetEnv() *conf {
	Conf.getSystem()
	Conf.getServer()
	Conf.getPostgresSql()
	Conf.getRedis()
	Conf.getLogsConfig()
	env.Fill(Conf)
	return Conf
}

func (c *conf) getSystem() {
	if c.System != nil {
		return
	}
	system := new(System)
	if err := env.Fill(system); err != nil {
		panic(err)
	}
	c.System = system
}

func (c *conf) getServer() {
	if c.Server != nil {
		return
	}
	server := new(Server)
	if err := env.Fill(server); err != nil {
		panic(err)
	}
	c.Server = server
}

func (c *conf) getPostgresSql() {
	if c.PostgresSQL != nil {
		return
	}
	postgreSql := new(PostgresSQL)
	if err := env.Fill(postgreSql); err != nil {
		panic(err)
	}
	c.PostgresSQL = postgreSql
}

func (c *conf) getRedis() {
	if c.Redis != nil {
		return
	}
	redis := new(Redis)
	if err := env.Fill(redis); err != nil {
		panic(err)
	}
	c.Redis = redis
}

func (c *conf) GetAddr() string {
	return c.Server.Ip + ":" + c.Server.Port
}

func (mode Mode) IsProd() bool {
	return strings.Contains(string(mode), "prod") || strings.Contains(string(mode), "release")
}

func (c *conf) getLogsConfig() {
	if c.Logs != nil {
		return
	}
	logsConfig := new(LogsConfig)
	if err := env.Fill(logsConfig); err != nil {
		panic(err)
	}
	c.Logs = logsConfig
}
