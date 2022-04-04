package controllers

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/db_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/response_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/utils"
	"github.com/gin-gonic/gin"
)

// Retrieves a list of poems that match a given query.
func FindPoems(c *gin.Context) {
	query := strings.ReplaceAll(strings.Trim(c.Query("q"), " "), "'", "")
	if len(query) < 3 {
		c.JSON(200, gin.H{"success": false, "message": "Query is too short."})
		return
	}
	if strings.Count(query, "\"")%2 != 0 {
		c.JSON(200, gin.H{"success": false, "message": "Unequal number of quotes."})
		return
	}
	pageSpec, err := utils.GetPageSpec(c, true)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	poems := []db_models.Poem{}
	err = db.Select(
		&poems,
		`SELECT * FROM poems WHERE
		to_tsvector(title || ' ' || text) @@ websearch_to_tsquery('english', $1)
		;`,
		query,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	pagePoems, err := utils.ExtractPage(poems, *pageSpec)
	pagePoemsObjs := make([]response_models.Poem, len(pagePoems))
	for i, pagePoem := range pagePoems {
		user := &db_models.User{}
		err = db.Get(
			user,
			"SELECT * FROM users WHERE id=$1;",
			pagePoem.UserId,
		)
		if err != nil {
			continue
		}
		repliesCount := 0
		err = db.Select(
			&repliesCount,
			"SELECT COUNT(*) AS replies_count FROM comments WHERE poem_id=$1 AND comment_id=$2",
			pagePoem.Id,
			"",
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		likesCount := 0
		err = db.Select(
			&likesCount,
			"SELECT COUNT(*) AS likes_count FROM poems_likes WHERE poem_id=$1;",
			pagePoems[i].Id,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		isLiked, isFollowingUser := false, false
		if authToken != nil {
			poemLike := &db_models.PoemLike{}
			err = db.Get(
				poemLike,
				"SELECT * FROM poems_likes WHERE poem_id=$1 AND user_id=$2;",
				pagePoems[i].Id,
				authToken.UserId,
			)
			isLiked = err == nil
			userFollowing := &db_models.UserFollowing{}
			err = db.Get(
				userFollowing,
				"SELECT * FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
				authToken.UserId,
				pagePoem.UserId,
			)
			isFollowingUser = err == nil
		}
		verses := []string{}
		err = json.Unmarshal([]byte(pagePoems[i].Text), &verses)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pagePoemsObjs[i] = response_models.Poem{
			Id: pagePoems[i].Id,
			User: response_models.UserMin{
				Id:             user.Id,
				Name:           user.Name,
				ProfilePhotoId: user.ProfilePhotoId,
				IsFollowing:    isFollowingUser,
			},
			Title:         pagePoems[i].Title,
			PublishedOn:   pagePoems[i].CreatedOn.Format(time.RFC3339),
			Verses:        verses,
			CommentsCount: repliesCount,
			LikesCount:    likesCount,
			IsLiked:       isLiked,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pagePoemsObjs,
		},
	)
}

// Retrieves a list of users that match a given query.
func FindPeople(c *gin.Context) {
	query := strings.ReplaceAll(strings.Trim(c.Query("q"), " "), "'", "")
	if len(query) < 3 {
		c.JSON(200, gin.H{"success": false, "message": "Query is too short."})
		return
	}
	if strings.Count(query, "\"")%2 != 0 {
		c.JSON(200, gin.H{"success": false, "message": "Unequal number of quotes."})
		return
	}
	pageSpec, err := utils.GetPageSpec(c, true)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	users := []db_models.User{}
	err = db.Select(
		&users,
		`SELECT * FROM users WHERE
		to_tsvector(name || ' ' || bio) @@ websearch_to_tsquery('english', $1)
		;`,
		query,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	pageUsers, err := utils.ExtractPage(users, *pageSpec)
	pageUsersObjs := make([]response_models.UserMin, len(pageUsers))
	for i, pageUser := range pageUsers {
		isFollowingUser := false
		if authToken != nil {
			userFollowing := &db_models.UserFollowing{}
			err = db.Get(
				userFollowing,
				"SELECT * FROM users_followings WHERE follower_id=$1 AND following_id=$2;",
				authToken.UserId,
				pageUser.Id,
			)
			isFollowingUser = err == nil
		}
		pageUsersObjs[i] = response_models.UserMin{
			Id:             pageUser.Id,
			Name:           pageUser.Name,
			ProfilePhotoId: pageUser.ProfilePhotoId,
			IsFollowing:    isFollowingUser,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pageUsersObjs,
		},
	)
}
