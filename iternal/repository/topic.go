package repository

import (
	"Arkadiy_Servis_authorization/config"
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/tools"
	"github.com/gofrs/uuid"
	"net/http"
)

type TopicRepo struct {
}

func NewTopicRepo() *TopicRepo {
	return &TopicRepo{}
}

func (tr *TopicRepo) Topics() ([]domain.Topic, error) {
	var topic []domain.Topic

	err := config.DB.
		Order("created-at").
		Preloads("Users").
		Find(&topic).
		Error

	if err != nil {
		return []domain.Topic{}, err
	}

	return topic, err
}

func (tr *TopicRepo) TopicMassages(topicID uuid.UUID) (domain.Topic, error) {
	var topic domain.Topic

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

func (tr *TopicRepo) JoinTopic(userID uuid.UUID, topicID string) (domain.User, error) {
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

func (tr *TopicRepo) CreateTopic(topic domain.Topic) tools.UserResult {
	err := tr.FindTopic("topic_name", topic.TopicName)
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

func (tr *TopicRepo) FindTopic(column string, find any) error {
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

func (tr *TopicRepo) FindUserInTopic(topicID uuid.UUID, userID uuid.UUID) error {
	var topic domain.Topic

	err := config.DB.
		Model(&domain.Topic{}).
		Preload("Users", "topic_id IN (?)", topicID).
		First(&topic).
		Error

	return err
}
