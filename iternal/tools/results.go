package tools

import "Arkadiy_Servis_authorization/iternal/domain"

type UserResult struct {
	Err     error
	Message string
	Status  int
	Result  domain.User
}
