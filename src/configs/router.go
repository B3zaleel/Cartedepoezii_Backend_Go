package configs

import (
	"github.com/gin-gonic/gin"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/controllers"
)

/**
 * Adds all required endpoints to the given engine.
 */
func AddEndpoints(ginEngine *gin.Engine) {
	ginEngine.GET("/", controllers.GetHome)
	ginEngine.GET("/api", controllers.GetHome)
	ginEngine.StaticFile("/favicon", "src/static/Logo.png")
	ginEngine.StaticFile("/favicon.ico", "src/static/Logo.png")
	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/", controllers.GetHome)
		// v1.GET("/profile-photo", controllers.GetProfilePhoto)

		v1.POST("/sign-in", controllers.SignIn)
		// v1.POST("/sign-up", controllers.SignUp)
		// v1.POST("/reset-password", controllers.RequestResetPassword)
		// v1.PUT("/reset-password", controllers.ResetPassword)

		v1.GET("/comment", controllers.GetComment)
		// v1.POST("/comment", controllers.AddComment)
		// v1.DELETE("/comment", controllers.RemoveComment)
		// v1.GET("/comments-of-poem", controllers.GetPoemComments)
		// v1.GET("/comment-replies", controllers.GetRepliesToComment)
		// v1.GET("/comments-by-user", controllers.GetUserComments)

		// v1.GET("/followers", controllers.GetFollowers)
		// v1.GET("/followings", controllers.GetFollowings)
		v1.PUT("/follow", controllers.ChangeConnection)

		v1.GET("/poem", controllers.GetPoem)
		// v1.POST("/poem", controllers.AddPoem)
		// v1.PUT("/poem", controllers.UpdatePoem)
		// v1.DELETE("/poem", controllers.RemovePoem)
		// v1.PUT("/like-poem", controllers.ChangePoemReaction)
		// v1.GET("/poems-user-created", controllers.GetPoemsUserCreated)
		// v1.GET("/poems-user-likes", controllers.GetPoemsUserLikes)
		// v1.GET("/poems-channel", controllers.GetPoemsForChannel)
		// v1.GET("/poems-explore", controllers.GetPoemsToExplore)

		// v1.GET("/search-poems", controllers.FindPoems)
		v1.GET("/search-people", controllers.FindPeople)

		v1.GET("/user", controllers.GetUser)
		// v1.PUT("/user", controllers.UpdateUser)
		// v1.DELETE("/user", controllers.RemoveUser)
	}
	ginEngine.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{"success": false, "message": "Page not found."})
	})
	ginEngine.NoMethod(func(c *gin.Context) {
		c.JSON(200, gin.H{"success": false, "message": "Method not found."})
	})
}
