package service

import (
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/repository"
	"Arkadiy_Servis_authorization/iternal/tools"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
)

type TopicService struct {
}

func NewTopicService() *TopicService {
	return &TopicService{}
}

var topicRepo = repository.NewTopicRepo()

func (ts *TopicService) GetTopics(page, limit int) ([]models.ResponseTopic, error) {
	var topics []models.ResponseTopic

	skip := page*limit - limit

	result, err := topicRepo.GetTopics(skip, limit)
	if err != nil {
		return []models.ResponseTopic{}, err
	}

	for _, topic := range result {
		topicEntity := models.ResponseTopic{
			Name:    topic.TopicName,
			About:   topic.About,
			Creator: topic.Creator,
			Users:   len(topic.Users),
		}
		topics = append(topics, topicEntity)
	}

	return topics, nil
}

func (ts *TopicService) GetTopic(topicID string) (models.ResponseTopic, error) {
	result, err := topicRepo.GetTopic(topicID)
	if err != nil {
		return models.ResponseTopic{}, err
	}

	topicEntity := models.ResponseTopic{
		Name:    result.TopicName,
		About:   result.About,
		Creator: result.Creator,
		Users:   len(result.Users),
	}

	return topicEntity, nil
}

func (ts *TopicService) TopicMassages(topicID, userID uuid.UUID, page, limit int) ([]models.RespMassage, error) {
	var respMassages []models.RespMassage
	skip := page*limit - limit

	err := topicRepo.FindUserInTopic(topicID, userID)
	if err != nil {
		return []models.RespMassage{}, err
	}

	result, err := topicRepo.TopicMassages(topicID, skip, limit)

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

func (ts *TopicService) JoinTopic(userID uuid.UUID, topicID string) tools.UserResult {

	result, err := topicRepo.JoinTopic(userID, topicID)
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

func (ts *TopicService) CreateTopic(topic models.Topic, path string) tools.UserResult {
	topicID, _ := uuid.NewV4()
	topicEntity := domain.Topic{
		TopicName:   topic.Name,
		About:       topic.About,
		Creator:     topic.Creator,
		PathToPhoto: path,
	}
	topicEntity.ID = topicID

	result := topicRepo.CreateTopic(topicEntity)

	return result
}

func (ts *TopicService) DeleteTopic(id string) tools.UserResult {
	uuidTopic, err := uuid.FromString(id)

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	result := topicRepo.DeleteTopic(uuidTopic)

	return result
}

func (ts *TopicService) ChangeTopic(topic models.Topic, path string, id string) error {
	uuidTopic, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	topicEntity := domain.Topic{
		TopicName:   topic.Name,
		About:       topic.About,
		Creator:     topic.Creator,
		PathToPhoto: path,
	}
	topicEntity.ID = uuidTopic

	err = topicRepo.ChangeTopic(topicEntity)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TopicService) DeleteUser(userID, topicID string) error {

	err := topicRepo.DeleteUser(userID, topicID)
	if err != nil {
		return err
	}

	return nil
}
