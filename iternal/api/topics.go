package api

import (
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/service"
	"Arkadiy_Servis_authorization/iternal/tools"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type TopicAPI struct {
}

func NewTopicAPI() *TopicAPI {
	return &TopicAPI{}
}

var topicService = service.NewTopicService()

// @Summary	получение топиков
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Success	200		{object}	[]models.ResponseTopic
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics [get]
func (ta *TopicAPI) Topics(c *gin.Context) {

	result, err := topicService.Topics()
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

// @Summary	получение сообщений в конкретном топике
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	id		path		string	true	"id топика"
// @Success	200		{int}		int
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/{id} [get]
func (ta *TopicAPI) TopicMassages(c *gin.Context) {

	topicID := c.Param("id")
	uuidMassage, err := uuid.FromString(topicID)

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result, err := topicService.TopicMassages(uuidMassage, claims.UserID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

// @Summary	присоединение к топику
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	topic_id			query		string	true	"id топика"
// @Success	200		{int}		int		200
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/join [get]
func (ta *TopicAPI) JoinTopic(c *gin.Context) {
	topicID := c.Query("topic_id")

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := topicService.JoinTopic(claims.UserID, topicID)

	c.Status(result.Status)
	defer c.Request.Body.Close()
}

// @Summary	присоединение к топику
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	data	query	models.Topic	true 	"заполните указанные поля"
// @Param	file	formData  	file	false 	"файл"
// @Success	200		{string}	string		"топик создан"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/create [get]
func (ta *TopicAPI) CreateTopic(c *gin.Context) {
	var topic models.Topic

	pathToLogo, err := tools.MultipartFormDataTopic(&topic, c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := topicService.CreateTopic(topic, pathToLogo)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}
