package controllers

import (
	"encoding/json"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/db_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/request_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/response_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// The maximum number of poems for the explore page.
	MAX_EXPLORE_PAGE_POEMS = 255
)

// Retrieves information about a given poem.
func GetPoem(c *gin.Context) {
	poemId := c.Query("id")
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	poem := &db_models.Poem{}
	err = db.Get(poem, "SELECT * FROM poems WHERE id=$1;", poemId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find poem."})
		return
	}
	user := &db_models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE id=$1;", poem.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
		return
	}
	comments := []db_models.Comment{}
	err = db.Select(
		&comments,
		"SELECT * FROM comments WHERE poem_id=$1 AND comment_id=$2",
		poemId,
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
		poemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	poemLike := &db_models.PoemLike{}
	err = db.Get(
		poemLike,
		"SELECT * FROM poems_likes WHERE poem_id=$1 AND user_id=$2;",
		poemId,
		authToken.UserId,
	)
	isLiked := err != nil
	verses := []string{}
	err = json.Unmarshal([]byte(poem.Text), &verses)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"id": poem.Id,
				"user": gin.H{
					"id":             user.Id,
					"name":           user.Name,
					"profilePhotoId": user.ProfilePhotoId,
				},
				"title":         poem.Title,
				"publishedOn":   poem.CreatedOn.UTC().Format(time.RFC3339),
				"verses":        verses,
				"commentsCount": len(comments),
				"likesCount":    likesCount,
				"isLiked":       isLiked,
			},
		},
	)
}

// Creates a new poem.
func AddPoem(c *gin.Context) {
	var jsonBody request_models.PoemAddForm
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
		c.JSON(200, gin.H{"success": false, "message": "User id and auth token are a mismatch."})
		return
	}
	currentTime := time.Now().UTC()
	versesTxt, err := json.Marshal(jsonBody.Verses)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(jsonBody.Title) > 256 {
		c.JSON(200, gin.H{"success": false, "message": "Title is too long."})
		return
	}
	if len(jsonBody.Verses) < 1 {
		c.JSON(200, gin.H{"success": false, "message": "At least 1 verse is needed."})
		return
	}
	for i := 0; i < len(jsonBody.Verses); i++ {
		if len(strings.Trim(jsonBody.Verses[i], " ")) < 1 {
			c.JSON(200, gin.H{"success": false, "message": "Some verses are too short."})
			return
		}
	}
	poemId := uuid.New().String()
	poem := &db_models.Poem{
		Id:        poemId,
		CreatedOn: currentTime,
		UpdatedOn: currentTime,
		UserId:    authToken.UserId,
		Title:     jsonBody.Title,
		Text:      string(versesTxt),
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`INSERT INTO poems (id, created_on, updated_on, user_id, title, text)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		poem.Id,
		poem.CreatedOn,
		poem.UpdatedOn,
		poem.UserId,
		poem.Title,
		poem.Text,
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
			"data": gin.H{
				"id":           poemId,
				"createdOn":    currentTime.UTC().Format(time.RFC3339),
				"repliesCount": 0,
				"likesCount":   0,
			},
		},
	)
}

// Edits an existing poem.
func UpdatePoem(c *gin.Context) {
	var jsonBody request_models.PoemUpdateForm
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
		c.JSON(200, gin.H{"success": false, "message": "User id and auth token are a mismatch."})
		return
	}
	currentTime := time.Now().UTC()
	poem := &db_models.Poem{}
	err = db.Get(poem, "SELECT * FROM poems WHERE id=$1;", jsonBody.PoemId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find poem."})
		return
	}
	if poem.UserId != jsonBody.UserId {
		c.JSON(200, gin.H{"success": false, "message": "You are not allowed to edit this poem."})
		return
	}
	versesTxt, err := json.Marshal(jsonBody.Verses)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(jsonBody.Title) > 256 {
		c.JSON(200, gin.H{"success": false, "message": "Title is too long."})
		return
	}
	if len(jsonBody.Verses) < 1 {
		c.JSON(200, gin.H{"success": false, "message": "At least 1 verse is needed."})
		return
	}
	for i := 0; i < len(jsonBody.Verses); i++ {
		if len(strings.Trim(jsonBody.Verses[i], " ")) < 1 {
			c.JSON(200, gin.H{"success": false, "message": "Some verses are too short."})
			return
		}
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`UPDATE poems SET
			updated_on=$1, title=$2, text=$3
			WHERE id=$4;`,
		currentTime,
		jsonBody.Title,
		string(versesTxt),
		jsonBody.PoemId,
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
			"data":    gin.H{},
		},
	)
}

// Deletes a poem.
func RemovePoem(c *gin.Context) {
	var jsonBody request_models.PoemDeleteForm
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
		c.JSON(200, gin.H{"success": false, "message": "User id and auth token are a mismatch."})
		return
	}
	poem := &db_models.Poem{}
	err = db.Get(
		poem,
		"SELECT * FROM poem WHERE id=$1;",
		jsonBody.PoemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Poem doesn't exist."})
		return
	} else if poem.UserId != authToken.UserId {
		c.JSON(200, gin.H{
			"success": false,
			"message": "Only the author of the poem can delete the poem.",
		})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		"DELETE FROM poems_likes WHERE poem_id=$1;",
		jsonBody.PoemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		"DELETE FROM comments WHERE poem_id=$1;",
		jsonBody.PoemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		"DELETE FROM poems WHERE id=$1;",
		jsonBody.PoemId,
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
			"data":    gin.H{},
		},
	)
}

// Toggles a user's reaction on a poem.
func ChangePoemReaction(c *gin.Context) {
	var jsonBody request_models.PoemLikeForm
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
	if jsonBody.UserId != authToken.UserId {
		c.JSON(200, gin.H{"success": false, "message": "User id and auth token are a mismatch."})
		return
	}
	poem := &db_models.Poem{}
	err = db.Get(
		poem,
		"SELECT * FROM poem WHERE id=$1;",
		jsonBody.PoemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Poem doesn't exist."})
		return
	}
	poemLike := &db_models.PoemLike{}
	err = db.Get(
		poemLike,
		"SELECT * FROM poems_likes WHERE user_id=$1 AND poem_id=$2;",
		jsonBody.UserId,
		jsonBody.PoemId,
	)
	if err == nil {
		// poemLike exists -> remove like
		tx, err := db.Begin()
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		_, err = tx.Exec(
			"DELETE FROM poems_likes WHERE user_id=$1 AND poem_id=$2;",
			jsonBody.UserId,
			jsonBody.PoemId,
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
		// poemLike doesn't exist -> create like
		currentTime := time.Now().UTC()
		connectionId := uuid.New().String()
		newPoemLike := &db_models.PoemLike{
			Id:        connectionId,
			UserId:    jsonBody.UserId,
			PoemId:    jsonBody.PoemId,
			CreatedOn: currentTime,
		}
		tx, err := db.Begin()
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		_, err = tx.Exec(
			`INSERT INTO poems_likes (id, user_id, poem_id, created_on)
			VALUES ($1, $2, $3, $4);`,
			newPoemLike.Id,
			newPoemLike.UserId,
			newPoemLike.PoemId,
			newPoemLike.CreatedOn,
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

// Retrieves poems created by the current user.
func GetPoemsUserCreated(c *gin.Context) {
	userId := c.Query("id")
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
	err = db.Select(&poems, "SELECT * FROM poems WHERE user_id=$1;", userId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(poems, func(i, j int) bool {
		return poems[i].CreatedOn.After(poems[j].CreatedOn)
	})
	pagePoems, err := utils.ExtractPage(poems, *pageSpec)
	pagePoemsObjs := make([]response_models.Poem, len(pagePoems))
	user := &db_models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE id=$1;", userId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
		return
	}
	for i := 0; i < len(pagePoems); i++ {
		repliesCount := 0
		err = db.Select(
			&repliesCount,
			"SELECT COUNT(*) AS replies_count FROM comments WHERE poem_id=$1 AND comment_id=$2",
			pagePoems[i].Id,
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
		isLiked := false
		if authToken != nil {
			poemLike := &db_models.PoemLike{}
			err = db.Get(
				poemLike,
				"SELECT * FROM poems_likes WHERE poem_id=$1 AND user_id=$2;",
				pagePoems[i].Id,
				authToken.UserId,
			)
			isLiked = err == nil
		}
		verses := []string{}
		err = json.Unmarshal([]byte(pagePoems[i].Text), verses)
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

// Retrieves poems liked by a given user.
func GetPoemsUserLikes(c *gin.Context) {
	userId := c.Query("id")
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
	poemLikes := []db_models.PoemLike{}
	err = db.Select(
		&poemLikes,
		"SELECT * FROM poems_likes WHERE user_id=$1;",
		userId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(poemLikes, func(i, j int) bool {
		return poemLikes[i].CreatedOn.After(poemLikes[j].CreatedOn)
	})
	pagePoemLikes, err := utils.ExtractPage(poemLikes, *pageSpec)
	pagePoemsObjs := make([]response_models.Poem, len(pagePoemLikes))
	for i := 0; i < len(pagePoemLikes); i++ {
		poem := &db_models.Poem{}
		err = db.Get(poem, "SELECT * FROM poems WHERE id=$1;", poem.Id)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find poem."})
			return
		}
		user := &db_models.User{}
		err = db.Get(user, "SELECT * FROM users WHERE id=$1;", userId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find poem creator."})
			return
		}
		repliesCount := 0
		err = db.Select(
			&repliesCount,
			"SELECT COUNT(*) AS replies_count FROM comments WHERE poem_id=$1 AND comment_id=$2",
			poem.Id,
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
			poem.Id,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		isLiked := false
		if authToken != nil {
			poemLike := &db_models.PoemLike{}
			err = db.Get(
				poemLike,
				"SELECT * FROM poems_likes WHERE poem_id=$1 AND user_id=$2;",
				poem.Id,
				authToken.UserId,
			)
			isLiked = err == nil
		}
		verses := []string{}
		err = json.Unmarshal([]byte(poem.Text), verses)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pagePoemsObjs[i] = response_models.Poem{
			Id: poem.Id,
			User: response_models.UserMin{
				Id:             user.Id,
				Name:           user.Name,
				ProfilePhotoId: user.ProfilePhotoId,
			},
			Title:         poem.Title,
			PublishedOn:   poem.CreatedOn.Format(time.RFC3339),
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

// Retrieves poems for a user's timeline or home section.
func GetPoemsForChannel(c *gin.Context) {
	pageSpec, err := utils.GetPageSpec(c, true)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	authToken, err := utils.DecodeAuthToken(c.Query("token"))
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	userFollowings := []db_models.UserFollowing{}
	err = db.Select(
		&userFollowings,
		"SELECT * FROM users_followings WHERE follower_id=$1;",
		authToken.UserId,
	)
	maxSize := math.MaxInt32
	usersCount := len(userFollowings) + 1
	poemsPerUser := int32(maxSize / usersCount)
	users := make(map[string]response_models.UserMin, usersCount)
	poems := []db_models.Poem{}
	user := &db_models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE id=$1;", authToken.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find user."})
		return
	}
	users[authToken.UserId] = response_models.UserMin{
		Id:             user.Id,
		Name:           user.Name,
		ProfilePhotoId: user.ProfilePhotoId,
		IsFollowing:    false,
	}
	_ = db.Select(
		&poems,
		"SELECT * FROM poems WHERE user_id=$1 LIMIT $2;",
		authToken.UserId,
		poemsPerUser,
	)
	for i := range userFollowings {
		userFollowingPoems := []db_models.Poem{}
		user := &db_models.User{}
		err = db.Get(
			user,
			"SELECT * FROM users WHERE id=$1;",
			userFollowings[i].FollowingId,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find user."})
			return
		}
		users[userFollowings[i].FollowingId] = response_models.UserMin{
			Id:             user.Id,
			Name:           user.Name,
			ProfilePhotoId: user.ProfilePhotoId,
			IsFollowing:    true,
		}
		err = db.Select(
			&userFollowingPoems,
			"SELECT * FROM poems WHERE user_id=$1 LIMIT $2;",
			userFollowings[i].FollowingId,
			poemsPerUser,
		)
		if err == nil {
			poems = append(poems, userFollowingPoems...)
		}
	}
	sort.SliceStable(poems, func(i, j int) bool {
		return poems[i].CreatedOn.After(poems[j].CreatedOn)
	})
	pagePoems, err := utils.ExtractPage(poems, *pageSpec)
	pagePoemsObjs := make([]response_models.Poem, len(pagePoems))
	for i, pagePoem := range pagePoems {
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
		isLiked := false
		if authToken != nil {
			poemLike := &db_models.PoemLike{}
			err = db.Get(
				poemLike,
				"SELECT * FROM poems_likes WHERE poem_id=$1 AND user_id=$2;",
				pagePoems[i].Id,
				authToken.UserId,
			)
			isLiked = err == nil
		}
		verses := []string{}
		err = json.Unmarshal([]byte(pagePoems[i].Text), verses)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pagePoemsObjs[i] = response_models.Poem{
			Id:            pagePoems[i].Id,
			User:          users[pagePoem.UserId],
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

// TODO: Retrieves poems a user can explore.
func GetPoemsToExplore(c *gin.Context) {
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
	userId := ""
	if authToken != nil {
		userId = authToken.UserId
	}
	poems := []db_models.Poem{}
	err = db.Select(
		&poems,
		`SELECT * FROM poems WHERE user_id NOT IN
			(SELECT following_id FROM users_followings WHERE follower_id=$1;)
			LIMIT $2;`,
		userId,
		MAX_EXPLORE_PAGE_POEMS,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(poems, func(i, j int) bool {
		return poems[i].CreatedOn.After(poems[j].CreatedOn)
	})
	pagePoems, err := utils.ExtractPage(poems, *pageSpec)
	pagePoemsObjs := make([]response_models.Poem, len(pagePoems))
	for i, pagePoem := range pagePoems {
		user := &db_models.User{}
		err = db.Get(user, "SELECT * FROM users WHERE id=$1;", pagePoem.UserId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find poem creator."})
			return
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
			pagePoem.Id,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		isLiked := false
		if authToken != nil {
			poemLike := &db_models.PoemLike{}
			err = db.Get(
				poemLike,
				"SELECT * FROM poems_likes WHERE poem_id=$1 AND user_id=$2;",
				pagePoem.Id,
				authToken.UserId,
			)
			isLiked = err == nil
		}
		verses := []string{}
		err = json.Unmarshal([]byte(pagePoem.Text), verses)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pagePoemsObjs[i] = response_models.Poem{
			Id: pagePoem.Id,
			User: response_models.UserMin{
				Id:             user.Id,
				Name:           user.Name,
				ProfilePhotoId: user.ProfilePhotoId,
			},
			Title:         pagePoem.Title,
			PublishedOn:   pagePoem.CreatedOn.Format(time.RFC3339),
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
