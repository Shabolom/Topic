package domain

import "github.com/gofrs/uuid"

type Like struct {
	Base
	MassageID uuid.UUID
	UserID    uuid.UUID
}

type DizLike struct {
	Base
	MassageID uuid.UUID
	UserID    uuid.UUID
}
