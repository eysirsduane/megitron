package service

import (
	"context"
	"fmt"
	"math/big"
	"megitron/pkg/common"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jinzhu/copier"
	"github.com/kslamph/tronlib/pb/api"
	"github.com/kslamph/tronlib/pkg/client"
	"github.com/kslamph/tronlib/pkg/resources"
	"github.com/kslamph/tronlib/pkg/signer"
	ttypes "github.com/kslamph/tronlib/pkg/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type TronService struct {
	grpc string
	db   *gorm.DB
}

func NewTronService(grpc string, db *gorm.DB) *TronService {
	return &TronService{grpc: grpc, db: db}
}

func (t *TronService) Something(oprikey, receiver string) (c *client.Client, siger *signer.PrivateKeySigner, from *ttypes.Address, to *ttypes.Address, err error) {
	c, err = client.NewClient(t.grpc)
	if err != nil {
		logx.Errorf("tron get client new clinet failed, err:%v", err)
		return
	}

	if oprikey != "" {
		key, ierr := signer.NewPrivateKeySigner(oprikey)
		if ierr != nil {
			logx.Errorf("tron get client new private key signer failed, err:%v", ierr)
			err = ierr
			return
		}

		siger = key
		from = key.Address()
	}

	if receiver != "" {
		to = ttypes.MustNewAddressFromBase58(receiver)
	}

	return
}

func (t *TronService) TransferTrxAndBroadcast(oprikey, receiver string, sun int64) (id, msg string, success bool, err error) {
	c, siger, from, to, err := t.Something(oprikey, receiver)
	if err != nil {
		return
	}
	defer c.Close()

	tx, err := c.Account().TransferTRX(context.Background(), from, to, sun)
	if err != nil {
		logx.Errorf("tron trx transfer create transfer failed, err:%v", err)
		return
	}
	if tx == nil {
		logx.Errorf("tron trx transfer create transfer failed, tx is nil")
		return
	}

	result, err := c.SignAndBroadcast(context.Background(), tx, client.DefaultBroadcastOptions(), siger)
	if err != nil {
		logx.Errorf("tron trx transfer sign and broadcast failed, txid:%v, err:%v", string(tx.GetTxid()), err)
		return
	}

	logx.Infof("tron trx transfer broadcast, result:%+v", result)

	if !result.Success || result.Code != api.Return_SUCCESS {
		return result.TxID, result.Message, false, nil
	}

	return result.TxID, result.Message, true, nil
}

func EncodeTransfer(to *ttypes.Address, amount *big.Int) ([]byte, error) {
	parsedABI, err := abi.JSON(strings.NewReader(`[{
        "constant": false,
        "inputs": [
            {"name":"_to", "type":"address"},
            {"name":"_value", "type":"uint256"}
        ],
        "name": "transfer",
        "outputs": [{"name":"","type":"bool"}],
        "type":"function"
    }]`))
	if err != nil {
		return nil, err
	}

	return parsedABI.Pack("transfer", to, amount)
}

func (t *TronService) ResourceDelegate(oprikey, receiver string, typo int32, amount int64) (txid, msg string, time int64, success bool, err error) {
	c, siger, from, to, err := t.Something(oprikey, receiver)
	if err != nil {
		return
	}
	defer c.Close()

	ctx := context.Background()

	tx, err := c.Resources().DelegateResource(ctx, from, to, amount, resources.ResourceType(typo), false)
	if err != nil || tx == nil {
		logx.Errorf("tron resource delegate create transfer failed, receiver:%v, tx:%+v, err:%v", receiver, tx, err)
		return
	}

	result, err := c.SignAndBroadcast(ctx, tx, client.DefaultBroadcastOptions(), siger)
	if err != nil {
		logx.Errorf("tron resource delegate sign and broadcast failed, receiver:%v, txid:%v, err:%v", receiver, string(tx.GetTxid()), err)
		return
	}

	logx.Infof("tron resource delegate broadcast, result:%+v", result)

	txid = result.TxID
	msg = result.Message
	time = tx.Transaction.RawData.Timestamp

	if !result.Success || result.Code != api.Return_SUCCESS {
		return
	}

	// drsc, err := c.Resources().GetDelegatedResourceV2(context.Background(), from, to)
	// if err != nil {
	// 	logx.Errorf("tron account delegated resource get failed, err:%v", err)
	// 	return
	// }
	success = true
	return
}

func temp() {
	// paramTypes := []string{"address", "uint256"}
	// params := []interface{}{
	// 	to.Base58(), // TRON address
	// 	big.NewInt(amount),
	// }
	// p := utils.NewABIProcessor(nil)
	// data, err := p.EncodeMethod("transfer", paramTypes, params)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// trc20address, err := ttypes.NewAddressFromBase58(trc20uaddress)
	// if err != nil {
	// 	logx.Errorf("tron resource delegate parse trc20usdt address failed, err:%v", err)
	// 	return
	// }

	// estimate, err := c.SmartContract().EstimateEnergy(ctx, from, trc20address, data, 0)
	// if err != nil {
	// 	logx.Errorf("tron resource delegate estimate energy failed, err:%v", err)
	// 	return
	// }
	// if estimate.Result.Result {
	// 	amount = estimate.EnergyRequired
	// 	// amount = 1.5*1_000_000
	// }
}

func (t *TronService) ResourceUnDelegate(oprikey, receiver string, typo int32, amount int64) (txid, msg string, time int64, success bool, err error) {
	if receiver == "" || amount == 0 {
		return
	}

	c, siger, from, to, err := t.Something(oprikey, receiver)
	if err != nil {
		return
	}
	defer c.Close()

	tx, err := c.Resources().UnDelegateResource(context.Background(), from, to, amount, resources.ResourceType(typo))
	if err != nil {
		logx.Errorf("tron resource undelegate create transfer failed, err:%v", err)
		return
	}
	if tx == nil {
		logx.Errorf("traon resource undelegate create transfer failed, tx is nil")
		return
	}

	result, err := c.SignAndBroadcast(context.Background(), tx, client.DefaultBroadcastOptions(), siger)
	if err != nil {
		logx.Errorf("tron resource undelegate sign and broadcast failed, txid:%v, err:%v", string(tx.GetTxid()), err)
		return
	}

	logx.Infof("tron resource undelegate broadcast, result:%+v", result)

	txid = result.TxID
	msg = result.Message
	time = tx.Transaction.RawData.Timestamp

	if !result.Success || result.Code != api.Return_SUCCESS {
		return
	}

	// drsc, err := c.Resources().GetDelegatedResourceV2(context.Background(), from, to)
	// if err != nil {
	// 	logx.Errorf("tron account delegated resource get failed, err:%v", err)
	// 	return
	// }
	success = true
	return
}

func (t *TronService) TronAccountResource(oprikey string) (resource *TronAccountResource, err error) {
	c, _, from, _, err := t.Something(oprikey, "")
	if err != nil {
		return
	}
	defer c.Close()

	rsc, err := c.Account().GetAccountResource(context.Background(), from)
	if err != nil {
		logx.Errorf("tron account resource get failed, err:%v", err)
		return
	}

	resource = &TronAccountResource{}
	copier.Copy(resource, rsc)

	logx.Infof("tron account resource, owner:%v, resource:%+v", from.Base58(), resource)

	return
}

func (t *TronService) TronAccountBalance(oprikey string) (amount float64, sun int64, err error) {
	c, _, from, _, err := t.Something(oprikey, "")
	if err != nil {
		return
	}
	defer c.Close()

	sun, err = c.Account().GetBalance(context.Background(), from)
	if err != nil {
		logx.Errorf("tron account balance get failed, err:%v", err)
		return
	}

	amount = common.Amount(sun)

	fmt.Printf("tron account balance, owner:%v, amount:%v, sun:%v", from.Base58(), amount, sun)

	return
}
