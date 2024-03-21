package api

import (
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/service"
	"Arkadiy_Servis_authorization/iternal/tools"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type UserApi struct {
}

func NewUserApi() *UserApi {
	return &UserApi{}
}

var userService = service.NewUserService()

func (ua *UserApi) Register(c *gin.Context) {
	var user models.Register

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.Register(user)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	err = tools.JWTCreator(result.Result, c)
	if err != nil {
		tools.CreateError(501, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

func (ua *UserApi) Login(c *gin.Context) {
	var user models.Register

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.Login(user)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	err = tools.JWTCreator(result.Result, c)
	if err != nil {
		tools.CreateError(501, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

func (ua *UserApi) GetUser(c *gin.Context) {
	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.GetUser(claims.UserID)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Result)
	defer c.Request.Body.Close()
}

func (ua *UserApi) JoinTopic(c *gin.Context) {
	topicID := c.Query("topic_id")

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.JoinTopic(claims.UserID, topicID)

	c.Status(result.Status)
	defer c.Request.Body.Close()
}

func (ua *UserApi) PostMassage(c *gin.Context) {
	var massage models.Massage

	massDirPath, topicID, claims, err := tools.MultipartFormDataMassage(&massage, c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := userService.PostMassage(massage, massDirPath, claims, topicID)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.Status(result.Status)
	defer c.Request.Body.Close()
}

func (ua *UserApi) DizLike(c *gin.Context) {

	massageID := c.Query("massage_id")
	uuidMassage, err := uuid.FromString(massageID)

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err, dizLikes := userService.DizLike(claims, uuidMassage)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, dizLikes)
	defer c.Request.Body.Close()
}

func (ua *UserApi) Like(c *gin.Context) {

	massageID := c.Query("massage_id")
	uuidMassage, err := uuid.FromString(massageID)

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err, likes := userService.Like(claims, uuidMassage)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, likes)
	defer c.Request.Body.Close()
}

func (ua *UserApi) TopicMassages(c *gin.Context) {

	topicID := c.Query("topic_id")
	uuidMassage, err := uuid.FromString(topicID)

	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result, err := userService.TopicMassages(uuidMassage, claims.UserID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}
