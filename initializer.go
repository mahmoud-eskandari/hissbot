package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

//Db gorm pool
var Db *gorm.DB
var conf *ini.File

//Redis redis connection pool
var Redis *redis.Client

//LoadDB Database Connection
func LoadDB() {
	connectionString := Str("database", "db_user", "root") + ":" + Str("database", "db_pass", "") + "@" + Str("database", "db_host", "") + "/" + Str("database", "db_name", "test") + "?charset=" + Str("database", "db_char", "utf8") + "&parseTime=True&loc=" + Str("database", "db_loc", "Local")
	var err error
	Db, err = gorm.Open("mysql", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Database Connection Error %v", err))
	}
	Db.DB().SetMaxIdleConns(20)
	Db.DB().SetMaxOpenConns(100)
	Db.DB().SetConnMaxLifetime(time.Minute)

	// Connect to redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Str("database", "redis_host", "localhost"), Int("database", "redis_port", 6379)),
		Password: Str("database", "redis_password", ""),
		DB:       0, // use default DB
	})
	_, err = Redis.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Redis Connection Error %v", err))
	}
}

//LoadConfig Load ini File
func LoadConfig() {
	ConfigPath := "./config.ini"
	if len(os.Args) > 1 {
		ConfigPath = os.Args[1]
	}
	err := errors.New("")
	conf, err = ini.Load(ConfigPath)
	if err != nil {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Printf("Fail to read file: %v at %s", err, dir)
		os.Exit(1)
	}

	// Hashtable
	HashTable = strings.Split(Str("hash", "strings", "bHoREzxvbn"), "")
	if len(HashTable) < 10 {
		panic("hash[strings] is small (min 10 char)")
	}
	HashSalt = int64(Int("hash", "int", 12345))
}

//CloseDB close database connections
func CloseDB() {
	_ = Db.Close()
	_ = Redis.Close()
}

//Int Get Integer Config data
func Int(section string, key string, defaultVal int) int {
	return conf.Section(section).Key(key).MustInt(defaultVal)
}

//Int64 Get Integer 64 Config data
func Int64(section string, key string, defaultVal int64) int64 {
	return conf.Section(section).Key(key).MustInt64(defaultVal)
}

//Str Get String Config data
func Str(section string, key string, defaultVal string) string {
	return conf.Section(section).Key(key).MustString(defaultVal)
}

//Bool Get Bool Config data
func Bool(section string, key string, defaultVal bool) bool {
	return conf.Section(section).Key(key).MustBool(defaultVal)
}

//SaveToRedis save data into redis db
func SaveToRedis(key string, data interface{}, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Redis.Set(key, b, ttl).Err()
}
