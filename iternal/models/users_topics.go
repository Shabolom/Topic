package models

import "github.com/gofrs/uuid"

type UsersTopics struct {
	UserID  uuid.UUID
	TopicID uuid.UUID
}
