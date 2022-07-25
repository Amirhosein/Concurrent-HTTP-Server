package model

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
)

type User struct {
	ID       uint64
	Username string
	Files    []string
	jwt.StandardClaims
}

func (s User) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

type UserRepo interface {
	Get(username string, password string) (User, error)
	Set(username string, password string) (User, error)
	Update(user User) error
}

type UserDB struct {
	DB *sql.DB
}

type UserCache struct {
	DB    UserDB
	Redis *redis.Client
}

func (s UserDB) Get(username string, password string) (User, error) {
	var user User

	if password == "" {
		err := s.DB.QueryRow("SELECT id, username, files FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Files)
		if err != nil {
			return user, err
		}
	}

	err := s.DB.QueryRow("SELECT id, username, files FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.ID, &user.Username, &user.Files)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserDB) Set(username string, password string) (User, error) {
	user := User{
		Username: username,
		Files:    make([]string, 0),
	}

	err := s.DB.QueryRow("INSERT INTO users (username, password, files) VALUES ($1, $2, $3) RETURNING id, username", username, password, pq.Array(user.Files)).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserDB) Update(user User) error {
	_, err := s.DB.Exec("UPDATE users SET files = $1 WHERE id = $2", pq.Array(user.Files), user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s UserCache) Get(username string, password string) (User, error) {
	var user User

	result, err := s.Redis.Get(username).Result()
	err1 := user.UnmarshalBinary([]byte(result))
	if err1 != nil {
		return user, err1
	}
	if err == redis.Nil {
		log.Println("Redis cache miss")

		user, err = s.DB.Get(username, password)
		if err != nil {
			return user, err
		}

		err = s.Redis.Set(username, user, time.Hour*24).Err()

		return user, err
	} else if err != nil {
		return user, err
	}

	log.Println("Redis cache hit")

	return user, nil
}

func (s UserCache) Set(username string, password string) (User, error) {
	var user User
	user, err := s.DB.Set(username, password)
	if err != nil {
		return user, err
	}

	err = s.Redis.Set(username, user, time.Hour*24).Err()
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s UserCache) Update(user User) error {
	err := s.DB.Update(user)
	if err != nil {
		return err
	}

	err = s.Redis.Set(user.Username, user, time.Hour*24).Err()
	if err != nil {
		return err
	}

	return nil
}
