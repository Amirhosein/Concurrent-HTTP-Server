package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

// initialize the database using postgres
const (
	DB_HOST     = "db"
	DB_PORT     = "5432"
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "concurrent-http-server"
)

// initialaize the redis client
const (
	REDIS_HOST     = "redis"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
)

func InitDB() *sql.DB {
	psqlInfo := fmt.Sprintf("user=%s "+
		"password=%s host=%s port=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		_, err = db.Exec("create database " + DB_NAME)
		if err != nil {
			log.Fatal(err)
		}
	}

	// check if the table exists
	var tableExists bool

	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		log.Fatal(err)
	}

	if !tableExists {
		_, err = db.Exec("CREATE TABLE users (id SERIAL PRIMARY KEY, username VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL, files TEXT[])")
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")

	return db
}

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: REDIS_PASSWORD,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the redis database")

	return client
}
