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
		authRequired.GET("/api/messages/topic/:id", topic.TopicMassages)
		authRequired.POST("/api/messages/topic/:id", massages.Post)
		authRequired.DELETE("/api/messages/:id", massages.Delete)
		authRequired.GET("/api/messages/rating/:id", massages.Rating)

		authRequired.GET("/api/users", user.GetUserSelf)

		authRequired.GET("/api/topics", topic.GetTopics)
		authRequired.GET("/api/topics/:id", topic.GetTopic)
		authRequired.GET("/api/topics/join/:id", topic.JoinTopic)
		//authRequired.GET("/api/topics/massages/:topic_id", topic.TopicMassages)
	}

	{
		authRequiredAdmin.PUT("/api/users/status", user.ChangeStatus)
		authRequiredAdmin.DELETE("api/users/:id", user.DeleteUser)
		authRequiredAdmin.PUT("api/users/permissions", user.SetPerm)
		authRequiredAdmin.GET("/api/users/all", user.GetUsers)
		authRequiredAdmin.GET("/api/users/:id", user.GetUser)

		authRequiredAdmin.DELETE("/api/messages/users_message/:id", massages.Delete)

		authRequiredAdmin.POST("/api/topics", topic.CreateTopic)
		authRequiredAdmin.DELETE("/api/topics/:id", topic.DeleteTopic)
		authRequiredAdmin.PUT("/api/topics/:id", topic.ChangeTopic)
		authRequiredAdmin.DELETE("/api/topics/:id/user/:user_id", topic.DeleteUser)
	}

	return r
}
