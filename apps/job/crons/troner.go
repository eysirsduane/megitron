package crons

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"megitron/pkg/common"
	"megitron/pkg/entity"
	"megitron/pkg/service"
	"net/http"
	"strconv"
	"sync"
	"time"

	alist "github.com/golanglibs/gocollections/list/arraylist"
	ttypes "github.com/kslamph/tronlib/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type Troner struct {
	grpc             string
	http             string
	rtrx2usdtapi     string
	monitoraddr      string
	trc20usdtaddress string
	ownerprivatekey  string

	trxmintime  uint64
	usdtmintime uint64

	excfgservice *service.ExchangeConfigService
	orderservice *service.OrderService
	tronservice  *service.TronService
}

var trxrunner *Runner[TrxTransaction]
var usdtrunner *Runner[UsdtTransaction]

func NewTroner(grpc, http, rtrx2usdtapi, monitoraddr, trc20usdtaddr string, prikey string, excfgservice *service.ExchangeConfigService, oservice *service.OrderService, tservice *service.TronService) (troner *Troner) {
	troner = &Troner{grpc: grpc, http: http, rtrx2usdtapi: rtrx2usdtapi, monitoraddr: monitoraddr, trc20usdtaddress: trc20usdtaddr, ownerprivatekey: prikey, excfgservice: excfgservice, orderservice: oservice, tronservice: tservice}
	troner.trxmintime = common.TimeNowMilli()
	troner.usdtmintime = common.TimeNowMilli()

	onceing := sync.OnceFunc(func() {
		initTrxRunner()
		initUsdtRunner()
	})

	onceing()

	return
}

func initTrxRunner() {
	trxrunner = &Runner[TrxTransaction]{Currency: string(common.CurrencyTypoTrx)}
	trxrunner.GetTxInfoHandler = func(tx *TrxTransaction) (txid, toaddr string, err error) {
		if len(tx.Ret) > 0 {
			ret := tx.Ret[0]
			if ret.ContractRet != common.StringSuccess {
				return "", "", errors.New("transfer not SUCCESS")
			}
		}

		if len(tx.RawData.Contract) > 0 {
			txid = tx.TxID
			if tx.RawData.Contract[0].Type != common.TronTransactionTypoTransferContract {
				return txid, "", errors.New("not TransferContract")
			}
			addr := tx.RawData.Contract[0].Parameter.Value.ToAddr
			if addr == "" {
				return txid, addr, errors.New("to address is empty")
			}

			address := ttypes.MustNewAddressFromHex(addr)
			toaddr = address.Base58()
			return
		}

		return "", "", errors.New("from or to address is empty")
	}
}

func initUsdtRunner() {
	usdtrunner = &Runner[UsdtTransaction]{Currency: string(common.CurrencyTypoUsdt)}
	usdtrunner.GetTxInfoHandler = func(tx *UsdtTransaction) (txid, toaddr string, err error) {
		if tx.Type != common.TronTransactionTypoTransfer {
			return "", "", errors.New("not transfer event")
		}
		if tx.From == "" || tx.To == "" {
			return "", "", errors.New("from or to address is empty")
		}

		return tx.TransactionId, tx.To, nil
	}
}

func (t *Troner) MonitorTrxTransaction() {
	if trxrunner.IsRunning {
		logx.Debugf("cron tron monitor trx transactions is running, wait for next round automatic.")
		return
	}

	trxrunner.IsRunning = true

	logx.Debugf("cron tron monitor trx transaction current is running.")

	url := fmt.Sprintf("%v/v1/accounts/%s/transactions?limit=200&only_to=true&only_confirmed=true&min_timestamp=%v&max_timestamp=%v&order_by=block_timestamp,asc&event_name=Transfer", t.http, t.monitoraddr, t.trxmintime, common.TimeNowMilli())
	lids, err := t.orderservice.FindLastHourDelegateTxIds()
	if err != nil {
		logx.Errorf("cron tron monitor trx transaction find delegate energy order last day's txids failed, err:%v", err)
		return
	}
	trxrunner.Url = url
	trxrunner.Existent = lids
	trxrunner.MonitorAddress = t.monitoraddr
	trxrunner.Transactions = make(chan *TrxTransaction)

	go trxrunner.Run()

	for {
		tx, ok := <-trxrunner.Transactions
		if !ok {
			logx.Debugf("cron tron monitor trx transaction current is finished.")
			break
		}

		contract := tx.RawData.Contract[0]
		sun := contract.Parameter.Value.Amount
		amount := common.Amount(sun)
		from := ttypes.MustNewAddressFromHex(contract.Parameter.Value.OwnerAddr)
		to := ttypes.MustNewAddressFromHex(contract.Parameter.Value.ToAddr)

		ecfg, err := t.excfgservice.GetExchangeConfig(string(common.ExchangeRateTrx2Energy), int64(amount), int64(amount))
		if err != nil {
			logx.Errorf("cron tron TRX get exchange config failed, TRX[%v], err:%v", amount, err)
			continue
		}

		entity := &entity.DelegateOrder{
			TransactionId:  tx.TxID,
			Typo:           int32(common.DelegateResourceTypoEnergy),
			Curreny:        string(common.CurrencyTypoTrx),
			From:           from.Base58(),
			To:             to.Base58(),
			FromHex:        from.Hex(),
			ToHex:          to.Hex(),
			ReceivedAmount: amount,
			ReceivedSun:    uint64(sun),
			DelegateAmount: int64(ecfg.Value),
			// Expires:        common.TimeInHourMilli(), //uint64(time.Now().Add(time.Second * 120).UnixMilli()),
		}
		err = t.orderservice.CreateDelegateOrder(entity)
		if err != nil {
			logx.Errorf("cron tron TRX create transaction db failed, err:%v", err)
			continue
		}

		logx.Infof("🎉🎉🎉 到账 TRX [%v][%v], from:%v, to:%v, time:%v, txid:%v", amount, sun, from.Base58(), to.Base58(), tx.RawData.Timestamp, tx.TxID)

		err = t.orderDelegate(entity, from, to, entity.DelegateAmount)
		if err != nil {
			logx.Errorf("cron tron order delegate failed, err:%v", err)
			continue
		}

		t.trxmintime = tx.BlockTimestamp
	}

	trxrunner.IsRunning = false
}

func (t *Troner) MonitorUsdtTransaction() {
	if usdtrunner.IsRunning {
		logx.Debugf("cron tron monitor usdt transaction is running, wait for next round automatic!")
		return
	}

	usdtrunner.IsRunning = true

	logx.Debugf("cron tron monitor usdt transaction current is running.")

	url := fmt.Sprintf("%v/v1/accounts/%s/transactions/trc20?limit=200&contract_address=%v&only_confirmed=true&min_timestamp=%v&max_timestamp=%v", t.http, t.monitoraddr, t.trc20usdtaddress, t.usdtmintime, common.TimeNowMilli())
	lids, err := t.orderservice.FindLastHourExchangeTxIds()
	if err != nil {
		logx.Errorf("cron tron monitor trx transaction find delegate energy order last day's txids failed, err:%v", err)
		return
	}
	usdtrunner.Url = url
	usdtrunner.Existent = lids
	usdtrunner.MonitorAddress = t.monitoraddr
	usdtrunner.Transactions = make(chan *UsdtTransaction)

	go usdtrunner.Run()

	for {
		tx, ok := <-usdtrunner.Transactions
		if !ok {
			logx.Debugf("cron tron monitor usdt transaction current is finished.")
			break
		}

		sun, err := strconv.ParseInt(tx.Value, 10, 64)
		if err != nil {
			logx.Errorf("cron tron usdt transaction parse amount failed, err:%v", err)
			continue
		}

		amount := common.Amount(sun)
		ecfg, err := t.excfgservice.GetExchangeConfig(string(common.ExchangeRateUsdt2Trx), 1, 1)
		if err != nil {
			logx.Errorf("cron tron USDT get exchange config failed, typo:%v, err:%v", common.ExchangeRateUsdt2Trx, err)
			continue
		}

		rate, err := common.GetTrx2UsdtRateFromHtx(t.rtrx2usdtapi)
		if err != nil {
			logx.Errorf("cron tron USDT get trx2usdt rate failed, err:%v", err)
			continue
		}

		eamount, erate := common.GetUsdt2TrxAmount(rate, amount, ecfg.Value)
		from := ttypes.MustNewAddressFromBase58(tx.From)
		to := ttypes.MustNewAddressFromBase58(tx.To)

		entity := &entity.ExchangeOrder{
			TransactionId:    tx.TransactionId,
			Typo:             string(common.ExchangeTypoUsdt2Trx),
			Curreny:          string(common.CurrencyTypoUsdt),
			From:             from.Base58(),
			To:               to.Base58(),
			FromHex:          from.Hex(),
			ToHex:            to.Hex(),
			ThenRate:         rate,
			ExchangeRate:     erate,
			ExchangeDiscount: ecfg.Value,
			ReceivedAmount:   amount,
			ReceivedSun:      sun,
			ExchangeAmount:   eamount,
			ExchangeSun:      common.Sun(eamount),
		}
		err = t.orderservice.CreateExchangeOrder(entity)
		if err != nil {
			logx.Errorf("cron tron exchange order create failed, err:%v", err)
			continue
		}

		logx.Infof("🏆🏆🏆 到账 USDT [%v][%v], from:%v, to:%v, time:%v, txid:%v", amount, sun, tx.From, tx.To, tx.BlockTimestamp, tx.TransactionId)

		err = t.orderExchange(entity, from, to, eamount)
		if err != nil {
			logx.Errorf("cron tron order exchange failed, err:%v", err)
			continue
		}

		t.usdtmintime = tx.BlockTimestamp
	}

	usdtrunner.IsRunning = false
}

func (t *Troner) orderExchange(order *entity.ExchangeOrder, from, to *ttypes.Address, amount float64) (err error) {
	balance, _, err := t.tronservice.TronAccountBalance(t.ownerprivatekey)
	if err != nil {
		logx.Errorf("cron tron get account balance failed, err:%v", err)
		return
	}

	if balance < amount {
		logx.Infof("🟠🟠🟠 余额不足, trx balance is insufficient! balance:%v, need:%v, from:%v", balance, amount, to)
		return
	}

	etxid, msg, success, err := t.tronservice.TransferTrxAndBroadcast(t.ownerprivatekey, from.Base58(), common.Sun(amount))
	if err != nil || !success {
		logx.Errorf("cron tron usdt2trx create transfer failed, TRX[%v], from:%v, to:%v, success:%v, msg:%v, err:%v", amount, to.Base58(), from.Base58(), success, msg, err)

		order.Status = int16(common.ExchangeOrderStatusPending)
		order.Time = common.TimeNowMilli()
		order.Description = fmt.Sprintf("兑换 TRX[%v] 失败, 创建 Tron Transaction 失败, success:%v, err:%v", amount, success, err)
		err = t.orderservice.UpdateExchangeOrder(order)
		if err != nil {
			logx.Errorf("cron tron exchange trx update order failed, status:%v, txid:%v, etxid:%v, err:%v", order.Status, order.TransactionId, etxid, err)
		}

		return
	}

	logx.Infof("🥇🥇🥇 兑出 TRX[%v], 收到 USDT[%v], from:%v, to:%v, etxid:%v", amount, order.ReceivedAmount, to.Base58(), from.Base58(), etxid)

	bill := &entity.ExchangeBill{}
	bill.UserId = order.UserId
	bill.OrderId = order.Id
	bill.TransactionId = etxid
	bill.Status = int16(common.ExchangeBillStatusSuccess)
	bill.Curreny = string(common.CurrencyTypoTrx)
	bill.From = to.Base58()
	bill.To = from.Base58()
	bill.FromHex = to.Hex()
	bill.ToHex = from.Hex()
	bill.ExchangedAmount = amount
	bill.ExchangedSun = common.Sun(amount)
	err = t.orderservice.UpdateExchangeBill(bill)
	if err != nil {
		logx.Errorf("cron tron exchange bill update failed, oid:%v, otxid:%v, btxid:%v, from:%v, to:%v, err:%v", order.Id, order.TransactionId, bill.TransactionId, bill.From, bill.To, err)
		return
	}

	ostatus := order.Status
	order.Time = common.TimeNowMilli()
	order.Status = int16(common.ExchangeOrderStatusFinished)
	err = t.orderservice.UpdateExchangeOrder(order)
	if err != nil {
		logx.Errorf("cron tron exchange order update failed, [%v]->[%v], oid:%v, txid:%v, err:%v", ostatus, order.Id, order.Status, order.TransactionId, err)
		return
	}

	return
}

func (t *Troner) orderDelegate(order *entity.DelegateOrder, from, to *ttypes.Address, amount int64) (err error) {
	dtxid, msg, btime, success, err := t.tronservice.ResourceDelegate(t.ownerprivatekey, from.Base58(), int32(common.DelegateResourceTypoEnergy), amount)
	if err != nil || !success {
		logx.Errorf("cron tron delegate energy failed, oid:%v, otxid:%v, success:%v, msg:%v, err:%v", order.Id, order.TransactionId, success, msg, err)

		order.Time = uint64(btime)
		if order.FailedTimes >= 3 {
			order.Status = uint16(common.DelegateOrderStatusError)
		} else {
			order.FailedTimes++
			order.Status = uint16(common.DelegateOrderStatusPending)
			order.Description = fmt.Sprintf("租出 能量[%v] 发生错误, TRX[%v] oid:%v, otxid:%v, success:%v, err:%v", amount, order.ReceivedAmount, order.Id, order.TransactionId, success, err)
		}
		err = t.orderservice.UpdateDelegateOrder(order)
		if err != nil {
			logx.Errorf("cron tron delegate energy update order failed, status:%v, oid:%v, otxid:%v, dtxid:%v, err:%v", order.Status, order.Id, order.TransactionId, dtxid, err)
		}

		return
	}

	logx.Infof("🥇🥇🥇 租出 TRON 能量[%v], 收到 TRX[%v], from:%v, to:%v, oid:%v, otxid:%v, time:%v, dtxid:%v", amount, order.ReceivedAmount, to.Base58(), from.Base58(), order.Id, order.TransactionId, btime, dtxid)

	bill := &entity.DelegateBill{}
	bill.UserId = order.UserId
	bill.OrderId = order.Id
	bill.TransactionId = dtxid
	bill.Status = int16(common.DelegateBillStatusSuccess)
	bill.Curreny = string(common.CurrencyTypoEnergy)
	bill.From = to.Base58()
	bill.To = from.Base58()
	bill.FromHex = to.Hex()
	bill.ToHex = from.Hex()
	bill.DelegatedAmount = amount
	err = t.orderservice.UpdateDelegateBill(bill)
	if err != nil {
		logx.Errorf("cron tron delegate bill update failed, oid:%v, otxid:%v, btxid:%v, from:%v, to:%v, err:%v", order.Id, order.TransactionId, bill.TransactionId, bill.From, bill.To, err)
		return
	}

	ostatus := order.Status
	order.Time = common.TimeNowMilli()
	order.WithdrawTime = common.TimeInHourMilli()
	order.Status = uint16(common.DelegateOrderStatusDelegated)
	err = t.orderservice.UpdateDelegateOrder(order)
	if err != nil {
		udtxid, _, _, success, err1 := t.tronservice.ResourceUnDelegate(t.ownerprivatekey, from.Base58(), int32(common.DelegateResourceTypoEnergy), amount)
		if err1 != nil || !success {
			logx.Errorf("cron tron delegate order update order to [%v] failed, and take back failed, oid:%v, otxid:%v, success:%v, dtxid:%v, err:%v", order.Status, order.Id, order.TransactionId, success, udtxid, err1)
			return
		}

		logx.Errorf("cron tron delegate order update failed, [%v]->[%v], from:%v, to:%v, oid:%v, otxid:%v, dtxid:%v, err:%v", ostatus, order.Status, from, to, order.Id, order.TransactionId, dtxid, err)
		return
	}

	return
}

func (t *Troner) WithdrawDelegatedOrders() {
	delegateds, err := t.orderservice.DelegatedOrders()
	if err != nil {
		logx.Errorf("cron tron delegated orders find failed, err:%v", err)
		return
	}

	for _, order := range delegateds {
		utxid, _, time, success, err := t.tronservice.ResourceUnDelegate(t.ownerprivatekey, order.From, order.Typo, order.DelegateAmount)
		if err != nil || !success {
			logx.Errorf("cron tron take back resource failed, oid:%v, txid:%v, utxid:%v, success:%v, time:%v, err:%v", order.Id, order.TransactionId, utxid, success, time, err)
			continue
		}

		from := ttypes.MustNewAddressFromBase58(order.From)
		to := ttypes.MustNewAddressFromBase58(order.To)

		drawal := &entity.DelegateWithdrawal{}
		drawal.UserId = order.UserId
		drawal.OrderId = order.Id
		drawal.TransactionId = utxid
		drawal.Typo = order.Typo
		drawal.From = from.Base58()
		drawal.To = to.Base58()
		drawal.FromHex = from.Hex()
		drawal.ToHex = to.Hex()
		drawal.UnDelegatedAmount = order.DelegateAmount
		err = t.orderservice.UpdateDelegateWithdrawal(drawal)
		if err != nil {
			logx.Errorf("cron tron get delegate takeback failed, oid:%v, otxid:%v, err:%v", order.Id, order.TransactionId, err)
			continue
		}

		order.Time = common.TimeNowMilli()
		order.Status = uint16(common.DelegateOrderStatusFinished)
		err = t.orderservice.UpdateDelegateOrder(order)
		if err != nil {
			logx.Errorf("cron tron expires delegated order batch update failed, err:%v", err)
			continue
		}

		logx.Infof("🟩🟩🟩 回收资源 [%v][%v], oid:%v, from:%v, to:%v, otxid:%v, udtxid:%v, success:%v, time:%v", drawal.UnDelegatedAmount, order.Typo, order.Id, order.From, order.To, order.TransactionId, drawal.TransactionId, success, time)
	}
}

func (t *Troner) ReDelegatePendingOrders() {
	pendings, _, err := t.orderservice.FindDelegateOrders(int16(common.DelegateOrderStatusPending))
	if err != nil {
		logx.Errorf("cron tron expires redelegated order failed, err:%v", err)
		return
	}

	for _, order := range pendings {
		from := ttypes.MustNewAddressFromBase58(order.From)
		to := ttypes.MustNewAddressFromBase58(order.To)
		err = t.orderDelegate(order, from, to, order.DelegateAmount)
		if err != nil {
			logx.Errorf("cron tron expires redelegated order batch update failed, err:%v", err)
		}
	}
}

func (t *Troner) ReExchangePendingOrders() {
	pendings, _, err := t.orderservice.FindExchangeOrders(int16(common.ExchangeOrderStatusPending))
	if err != nil {
		logx.Errorf("cron tron reexchange delegated order failed, err:%v", err)
		return
	}

	for _, order := range pendings {
		from := ttypes.MustNewAddressFromBase58(order.From)
		to := ttypes.MustNewAddressFromBase58(order.To)
		err = t.orderExchange(order, from, to, order.ExchangeAmount)
		if err != nil {
			logx.Errorf("cron tron reexchange delegated order batch update failed, err:%v", err)
		}
	}
}

type Runner[T TrxTransaction | UsdtTransaction] struct {
	IsRunning        bool
	Url              string
	Currency         string
	MonitorAddress   string
	MinTimestamp     int64
	MaxTimestamp     int64
	Existent         []string
	Transactions     chan *T
	Wrapper          struct{ Data []T }
	GetTxInfoHandler func(tx *T) (txid, toaddr string, err error)
}

func (r *Runner[T]) Run() {
	r.MinTimestamp = time.Now().AddDate(0, 0, -1).UnixMilli()
	r.MaxTimestamp = time.Now().UnixMilli()

	resp, err := http.Get(r.Url)
	if err != nil {
		logx.Infof("cron tron %v transaction listen request failed, err:%v", r.Currency, err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	wrapper := &struct{ Data []*T }{}
	if err := json.Unmarshal(body, wrapper); err != nil {
		logx.Infof("cron tron %v transaction listen unmarshal json failed, err:%v", r.Currency, err)
		return
	}

	if len(wrapper.Data) > 0 {
		llist := alist.New(r.Existent...)
		for _, tx := range wrapper.Data {
			if tx != nil {
				txid, addr, err := r.GetTxInfoHandler(tx)
				if err != nil {
					// logx.Infof("cron tron %v transaction is invalid, txid:%v, err:%v", r.Currency, txid, err)
					continue
				}

				if addr == r.MonitorAddress {
					if llist.Contains(txid) {
						continue
					}

					r.Transactions <- tx
				}
			}
		}
	}

	close(r.Transactions)
}

type Subscription struct {
	Topic   string `json:"topic"`
	Address string `json:"address"`
}

type TrxTransaction struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
		Fee         int64  `json:"fee"`
	}
	TxID           string `json:"txID"`
	BlockTimestamp uint64 `json:"block_timestamp"`
	RawData        struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Amount    int64  `json:"amount"`
					OwnerAddr string `json:"owner_address"`
					ToAddr    string `json:"to_address"`
				} `json:"value"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		Timestamp int64 `json:"timestamp"`
	} `json:"raw_data"`
}

type UsdtTransaction struct {
	TransactionId string `json:"transaction_id"`
	TokenInfo     struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int32  `json:"decimals"`
		Name     string `json:"name"`
	} `json:"token_info"`
	BlockTimestamp uint64 `json:"block_timestamp"`
	From           string `json:"from"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Value          string `json:"value"`
}
