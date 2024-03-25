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
// @Param	page		query		string	true	"количество страниц"
// @Param	limit		query		string	true	"количество элементов на странице"
// @Success	200		{object}	[]models.ResponseTopic
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics [get]
func (ta *TopicAPI) GetTopics(c *gin.Context) {
	page, limit, err := tools.GetQueryPagination(c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result, err := topicService.GetTopics(page, limit)
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
// @Param	page	query		string	true	"страница"
// @Param	limit	query		string	true	"колличество элементов на странице"
// @Success	200		{int}		int
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/messages/topic/{id} [get]
func (ta *TopicAPI) TopicMassages(c *gin.Context) {
	uuidMassage, err := uuid.FromString(c.Param("id"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	page, limit, err := tools.GetQueryPagination(c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result, err := topicService.TopicMassages(uuidMassage, claims.UserID, page, limit)
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
// @Param	id		path		string	true	"id топика"
// @Success	200		{int}		int		200
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/join/{id} [get]
func (ta *TopicAPI) JoinTopic(c *gin.Context) {
	topicID := c.Param("id")

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

// @Summary	создание топика
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
// @Router	/api/topics [post]
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

// @Summary	удаление топика
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	id		query		string	true	"id"
// @Success	200		{string}	string		"топик удален"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/{id} [delete]
func (ta *TopicAPI) DeleteTopic(c *gin.Context) {
	id := c.Query("id")

	result := topicService.DeleteTopic(id)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

// @Summary	изменение топика
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	data	query	models.Topic	true 	"заполните указанные поля"
// @Param	file	formData  	file	false 	"файл"
// @Success	200		{string}	string		"топик обновлен"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/{id} [put]
func (ta *TopicAPI) ChangeTopic(c *gin.Context) {
	var topic models.Topic
	id := c.Query("name")

	pathToLogo, err := tools.MultipartFormDataTopic(&topic, c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = topicService.ChangeTopic(topic, pathToLogo, id)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, "топик удален")
	defer c.Request.Body.Close()
}

// @Summary	отображение выбранного топикоа
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	id		path  		string	true 	"id"
// @Success	200		{object}	models.ResponseTopic
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/{id} [get]
func (ta *TopicAPI) GetTopic(c *gin.Context) {
	topicID := c.Param("id")

	result, err := topicService.GetTopic(topicID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

// @Summary	удаление пользователя из топмка
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Topic
// @Param	id		path  		string	true 	"id"
// @Param	user_id		path  		string	true 	"user_id"
// @Success	200		{int}		200
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/topics/{id} [get]
func (ta *TopicAPI) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	topicID := c.Param("id")

	err := topicService.DeleteUser(userID, topicID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.Status(http.StatusOK)
	defer c.Request.Body.Close()
}
