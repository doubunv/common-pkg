package dbM

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type SelectDbConfig struct {
	DbKey      string
	DataSource string
}

type SelectDb struct {
	dbMap map[string]*gorm.DB
}

func NewSelectDb(config []SelectDbConfig) *SelectDb {
	res := &SelectDb{
		dbMap: make(map[string]*gorm.DB),
	}
	for _, v := range config {
		db := MySqlConnect(v.DataSource)
		if db.Error != nil {
			panic(db.Error)
		}
		res.dbMap[v.DbKey] = db
	}

	return res
}

func (s *SelectDb) GetDb(ctx context.Context, dbKey string) *gorm.DB {
	if db, ok := s.dbMap[dbKey]; ok {
		if db.Error != nil {
			logc.Error(context.Background(), fmt.Sprintf("Business[%s] DB is invalid, %s", dbKey, db.Error))
			panic(fmt.Sprintf("Business[%s] DB is invalid, %s", dbKey, db.Error))
		}
		return db
	}
	logc.Error(context.Background(), fmt.Sprintf("Business[%s] no DB", dbKey))
	panic(fmt.Sprintf("Business[%s] no DB", dbKey))
}

func (s *SelectDb) GetDbAll() map[string]*gorm.DB {
	return s.dbMap
}
