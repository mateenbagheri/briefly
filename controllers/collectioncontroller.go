package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/models"
)

func GetAllCollections(c *gin.Context) {
	var collections []models.Collection

	results, err := Mysql.Query(`
		SELECT
			C.collectionID,
			C.collectionName,
			C.userID
		FROM briefly.collections AS C;
	`)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Could not get all the collections from database",
			"error":   string(err.Error()),
		})
		return
	}

	for results.Next() {
		var collection models.Collection

		err := results.Scan(
			&collection.CollectionID,
			&collection.CollectionName,
			&collection.UserID,
		)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not match collection type with body",
				"error":   string(err.Error()),
			})
			return
		}

		collections = append(collections, collection)
	}

	results.Close()

	c.IndentedJSON(http.StatusOK, collections)
}

func GetCollectionByID(c *gin.Context) {
	id := c.Param("CollectionID")
	var collection models.Collection

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter CollectionID was not found in request",
		})
		return
	}

	result, err := Mysql.Query("SELECT * FROM briefly.collections WHERE CollectionID = ?;", id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error runnin sql select query.",
			"error":   string(err.Error()),
		})
		return
	}

	if result.Next() {

		err = result.Scan(&collection.CollectionID, &collection.CollectionName, &collection.UserID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not match database collection type with body",
				"error":   string(err.Error()),
			})
			return
		}
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no collection with given ID was found.",
			"error":   string(err.Error()),
		})
		return
	}
	c.IndentedJSON(http.StatusOK, collection)
}

func DeleteCollectionByID(c *gin.Context) {
	id := c.Param("CollectionID")

	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	result, err := Mysql.Query("DELETE FROM collections WHERE CollectionID = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Could not delete collection with this id from database",
			"error":   string(err.Error()),
		})
		return
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
}

func CreateCollection(c *gin.Context) {
	var newCollection models.Collection

	if err := c.BindJSON(&newCollection); err != nil {
		return
	}

	stmt, err := Mysql.Prepare("INSERT INTO collections SET collectionName=?, userID=?;")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error in creating insert into statement",
			"error":   string(err.Error()),
		})
		return
	}

	_, err = stmt.Exec(newCollection.CollectionName, newCollection.UserID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed inserting the data into database.",
			"error":   string(err.Error()),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, newCollection)
}

func EditCollectionByID(c *gin.Context) {
	var body struct {
		CollectionID   int    `json:"CollectionID" binding:"required"`
		CollectionName string `json:"CollectionName" binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "could not bind collection to collection model",
			"error":   string(err.Error()),
		})
		return
	}

	id := c.Param("CollectionID")

	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	stmt, err := Mysql.Prepare(`
		UPDATE collections 
		SET 
			collectionName=?
		WHERE collectionID=?;
	`)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error in creating update statement",
			"error":   string(err.Error()),
		})
		return
	}

	_, err = stmt.Exec(body.CollectionName, id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error updating collection data in database",
			"error":   string(err.Error()),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "collection successfully edited",
	})
}

func GetUserCollections(c *gin.Context) {
	var collections []models.Collection
	userID := c.Param("UserID")

	results, err := Mysql.Query(
		`
		SELECT
			collectionID,
			collectionName,
			userID
		FROM briefly.Collections
		WHERE userID = ?
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
		var currentCollection models.Collection

		err := results.Scan(
			&currentCollection.CollectionID,
			&currentCollection.CollectionName,
			&currentCollection.UserID,
		)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "could not match collection type with body",
				"error":   string(err.Error()),
			})
			return
		}

		collections = append(collections, currentCollection)

	}

	results.Close()

	c.IndentedJSON(http.StatusOK, collections)
}
