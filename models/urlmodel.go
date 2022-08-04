package models

type Url struct {
	LinkID       int    `json:"LinkID"`
	ShortenedUrl string `json:"ShortenedUrl"`
	MainUrl      string `json:"MainUrl" validate:"required"`
	ExpDate      string `json:"ExpDate" validate:"required"`
	CollectionID string `json:"CollectionID" validate:"required"`
	HitNumbers   string `json:"HitNumbers"`
}
