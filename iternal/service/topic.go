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

func (ts *TopicService) Topics() ([]models.ResponseTopic, error) {
	var topics []models.ResponseTopic

	result, err := topicRepo.Topics()
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

func (ts *TopicService) TopicMassages(topicID uuid.UUID, userID uuid.UUID) ([]models.RespMassage, error) {
	var respMassages []models.RespMassage

	err := topicRepo.FindUserInTopic(topicID, userID)
	if err != nil {
		return []models.RespMassage{}, err
	}

	result, err := topicRepo.TopicMassages(topicID)

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
