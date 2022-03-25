package configs

import (
	"github.com/gin-gonic/gin"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/controllers"
)

/**
 * Adds all required endpoints to the given engine.
 */
func AddEndpoints(ginEngine *gin.Engine) {
	authController := controllers.Authentication{}
	commentController := controllers.Comment{}
	connectionController := controllers.Connection{}
	homeController := controllers.Home{}
	poemController := controllers.Poem{}
	searchController := controllers.Search{}
	userController := controllers.User{}
	ginEngine.GET("/", homeController.GetHome())
	ginEngine.GET("/api", homeController.GetHome())
	ginEngine.StaticFile("/favicon", "src/static/Logo.png")
	ginEngine.StaticFile("/favicon.ico", "src/static/Logo.png")
	v1 := ginEngine.Group("/api/v1")
	{
		v1.GET("/", homeController.GetHome())
		// v1.GET("/profile-photo", homeController.GetProfilePhoto())

		v1.POST("/sign-in", authController.SignIn())
		// v1.POST("/sign-up", authController.SignUp())
		// v1.POST("/reset-password", authController.RequestResetPassword())
		// v1.PUT("/reset-password", authController.ResetPassword())

		v1.GET("/comment", commentController.GetComment())
		// v1.POST("/comment", commentController.AddComment())
		// v1.DELETE("/comment", commentController.RemoveComment())
		// v1.GET("/comments-of-poem", commentController.GetPoemComments())
		// v1.GET("/comment-replies", commentController.GetRepliesToComment())
		// v1.GET("/comments-by-user", commentController.GetUserComments())

		// v1.GET("/followers", connectionController.GetFollowers())
		// v1.GET("/followings", connectionController.GetFollowings())
		v1.PUT("/follow", connectionController.ChangeConnection())

		v1.GET("/poem", poemController.GetPoem())
		// v1.POST("/poem", poemController.AddPoem())
		// v1.PUT("/poem", poemController.UpdatePoem())
		// v1.DELETE("/poem", poemController.RemovePoem())
		// v1.PUT("/like-poem", poemController.ChangePoemReaction())
		// v1.GET("/poems-user-created", poemController.GetPoemsUserCreated())
		// v1.GET("/poems-user-likes", poemController.GetPoemsUserLikes())
		// v1.GET("/poems-channel", poemController.GetPoemsForChannel())
		// v1.GET("/poems-explore", poemController.GetPoemsToExplore())

		// v1.GET("/search-poems", searchController.FindPoems())
		v1.GET("/search-people", searchController.FindPeople())

		v1.GET("/user", userController.GetUser())
		// v1.PUT("/user", userController.UpdateUser())
		// v1.DELETE("/user", userController.RemoveUser())
	}
	ginEngine.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{"success": false, "message": "Page not found."})
	})
	ginEngine.NoMethod(func(c *gin.Context) {
		c.JSON(200, gin.H{"success": false, "message": "Method not found."})
	})
}
