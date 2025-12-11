package entity

type DelegateOrder struct {
	Id            int64  `gorm:"type:BIGSERIAL;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Typo    int32  `gorm:""`
	Status  uint16 `gorm:""`
	Curreny string `gorm:"size:16"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    uint64  `gorm:""`

	From    string `gorm:"size:255"`
	To      string `gorm:"size:255"`
	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	DelegateAmount int64 `gorm:""`

	Time         uint64 `gorm:""`
	Expires      uint64 `gorm:""`
	WithdrawTime uint64
	FailedTimes  uint32

	Description string `gorm:"size:2047"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}

type ExchangeOrder struct {
	Id            int64  `gorm:"type:bigserial;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Typo    string `gorm:"size:16"`
	Status  int16  `gorm:""`
	Curreny string `gorm:"size:16"`

	ReceivedAmount float64 `gorm:""`
	ReceivedSun    int64   `gorm:""`

	From    string `gorm:"size:255"`
	To      string `gorm:"size:255"`
	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	ThenRate         float64
	ExchangeRate     float64
	ExchangeDiscount float64
	ExchangeAmount   float64 `gorm:""`
	ExchangeSun      int64   `gorm:""`

	Time    uint64 `gorm:""`
	Expires uint64 `gorm:""`

	Description string `gorm:"size:1023"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}
