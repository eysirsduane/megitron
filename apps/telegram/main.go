package main

import (
	"flag"
	"fmt"
	"megitron/apps/telegram/bots"
	"megitron/pkg/entity"
	"megitron/pkg/model"
	"megitron/pkg/service"

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

	excfgservice := service.NewExchangeConfigService(db)
	bot := bots.NewBot(cfg.Bot.Token, cfg.Bot.Service, cfg.Tron.MonitorAddress, cfg.Tron.Trx2UsdtRateApi, excfgservice)
	defer bot.Stop()

	starting := "Starting telegram bot..."
	fmt.Println(starting)
	logx.Info(starting)

	bot.Start()
}
