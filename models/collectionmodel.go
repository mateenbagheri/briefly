package models

type collection struct {
	CollectionID   int    `json:"CollectionID"`
	UserID         int    `json:"UserID"`
	CollectionName string `json:"CollectionName" validate:"required"`
}
