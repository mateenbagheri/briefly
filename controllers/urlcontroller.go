package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/models"
	"github.com/mateenbagheri/briefly/scripts"
)

func CreateURL(c *gin.Context) {
	var body struct {
		MainUrl      string         `json:"MainUrl" binding:"required"`
		ExpDate      sql.NullString `json:"ExpDate" binding:"required"`
		CollectionID sql.NullInt64  `json:"CollectionID" binding:"required"`
	}

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
					collectionID=?;
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

			_, err = stmt.Exec(
				url.MainUrl,
				url.ShortenedUrl,
				url.ExpDate,
				url.CollectionID,
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
			url.ExpDate = tmpExpDate.Format("2006-01-02")
			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					collectionID=?,
					expDate=?;
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
		url.MainUrl = body.MainUrl
		if body.ExpDate.Valid {
			url.ExpDate = body.ExpDate.String

			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					expDate =?;
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
			url.ExpDate = tmpExpDate.Format("2006-01-02")
			stmt, err = Mysql.Prepare(
				`
				INSERT INTO links 
				SET 
					link=?, 
					shortened=?, 
					expDate=?;
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

func GetURLByShortened(c *gin.Context) {
	var url models.UrlAlt

	var now string
	currentTime := time.Now()
	now = currentTime.Format("2006-01-02")

	shortened := c.Param("ShortenedUrl")

	if shortened == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "ShortenedUrl not found in request params",
		})
		return
	}

	result, err := Mysql.Query(
		`
		SELECT 
			linkID,
			link,
			shortened,
			expDate,
			collectionID
		FROM links 
		WHERE shortened=?
		LIMIT 1; 
		`, shortened,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error running select sql query",
			"error":   string(err.Error()),
		})
		return
	}

	if result.Next() {
		err = result.Scan(
			&url.LinkID,
			&url.MainUrl,
			&url.ShortenedUrl,
			&url.ExpDate,
			&url.CollectionID,
		)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not match database url type with body",
				"error":   string(err.Error()),
			})
			return
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no url with given shortened version was found.",
			"error":   string(err.Error()),
		})
		return
	}

	// TODO :: TEST
	if url.ExpDate < now {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "the expiration date of this url is passed",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, url)
}

func GetCollectionURLs(c *gin.Context) {
	var urls []models.Url

	id := c.Param("CollectionID")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "could not find CollectionID in parameters",
			"err":     "CollectionID == \"\"",
		})
		return
	}

	results, err := Mysql.Query(
		`
		SELECT
			collectionID,
			expDate,
			link,
			linkID,
			shortened
		FROM Links
		WHERE collectionID = ?
		`, id,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error running select query",
			"error":   string(err.Error()),
		})
	}

	for results.Next() {
		var url models.Url

		err := results.Scan(
			&url.CollectionID,
			&url.ExpDate,
			&url.MainUrl,
			&url.LinkID,
			&url.ShortenedUrl,
		)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not match url type with body",
				"error":   string(err.Error()),
			})
		}

		urls = append(urls, url)
	}

	c.IndentedJSON(http.StatusOK, urls)
}

func GetUserURLs(c *gin.Context) {
	var urls []models.UrlAlt

	userID := c.Param("UserID")

	if userID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "No UserID parameter was found in request.",
		})
		return
	}

	results, err := Mysql.Query(
		`
		SELECT 
			L.collectionID,
			L.expDate,
			L.link,
			L.linkID,
			L.shortened
		FROM Links AS L
			INNER JOIN Collections AS C
				ON L.collectionID = C.collectionID
			INNER JOIN Users AS U
				ON U.userID = C.userID
		WHERE U.userID = ?
		`, userID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error running select query",
			"error":   string(err.Error()),
		})
	}

	for results.Next() {
		var url models.UrlAlt

		err := results.Scan(
			&url.CollectionID,
			&url.ExpDate,
			&url.MainUrl,
			&url.LinkID,
			&url.ShortenedUrl,
		)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not match url type with body.",
				"error":   string(err.Error()),
			})
		}

		urls = append(urls, url)
	}
	c.IndentedJSON(http.StatusOK, urls)
}

func DeleteURLByID(c *gin.Context) {
	id := c.Param("LinkID")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Could not find LinkID parameter in request.",
		})
		return
	}

	result, err := Mysql.Exec(`
		DELETE FROM Links
		WHERE linkID = ?
		`, id,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "could not run DELETE query to remove the URL.",
			"error":   err.Error(),
		})
		return
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return
	}

	if int(count) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Could not find any URLs with given ID",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "URL has been deleted successfully",
		"result":  result,
	})
}
