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

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

var userRepo = repository.NewUserRepo()

func (us *UserService) Register(user models.Register) tools.UserResult {
	userID, _ := uuid.NewV4()

	password, err := tools.HashPassword(user.Password)
	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	userEntity := domain.User{
		Status:      "not confirmed",
		Login:       user.Login,
		Password:    password,
		Permissions: 0,
	}
	userEntity.ID = userID

	result := userRepo.Register(userEntity)

	return result
}

func (us *UserService) Login(user models.Register) tools.UserResult {

	result, err := userRepo.FindOnUser("login", user.Login)
	if err != nil {
		return tools.UserResult{
			Message: "не верный логин или пароль",
			Status:  http.StatusBadRequest,
		}
	}

	if !tools.CheckPasswordHash(user.Password, result.Password) {
		return tools.UserResult{
			Err:     nil,
			Message: "не верный логин или пароль",
			Status:  http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Message: "вы успешно авторезированны",
		Status:  http.StatusOK,
		Result:  result,
	}
}

func (us *UserService) ChangeStatus(status models.Status) tools.UserResult {

	if status.Confirm {
		result := userRepo.ChangeStatus("confirm", status.UserID)
		return result
	}

	result := userRepo.ChangeStatus("not confirm", status.UserID)
	return result
}

func (us *UserService) DeleteUser(userID string) tools.UserResult {
	result := userRepo.DeleteUser(userID)
	return result
}

func (us *UserService) SetPerm(userPerm models.Permissions) tools.UserResult {
	if userPerm.Perm > 3 || userPerm.Perm < 0 {
		return tools.UserResult{
			Message: "введите права доступа от 0 до 3",
			Status:  http.StatusBadRequest,
		}
	}
	result := userRepo.SetPerm(userPerm)
	return result
}

func (us *UserService) GetUsers(page, limit int) ([]domain.User, error) {
	skip := page*limit - limit

	if page == 0 || limit == 0 {
		return []domain.User{}, errors.New("задайте необходимое количество элементов и страниц в param")
	}

	result, err := userRepo.GetUsers(limit, skip)
	if err != nil {
		return []domain.User{}, err
	}
	return result, nil
}

func (us *UserService) GetUser(userID string) (domain.User, error) {
	result, err := userRepo.GetUser(userID)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}
