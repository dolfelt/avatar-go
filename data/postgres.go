package data

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// PostgresDB wraps the gorm DB interface
type PostgresDB struct {
	Gorm *gorm.DB
}

// Connect begins the connection with the database
func (p *PostgresDB) Connect() error {
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
		return err
	}

	p.Gorm = conn

	return nil
}

func (p *PostgresDB) FindByHash(hash string) (*Avatar, error) {
	var avatar Avatar

	p.Gorm.Find(&avatar, "hash = ?", hash)
	if len(avatar.Hash) == 0 {
		return nil, fmt.Errorf("Cannot find avatar with hash %s", hash)
	}

	return &avatar, nil
}

func (p *PostgresDB) Save(a *Avatar) error {
	res := p.Gorm.Save(a)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		res = p.Gorm.Create(a)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (p *PostgresDB) Migrate() error {
	p.Gorm.AutoMigrate(&Avatar{})
	return nil
}
