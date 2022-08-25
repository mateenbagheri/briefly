package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mateenbagheri/briefly/configs"
	"github.com/mateenbagheri/briefly/controllers"
	"github.com/mateenbagheri/briefly/models"
)

func RequireAuth(c *gin.Context) {
	confs, _ := configs.LoadConfig()

	tokenStringList, ok := c.Request.Header["Authorization"]

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Authorization header is required",
		})
	}

	tokenString := tokenStringList[0]

	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// decode/validate
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(confs.JWT.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// find the user with token sub
		var user models.User
		id := claims["sub"]

		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "no sub found in given JWT claim",
			})
			return
		}

		result, err := controllers.Mysql.Query("SELECT * FROM briefly.users WHERE userID = ?;", id)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not select user info from database.",
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
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "could not scan user info from result set.",
					"error":   string(err.Error()),
				})
				return
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// attach request
		c.Set("UserID", user.UserID)
		// continiue
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
