package main

import (
	"flag"
	"fmt"
	"megitron/apps/job/crons"
	"megitron/pkg/entity"
	"megitron/pkg/model"
	"megitron/pkg/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
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

	orderservice := service.NewOrderService(db)
	excfgservice := service.NewExchangeConfigService(db)
	tronservice := service.NewTronService(cfg.Tron.GrpcNetwork, db)

	troner := crons.NewTroner(cfg.Tron.GrpcNetwork, cfg.Tron.HttpNetwork, cfg.Tron.Trx2UsdtRateApi, cfg.Tron.MonitorAddress, cfg.Tron.TRC20USDTAddress, cfg.Tron.OwnerPrivateKey, excfgservice, orderservice, tronservice)
	mgr := crons.NewCronManager(troner)
	defer mgr.Stop()

	starting := "Starting cron job..."
	fmt.Println(starting)
	logx.Info(starting)

	mgr.Start()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigchan
	logx.Infof("收到信号:%s, 准备退出...", sig)
}
