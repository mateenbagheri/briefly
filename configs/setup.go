package configs

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var confs Config

func ConnectRedis() (*redis.Client, error) {
	confs, _ = LoadConfig()

	fmt.Println(confs.Redis)

	fmt.Println()
	client := redis.NewClient(&redis.Options{
		Addr:     confs.Redis.Address,
		Password: confs.Redis.Password,
		DB:       confs.Redis.DB,
	})

	if err := client.Ping().Err(); err != nil {
		log.Fatal("Could not connect to redis database.")
		return nil, err
	}
	fmt.Println("Connected to redis database.")

	return client, nil
}

func ConnectMySQL() (db *sql.DB) {
	confs, _ = LoadConfig()

	address := confs.MySql.Address
	driver := confs.MySql.Driver
	user := confs.MySql.Username
	password := confs.MySql.Address
	schema := confs.MySql.Schema
	port := confs.MySql.Port

	db, err := sql.Open(driver, user+":"+password+"@tcp("+address+":"+port+")/"+schema)

	if err != nil {
		log.Fatal("Could not connect to mysql database.")
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL Database.")

	return db
}
