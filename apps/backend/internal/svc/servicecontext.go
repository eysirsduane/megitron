// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"megitron/pkg/model"
	"megitron/pkg/service"
)

type ServiceContext struct {
	Config       model.Config
	TronService  *service.TronService
	OrderService *service.OrderService
	ExcfgService *service.ExchangeConfigService
}

func NewServiceContext(c model.Config, tron *service.TronService, order *service.OrderService, excfg *service.ExchangeConfigService) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		TronService:  tron,
		OrderService: order,
		ExcfgService: excfg,
	}
}
