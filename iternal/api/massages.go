package api

import (
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/service"
	"Arkadiy_Servis_authorization/iternal/tools"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type MassagesAPI struct{}

func NewMassagesAPI() *MassagesAPI {
	return &MassagesAPI{}
}

var massagesService = service.NewServiceMassages()

// @Summary	отправка сообщения в определенный топик
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Massages
// @Param	id			path     	string	true	"id сообщения"
// @Param	data		query		models.Massage	true	"ввести сообщение"
// @Param	file		formData  		file	false	"ввести сообщение"
// @Success	200			{string}	string "200"
// @Failure	400			{object}	models.Error
// @Failure	500			{object}	models.Error
// @Failure	404			{object}	models.Error
// @Failure	409			{object}	models.Error
// @Router	/api/massages/topic/{id} [post]
func (ma *MassagesAPI) Post(c *gin.Context) {
	var massage models.Massage

	massDirPath, topicID, claims, err := tools.MultipartFormDataMassage(&massage, c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := massagesService.Post(massage, massDirPath, claims, topicID)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.Status(result.Status)
	defer c.Request.Body.Close()
}

// @Summary	поставить дизлайк
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Massages
// @Param	id		path		string	true	"введите id сообщения"
// @Param	rating	query		string	true	"введите 1 если лайк 0 если диз лайк"
// @Success	200		{int}		int
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/massages/rating/{id} [get]
func (ma *MassagesAPI) Rating(c *gin.Context) {

	massageID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	rating, err := strconv.Atoi(c.Query("rating"))
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

	if rating == 1 {

		err2, likes := massagesService.Like(claims, massageID)
		if err2 != nil {
			tools.CreateError(http.StatusBadRequest, err2, c)
			log.WithField("component", "rest").Warn(err2)
			return
		}

		c.JSON(http.StatusOK, likes)
		defer c.Request.Body.Close()
	}

	err, dizLikes := massagesService.DizLike(claims, massageID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, dizLikes)
	defer c.Request.Body.Close()
}

// @Summary	удаление сообщения
// @Security ApiKeyAuth
// @Accept	json
// @Tags	Massages
// @Param	id		path		string	true	"введите id сообщения"
// @Success	200		{int}		int
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/massages/{id} [delete]
// @Router  /api/messages/users_message/{id} [delete]
func (ma *MassagesAPI) Delete(c *gin.Context) {
	id := c.Param("id")

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	if claims.UserPerm == 3 {
		err := massagesService.Delete(id, claims)
		if err != nil {
			tools.CreateError(http.StatusBadRequest, err, c)
			log.WithField("component", "rest").Warn(err)
			return
		}
		c.Status(http.StatusOK)
		return
	}

	err = massagesService.Delete(id, claims)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.Status(http.StatusOK)
	defer c.Request.Body.Close()
}
