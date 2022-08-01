package configs

import (
	"fmt"

	"github.com/go-redis/redis"
)

func ConnectRedis() {
	confs, _ := LoadConfig()

	fmt.Println(confs.Redis)

	fmt.Println()
	client := redis.NewClient(&redis.Options{
		Addr:     confs.Redis.Address,
		Password: confs.Redis.Password,
		DB:       confs.Redis.DB,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func ConnectMySQL() {

}
