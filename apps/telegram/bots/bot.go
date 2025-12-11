package bots

import (
	"fmt"
	"math"
	"megitron/pkg/common"
	"megitron/pkg/service"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/telebot.v4"
)

type Bot struct {
	bot             *telebot.Bot
	token           string
	service         string
	receiveAddress  string
	rateApiTrx2Usdt string
	excfgservice    *service.ExchangeConfigService
}

func NewBot(token, service, receiveaddr, rateapi string, excfgservice *service.ExchangeConfigService) (bot *Bot) {
	return &Bot{excfgservice: excfgservice, token: token, service: service, receiveAddress: receiveaddr, rateApiTrx2Usdt: rateapi}
}

func (b *Bot) Init() (err error) {
	return
}

func (b *Bot) Stop() {
	b.bot.Stop()
}

var (
	menu       = &telebot.ReplyMarkup{ResizeKeyboard: true}
	btntrx     = menu.Text("💹闪兑TRX")
	btnenergy  = menu.Text("🔋闪租能量")
	btnservice = menu.Text("🧑‍💼联系客服")
)

func (b *Bot) Start() {
	pref := telebot.Settings{
		Token:  b.token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		panic(err)
	}

	b.bot = bot

	err = b.bot.SetCommands([]telebot.Command{
		{Text: "start", Description: "开始使用"},
		{Text: "usdt2trx", Description: "💹闪兑TRX"},
		{Text: "trx2energy", Description: "🔋闪租能量"},
		{Text: "service", Description: "🧑‍💼联系客服"},
	})
	if err != nil {
		panic(err)
	}

	menu.Reply(
		menu.Row(btntrx, btnenergy, btnservice),
	)
	b.bot.Handle("/start", func(c telebot.Context) error {
		return c.Send("Hello!", menu)
	})

	b.bot.Handle("/usdt2trx", b.handlerForUsdt2Trx)
	b.bot.Handle(&btntrx, b.handlerForUsdt2Trx)

	b.bot.Handle("/trx2energy", b.handlerForTrx2Energy)
	b.bot.Handle(&btnenergy, b.handlerForTrx2Energy)

	b.bot.Handle("/service", b.handlerForService)
	b.bot.Handle(&btnservice, b.handlerForService)

	b.bot.Start()
}

func (b *Bot) handlerForUsdt2Trx(c telebot.Context) (err error) {
	cfg, err := b.excfgservice.GetExchangeConfig(string(common.ExchangeTypoUsdt2Trx), 1.0, 1.0)
	if err != nil {
		logx.Errorf("telegram bot get exchange config failed, err:%v", err)
		return
	}
	rate, err := common.GetTrx2UsdtRateFromHtx(b.rateApiTrx2Usdt)
	if err != nil {
		logx.Errorf("telegram bot init get trx2usdt rate failed, err:%v", err)
	}

	one, _ := common.GetUsdt2TrxAmount(rate, 1, cfg.Value)

	line := "💹24小时自动兑换💹 地址:\n"
	line += "【点击自动复制】\n"
	line += "➖➖➖➖➖➖➖➖➖➖➖➖\n"
	line += fmt.Sprintf("```%v```\n", b.receiveAddress)
	line += "➖➖➖➖➖➖➖➖➖➖➖➖\n"
	line += "当前汇率：\n"
	line += fmt.Sprintf(`1 USDT \= %v TRX%v`, convertSpecialChars(one), "\n")
	line += fmt.Sprintf(`10 USDT \= %v TRX%v`, convertSpecialChars(one*10), "\n")
	line += fmt.Sprintf(`100 USDT \= %v TRX%v`, convertSpecialChars(one*100), "\n")
	line += fmt.Sprintf(`1000 USDT \= %v TRX%v`, convertSpecialChars(one*1000), "\n\n")
	line += "💹进U即兑, 全自动返TRX, 1U起兑\n"
	line += "❌请勿使用交易所或中心化钱包转账\n"
	line += fmt.Sprintf("💹如有老板需要用交易所转账, 请联系客服: %v\n", b.service)

	menu.Reply(
		menu.Row(btntrx, btnenergy, btnservice),
	)

	c.Respond()
	return c.Send(line, telebot.ModeMarkdownV2, menu)
}

func (b *Bot) handlerForTrx2Energy(c telebot.Context) error {
	line := "🔋1小时能量闪租🔋 地址:\n"
	line += "【点击自动复制】\n"
	line += "➖➖➖➖➖➖➖➖➖➖➖➖\n"
	line += fmt.Sprintf("```%v```\n", b.receiveAddress)
	line += "➖➖➖➖➖➖➖➖➖➖➖➖\n"
	line += "租用能量, 转账无需TRX消耗, 0手续费！\n"
	line += "1小时能量闪租, 转U不扣手续费\n"
	line += fmt.Sprintf(`2 TRX \= 1笔 %v 能量%v`, convertSpecialChars(64285), "\n")
	line += fmt.Sprintf(`4 TRX \= 1笔 %v 能量%v`, convertSpecialChars(130285), "\n\n")
	line += "🔋进TRX即到账能量, 全自动, 2TRX起租\n"
	line += "🔋向无U地址转账, 需要双倍能量\n"
	line += "❌请在1小时内使用, 否则过期收回\n"
	line += fmt.Sprintf("🔋如有老板需要用大量能量, 请联系客服: %v\n", b.service)

	menu.Reply(
		menu.Row(btntrx, btnenergy, btnservice),
	)

	c.Respond()
	return c.Send(line, telebot.ModeMarkdownV2, menu)
}

func (b *Bot) handlerForService(c telebot.Context) error {
	line := "🛡️🛡️ 大平台 🔒 更可靠 🔒 秒到账 🛡️🛡️\n"
	line += "\n"
	line += fmt.Sprintf("客服: 💬%v💬\n", b.service)

	menu.Reply(
		menu.Row(btntrx, btnenergy, btnservice),
	)

	c.Respond()
	return c.Send(line, telebot.ModeMarkdownV2, menu)
}

func convertSpecialChars(val float64) (r string) {
	s := fmt.Sprintf("%v", math.Floor(val*100)/100)
	return strings.ReplaceAll(s, ".", `\.`)
}
