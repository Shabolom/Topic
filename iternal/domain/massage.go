package domain

import "github.com/gofrs/uuid"

type Massage struct {
	Base
	TopicID      uuid.UUID
	Text         string `gorm:"colum:text; type:text"`
	UserFilePath string `gorm:"colum:user_file_path; type:text"`
	UserID       uuid.UUID
	Likes        []Like
	DizLikes     []DizLike
}
