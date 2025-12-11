package entity

type User struct {
	Id int64 `gorm:"type:bigserial;primaryKey;autoIncrement"`

	DisplayId string `gorm:"size:16;uniqueIndex"`
	Nickname  string `gorm:"size:255;uniqueIndex"`

	Username string `gorm:"size:255;uniqueIndex"`
	Password string `gorm:"size:255"`

	Telegram string `gorm:"size:255"`

	Contacts `gorm:"embedded"`
	TimeAts  `gorm:"embedded"`
}
