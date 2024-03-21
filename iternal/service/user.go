package service

import (
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/repository"
	"Arkadiy_Servis_authorization/iternal/tools"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
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

func (us *UserService) GetUser(userID uuid.UUID) tools.UserResult {

	result, err := userRepo.FindOnUser("id", userID)
	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Status: http.StatusOK,
		Result: result,
	}
}

func (us *UserService) JoinTopic(userID uuid.UUID, topicID string) tools.UserResult {

	result, err := userRepo.JoinTopic(userID, topicID)
	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Status: http.StatusOK,
		Result: result,
	}
}

func (us *UserService) PostMassage(massage models.Massage, massPaths []string, claims tools.Claims, topicID string) tools.UserResult {
	massageId, _ := uuid.NewV4()
	topicUUID, _ := uuid.FromString(topicID)
	formatPaths := strings.Join(massPaths, "(space)")
	fmt.Println(formatPaths)
	userMassageEntity := domain.Massage{
		TopicID:      topicUUID,
		Text:         massage.Text,
		UserFilePath: formatPaths,
		UserID:       claims.UserID,
	}
	userMassageEntity.ID = massageId

	result := userRepo.PostMassage(userMassageEntity)

	return result
}

func (us *UserService) DizLike(claims tools.Claims, massageID uuid.UUID) (error, int) {

	dizID, _ := uuid.NewV4()
	dizLikeEntity := domain.DizLike{
		MassageID: massageID,
		UserID:    claims.UserID,
	}
	dizLikeEntity.ID = dizID

	result, err := userRepo.FindDiz(massageID, claims.UserID)
	if err != nil {
		return err, 0
	}

	if len(result.DizLikes) == 0 && len(result.Likes) == 0 {

		err2 := userRepo.CreateDizLike(dizLikeEntity)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}
		return nil, len(result.DizLikes)
	} else if len(result.DizLikes) == 1 && len(result.Likes) == 0 {

		err2 := userRepo.DeleteDizLike(massageID)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}
		return nil, len(result.DizLikes)
	} else if len(result.DizLikes) == 0 && len(result.Likes) == 1 {

		err2 := userRepo.DeleteLike(massageID)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}

		err2 = userRepo.CreateDizLike(dizLikeEntity)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}

		return nil, len(result.DizLikes)
	}

	return errors.New("как-то занесли и лайк и диз лайк что не возможно"), 0
}

func (us *UserService) Like(claims tools.Claims, massageID uuid.UUID) (error, int) {

	dizID, _ := uuid.NewV4()
	likeEntity := domain.Like{
		MassageID: massageID,
		UserID:    claims.UserID,
	}
	likeEntity.ID = dizID

	result, err := userRepo.FindDiz(massageID, claims.UserID)
	if err != nil {
		return err, 0
	}

	if len(result.DizLikes) == 0 && len(result.Likes) == 0 {

		err2 := userRepo.CreateLike(likeEntity)
		if err2 != nil {
			return err2, len(result.Likes)
		}
		return nil, len(result.Likes)
	} else if len(result.DizLikes) == 0 && len(result.Likes) == 1 {

		err2 := userRepo.DeleteLike(massageID)
		if err2 != nil {
			return err2, len(result.Likes)
		}
		return nil, len(result.Likes)
	} else if len(result.DizLikes) == 1 && len(result.Likes) == 0 {

		err2 := userRepo.DeleteDizLike(massageID)
		if err2 != nil {
			return err2, len(result.Likes)
		}

		err2 = userRepo.CreateLike(likeEntity)
		if err2 != nil {
			return err2, len(result.Likes)
		}

		return nil, len(result.Likes)
	}

	return errors.New("как-то занесли и лайк и диз лайк что не возможно"), 0
}

func (us *UserService) TopicMassages(topicID uuid.UUID, userID uuid.UUID) ([]models.RespMassage, error) {
	var respMassages []models.RespMassage

	err := userRepo.FindUserInTopic(topicID, userID)
	if err != nil {
		return []models.RespMassage{}, err
	}

	result, err := userRepo.TopicMassages(topicID)

	for _, massage := range result.Massages {
		paths := strings.Split(massage.UserFilePath, "(space)")

		respMassageEntity := models.RespMassage{
			Text:      massage.Text,
			FilesPath: paths,
			UserID:    massage.UserID,
			Likes:     len(massage.Likes),
			DizLikes:  len(massage.DizLikes),
		}
		respMassages = append(respMassages, respMassageEntity)
	}

	if err != nil {
		return []models.RespMassage{}, err
	}

	return respMassages, nil
}
