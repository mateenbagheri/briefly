package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/models"
	"github.com/mateenbagheri/briefly/scripts"
)

func CreateURL(c *gin.Context) {
	var body struct {
		MainUrl      string         `json:"MailUrl" binding:"required"`
		ExpDate      sql.NullString `json:"ExpDate" binding:"required"`
		CollectionID sql.NullInt64  `json:"CollectionID" binding:"required"`
	}

	// TODO :: TEST THIS MODULE THOROUGHLY
	var url models.Url

	err := c.BindJSON(&body)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not bind body to structure",
			"error":   string(err.Error()),
		})
		return
	}

	var stmt *sql.Stmt

	url.HitNumbers = 0
	url.MainUrl = body.MainUrl
	url.ShortenedUrl = scripts.ShortenUrl(body.MainUrl)

	if body.CollectionID.Valid {
		url.CollectionID = body.CollectionID.Int64
		if body.ExpDate.Valid {
			url.ExpDate = body.ExpDate.String
			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					expDate=?,
					collectionID=?,
					hitNumbers=?;
				`,
			)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "could not prepare insert statement [type 1]",
					"error":   string(err.Error()),
				})
				return
			}

			// TODO :: EXEC STATEMENT
			_, err = stmt.Exec(
				url.MainUrl,
				url.ShortenedUrl,
				url.ExpDate,
				url.CollectionID,
				url.HitNumbers,
			)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "error inserting data into database type[1]",
					"error":   string(err.Error()),
				})
				return
			}

		} else {
			tmpExpDate := time.Now()
			tmpExpDate = tmpExpDate.AddDate(0, 1, 0)
			url.ExpDate = tmpExpDate.Format("2006-02-01")
			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					collectionID=?,
					expDate=?,
					hitNumbers=?;
				`,
			)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "could not prepare insert statement [type 2]",
					"error":   string(err.Error()),
				})
				return
			}

			_, err = stmt.Exec(
				url.MainUrl,
				url.ShortenedUrl,
				url.CollectionID,
				url.ExpDate,
				url.HitNumbers,
			)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "error inserting data into database type[2]",
					"error":   string(err.Error()),
				})
				return
			}
		}

	} else {
		url.HitNumbers = 0
		url.MainUrl = body.MainUrl
		if body.ExpDate.Valid {
			url.ExpDate = body.ExpDate.String

			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					expDate =?,
					hitNumbers=?;
				`,
			)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "could not prepare insert statement [type 3]",
					"error":   string(err.Error()),
				})
				return
			}

			_, err = stmt.Exec(
				url.MainUrl,
				url.ShortenedUrl,
				url.ExpDate,
				url.HitNumbers,
			)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "error inserting data into database type[3]",
					"error":   string(err.Error()),
				})
				return
			}
		} else {
			tmpExpDate := time.Now()
			tmpExpDate = tmpExpDate.AddDate(0, 1, 0)
			url.ExpDate = tmpExpDate.Format("2006-02-01")
			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					expDate=?,
					hitNumbers=?;
				`,
			)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "could not prepare insert statement [type 4]",
					"error":   string(err.Error()),
				})
				return
			}

			_, err = stmt.Exec(
				url.MainUrl,
				url.ShortenedUrl,
				url.ExpDate,
				url.HitNumbers,
			)

			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "error inserting data into database type[4]",
					"error":   string(err.Error()),
				})
				return
			}

		}

	}
	c.IndentedJSON(http.StatusOK, url)
}
