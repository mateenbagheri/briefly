package models

import "database/sql"

type User struct {
	UserID     int            `json:"UserID"`
	Name       string         `json:"Name" binding:"required"`
	FamilyName string         `json:"FamilyName" binding:"required"`
	Password   string         `json:"Password" binding:"required"`
	Email      string         `json:"Email" binding:"required"`
	Salt       sql.NullString `json:"Salt"`
}
