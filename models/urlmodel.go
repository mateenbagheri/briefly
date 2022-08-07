package models

type Url struct {
	LinkID       int    `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl"`
	MainUrl      string `json:"MainUrl" binding:"required"`
	ExpDate      string `json:"ExpDate" binding:"required"`
	CollectionID string `json:"CollectionID" binding:"required"`
	HitNumbers   string `json:"HitNumbers"`
}
