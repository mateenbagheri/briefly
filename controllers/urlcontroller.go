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
	var body models.CreateUrlBody

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

	tmpCreateDate := time.Now()

	url.MainUrl = body.MainUrl
	url.ShortenedUrl = scripts.ShortenUrl(body.MainUrl)
	url.CreateDate = tmpCreateDate.Format("2006-01-02")

	// checking if this url is already shortened in our system
	statement := `
		SELECT 
			linkID,
			shortened, 
			createDate,
			expDate
		FROM links 
		WHERE shortened = ?;
	`

	err = Mysql.QueryRow(statement, url.ShortenedUrl).Scan(
		&url.LinkID,
		&url.ShortenedUrl,
		&url.CreateDate,
		&url.ExpDate,
	)

	if err == nil {
		c.IndentedJSON(http.StatusOK, url)
		return
	}

	// Inserting URL add result into database
	if body.ExpDate.Valid {
		url.ExpDate = body.ExpDate.String
		stmt, err = Mysql.Prepare(
			`
			INSERT INTO links 
			SET 
				link=?, 
				shortened=?, 
				expDate=?,
				createDate=?;
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

		res, err := stmt.Exec(
			url.MainUrl,
			url.ShortenedUrl,
			url.ExpDate,
			url.CreateDate,
		)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "error inserting data into database type[1]",
				"error":   string(err.Error()),
			})
			return
		}

		lid, err := res.LastInsertId()

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not retrieve last inserted data",
				"error":   string(err.Error()),
			})
			return
		}

		url.LinkID = lid
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
				expDate=?,
				createDate=?;
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

		res, err := stmt.Exec(
			url.MainUrl,
			url.ShortenedUrl,
			url.ExpDate,
			url.CreateDate,
		)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "error inserting data into database type[2]",
				"error":   string(err.Error()),
			})
			return
		}

		lid, err := res.LastInsertId()

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not retrieve last inserted data",
				"error":   string(err.Error()),
			})
			return
		}

		url.LinkID = lid
	}

	c.IndentedJSON(http.StatusOK, url)
}

func GetURLByShortened(c *gin.Context) {
	var url models.DBUrl

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
			L.linkID,
			L.link,
			L.shortened,
			L.expDate,
			L.createDate
		FROM links AS L
		WHERE shortened=?;
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
			&url.CreateDate,
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

	if url.ExpDate < now {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "the expiration date of this url is passed",
		})
		return
	}

	stmt, err := Mysql.Prepare(`
		INSERT INTO linkhits 
		SET linkID=?, hitDate=?;
	`)

	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error in creating insert into statement",
			"error":   string(err.Error()),
		})
		return
	}

	_, err = stmt.Exec(url.LinkID, now)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed inserting the data into database.",
			"error":   string(err.Error()),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, url)
}

func GetCollectionURLs(c *gin.Context) {
	var urls []models.ReportUrl

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
			LC.collectionID,
			L.expDate,
			L.createDate,
			L.link,
			L.linkID,
			L.shortened,
			(
				SELECT COUNT(LH.hitID)
				FROM linkhits AS LH
				WHERE LH.linkID = LH.linkID
			) AS hitNumber
		FROM links AS L
			INNER JOIN collectionlinks AS LC
				ON LC.linkID = L.linkID
		WHERE LC.collectionID = ?
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
		var url models.ReportUrl

		err := results.Scan(
			&url.CollectionID,
			&url.ExpDate,
			&url.CreateDate,
			&url.MainUrl,
			&url.LinkID,
			&url.ShortenedUrl,
			&url.HitNumber,
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
			C.collectionID,
			L.expDate,
			L.createDate,
			L.link,
			L.linkID,
			L.shortened
		FROM links AS L
			INNER JOIN collectionlinks AS LC
				ON LC.linkID = L.linkID
			INNER JOIN collections AS C
				ON LC.collectionID = C.collectionID
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
			&url.CreateDate,
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
		DELETE FROM links
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

func GetUserUrlHistory(c *gin.Context) {
	userID := c.Param("UserID")
	var urls []models.HistoryUrl

	results, err := Mysql.Query(
		`
		SELECT
			C.collectionID,
			C.collectionName,
			L.expDate,
			L.createDate,
			L.link,
			L.linkID,
			L.shortened,
			LH.hitDate
		FROM briefly.links AS L
			INNER JOIN briefly.collectionlinks AS LC
				ON LC.linkID = L.linkID
			INNER JOIN briefly.collections AS C
				ON C.collectionID = LC.collectionID
			INNER JOIN linkhits AS LH
				ON LH.linkID = L.linkID
		WHERE C.userID = ?
		ORDER BY L.createDate DESC
		`, userID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error running select sql query",
			"error":   string(err.Error()),
		})
		return
	}

	for results.Next() {
		var url models.HistoryUrl

		err := results.Scan(
			&url.CollectionID,
			&url.CollectionName,
			&url.ExpDate,
			&url.CreateDate,
			&url.MainUrl,
			&url.LinkID,
			&url.ShortenedUrl,
			&url.HitDate,
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

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": urls,
	})
}

func GetUserUrlReport(c *gin.Context) {
	userID := c.Param("UserID")
	var urls []models.ReportUrl

	result, err := Mysql.Query(
		`
		SELECT
			L.linkID,
			L.shortened,
			L.link,
			C.collectionName,
			C.collectionID,
			L.expDate,
			L.createDate,
			(
				SELECT COUNT(LH.hitID)
				FROM linkhits AS LH
				WHERE linkID = L.linkID
			) AS hitNumber
		FROM links AS L
			INNER JOIN collectionlinks AS CL
				ON CL.linkID = L.linkID
			INNER JOIN collections AS C
				ON C.collectionID = CL.collectionID
		WHERE C.userID = ?
		ORDER BY hitNumber DESC
		`, userID,
	)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "could not select required data from database",
			"error":   string(err.Error()),
		})
	}

	for result.Next() {
		var url models.ReportUrl

		result.Scan(
			&url.LinkID,
			&url.ShortenedUrl,
			&url.MainUrl,
			&url.CollectionName,
			&url.CollectionID,
			&url.ExpDate,
			&url.CreateDate,
			&url.HitNumber,
		)

		urls = append(urls, url)
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": urls,
	})
}
