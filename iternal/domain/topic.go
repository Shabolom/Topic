package domain

type Topic struct {
	Base
	Massages    []Massage
	PathToPhoto string `gorm:"colum:path_to_photo; type:text"`
	TopicName   string `gorm:"colum:topic_name; type:text"`
	About       string `gorm:"colum:about; type:text"`
	Creator     string `gorm:"colum:creator; type:text"`
	Users       []User `gorm:"many2many:user_topic;"`
}
