package controllers

import (
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/db_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/request_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/utils"
	imagekit "github.com/B3zaleel/imagekit-go"
	"github.com/gin-gonic/gin"
)

// Retrieves information about a given user.
func GetUser(c *gin.Context) {
	userId := c.Query("id")
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	user := &db_models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE id=$1;", userId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find user."})
		return
	}
	isFollowingUser := false
	userEmail := ""
	if authToken != nil {
		userFollowing := &db_models.UserFollowing{}
		err = db.Get(
			userFollowing,
			"SELECT * FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
			authToken.UserId,
			userId,
		)
		isFollowingUser = err == nil
		if authToken.UserId == userId {
			userEmail = user.Email
		}
	}
	poemsCount := 0
	db.Select(
		&poemsCount,
		"SELECT COUNT(*) AS poems_count FROM poems WHERE user_id=$1;",
		userId,
	)
	poemLikesCount := 0
	db.Select(
		&poemLikesCount,
		"SELECT COUNT(*) AS poems_likes_count FROM poems_likes WHERE user_id=$1;",
		userId,
	)
	commentsCount := 0
	db.Select(
		&commentsCount,
		"SELECT COUNT(*) AS comments_count FROM comments WHERE user_id=$1;",
		userId,
	)
	followersCount := 0
	db.Select(
		&followersCount,
		`SELECT COUNT(*) AS followers_count FROM users_followings
		WHERE following_id=$1;`,
		userId,
	)
	followingsCount := 0
	db.Select(
		&followingsCount,
		`SELECT COUNT(*) AS followings_count FROM users_followings
		WHERE follower_id=$1;`,
		userId,
	)
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"id":              user.Id,
				"joined":          user.CreatedOn.UTC().Format(time.RFC3339),
				"name":            user.Name,
				"email":           userEmail,
				"bio":             user.Bio,
				"profilePhotoId":  user.ProfilePhotoId,
				"followersCount":  followersCount,
				"followingsCount": followingsCount,
				"poemsCount":      poemsCount,
				"likesCount":      poemLikesCount,
				"commentsCount":   commentsCount,
				"isFollowing":     isFollowingUser,
			},
		},
	)
}

// Updates a user's information.
func UpdateUser(c *gin.Context) {
	var jsonBody request_models.UserUpdateForm
	currentTime := time.Now().UTC()
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(jsonBody.AuthToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if jsonBody.UserId != authToken.UserId {
		c.JSON(
			200,
			gin.H{
				"success": false,
				"message": "User id and auth token are a mismatch.",
			},
		)
		return
	}
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(jsonBody.Name) > 64 {
		c.JSON(200, gin.H{"success": false, "message": "Name is too long."})
		return
	}
	if len(jsonBody.Bio) > 384 {
		c.JSON(200, gin.H{"success": false, "message": "Bio is too long."})
		return
	}
	if _, err := mail.ParseAddress(jsonBody.Email); err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Invalid email."})
		return
	}
	db, err := utils.GetDBConnection()
	user := &db_models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE id=$1;", jsonBody.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find user with id."})
		return
	}
	// check and update user's profile photo
	imgKit := &imagekit.ImageKit{
		PublicKey:   os.Getenv("IMG_CDN_PUB_KEY"),
		PrivateKey:  os.Getenv("IMG_CDN_PRI_KEY"),
		UrlEndpoint: os.Getenv("IMG_CDN_URL_EPT"),
	}
	profilePhotoId := user.ProfilePhotoId
	if jsonBody.RemoveProfilePhoto {
		if len(profilePhotoId) > 0 {
			err = imgKit.DeleteFile(profilePhotoId)
			if err != nil {
				c.JSON(200, gin.H{"success": false, "message": err.Error()})
				return
			}
			profilePhotoId = ""
		}
	} else if len(strings.Trim(jsonBody.ProfilePhoto, " ")) > 0 {
		if len(profilePhotoId) > 0 {
			err = imgKit.DeleteFile(profilePhotoId)
			if err != nil {
				c.JSON(200, gin.H{"success": false, "message": err.Error()})
				return
			}
		}
		destFolder := imagekit.String("cartedepoezii/profile_photos/")
		isPrivateFile := imagekit.Bool(false)
		fileDetails, err := imgKit.Upload(
			jsonBody.ProfilePhoto,
			strings.ReplaceAll(user.Id, "-", ""),
			&imagekit.FileOptions{
				Folder:        &destFolder,
				IsPrivateFile: &isPrivateFile,
			},
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		} else {
			profilePhotoId = string(*fileDetails.FileId)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`UPDATE users
		SET name=$1, email=$2, bio=$3, profile_photo_id=$4, updated_on=$5
		WHERE id=$6;`,
		jsonBody.Name,
		jsonBody.Email,
		jsonBody.Bio,
		profilePhotoId,
		currentTime.UTC().Format(time.RFC3339),
		jsonBody.UserId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	newAuthToken := &utils.AuthToken{
		UserId:     user.Id,
		Email:      jsonBody.Email,
		SecureText: user.PasswordHash,
		Expires:    time.Now().UTC(),
	}
	token, err := utils.EncodeAuthToken(newAuthToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"authToken":      token,
				"profilePhotoId": profilePhotoId,
			},
		},
	)
}

// Deletes a user's account.
func RemoveUser(c *gin.Context) {
	var jsonBody request_models.UserDeleteForm
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(jsonBody.AuthToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if jsonBody.UserId != authToken.UserId {
		c.JSON(
			200,
			gin.H{
				"success": false,
				"message": "User id and auth token are a mismatch.",
			},
		)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	// remove user's poem likes
	_, err = tx.Exec(
		"DELETE FROM poems_likes WHERE user_id=$1;",
		jsonBody.UserId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	// remove user's connections
	_, err = tx.Exec(
		"DELETE FROM users_followings WHERE follower_id=$1 OR following_id=$1;",
		jsonBody.UserId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	// remove replies to user's comments
	_, err = tx.Exec(
		`DELETE FROM comments WHERE comment_id IN (
			SELECT * FROM comments WHERE user_id=$1 AND comment_id=''
		);`,
		jsonBody.UserId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	// remove user's comments
	_, err = tx.Exec("DELETE FROM comments WHERE user_id=$1;", jsonBody.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec("DELETE FROM poems WHERE user_id=$1;", jsonBody.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	// remove user's record
	_, err = tx.Exec("DELETE FROM users WHERE id=$1;", jsonBody.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    gin.H{},
		},
	)
}
