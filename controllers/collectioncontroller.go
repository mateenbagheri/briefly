package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/models"
)

// TODO :: ADD ERROR HANDLING TO THIS MODULE!

func GetAllCollections(c *gin.Context) {
	var collections []models.Collection

	results, err := mysql.Query("SELECT * FROM briefly.collections;")

	if err != nil {
		log.Fatal("Could not get all the collections from database")
	}

	for results.Next() {
		var collection models.Collection

		err := results.Scan(&collection.CollectionID, &collection.CollectionName, &collection.UserID)

		if err != nil {
			log.Fatal(err)
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
		log.Fatal("CollectionID parameter was not found")
	}

	result, err := mysql.Query("SELECT * FROM briefly.collections WHERE CollectionID = ?;", id)

	if err != nil {
		log.Fatal(err)
	}

	if result.Next() {

		err = result.Scan(&collection.CollectionID, &collection.CollectionName, &collection.UserID)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Collection Not Found")
	}
	c.IndentedJSON(http.StatusOK, collection)
}

func DeleteCollectionByID(c *gin.Context) {
	id := c.Param("CollectionID")

	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	result, err := mysql.Query("DELETE FROM collections WHERE CollectionID = ?", id)
	if err != nil {
		log.Fatal("Could not delete collection with this id from database")
	} else {
		c.IndentedJSON(http.StatusOK, result)
	}
}

func CreateCollection(c *gin.Context) {
	var newCollection models.Collection

	if err := c.BindJSON(&newCollection); err != nil {
		return
	}

	stmt, err := mysql.Prepare("INSERT INTO collections SET collectionName=?, userID=?;")

	if err != nil {
		log.Fatal("Error Preparing Insert Statement")
	}

	_, err = stmt.Exec(newCollection.CollectionName, newCollection.UserID)

	if err != nil {
		log.Fatal("Error inserting data into database")
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

	stmt, err := mysql.Prepare("UPDATE collections SET collectionName=?, userID=? WHERE collectionID=?;")

	if err != nil {
		log.Fatal("Error Preparing Update Statement")
	}

	_, err = stmt.Exec(collection.CollectionName, collection.UserID, id)

	if err != nil {
		log.Fatal("Error updating data in database")
	}

	c.IndentedJSON(http.StatusNoContent, "")
}
