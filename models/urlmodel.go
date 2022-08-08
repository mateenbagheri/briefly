package models

type Url struct {
	LinkID       int64  `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl" binding:"required"`
	MainUrl      string `json:"MainUrl" binding:"required"`
	ExpDate      string `json:"ExpDate" binding:"required"`
	HitNumbers   int64  `json:"HitNumbers" binding:"required"`
	CollectionID int64  `json:"CollectionID"`
}
