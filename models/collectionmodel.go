package models

type Collection struct {
	CollectionID   int    `json:"CollectionID"`
	UserID         int    `json:"UserID" binding:"required"`
	CollectionName string `json:"CollectionName" binding:"required"`
}

type AddUrlToCollectionBody struct {
	LinkID       int64 `json:"LinkID" binding:"required"`
	CollectionID int   `json:"CollectionID" binding:"required"`
}

type EditCollectionByIDBody struct {
	CollectionID   int    `json:"CollectionID" binding:"required"`
	CollectionName string `json:"CollectionName" binding:"required"`
}
