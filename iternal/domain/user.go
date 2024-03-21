package domain

type User struct {
	Base
	Status      string `gorm:"colum:status; type:text"`
	Login       string `gorm:"colum:login; type:text"`
	Password    string `gorm:"colum:password; type:text"`
	Permissions int    `gorm:"colum:permissions; type:int"`
	Massages    []Massage
	Topics      []Topic `gorm:"many2many:user_topic;"`
}
