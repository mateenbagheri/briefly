package models

type User struct {
	UserID     int    `json:"UserID"`
	Name       string `json:"Name" binding:"required"`
	FamilyName string `json:"FamilyName" binding:"required"`
	Password   string `json:"Password" binding:"required"`
	Email      string `json:"Email" validate:"required"`
	Salt       string `json:"Salt"`
}
