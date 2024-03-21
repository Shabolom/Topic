package models

import "github.com/gofrs/uuid"

type RespMassage struct {
	Text      string    `json:"text"`
	FilesPath []string  `json:"userFilePath"`
	UserID    uuid.UUID `json:"userID"`
	Likes     int       `json:"likes"`
	DizLikes  int       `json:"dizLikes"`
}
