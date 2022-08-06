package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get email/pass off request body
	var user models.User

	if c.BindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read users from body",
		})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
	}

	// create the user
	user.Password = string(hash)

	stmt, err := mysql.Prepare(
		`
		INSERT INTO users 
		SET 
			userID=?,
			name=?,
			familyName=?,
			password=?,
			email=?
		;
		`,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Could not prepare user insert statement",
			"detail": err.Error(),
		})
		return
	}

	_, err = stmt.Exec(user.UserID,
		user.Name,
		user.FamilyName,
		user.Password,
		user.Email,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Could not insert new user in database",
			"detail": err.Error(),
		})
		return
	}

	// respond
	c.IndentedJSON(http.StatusOK, "")
}

func Login(c *gin.Context) {
	// get email/pass off request body

	// lookup requested user

	// compare sent in pass with saved user pass hash

	//generate a jwt token

	// send it back
}
