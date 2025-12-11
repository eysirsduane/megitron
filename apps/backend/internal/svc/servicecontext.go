// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"megitron/pkg/model"
	"megitron/pkg/service"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      model.Config
	DB          *gorm.DB
	TronService *service.TronService
}

func NewServiceContext(c model.Config, db *gorm.DB, tron *service.TronService) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		DB:          db,
		TronService: tron,
	}
}
