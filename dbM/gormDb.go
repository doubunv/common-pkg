package dbM

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func MySqlConnect(conf string) *gorm.DB {
	var (
		err error
		res *gorm.DB
	)
	res, err = gorm.Open(mysql.Open(conf), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}

	sqlDB, err := res.DB()
	if err != nil {
		panic("mysql connect err," + conf + "," + err.Error())
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logx.Info("mysql connect success")
	return res
}

func RedisConnect(redisConf *redis.Options) *redis.Client {
	return redis.NewClient(redisConf)
}
