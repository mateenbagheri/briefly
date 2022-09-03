package models

import "database/sql"

type Url struct {
	LinkID       int64  `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl" binding:"required"`
	MainUrl      string `json:"MainUrl" binding:"required"`
	ExpDate      string `json:"ExpDate" binding:"required"`
	CreateDate   string `json:"CreateDate" binding:"required"`
	CollectionID int64  `json:"CollectionID"`
}

/*
	A map-like data structure supposed to be matched with
	table <links> in MySQL database.
	We use this type when we have no need of collection mixed
	with URL related data.
*/
type DBUrl struct {
	LinkID       int64  `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl" binding:"required"`
	MainUrl      string `json:"MainUrl" binding:"required"`
	ExpDate      string `json:"ExpDate" binding:"required"`
	CreateDate   string `json:"CreateDate" binding:"required"`
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
	CreateDate   string        `json:"CreateDate" binding:"required"`
	CollectionID sql.NullInt64 `json:"CollectionID"`
}

/*
	CreateUrlBody struct type is the type used for getting
	data in CreateURL function in urlcontroller.go
*/
type CreateUrlBody struct {
	MainUrl      string         `json:"MainUrl" binding:"required"`
	ExpDate      sql.NullString `json:"ExpDate" binding:"required"`
	CollectionID sql.NullInt64  `json:"CollectionID" binding:"required"`
}

/*
	ReportUrl is the url output type we use when we want
	to extract a report from our program.
	This type is basically like type url but also includes
	hitnumbers of the specified url.
*/
type ReportUrl struct {
	LinkID         int64  `json:"LinkID"`
	ShortenedUrl   string `json:"ShortenedUrl" binding:"required"`
	MainUrl        string `json:"MainUrl" binding:"required"`
	CollectionName string `json:"CollectionName"`
	CollectionID   int64  `json:"CollectionID"`
	ExpDate        string `json:"ExpDate" binding:"required"`
	CreateDate     string `json:"CreateDate" binding:"required"`
	HitNumber      int64  `json:"hitnumber" binding:"required"`
}

/*
	HistoryUrl Is the url type used for /user/:UserID/history
	which is supposed to show us a timeline of user URL related
	activity.
	The main difference between this type and other data types
	is containing HitDate info.
*/
type HistoryUrl struct {
	LinkID         int64  `json:"LinkID"`
	ShortenedUrl   string `json:"ShortenedUrl" binding:"required"`
	MainUrl        string `json:"MainUrl" binding:"required"`
	CollectionName string `json:"CollectionName"`
	CollectionID   int64  `json:"CollectionID"`
	ExpDate        string `json:"ExpDate" binding:"required"`
	CreateDate     string `json:"CreateDate" binding:"required"`
	HitDate        string `json:"HitDate"`
}
