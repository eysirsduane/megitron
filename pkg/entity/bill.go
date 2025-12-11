package entity

type DelegateBill struct {
	Id            int64  `gorm:"type:bigserial;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	OrderId       int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Status  int16  `gorm:""`
	Curreny string `gorm:"size:16"`

	From    string `gorm:"size:255"`
	To      string `gorm:"size:255"`
	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	DelegatedAmount int64 `gorm:""`

	Time        uint64 `gorm:""`
	Description string `gorm:"size:1023"`

	TimeAts `gorm:"embedded"`
}

type ExchangeBill struct {
	Id            int64  `gorm:"type:bigserial;primaryKey;autoIncrement"`
	UserId        int64  `gorm:""`
	OrderId       int64  `gorm:""`
	TransactionId string `gorm:"size:255;uniqueIndex"`

	Status  int16  `gorm:""`
	Curreny string `gorm:"size:16"`

	From    string `gorm:"size:255"`
	To      string `gorm:"size:255"`
	FromHex string `gorm:"size:255"`
	ToHex   string `gorm:"size:255"`

	ExchangedAmount float64 `gorm:""`
	ExchangedSun    int64   `gorm:""`

	Time        uint64 `gorm:""`
	Description string `gorm:"size:1023"`

	TimeAts `gorm:"embedded"`
}
