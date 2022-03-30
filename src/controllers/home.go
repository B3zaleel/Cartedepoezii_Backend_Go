package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
	imagekit "github.com/B3zaleel/imagekit-go"
)

// Retrieves the welcome page.
func GetHome(c *gin.Context) {
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"message": "Welcome to the Cartedepoezii API.",
			},
		},
	)
}

// Retrieves a user's profile photo.
func GetProfilePhoto(c *gin.Context) {
	imgKit := imagekit.ImageKit{
		PrivateKey: os.Getenv("IMG_CDN_PRI_KEY"),
		PublicKey: os.Getenv("IMG_CDN_PUB_KEY"),
		UrlEndpoint: os.Getenv("IMG_CDN_URL_EPT"),
	}
	imgId := c.DefaultQuery("imgId", "")
	fileDetails, err := imgKit.GetFileDetails(imgId)
	if err != nil {
		c.JSON(
			200,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)
	} else {
		c.JSON(
			200,
			gin.H{
				"success": true,
				"data": gin.H{
					"url": fileDetails.Url,
				},
			},
		)
	}
}
