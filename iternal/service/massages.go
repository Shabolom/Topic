package service

import (
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/repository"
	"Arkadiy_Servis_authorization/iternal/tools"
	"errors"
	"github.com/gofrs/uuid"
	"strings"
)

type ServiceMassages struct {
}

func NewServiceMassages() *ServiceMassages {
	return &ServiceMassages{}
}

var massageRepo = repository.NewMassagesRepo()

func (sm *ServiceMassages) Post(massage models.Massage, massPaths []string, claims tools.Claims, topicID string) tools.UserResult {
	massageId, _ := uuid.NewV4()
	topicUUID, _ := uuid.FromString(topicID)
	formatPaths := strings.Join(massPaths, "(space)")

	userMassageEntity := domain.Massage{
		TopicID:      topicUUID,
		Text:         massage.Text,
		UserFilePath: formatPaths,
		UserID:       claims.UserID,
	}
	userMassageEntity.ID = massageId

	result := massageRepo.Post(userMassageEntity)

	return result
}

func (sm *ServiceMassages) TopicMassages(topicID uuid.UUID, userID uuid.UUID) ([]models.RespMassage, error) {
	var respMassages []models.RespMassage

	err := massageRepo.FindUserInTopic(topicID)
	if err != nil {
		return []models.RespMassage{}, err
	}

	result, err := massageRepo.TopicMassages(topicID)

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

func (sm *ServiceMassages) DizLike(claims tools.Claims, massageID uuid.UUID) (error, int) {

	dizID, _ := uuid.NewV4()
	dizLikeEntity := domain.DizLike{
		MassageID: massageID,
		UserID:    claims.UserID,
	}
	dizLikeEntity.ID = dizID

	result, err := massageRepo.FindDiz(massageID, claims.UserID)
	if err != nil {
		return err, 0
	}

	if len(result.DizLikes) == 0 && len(result.Likes) == 0 {

		err2 := massageRepo.CreateDizLike(dizLikeEntity)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}
		return nil, len(result.DizLikes)

	} else if len(result.DizLikes) == 1 && len(result.Likes) == 0 {

		err2 := massageRepo.DeleteDizLike(massageID)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}
		return nil, len(result.DizLikes)

	} else if len(result.DizLikes) == 0 && len(result.Likes) == 1 {

		err2 := massageRepo.DeleteLike(massageID)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}

		err2 = massageRepo.CreateDizLike(dizLikeEntity)
		if err2 != nil {
			return err2, len(result.DizLikes)
		}

		return nil, len(result.DizLikes)
	}

	return errors.New("как-то занесли и лайк и диз лайк что не возможно"), 0
}

func (sm *ServiceMassages) Like(claims tools.Claims, massageID uuid.UUID) (error, int) {

	dizID, _ := uuid.NewV4()
	likeEntity := domain.Like{
		MassageID: massageID,
		UserID:    claims.UserID,
	}
	likeEntity.ID = dizID

	result, err := massageRepo.FindDiz(massageID, claims.UserID)
	if err != nil {
		return err, 0
	}

	if len(result.DizLikes) == 0 && len(result.Likes) == 0 {

		err2 := massageRepo.CreateLike(likeEntity)
		if err2 != nil {
			return err2, len(result.Likes)
		}
		return nil, len(result.Likes)
	} else if len(result.DizLikes) == 0 && len(result.Likes) == 1 {

		err2 := massageRepo.DeleteLike(massageID)
		if err2 != nil {
			return err2, len(result.Likes)
		}
		return nil, len(result.Likes)
	} else if len(result.DizLikes) == 1 && len(result.Likes) == 0 {

		err2 := massageRepo.DeleteDizLike(massageID)
		if err2 != nil {
			return err2, len(result.Likes)
		}

		err2 = massageRepo.CreateLike(likeEntity)
		if err2 != nil {
			return err2, len(result.Likes)
		}

		return nil, len(result.Likes)
	}

	return errors.New("как-то занесли и лайк и диз лайк что не возможно"), 0
}

func (sm *ServiceMassages) Delete(massageID string, claims tools.Claims) error {

	err := massageRepo.Delete(massageID, claims)
	if err != nil {
		return err
	}

	return nil
}
