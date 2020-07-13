package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Db gorm pool
var Db *gorm.DB

//Redis redis connection pool
var Redis *redis.Client

//LoadDB Database Connection
func initConnections() {
	// Hashtable
	HashTable = strings.Split(Str("HASH_TABLE", "bHoREzxvbn"), "")
	if len(HashTable) < 10 {
		panic("hash[strings] is small (min 10 char)")
	}
	HashSalt = int64(Int("RANDOM_INT", 12345))

	connectionString := Str("DB_USER", "root") + ":" + Str("DB_PASSWORD", "") + "@" + Str("DB_HOST", "") + "/" + Str("DB", "test") + "?charset=utf8mb4&parseTime=True&loc=" + Str("DB_LOCATION", "Local")
	var err error
	Db, err = gorm.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("database connection error %v", err)
	}
	Db.DB().SetMaxIdleConns(20)
	Db.DB().SetMaxOpenConns(100)
	Db.DB().SetConnMaxLifetime(time.Minute)
	// migrate models
	migrate()
	// Connect to redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Str("REDIS_HOST", "localhost"), Int("REDIS_PORT", 6379)),
		Password: Str("REDIS_PASSWORD", ""),
		DB:       0, // use default DB
	})
	_, err = Redis.Ping().Result()
	if err != nil {
		log.Fatalf("redis connection error %v", err)
	}
}

//CloseDB close database connections
func CloseDB() {
	_ = Db.Close()
	_ = Redis.Close()
}

//Int Get Integer Config data
func Int(key string, defaultVal int) int {
	data := os.Getenv(key)
	if len(data) == 0 {
		return defaultVal
	}
	out, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return defaultVal
	}
	return int(out)
}

//Str Get String Config data
func Str(key string, defaultVal string) string {
	data := os.Getenv(key)
	if len(data) == 0 {
		return defaultVal
	}
	return data
}

//Bool Get Bool Config data
func Bool(key string, defaultVal bool) bool {
	data := os.Getenv(key)
	if len(data) == 0 {
		return defaultVal
	}
	if strings.ToLower(data) == "true" {
		return true
	}
	return false
}

//SaveToRedis save data into redis db
func SaveToRedis(key string, data interface{}, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Redis.Set(key, b, ttl).Err()
}
