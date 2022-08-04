package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/configs"
	"github.com/mateenbagheri/briefly/models"
)

// TODO :: ADD ERROR HANDLING TO THIS MODULE!

var mysql *sql.DB = configs.ConnectMySQL()

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

}
