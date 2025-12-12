package common

const (
	StringSuccess = "SUCCESS"
)

const (
	TronTransactionTypoTransfer                 string = "Transfer"
	TronTransactionTypoTransferContract         string = "TransferContract"
	TronTransactionTypoDelegateResourceContract string = "DelegateResourceContract"
)

type CurrenyTypo string

const (
	CurrencyTypoTrx    CurrenyTypo = "trx"
	CurrencyTypoUsdt   CurrenyTypo = "usdt"
	CurrencyTypoEnergy CurrenyTypo = "energy"
)

type ExchangeRateTypo string

const (
	ExchangeRateTrx2Energy ExchangeRateTypo = "trx2energy"
	ExchangeRateUsdt2Trx   ExchangeRateTypo = "usdt2trx"
)

type DelegateResourceTypo uint32

const (
	DelegateResourceTypoEnergy    DelegateResourceTypo = 1
	DelegateResourceTypoBindWidth DelegateResourceTypo = 0
)

type ExchangeTypo string

const (
	ExchangeTypoUsdt2Trx ExchangeTypo = "usdt2trx"
	ExchangeTypoTrx2Usdt ExchangeTypo = "trx2usdt"
)

type DelegateOrderStatus string

const (
	DelegateOrderStatusCreated   DelegateOrderStatus = "已创建"
	DelegateOrderStatusPending   DelegateOrderStatus = "已挂起"
	DelegateOrderStatusExpired   DelegateOrderStatus = "已过期"
	DelegateOrderStatusCanceled  DelegateOrderStatus = "已取消"
	DelegateOrderStatusDelegated DelegateOrderStatus = "已委托"
	DelegateOrderStatusError     DelegateOrderStatus = "错误"
	DelegateOrderStatusFinished  DelegateOrderStatus = "已完成"
)

type DelegateBillStatus string

const (
	DelegateBillStatusCreated  DelegateBillStatus = "已创建"
	DelegateBillStatusPending  DelegateBillStatus = "已挂起"
	DelegateBillStatusExpired  DelegateBillStatus = "已过期"
	DelegateBillStatusCanceled DelegateBillStatus = "已取消"
	DelegateBillStatusPaid     DelegateBillStatus = "已委托"
	DelegateBillStatusError    DelegateBillStatus = "错误"
	DelegateBillStatusSuccess  DelegateBillStatus = "已完成"
)

type ExchangeOrderStatus string

const (
	ExchangeOrderStatusCreated   ExchangeOrderStatus = "已创建"
	ExchangeOrderStatusPending   ExchangeOrderStatus = "已挂起"
	ExchangeOrderStatusExpired   DelegateOrderStatus = "已过期"
	ExchangeOrderStatusCanceled  ExchangeOrderStatus = "已取消"
	ExchangeOrderStatusExchanged ExchangeOrderStatus = "已委托"
	ExchangeOrderStatusError     DelegateOrderStatus = "错误"
	ExchangeOrderStatusFinished  ExchangeOrderStatus = "已完成"
)

type ExchangeBillStatus string

const (
	ExchangeBillStatusCreated  ExchangeBillStatus = "已创建"
	ExchangeBillStatusPending  ExchangeBillStatus = "已挂起"
	ExchangeBillStatusExpired  ExchangeBillStatus = "已过期"
	ExchangeBillStatusCanceled ExchangeBillStatus = "已取消"
	ExchangeBillStatusPaid     ExchangeBillStatus = "已委托"
	ExchangeBillStatusError    ExchangeBillStatus = "错误"
	ExchangeBillStatusSuccess  ExchangeBillStatus = "已完成"
)
