package repository

import (
	"Arkadiy_Servis_authorization/config"
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/tools"
	"github.com/gofrs/uuid"
	uuid2 "github.com/google/uuid"
	"net/http"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) Register(user domain.User) tools.UserResult {
	if _, err := ur.FindOnUser("login", user.Login); err == nil {
		return tools.UserResult{
			Message: "такой пользователь уже существует",
			Status:  http.StatusBadRequest,
		}
	}
	err := config.DB.
		Create(&user).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Err:     nil,
		Message: "успешно зарегестрированны",
		Status:  http.StatusCreated,
		Result:  user,
	}
}

func (ur *UserRepo) FindOnUser(what string, find any) (domain.User, error) {
	var user domain.User

	err := config.DB.
		Where(what+" =?", find).
		First(&user).
		Error

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) ChangeStatus(status string, userID uuid.UUID) tools.UserResult {
	var user domain.User
	err := config.DB.Model(&user).
		Where("id =?", userID).
		Update("status", status).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Message: "статус изменен",
		Status:  http.StatusOK,
	}
}

func (ur *UserRepo) DeleteUser(userStrID string) tools.UserResult {
	var user domain.User
	userID, _ := uuid2.Parse(userStrID)

	err := config.DB.
		Where("id =?", userID).
		Delete(&user).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Message: "пользователь удален",
		Status:  http.StatusOK,
	}
}

func (ur *UserRepo) SetPerm(userPerm models.Permissions) tools.UserResult {
	var user domain.User

	err := config.DB.Model(&user).
		Where("id =?", userPerm.UserID).
		Update("permissions", userPerm.Perm).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Message: "права изменены",
		Status:  http.StatusOK,
	}
}

func (ur *UserRepo) GetUsers(limit, skip int) ([]domain.User, error) {
	var user []domain.User

	err := config.DB.
		Order("id asc").
		Limit(limit).
		Offset(skip).
		Find(&user).
		Error

	if err != nil {
		return []domain.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) GetUser(stringUserID string) (domain.User, error) {
	var user domain.User

	userID, _ := uuid2.Parse(stringUserID)
	err := config.DB.
		Where("id =?", userID).
		First(&user).
		Error

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
