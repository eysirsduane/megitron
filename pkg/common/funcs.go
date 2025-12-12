package common

import (
	"encoding/json"
	"io"
	"math"
	"megitron/pkg/biz"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

func Offset(page, limit int) (offset int) {
	return (page - 1) * limit
}

func GetTrx2UsdtRateFromHtx(url string) (rate float64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		logx.Infof("funcs get trx2usdt rate http request failed, err:%v", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	wrapper := &struct {
		Status string `json:"status"`
		Tick   struct {
			Open float64 `json:"open"`
		} `json:"tick"`
	}{}

	if err = json.Unmarshal(body, wrapper); err != nil {
		logx.Infof("funcs get trx2usdt rate unmarshal json failed, err:%v", err)
		return
	}

	rate = wrapper.Tick.Open
	if rate >= 1 || rate <= 0 {
		logx.Errorf("funcs get trx2usdt rate value is incorrect, rate:%v", rate)
		err = biz.ExchangeTrx2UsdtRateIncorrect
		return
	}

	return
}

func GetTrx2UsdtExchangeRate(rate float64) (erate float64) {
	return math.Ceil(rate*100) / 100
}

func GetUsdt2TrxAmount(rate float64, amount float64, discount float64) (eamount, erate float64) {
	erate = GetTrx2UsdtExchangeRate(rate)
	eamount = math.Floor(((amount/erate)*discount)*100) / 100

	return
}
