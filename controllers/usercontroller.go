package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mateenbagheri/briefly/configs"
	"github.com/mateenbagheri/briefly/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get email/pass off request body
	var user models.User

	err := c.BindJSON(&user)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed to read users from body",
			"error":   string(err.Error()),
		})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed to hash password",
			"error":   string(err.Error()),
		})
		return
	}

	// create the user
	user.Password = string(hash)

	stmt, err := Mysql.Prepare(
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not prepare user insert statement",
			"error":   string(err.Error()),
		})
		return
	}

	_, err = stmt.Exec(
		user.UserID,
		user.Name,
		user.FamilyName,
		user.Password,
		user.Email,
	)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not insert new user in database",
			"error":   string(err.Error()),
		})
		return
	}

	// respond
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":http.StatusOK,
		"message": "The user has been added successfully",
	})
}

func Login(c *gin.Context) {
	// get email/pass off request body
	var body struct {
		Email    string `json:"Email" binding:"required"`
		Password string `json:"Password" binding:"required"`
	}

	var user models.User

	err := c.BindJSON(&body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not bind body to structure",
			"error":   string(err.Error()),
		})
	}

	// lookup requested user
	result, err := Mysql.Query(
		`SELECT * FROM users WHERE email=?`, body.Email,
	)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not create the SELECT statement",
			"error":   string(err.Error()),
		})
		return
	}

	if result.Next() {
		err = result.Scan(
			&user.UserID,
			&user.Name,
			&user.FamilyName,
			&user.Password,
			&user.Salt,
			&user.Email,
		)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "SELECT reusults differ from Struct",
				"error":   string(err.Error()),
			})
			return
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Invalid Email Address. No Such Email Exists",
			"error":   string(err.Error()),
		})
		return
	}

	// compare sent in pass with saved user pass hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Wrong Password. Try Again!",
		})
		return
	}

	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	confs, _ := configs.LoadConfig()

	tokenString, err := token.SignedString([]byte(confs.JWT.JWTSecret))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not create token",
			"error":   string(err.Error()),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("authorization", tokenString, 3600*24*30, "", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{})

}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": user,
	})
}
