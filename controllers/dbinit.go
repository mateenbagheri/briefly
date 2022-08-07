package controllers

import (
	"database/sql"

	"github.com/mateenbagheri/briefly/configs"
)

var Mysql *sql.DB = configs.ConnectMySQL()
