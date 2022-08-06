package controllers

import (
	"database/sql"

	"github.com/mateenbagheri/briefly/configs"
)

var mysql *sql.DB = configs.ConnectMySQL()
