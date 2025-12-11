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

type DelegateOrderStatus uint16

const (
	DelegateOrderStatusCreated   DelegateOrderStatus = 0
	DelegateOrderStatusPending   DelegateOrderStatus = 100
	DelegateOrderStatusExpired   DelegateOrderStatus = 300
	DelegateOrderStatusCanceled  DelegateOrderStatus = 400
	DelegateOrderStatusDelegated DelegateOrderStatus = 700
	DelegateOrderStatusError     DelegateOrderStatus = 800
	DelegateOrderStatusFinished  DelegateOrderStatus = 900
)

type DelegateBillStatus uint16

const (
	DelegateBillStatusCreated  DelegateBillStatus = 0
	DelegateBillStatusPending  DelegateBillStatus = 100
	DelegateBillStatusExpired  DelegateBillStatus = 300
	DelegateBillStatusCanceled DelegateBillStatus = 400
	DelegateBillStatusPaid     DelegateBillStatus = 700
	DelegateBillStatusError    DelegateBillStatus = 800
	DelegateBillStatusSuccess  DelegateBillStatus = 900
)

type ExchangeOrderStatus uint16

const (
	ExchangeOrderStatusCreated   ExchangeOrderStatus = 0
	ExchangeOrderStatusPending   ExchangeOrderStatus = 100
	ExchangeOrderStatusExpired   DelegateOrderStatus = 300
	ExchangeOrderStatusCanceled  ExchangeOrderStatus = 400
	ExchangeOrderStatusExchanged ExchangeOrderStatus = 700
	ExchangeOrderStatusError     DelegateOrderStatus = 800
	ExchangeOrderStatusFinished  ExchangeOrderStatus = 900
)

type ExchangeBillStatus uint16

const (
	ExchangeBillStatusCreated  ExchangeBillStatus = 0
	ExchangeBillStatusPending  ExchangeBillStatus = 100
	ExchangeBillStatusExpired  ExchangeBillStatus = 300
	ExchangeBillStatusCanceled ExchangeBillStatus = 400
	ExchangeBillStatusPaid     ExchangeBillStatus = 700
	ExchangeBillStatusError    ExchangeBillStatus = 800
	ExchangeBillStatusSuccess  ExchangeBillStatus = 900
)
