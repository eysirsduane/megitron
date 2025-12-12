// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"

	"megitron/pkg/biz"
	"megitron/pkg/entity"
	"megitron/pkg/model"
	"megitron/pkg/service"

	"megitron/apps/backend/internal/handler"
	"megitron/apps/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "../../etc/megitron.dev.yaml", "the config file")

func main() {
	flag.Parse()

	var cfg model.Config
	conf.MustLoad(*configFile, &cfg)
	logx.MustSetup(cfg.Log)
	defer logx.Close()

	db, err := entity.NewGormDB(&cfg)
	if err != nil {
		logx.Errorf("init gorm db failed, err:%v", err)
		panic(err)
	}

	server := rest.MustNewServer(cfg.RestConf)
	defer server.Stop()

	orderservice := service.NewOrderService(db)
	excfgservice := service.NewExchangeConfigService(db)
	tronservice := service.NewTronService(cfg.Tron.GrpcNetwork, db)
	ctx := svc.NewServiceContext(cfg, tronservice, orderservice, excfgservice)
	handler.RegisterHandlers(server, ctx)

	httpx.SetOkHandler(biz.OkHandler)
	httpx.SetErrorHandlerCtx(biz.ErrHandler(cfg.Name))

	starting := fmt.Sprintf("Starting http server %s at %s:%d ...", cfg.Name, cfg.Host, cfg.Port)
	fmt.Println(starting)
	logx.Info(starting)

	server.Start()
}
