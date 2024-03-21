package routes

import (
	"Arkadiy_Servis_authorization/iternal/api"
	"Arkadiy_Servis_authorization/iternal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	user := api.NewUserApi()
	admin := api.NewAdminApi()

	authRequired := r.Group("/")
	authRequired.Use(middleware.Authorization())
	authRequired.Use(middleware.Timer())

	authRequiredAdmin := r.Group("/")
	authRequiredAdmin.Use(middleware.AdminAuthorization())
	authRequiredAdmin.Use(middleware.Timer())

	r.POST("/api/register_user", user.Register)
	r.POST("/api/login_user_user", user.Login)

	{
		authRequired.GET("/api/get_self_info", user.GetUser)
		authRequired.GET("/api/join_topic", user.JoinTopic)
		authRequired.POST("/api/post_massage", user.PostMassage)
		authRequired.GET("/api/diz_like", user.DizLike)
		authRequired.GET("/api/like", user.Like)
		authRequired.GET("/api/all_massages_topic", user.TopicMassages)
	}

	{
		authRequiredAdmin.POST("/api/change_user_status", admin.ChangeStatus)
		authRequiredAdmin.POST("/api/create_topic", admin.CreateTopic)
		authRequiredAdmin.DELETE("api/user_delete", admin.DeleteUser)
		authRequiredAdmin.POST("api/user_permissions", admin.SetPerm)
		authRequiredAdmin.GET("/api/get_users", admin.GetUsers)
		authRequiredAdmin.GET("/api/get_user/user_id/:id", admin.GetUser)
	}

	return r
}
