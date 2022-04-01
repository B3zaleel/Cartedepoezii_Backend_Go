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

// Retrieves information about a given comment.
func GetComment(c *gin.Context) {
	commentId := c.Query("id")
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	comment := &db_models.Comment{}
	err = db.Get(comment, "SELECT * FROM comments WHERE id=$1;", commentId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find comment."})
		return
	}
	user := &db_models.User{}
	err = db.Get(user, "SELECT * FROM users WHERE id=$1;", comment.UserId)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
		return
	}
	replies := []db_models.Comment{}
	db.Select(&replies, "SELECT * FROM comments WHERE comment_id=$1;", commentId)
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"id": comment.Id,
				"user": gin.H{
					"id":             user.Id,
					"name":           user.Name,
					"profilePhotoId": user.ProfilePhotoId,
				},
				"createdOn":    comment.CreatedOn.UTC().Format(time.RFC3339),
				"text":         comment.Text,
				"poemId":       comment.PoemId,
				"repliesCount": len(replies),
				"replyTo":      comment.CommentId,
			},
		},
	)
}

// Creates a new comment.
func AddComment(c *gin.Context) {
	var jsonBody request_models.CommentAddForm
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
	commentId := uuid.New().String()
	comment := &db_models.Comment{
		Id:        commentId,
		UserId:    authToken.UserId,
		PoemId:    jsonBody.PoemId,
		CommentId: commentId,
		Text:      jsonBody.Text,
		CreatedOn: currentTime,
	}
	if len(jsonBody.Text) > 384 {
		c.JSON(200, gin.H{"success": false, "message": "Name is too long."})
		return
	}
	user := &db_models.User{}
	err = db.Get(
		user,
		"SELECT * FROM comments WHERE id=$1 AND poem_id=$2;",
		jsonBody.ReplyTo,
		jsonBody.PoemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
		return
	}
	poem := &db_models.Poem{}
	err = db.Get(
		poem,
		"SELECT * FROM poems WHERE id=$1;",
		jsonBody.PoemId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find poem with id."})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`INSERT INTO comments (id, user_id, poem_id, comment_id, text, created_on)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		comment.Id,
		comment.UserId,
		comment.PoemId,
		comment.CommentId,
		comment.Text,
		comment.CreatedOn,
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
				"id":           commentId,
				"createdOn":    comment.CreatedOn.UTC().Format(time.RFC3339),
				"poemId":       jsonBody.PoemId,
				"replyTo":      jsonBody.ReplyTo,
				"repliesCount": 0,
			},
		},
	)
}

// Deletes a comment made by a user.
func RemoveComment(c *gin.Context) {
	var jsonBody request_models.CommentDeleteForm
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
	comment := &db_models.Comment{}
	err = db.Get(
		comment,
		"SELECT * FROM comments WHERE id=$1;",
		jsonBody.CommentId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
		return
	} else if comment.UserId != authToken.UserId {
		c.JSON(200, gin.H{
			"success": false,
			"message": "Only the author of the comment can delete the comment.",
		})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		"DELETE FROM comments WHERE comment_id=$1 OR id=$1;",
		jsonBody.CommentId,
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

// Retrieves all comments made directly under a poem.
func GetPoemComments(c *gin.Context) {
	poemId := c.Query("id")
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
	sort.SliceStable(comments, func(i, j int) bool {
		return comments[i].CreatedOn.After(comments[j].CreatedOn)
	})
	pageComments, err := utils.ExtractPage(comments, *pageSpec)
	pageCommentObjs := make([]response_models.Comment, len(pageComments))
	for i := 0; i < len(pageComments); i++ {
		user := &db_models.User{}
		err = db.Get(user, "SELECT * FROM users WHERE id=$1;", pageComments[i].UserId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
			return
		}
		replies := []db_models.Comment{}
		err = db.Select(
			&replies,
			"SELECT * FROM comments WHERE poem_id=$1 AND comment_id=$2",
			poemId,
			pageComments[i].Id,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pageCommentObjs[i] = response_models.Comment{
			Id: pageComments[i].Id,
			User: response_models.UserMin{
				Id:             user.Id,
				Name:           user.Name,
				ProfilePhotoId: user.ProfilePhotoId,
			},
			CreatedOn:    pageComments[i].CreatedOn.Format(time.RFC3339),
			Text:         pageComments[i].Text,
			PoemId:       pageComments[i].PoemId,
			RepliesCount: len(replies),
			ReplyTo:      pageComments[i].CommentId,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pageCommentObjs,
		},
	)
}

// Retrieves all comments made directly under another comment.
func GetRepliesToComment(c *gin.Context) {
	commentId := c.Query("id")
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
	comment := db_models.Comment{}
	err = db.Select(
		&comment,
		"SELECT * FROM comments WHERE id=$1",
		commentId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	comments := []db_models.Comment{}
	err = db.Select(
		&comments,
		"SELECT * FROM comments WHERE comment_id=$1",
		commentId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(comments, func(i, j int) bool {
		return comments[i].CreatedOn.After(comments[j].CreatedOn)
	})
	pageReplies, err := utils.ExtractPage(comments, *pageSpec)
	pageRepliesObjs := make([]response_models.Comment, len(pageReplies))
	for i := 0; i < len(pageReplies); i++ {
		user := &db_models.User{}
		err = db.Get(user, "SELECT * FROM users WHERE id=$1;", pageReplies[i].UserId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
			return
		}
		replies := []db_models.Comment{}
		err = db.Select(
			&replies,
			"SELECT * FROM comments WHERE poem_id=$1 AND comment_id=$2",
			comment.PoemId,
			pageReplies[i].Id,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pageRepliesObjs[i] = response_models.Comment{
			Id: pageReplies[i].Id,
			User: response_models.UserMin{
				Id:             user.Id,
				Name:           user.Name,
				ProfilePhotoId: user.ProfilePhotoId,
			},
			CreatedOn:    pageReplies[i].CreatedOn.Format(time.RFC3339),
			Text:         pageReplies[i].Text,
			PoemId:       pageReplies[i].PoemId,
			RepliesCount: len(replies),
			ReplyTo:      pageReplies[i].CommentId,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pageRepliesObjs,
		},
	)
}

// Retrieves all comments made by a user.
func GetUserComments(c *gin.Context) {
	userId := c.Query("id")
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
	user := db_models.User{}
	err = db.Select(
		&user,
		"SELECT * FROM users WHERE id=$1",
		userId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	comments := []db_models.Comment{}
	err = db.Select(
		&comments,
		"SELECT * FROM comments WHERE user_id=$1",
		userId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	sort.SliceStable(comments, func(i, j int) bool {
		return comments[i].CreatedOn.After(comments[j].CreatedOn)
	})
	pageReplies, err := utils.ExtractPage(comments, *pageSpec)
	pageRepliesObjs := make([]response_models.Comment, len(pageReplies))
	for i := 0; i < len(pageReplies); i++ {
		poem := &db_models.Poem{}
		err = db.Get(poem, "SELECT * FROM poems WHERE id=$1;", pageReplies[i].PoemId)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": "Failed to find comment creator."})
			return
		}
		replies := []db_models.Comment{}
		err = db.Select(
			&replies,
			"SELECT * FROM comments WHERE poem_id=$1 AND comment_id=$2",
			poem.Id,
			pageReplies[i].Id,
		)
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		pageRepliesObjs[i] = response_models.Comment{
			Id: pageReplies[i].Id,
			User: response_models.UserMin{
				Id:             user.Id,
				Name:           user.Name,
				ProfilePhotoId: user.ProfilePhotoId,
			},
			CreatedOn:    pageReplies[i].CreatedOn.Format(time.RFC3339),
			Text:         pageReplies[i].Text,
			PoemId:       pageReplies[i].PoemId,
			RepliesCount: len(replies),
			ReplyTo:      pageReplies[i].CommentId,
		}
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data":    pageRepliesObjs,
		},
	)
}
