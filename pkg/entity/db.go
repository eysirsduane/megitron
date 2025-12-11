package entity

import (
	"fmt"
	"megitron/pkg/model"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg *model.Config) (db *gorm.DB, err error) {
	// envs := getEnvs([]string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_PORT", "DB_NAME", "DB_TIMEZONE"})
	dsn := fmt.Sprintf("host=%v user=%v password=%v port=%v dbname=%v sslmode=disable TimeZone=%v", cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Port, cfg.DB.Name, cfg.DB.Timezone)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		logx.Errorf("init gorm engine failed, err:%v", err)
		panic(err)
	}
	err = db.AutoMigrate(&ExchangeConfig{}, &ExchangeOrder{}, &ExchangeBill{}, &DelegateOrder{}, &DelegateBill{}, &DelegateWithdrawal{})
	if err != nil {
		logx.Errorf("init gorm create tables failed, err:%v", err)
		panic(err)
	}

	return
}

func getEnvs(keys []string) (vals []string) {
	vals = []string{}
	for _, key := range keys {
		vals = append(vals, os.Getenv(key))
	}

	return
}
