package router

import (
	"final-project/controllers"
	"final-project/middlewares"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.PUT("/:userId", middlewares.Authentication(), middlewares.ProfileAuthorization(), controllers.UserUpdate)
		userRouter.DELETE("/", middlewares.Authentication(), middlewares.ProfileAuthorization(), controllers.UserDelete)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", controllers.PhotoCreate)
		photoRouter.GET("/", controllers.PhotoGetAll)
		photoRouter.GET("/:photoId", controllers.PhotoGetByID)
		photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(), controllers.PhotoUpdate)
		photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(), controllers.PhotoDelete)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.POST("/", controllers.CommentCreate)
		commentRouter.GET("/", controllers.CommentList)
		commentRouter.GET("/:commentId", controllers.CommentByID)
		commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), controllers.CommentUpdate)
		commentRouter.DELETE("/:commentId", middlewares.CommentAuthorization(), controllers.CommentDelete)
	}

	socialmediasRouter := r.Group("/socialmedias")
	{
		socialmediasRouter.Use(middlewares.Authentication())
		socialmediasRouter.POST("/", controllers.SocialMediaCreate)
		socialmediasRouter.GET("/", controllers.SocialMediaList)
		socialmediasRouter.GET("/:socialMediaId", controllers.GetSocialMediaByID) 
		socialmediasRouter.PUT("/:socialMediaId", middlewares.SocialMediaAuthorization(), controllers.SocialMediaUpdate)
		socialmediasRouter.DELETE("/:socialMediaId", middlewares.SocialMediaAuthorization(), controllers.SocialMediaDelete)
		
	}

	return r
}
