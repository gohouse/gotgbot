package config

import (
	"github.com/go-redis/redis/v7"
	"github.com/gohouse/gorose/v2"
	"github.com/sirupsen/logrus"
)
var conf *ConfigOption
type BotOption struct {
	Ticker int	// 多久自动发一次(min),default:10
}
type ConfigOption struct {
	TgOption  TgOption
	Mysql     gorose.Config
	Redis     redis.Options
	LogOption logrus.Logger
	BotOption BotOption
}

func Init(conf2 *ConfigOption) {
	//logInit(&conf2.LogOption)
	if conf2.Redis.Addr != "" {
		redisInit(&conf2.Redis)
	}
	if conf2.Mysql.Dsn != "" {
		mysqlInit(&conf2.Mysql)
	}
	tgbotapiInit(&conf2.TgOption)
	conf = conf2
}

func Config() *ConfigOption {
	return conf
}