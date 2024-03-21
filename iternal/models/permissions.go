package models

import "github.com/gofrs/uuid"

type Permissions struct {
	UserID uuid.UUID `json:"userID"`
	Perm   int       `json:"perm"`
}
