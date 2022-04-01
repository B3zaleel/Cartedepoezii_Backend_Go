package controllers

import (
	"sort"
	"time"

	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/db_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/request_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/response_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Retrieves a user's followers.
func GetFollowers(c *gin.Context) {
	userId := c.Query("id")
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	pageSpec, err := utils.GetPageSpec(c, true)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	user_followers := []db_models.UserFollowing{}
	err = db.Select(
		&user_followers,
		"SELECT * FROM users_followings WHERE following_id=$1;",
		userId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(user_followers, func(i, j int) bool {
		return user_followers[i].CreatedOn.After(user_followers[j].CreatedOn)
	})
	pageFollowers, err := utils.ExtractPage(user_followers, *pageSpec)
	pageFollowersObjs := make([]response_models.UserMin, len(pageFollowers))
	for i := 0; i < len(pageFollowers); i++ {
		user := &db_models.User{}
		err = db.Get(user, "SELECT * FROM users WHERE id=$1;", pageFollowers[i].FollowerId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find follower."})
			return
		}
		isFollowingUser := false
		if authToken != nil {
			user_following := &db_models.UserFollowing{}
			err = db.Get(
				user_following,
				"SELECT * FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
				authToken.UserId,
				user.Id,
			)
			if err == nil {
				if user_following.FollowerId == authToken.UserId && user_following.FollowingId == user.Id {
					isFollowingUser = true
				}
			}
		}
		pageFollowersObjs[i] = response_models.UserMin{
			Id:             user.Id,
			Name:           user.Name,
			ProfilePhotoId: user.ProfilePhotoId,
			IsFollowing:    isFollowingUser,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pageFollowersObjs,
		},
	)
}

// Retrieves users followed by a given user.
func GetFollowings(c *gin.Context) {
	userId := c.Query("id")
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	pageSpec, err := utils.GetPageSpec(c, true)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	user_followings := []db_models.UserFollowing{}
	err = db.Select(
		&user_followings,
		"SELECT * FROM users_followings WHERE follower_id=$1;",
		userId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(user_followings, func(i, j int) bool {
		return user_followings[i].CreatedOn.After(user_followings[j].CreatedOn)
	})
	pageFollowings, err := utils.ExtractPage(user_followings, *pageSpec)
	pageFollowingsObjs := make([]response_models.UserMin, len(pageFollowings))
	for i := 0; i < len(pageFollowings); i++ {
		user := &db_models.User{}
		err = db.Get(user, "SELECT * FROM users WHERE id=$1;", pageFollowings[i].FollowingId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find following."})
			return
		}
		isFollowingUser := false
		if authToken != nil {
			user_following := &db_models.UserFollowing{}
			err = db.Get(
				user_following,
				"SELECT * FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
				authToken.UserId,
				user.Id,
			)
			if err == nil {
				if user_following.FollowerId == authToken.UserId && user_following.FollowingId == user.Id {
					isFollowingUser = true
				}
			}
		}
		pageFollowingsObjs[i] = response_models.UserMin{
			Id:             user.Id,
			Name:           user.Name,
			ProfilePhotoId: user.ProfilePhotoId,
			IsFollowing:    isFollowingUser,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pageFollowingsObjs,
		},
	)
}

// Toggles the connection between two users.
func ChangeConnection(c *gin.Context) {
	var jsonBody request_models.ConnectionForm
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(jsonBody.AuthToken)
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if jsonBody.FollowId == authToken.UserId {
		c.JSON(200, gin.H{"success": false, "message": "You cannot follow yourself."})
		return
	}
	if jsonBody.UserId != authToken.UserId {
		c.JSON(200, gin.H{"success": false, "message": "User id and auth token are a mismatch."})
		return
	}
	user_following := db_models.UserFollowing{}
	err = db.Get(
		&user_following,
		"SELECT * FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
		jsonBody.UserId,
		jsonBody.FollowId,
	)
	if err != nil {
		// remove connection
		tx, err := db.Begin()
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		_, err = tx.Exec(
			"DELETE FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
			jsonBody.UserId,
			jsonBody.FollowId,
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
		c.JSON(
			200,
			gin.H{
				"success": true,
				"data":    gin.H{"status": false},
			},
		)
	} else {
		// create connection
		currentTime := time.Now().UTC()
		connectionId := uuid.New().String()
		user_following := &db_models.UserFollowing{
			Id:          connectionId,
			FollowerId:  jsonBody.UserId,
			FollowingId: jsonBody.FollowId,
			CreatedOn:   currentTime,
		}
		tx, err := db.Begin()
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		_, err = tx.Exec(
			`INSERT INTO users_followings (id, follower_id, following_id, created_on)
			VALUES ($1, $2, $3, $4);`,
			user_following.Id,
			user_following.FollowerId,
			user_following.FollowingId,
			user_following.CreatedOn,
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
		c.JSON(
			200,
			gin.H{
				"success": true,
				"data":    gin.H{"status": true},
			},
		)
	}
}
