package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/models"
)

func GetAllCollections(c *gin.Context) {
	var collections []models.Collection

	results, err := Mysql.Query("SELECT * FROM briefly.collections;")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Could not get all the collections from database",
			"error":   err.Error(),
		})
		return
	}

	for results.Next() {
		var collection models.Collection

		err := results.Scan(&collection.CollectionID, &collection.CollectionName, &collection.UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not match collection type with body",
				"error":   err.Error(),
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Parameter CollectionID was not found in request",
		})
		return
	}

	result, err := Mysql.Query("SELECT * FROM briefly.collections WHERE CollectionID = ?;", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error runnin sql select query.",
			"error":   err.Error(),
		})
		return
	}

	if result.Next() {

		err = result.Scan(&collection.CollectionID, &collection.CollectionName, &collection.UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "could not match database collection type with body",
				"error":   err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "no collection with given ID was found.",
			"error":   err.Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Could not delete collection with this id from database",
			"error":   err.Error(),
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error in creating insert into statement",
			"error":   err.Error(),
		})
		return
	}

	_, err = stmt.Exec(newCollection.CollectionName, newCollection.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed inserting the data into database.",
			"error":   err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, newCollection)
}

func EditCollectionByID(c *gin.Context) {
	var collection models.Collection

	if err := c.BindJSON(&collection); err != nil {
		return
	}

	id := c.Param("CollectionID")

	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	stmt, err := Mysql.Prepare("UPDATE collections SET collectionName=?, userID=? WHERE collectionID=?;")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error in creating update statement",
			"error":   err.Error(),
		})
		return
	}

	_, err = stmt.Exec(collection.CollectionName, collection.UserID, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "error updating collection data in database",
			"error":   err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}
