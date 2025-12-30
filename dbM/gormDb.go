package dbM

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func MySqlConnectV2(conf string, logLevel int) *gorm.DB {
	var (
		err error
		res *gorm.DB
	)
	res, err = gorm.Open(mysql.Open(conf), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
	})
	if err != nil {
		panic(err.Error())
	}

	sqlDB, err := res.DB()
	if err != nil {
		panic("mysql connect err," + conf + "," + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(80) // 起步值，按实例数/DB上限再调
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	logx.Info("mysql connect success")
	return res
}

func RedisConnect(redisConf *redis.Options) *redis.Client {
	return redis.NewClient(redisConf)
}
