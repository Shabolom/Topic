package repository

import (
	"Arkadiy_Servis_authorization/config"
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/tools"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
)

type MassagesRepo struct {
}

func NewMassagesRepo() *MassagesRepo {
	return &MassagesRepo{}
}

func (mr *MassagesRepo) Post(massage domain.Massage) tools.UserResult {

	err := config.DB.
		Create(&massage).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Err:    nil,
		Status: http.StatusOK,
	}
}

func (mr *MassagesRepo) TopicMassages(topicID uuid.UUID) (domain.Topic, error) {
	var topic domain.Topic
	fmt.Println(topicID)
	err := config.DB.
		Model(&domain.Topic{}).
		Order("created-at").
		Preload("Massages", "topic_id IN (?)", topicID).
		Preload("Massages.Likes").
		Preload("Massages.DizLikes").
		Find(&topic).
		Error

	if err != nil {
		return domain.Topic{}, err
	}

	return topic, err
}

func (mr *MassagesRepo) FindUserInTopic(topicID uuid.UUID) error {
	var topic domain.Topic

	err := config.DB.
		Model(&domain.Topic{}).
		Preload("Users", "topic_id IN (?)", topicID).
		First(&topic).
		Error

	return err
}

func (mr *MassagesRepo) JoinTopic(userID uuid.UUID, topicID string) (domain.User, error) {
	id, _ := uuid.FromString(topicID)

	err := config.DB.
		Model(&domain.Topic{
			Base: domain.Base{
				ID: id,
			},
		}).
		Association("Users").
		Append(domain.User{
			Base: domain.Base{
				ID: userID,
			},
		}).
		Error

	if err != nil {
		return domain.User{}, err
	}

	return domain.User{}, nil
}

func (mr *MassagesRepo) CreateTopic(topic domain.Topic) tools.UserResult {
	err := mr.FindTopic("topic_name", topic.TopicName)
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

func (mr *MassagesRepo) FindDiz(massageID uuid.UUID, userID uuid.UUID) (domain.Massage, error) {
	var massage domain.Massage

	err := config.DB.
		Model(domain.Massage{}).
		Where("id = ?", massageID).
		Preload("DizLikes", "user_id IN (?)", userID).
		Preload("Likes", "user_id IN (?)", userID).
		//Preload("DizLikes").
		//Preload("Likes").
		Find(&massage).
		Error

	if err != nil {
		return domain.Massage{}, err
	}

	return massage, nil
}

func (mr *MassagesRepo) CreateDizLike(dizLike domain.DizLike) error {

	err := config.DB.
		Create(&dizLike).
		Error

	return err
}

func (mr *MassagesRepo) DeleteDizLike(massageID uuid.UUID) error {
	var dizLike domain.DizLike

	err := config.DB.
		Where("massage_id =?", massageID).
		Delete(&dizLike).
		Error

	return err
}

func (mr *MassagesRepo) DeleteLike(massageID uuid.UUID) error {
	var like domain.Like

	err := config.DB.
		Where("massage_id =?", massageID).
		Delete(&like).
		Error

	return err
}

func (mr *MassagesRepo) FindTopic(column string, find any) error {
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

func (mr *MassagesRepo) CreateLike(like domain.Like) error {

	err := config.DB.
		Create(&like).
		Error

	return err
}

func (mr *MassagesRepo) Delete(massageID string, claims tools.Claims) error {

	if claims.UserPerm == 3 {
		err := config.DB.
			Where("id =?", massageID).
			Delete(&domain.Massage{}).
			Error
		return err
	}

	err := config.DB.
		Where("id =? AND user_id =?", massageID, claims.UserID).
		Delete(&domain.Massage{}).
		Error

	return err
}
