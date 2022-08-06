package models

type Collection struct {
	CollectionID   int    `json:"CollectionID"`
	UserID         int    `json:"UserID" binding:"required"`
	CollectionName string `json:"CollectionName" binding:"required"`
}
