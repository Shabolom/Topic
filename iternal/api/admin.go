package api

import (
	"Arkadiy_Servis_authorization/iternal/models"
	"Arkadiy_Servis_authorization/iternal/service"
	"Arkadiy_Servis_authorization/iternal/tools"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type AdminApi struct {
}

func NewAdminApi() *AdminApi {
	return &AdminApi{}
}

var adminService = service.NewAdminService()

func (aa *AdminApi) ChangeStatus(c *gin.Context) {
	var status models.Status

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &status)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := adminService.ChangeStatus(status)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

func (aa *AdminApi) DeleteUser(c *gin.Context) {
	userID := c.Query("user_id")

	result := adminService.DeleteUser(userID)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}
	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

func (aa *AdminApi) SetPerm(c *gin.Context) {
	var perm models.Permissions

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	err = json.Unmarshal(data, &perm)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := adminService.SetPerm(perm)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

func (aa *AdminApi) GetUsers(c *gin.Context) {

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result, err := adminService.GetUsers(page, limit)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

func (aa *AdminApi) GetUser(c *gin.Context) {
	userID := c.Param("id")

	result, err := adminService.GetUser(userID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

func (aa *AdminApi) CreateTopic(c *gin.Context) {
	var topic models.Topic

	pathToLogo, err := tools.MultipartFormDataTopic(&topic, c)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result := adminService.CreateTopic(topic, pathToLogo)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}
