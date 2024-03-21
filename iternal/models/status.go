package models

import (
	"github.com/gofrs/uuid"
)

type Status struct {
	UserID  uuid.UUID `json:"userID"`
	Confirm bool      `json:"confirm"`
}
