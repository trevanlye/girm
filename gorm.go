package girm

import (
	"fmt"

	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	maxOprCount = 500
)

func bulkOperate[T any](es []*T, f func(any) *gorm.DB) error {
	if len(es) < maxOprCount {
		return f(es).Error
	}

	segs := lo.Chunk(es, maxOprCount)
	for _, seg := range segs {
		if err := f(seg).Error; err != nil {
			return err
		}
	}
	return nil
}

func NewSqlite(dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewMySQL(ipaddr, dbName, userName, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, password, ipaddr, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

