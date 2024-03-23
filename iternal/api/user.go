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

type UserApi struct {
}

func NewUserApi() *UserApi {
	return &UserApi{}
}

var userService = service.NewUserService()

// @Summary	регистрация пользователя с выдачей токена
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	ввод	body		models.Register	true	"ввести логин и пароль"
// @Success	200		{string}	string	"вы зарегестрировались"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/users/register [post]
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

// @Summary	авторизация пользователя
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	ввод	body		models.Register	true	"ввести логин и пароль"
// @Success	200		{string}	string	"вы успешно авторизировались"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/users/login [post]
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

// @Summary	изменение статуса пользователя
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	ввод	body		models.Status	true	"измените статус"
// @Success	200		{string}	string	"статус изменен"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/users/status [post]
func (ua *UserApi) ChangeStatus(c *gin.Context) {
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

	result := userService.ChangeStatus(status)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

// @Summary	удаление пользователя
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	user_id	query    	string	true	"введите id_user"
// @Success	200		{string}	string	"пользователь удален"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router /api/users/delete [delete]
func (ua *UserApi) DeleteUser(c *gin.Context) {
	userID := c.Query("user_id")

	result := userService.DeleteUser(userID)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}
	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

// @Summary	выдача прав
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	ввод	body		models.Permissions	true	"выдайте права от 0 до 3"
// @Success	200		{string}	string	"права выданы"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router /api/users/permission [post]
func (ua *UserApi) SetPerm(c *gin.Context) {
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

	result := userService.SetPerm(perm)
	if result.Err != nil {
		tools.CreateError(result.Status, result.Err, c)
		log.WithField("component", "rest").Warn(result.Err)
		return
	}

	c.JSON(result.Status, result.Message)
	defer c.Request.Body.Close()
}

// @Summary	выдача прав
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	page	query		string	true	"страница"
// @Param	limit	query		string	true	"количество строк"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/users/get [get]
func (ua *UserApi) GetUsers(c *gin.Context) {

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

	result, err := userService.GetUsers(page, limit)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

// @Summary	выдача прав
// @Security ApiKeyAuth
// @Produce	json
// @Accept	json
// @Tags	User
// @Param	id	path		string	true	"user_id"
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/users/get/{id} [post]
func (ua *UserApi) GetUser(c *gin.Context) {
	userID := c.Param("id")

	result, err := userService.GetUser(userID)
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}

// @Summary	выдача прав
// @Security ApiKeyAuth
// @Produce	json
// @Tags	User
// @Failure	400		{object}	models.Error
// @Failure	500		{object}	models.Error
// @Failure	404		{object}	models.Error
// @Failure	409		{object}	models.Error
// @Router	/api/users/get [get]
func (ua *UserApi) GetUserSelf(c *gin.Context) {
	claims, err := tools.ParsTokenClaims(c.Request.Header.Get("Authorization"))
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	result, err := userService.GetUser(claims.UserID.String())
	if err != nil {
		tools.CreateError(http.StatusBadRequest, err, c)
		log.WithField("component", "rest").Warn(err)
		return
	}

	c.JSON(http.StatusOK, result)
	defer c.Request.Body.Close()
}
