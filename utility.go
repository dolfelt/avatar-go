package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"

	_ "github.com/lib/pq"
)

// DefaultSizes is the list of sizes
var DefaultSizes = map[string]int{
	"small":    128,
	"medium":   256,
	"large":    512,
	"original": 1024,
}

// DefaultSizeKeys is the list of acceptable size strings
func DefaultSizeKeys() []string {
	keys := make([]string, 0, len(DefaultSizes))
	for k := range DefaultSizes {
		keys = append(keys, k)
	}
	return keys
}

// AppError handles generic application errors
type AppError struct {
	msg string
}

func (e *AppError) Error() string {
	return e.msg
}

// Configuration loads config from JSON
type Configuration struct {
	AWS           awsConfig `json:"aws"`
	DB            dbConfig  `json:"db"`
	DefaultAvatar Avatar    `json:"default_avatar"`
	JwtKey        string    `json:"jwt_key"`
}

type awsConfig struct {
	AccessKey    string `json:"key"`
	AccessSecret string `json:"secret"`
	Bucket       string `json:"bucket"`
}

type dbConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

// SigningKey returns the byte value of the string signing key
func (config Configuration) SigningKey() []byte {
	key := make([]byte, len(config.JwtKey))
	copy(key[:], config.JwtKey)

	return key
}

// LoadConfig loads external configuration file
func LoadConfig(path string) Configuration {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}

	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}

	return config
}

/*
	Database Connections
*/

func (config dbConfig) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://"+config.User+":"+url.QueryEscape(config.Password)+"@"+config.Host+":"+config.Port+"/"+config.Database+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Additional special functions that are lacking in Go

// MinInt is a shim for determining the minimum for integers rather than floats
func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}
