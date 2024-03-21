package service

import (
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/repository"
	"Arkadiy_Servis_authorization/iternal/tools"
	"errors"
	"github.com/gofrs/uuid"
	"net/http"
)

type AdminService struct {
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

var adminRepo = repository.NewAdminRepo()

func (as *AdminService) ChangeStatus(status models.Status) tools.UserResult {

	if status.Confirm {
		result := adminRepo.ChangeStatus("confirm", status.UserID)
		return result
	}

	result := adminRepo.ChangeStatus("not confirm", status.UserID)
	return result
}

func (as *AdminService) DeleteUser(userID string) tools.UserResult {
	result := adminRepo.DeleteUser(userID)
	return result
}

func (as *AdminService) SetPerm(userPerm models.Permissions) tools.UserResult {
	if userPerm.Perm > 3 || userPerm.Perm < 0 {
		return tools.UserResult{
			Message: "введите права доступа от 0 до 3",
			Status:  http.StatusBadRequest,
		}
	}
	result := adminRepo.SetPerm(userPerm)
	return result
}

func (as *AdminService) GetUsers(page, limit int) ([]domain.User, error) {
	skip := page*limit - limit

	if page == 0 || limit == 0 {
		return []domain.User{}, errors.New("задайте необходимое количество элементов и страниц в param")
	}

	result, err := adminRepo.GetUsers(limit, skip)
	if err != nil {
		return []domain.User{}, err
	}
	return result, nil
}

func (as *AdminService) GetUser(userID string) (domain.User, error) {
	result, err := adminRepo.GetUser(userID)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}

func (as *AdminService) CreateTopic(topic models.Topic, path string) tools.UserResult {
	topicID, _ := uuid.NewV4()
	topicEntity := domain.Topic{
		TopicName:   topic.Name,
		About:       topic.About,
		Creator:     topic.Creator,
		PathToPhoto: path,
	}
	topicEntity.ID = topicID

	result := adminRepo.CreateTopic(topicEntity)

	return result
}
