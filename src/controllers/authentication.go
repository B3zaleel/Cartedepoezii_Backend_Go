package controllers

import (
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/db_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/request_models"
	"github.com/B3zaleel/Cartedepoezii_Backend_Go/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Signs in a user.
func SignIn(c *gin.Context) {
	var jsonBody request_models.SignInForm
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if _, err := mail.ParseAddress(jsonBody.Email); err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Invalid email."})
		return
	}
	db, err := utils.GetDBConnection()
	user := &db_models.User{}
	err = db.Get(
		user,
		"SELECT * FROM users WHERE email=$1;",
		jsonBody.Email,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Invalid email and/or password."})
		return
	}
	maxSignInAttempts, err := strconv.Atoi(os.Getenv("APP_MAX_SIGNIN_TRIES"))
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	isValid, err := utils.IsValidPassword(jsonBody.Password, user.PasswordHash)
	if err != nil || !isValid {
		if !isValid {
			tx, err := db.Begin()
			if err != nil {
				c.JSON(200, gin.H{"success": false, "message": err.Error()})
				return
			}
			if user.SignInAttempts+1 == maxSignInAttempts {
				// make account inactive
				_, err = tx.Exec(
					`UPDATE users SET is_active=$1 WHERE id=$2;`,
					false,
					user.Id,
				)
				if err != nil {
					c.JSON(200, gin.H{"success": false, "message": err.Error()})
					return
				}
				// TODO: send an account locked message
			} else {
				// increase attempts
				_, err = tx.Exec(
					`UPDATE users SET sign_in_attempts=$1 WHERE id=$2;`,
					user.SignInAttempts+1,
					user.Id,
				)
				if err != nil {
					c.JSON(200, gin.H{"success": false, "message": err.Error()})
					return
				}
			}
			err = tx.Commit()
			if err != nil {
				c.JSON(200, gin.H{"success": false, "message": err.Error()})
				return
			}
		}
		c.JSON(200, gin.H{"success": false, "message": "Invalid email and/or password."})
		return
	}
	if !user.IsActive {
		c.JSON(200, gin.H{"success": false, "message": "This account is locked."})
		return
	}
	if user.SignInAttempts > 1 || len(user.AccountResetToken) > 0 {
		// reset SignInAttempts to 1 and empty AccountResetToken
		tx, err := db.Begin()
		if err != nil {
			c.JSON(200, gin.H{"success": false, "message": err.Error()})
			return
		}
		_, err = tx.Exec(
			`UPDATE users SET sign_in_attempts=$1, account_reset_token=$2 WHERE id=$3;`,
			1,
			"",
			user.Id,
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
	}
	authToken := &utils.AuthToken{
		UserId:     user.Id,
		Email:      jsonBody.Email,
		SecureText: user.PasswordHash,
		Expires:    time.Now().UTC(),
	}
	token, err := utils.EncodeAuthToken(authToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"userId":    user.Id,
				"name":      user.Name,
				"authToken": token,
			},
		},
	)
}

// Creates a new user.
func SignUp(c *gin.Context) {
	var jsonBody request_models.SignUpForm
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(jsonBody.Name) > 64 {
		c.JSON(200, gin.H{"success": false, "message": "Name is too long."})
		return
	}
	if len(jsonBody.Password) < 8 {
		c.JSON(200, gin.H{"success": false, "message": "Password is too short."})
		return
	}
	if _, err := mail.ParseAddress(jsonBody.Email); err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Invalid email."})
		return
	}
	currentTime := time.Now().UTC()
	userId := uuid.New().String()
	pwdHash, err := utils.GenerateHash(jsonBody.Password)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`INSERT INTO users(
			id, created_on, updated_on, email, name, password_hash
		)
		VALUES($1, $2, $3, $4, $5, $6);`,
		userId,
		currentTime.Format(time.RFC3339),
		currentTime.Format(time.RFC3339),
		jsonBody.Email,
		jsonBody.Name,
		pwdHash,
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
	authToken := &utils.AuthToken{
		UserId:     userId,
		Email:      jsonBody.Email,
		SecureText: pwdHash,
		Expires:    currentTime,
	}
	token, err := utils.EncodeAuthToken(authToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"userId":    userId,
				"name":      jsonBody.Name,
				"authToken": token,
			},
		},
	)
	// TODO: Send a welcome message
}

// TODO: Creates a password reset token for a user.
func RequestResetPassword(c *gin.Context) {
	var jsonBody request_models.PasswordResetRequestForm
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if _, err := mail.ParseAddress(jsonBody.Email); err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Invalid email."})
		return
	}
	db, err := utils.GetDBConnection()
	user := &db_models.User{}
	err = db.Get(
		user,
		"SELECT * FROM users WHERE email=$1;",
		jsonBody.Email,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to create reset token."})
		return
	}
	resetToken := &utils.ResetToken{
		UserId:  user.Id,
		Email:   user.Email,
		Message: "password_reset",
		Expires: time.Now().UTC(),
	}
	resetTokenStr, err := utils.EncodeResetToken(resetToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`UPDATE users SET account_reset_token=$1 WHERE id=$2;`,
		resetTokenStr,
		user.Id,
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
	// TODO: Send a password reset message
}

// Resets a user's password
func ResetPassword(c *gin.Context) {
	var jsonBody request_models.PasswordResetForm
	err := c.ShouldBindJSON(&jsonBody)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(jsonBody.Password) < 8 {
		c.JSON(200, gin.H{"success": false, "message": "Password is too short."})
		return
	}
	if _, err := mail.ParseAddress(jsonBody.Email); err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Invalid email."})
		return
	}
	resetToken, err := utils.DecodeResetToken(jsonBody.ResetToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	if resetToken.Message != "password_reset" {
		c.JSON(200, gin.H{"success": false, "message": "Invalid reset token."})
		return
	}
	if jsonBody.Email != resetToken.Email {
		c.JSON(
			200,
			gin.H{
				"success": false,
				"message": "User email and reset token are a mismatch.",
			},
		)
		return
	}
	currentTime := time.Now().UTC()
	db, err := utils.GetDBConnection()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	user := &db_models.User{}
	err = db.Get(
		user,
		"SELECT * FROM users WHERE id=$1;",
		resetToken.UserId,
	)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": "Failed to find user with id."})
		return
	}
	if len(user.AccountResetToken) == 0 {
		c.JSON(200, gin.H{"success": false, "message": "Invalid reset token."})
		return
	}
	pwdHash, err := utils.GenerateHash(jsonBody.Password)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	tx, err := db.Begin()
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	_, err = tx.Exec(
		`UPDATE users SET account_reset_token=$1, is_active=$2
		password_hash=$3, sign_in_attempts=$4, updated_on=$5
		WHERE user_id=$6;`,
		"",
		true,
		pwdHash,
		1,
		currentTime.Format(time.RFC3339),
		resetToken.UserId,
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
	authToken := &utils.AuthToken{
		UserId:     user.Id,
		Email:      jsonBody.Email,
		SecureText: pwdHash,
		Expires:    currentTime,
	}
	token, err := utils.EncodeAuthToken(authToken)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(
		200,
		gin.H{
			"success": true,
			"data": gin.H{
				"userId":    user.Id,
				"name":      user.Name,
				"authToken": token,
			},
		},
	)
	// TODO: Send a password changed message
}
