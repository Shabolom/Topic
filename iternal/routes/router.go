package routes

import (
	"Arkadiy_Servis_authorization/iternal/api"
	"Arkadiy_Servis_authorization/iternal/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	user := api.NewUserApi()
	topic := api.NewTopicAPI()
	massages := api.NewMassagesAPI()

	authRequired := r.Group("/")
	authRequired.Use(middleware.Authorization())
	authRequired.Use(middleware.Timer())

	authRequiredAdmin := r.Group("/")
	authRequiredAdmin.Use(middleware.AdminAuthorization())
	authRequiredAdmin.Use(middleware.Timer())

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/users/register", user.Register)
	r.POST("/api/users/login", user.Login)

	{
		authRequired.POST("/api/massages/topic/:id", massages.PostMassage)
		authRequired.GET("/api/massages/rating/:id", massages.Rating)

		authRequired.GET("/api/users/get", user.GetUserSelf)

		authRequired.GET("/api/topics/join", topic.JoinTopic)
		authRequired.GET("/api/topics", topic.Topics)
		authRequired.GET("/api/topics/:id", topic.TopicMassages)
		//authRequired.GET("/api/topics/massages/:topic_id", topic.TopicMassages)
	}

	{
		authRequiredAdmin.POST("/api/users/status", user.ChangeStatus)
		authRequiredAdmin.DELETE("api/users/delete", user.DeleteUser)
		authRequiredAdmin.POST("api/users/permissions", user.SetPerm)
		authRequiredAdmin.GET("/api/users/get_all", user.GetUsers)
		authRequiredAdmin.GET("/api/users/get/:id", user.GetUser)

		authRequiredAdmin.POST("/api/topics/create", topic.CreateTopic)
	}

	return r
}
