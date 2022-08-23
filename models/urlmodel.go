package models

import "database/sql"

type Url struct {
	LinkID       int64  `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl" binding:"required"`
	MainUrl      string `json:"MainUrl" binding:"required"`
	ExpDate      string `json:"ExpDate" binding:"required"`
	CollectionID int64  `json:"CollectionID"`
}

/*
	UrlAlt is an alternative url struct format which
	can be used in cases we consider a chance of having
	a null value as result of our select query from
	collection table.
*/
type UrlAlt struct {
	LinkID       int64         `json:"LinkID"`
	ShortenedUrl string        `json:"ShortenedUrl" binding:"required"`
	MainUrl      string        `json:"MainUrl" binding:"required"`
	ExpDate      string        `json:"ExpDate" binding:"required"`
	CollectionID sql.NullInt64 `json:"CollectionID"`
}
