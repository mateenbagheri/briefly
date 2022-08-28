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

/*
	ReportUrl is the url output type we use when we want
	to extract a report from our program.
	This type is basically like type url but also includes
	hitnumbers of the specified url.
*/
type ReportUrl struct {
	LinkID       int64  `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl" binding:"required"`
	MainUrl      string `json:"MainUrl" binding:"required"`
	ExpDate      string `json:"ExpDate" binding:"required"`
	HitNumber    int64  `json:"hitnumber" binding:"required"`
	CollectionID int64  `json:"CollectionID"`
}
