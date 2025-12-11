package service

import (
	"megitron/pkg/biz"
	"megitron/pkg/common"
	"megitron/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ExchangeConfigService struct {
	db *gorm.DB
}

func NewExchangeConfigService(db *gorm.DB) *ExchangeConfigService {
	return &ExchangeConfigService{db: db}
}

func (t *ExchangeConfigService) GetExchangeConfig(typo string, rfrom, rto int64) (cfg *entity.ExchangeConfig, err error) {
	cfg = &entity.ExchangeConfig{}

	err = t.db.Model(&entity.ExchangeConfig{}).Where("typo = ? and  range_from >= ? and range_to <= ?", typo, rfrom, rto).First(cfg).Error
	if err != nil {
		logx.Errorf("database exchange config get failed, typo:%v, rfrom:%v, rto:%v, err:%v", typo, rfrom, rto, err)
		err = biz.ExchangeConfigGetFailed
		return
	}

	if cfg.Value == 0 {
		logx.Errorf("database exchange config value is 0, typo:%v, rfrom:%v, rto:%v, err:%v", typo, rfrom, rto, err)
		err = biz.ExchangeConfigValueInvalid
		return
	}
	if typo == string(common.ExchangeTypoUsdt2Trx) {
		if cfg.Value >= 1 {
			logx.Errorf("database exchange config usdt2trx value invalid, value:%v", cfg.Value)
			err = biz.ExchangeConfigValueInvalid
			return
		}
	}

	return
}
