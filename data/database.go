package data

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// DB wraps the gorm DB interface
type DB struct {
	*gorm.DB
}

// Connect begins the connection with the database
func Connect() (*DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("DBUser"),
		url.QueryEscape(viper.GetString("DBPassword")),
		viper.GetString("DBHost"),
		viper.GetString("DBPort"),
		viper.GetString("DBDatabase"),
	)
	conn, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}
