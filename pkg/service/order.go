package service

import (
	"megitron/pkg/biz"
	"megitron/pkg/common"
	"megitron/pkg/entity"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (t *OrderService) CreateDelegateOrder(tx *entity.DelegateOrder) (err error) {
	err = t.db.Create(tx).Error
	if err != nil {
		logx.Errorf("database delegate order insert failed, err:%v", err)
		err = biz.DelegateOrderCreateFailed
		return
	}

	return
}

func (t *OrderService) CreateExchangeOrder(tx *entity.ExchangeOrder) (err error) {
	err = t.db.Create(tx).Error
	if err != nil {
		logx.Errorf("database exchange order insert failed, err:%v", err)
		err = biz.ExchangeOrderCreateFailed
		return
	}

	return
}

func (t *OrderService) UpdateExchangeOrder(tx *entity.ExchangeOrder) (err error) {
	err = t.db.Save(tx).Error
	if err != nil {
		logx.Errorf("database exchange order update failed, err:%v", err)
		err = biz.ExchangeOrderUpdateFailed
		return
	}

	return
}

func (t *OrderService) UpdateExchangeBill(bill *entity.ExchangeBill) (err error) {
	err = t.db.Save(bill).Error
	if err != nil {
		logx.Errorf("database exchange bill update failed, err:%v", err)
		err = biz.ExchangeBillUpdateFailed
		return
	}

	return
}

func (t *OrderService) UpdateDelegateBill(bill *entity.DelegateBill) (err error) {
	err = t.db.Model(&entity.DelegateBill{}).Save(bill).Error
	if err != nil {
		logx.Errorf("database delegate bill create failed, err:%v", err)
		err = biz.DelegateBillCreateFailed
		return
	}

	return
}

func (t *OrderService) UpdateDelegateOrder(order *entity.DelegateOrder) (err error) {
	err = t.db.Save(order).Error
	if err != nil {
		logx.Errorf("database delegate order update failed, order:%+v, err:%v", order, err)
		err = biz.DelegateOrderUpdateFailed
	}

	return
}

func (t *OrderService) ExpiresDelegateOrder() (expireds []*entity.DelegateOrder, err error) {
	err = t.db.Model(&entity.DelegateOrder{}).Where("status = ? and expires < ? ", int16(common.DelegateOrderStatusDelegated), common.TimeNowMilli()).Update("status", int16(common.DelegateOrderStatusExpired)).Error
	if err != nil {
		logx.Errorf("database delegate order expires failed, err:%v", err)
		err = biz.DelegateOrderExpiresFailed
		return
	}

	expireds = make([]*entity.DelegateOrder, 0)
	err = t.db.Model(&entity.DelegateOrder{}).Where("status = ? and delegate_amount > 0", int16(common.DelegateOrderStatusExpired)).Find(&expireds).Error
	if err != nil {
		logx.Errorf("database delegate order expires find failed, err:%v", err)
		err = biz.DelegateOrderExpiresFailed
		return
	}

	return
}

func (t *OrderService) DelegatedOrders() (delegateds []*entity.DelegateOrder, err error) {
	delegateds = make([]*entity.DelegateOrder, 0)

	err = t.db.Model(&entity.DelegateOrder{}).Where("status = ? and withdraw_time < ? and delegate_amount > 0", int16(common.DelegateOrderStatusDelegated), common.TimeNowMilli()).Find(&delegateds).Error
	if err != nil {
		logx.Errorf("database delegate order find delegated orders failed, err:%v", err)
		err = biz.DelegateOrderFindFailed
		return
	}

	return
}

func (t *OrderService) FindDelegateOrders(status int16) (orders []*entity.DelegateOrder, count int64, err error) {
	orders = make([]*entity.DelegateOrder, 0)
	err = t.db.Model(&entity.DelegateOrder{}).Where("status = ?", common.DelegateOrderStatus(status)).Find(&orders).Error
	if err != nil {
		logx.Errorf("database delegate order find failed, err:%v", err)
		err = biz.DelegateOrderFindFailed
		return
	}

	count = int64(len(orders))

	return
}

func (t *OrderService) FindExchangeOrders(status int16) (orders []*entity.ExchangeOrder, count int64, err error) {
	orders = make([]*entity.ExchangeOrder, 0)
	err = t.db.Model(&entity.ExchangeOrder{}).Where("status = ?", common.ExchangeOrderStatus(status)).Find(&orders).Error
	if err != nil {
		logx.Errorf("database exchange order find failed, err:%v", err)
		err = biz.ExchangeOrderFindFailed
		return
	}

	count = int64(len(orders))

	return
}

func (t *OrderService) ExpiresExchangeOrder() (err error) {
	err = t.db.Model(&entity.ExchangeOrder{}).Where("status = ? and expires < ?", common.ExchangeOrderStatus(common.ExchangeOrderStatusCreated), common.TimeNowMilli()).Update("status", int16(common.ExchangeOrderStatusExpired)).Error
	if err != nil {
		logx.Errorf("database exchange order expires failed, err:%v", err)
		err = biz.ExchangeOrderExpiresFailed
		return
	}

	return
}

func (t *OrderService) FindLastHourDelegateTxIds() (lids []string, err error) {
	lids = make([]string, 0)
	err = t.db.Model(&entity.DelegateOrder{}).Where("created_at > ?", common.TimeLastHourMilli()).Select("transaction_id").Find(&lids).Error
	if err != nil {
		logx.Errorf("database delegate order find last day's failed, err:%v", err)
		err = biz.DelegateOrderFindFailed
		return
	}

	return
}

func (t *OrderService) FindLastHourExchangeTxIds() (lids []string, err error) {
	lids = make([]string, 0)
	err = t.db.Model(&entity.ExchangeOrder{}).Where("created_at > ?", common.TimeLastHourMilli()).Select("transaction_id").Find(&lids).Error
	if err != nil {
		logx.Errorf("database exchange order find last day's failed, err:%v", err)
		err = biz.DelegateOrderFindFailed
		return
	}

	return
}

func (t *OrderService) FirstDelegateBill(oid int64) (bill *entity.DelegateBill, err error) {
	bill = &entity.DelegateBill{}
	err = t.db.Model(&entity.DelegateBill{}).Where("order_id = ?", oid).First(bill).Error
	if err != nil {
		logx.Errorf("database delegate bill get failed, oid:%v, err:%v", oid, err)
		err = biz.DelegateBillFindFailed
		return
	}

	return
}

func (t *OrderService) FirstDelegateWithdrawal(oid int64) (back *entity.DelegateWithdrawal, err error) {
	back = &entity.DelegateWithdrawal{}
	err = t.db.Model(&entity.DelegateWithdrawal{}).Where("order_id = ?", oid).First(back).Error
	if err != nil {
		logx.Errorf("database delegate takeback get failed, oid:%v, err:%v", oid, err)
		err = biz.DelegateBillFindFailed
		return
	}

	return
}

func (t *OrderService) UpdateDelegateWithdrawal(drawal *entity.DelegateWithdrawal) (err error) {
	err = t.db.Save(drawal).Error
	if err != nil {
		logx.Errorf("database delegate takeback update failed, oid:%v, err:%v", drawal.OrderId, err)
		err = biz.DelegateBillFindFailed
		return
	}

	return
}
