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

type AdminRepo struct {
}

func NewAdminRepo() *AdminRepo {
	return &AdminRepo{}
}

func (ar *AdminRepo) ChangeStatus(status string, userID uuid.UUID) tools.UserResult {
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

func (ar *AdminRepo) DeleteUser(userStrID string) tools.UserResult {
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

func (ar *AdminRepo) SetPerm(userPerm models.Permissions) tools.UserResult {
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

func (ar *AdminRepo) GetUsers(limit, skip int) ([]domain.User, error) {
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

func (ar *AdminRepo) GetUser(stringUserID string) (domain.User, error) {
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

func (ar *AdminRepo) CreateTopic(topic domain.Topic) tools.UserResult {
	err := ar.FindTopic("topic_name", topic.TopicName)
	if err == nil {
		return tools.UserResult{
			Err:     nil,
			Message: "такой топик уже существует",
			Status:  http.StatusBadRequest,
		}
	}
	err = config.DB.
		Create(&topic).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Message: "топик создан",
		Status:  http.StatusCreated,
	}
}

func (ar *AdminRepo) FindTopic(column string, find any) error {
	var topics domain.Topic

	err := config.DB.
		Where(column+"=?", find).
		First(&topics).
		Error
	if err != nil {
		return err
	}

	return err
}
