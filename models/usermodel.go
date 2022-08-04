package models

type User struct {
	UserID     int    `json:"UserID"`
	Name       string `json:"Name" validate:"required"`
	FamilyName string `json:"FamilyName" validate:"required"`
	Password   string `json:"Password" validate:"required"`
	Salt       string `json:"Salt" validate:"required"`
	Email      string `json:"Email" validate:"required"`
}
