package data

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// PostgresDB wraps the gorm DB interface
type PostgresDB struct {
	Gorm *gorm.DB
}

type AvatarPostgres struct {
	Avatar
	Sizes string `gorm:"column:sizes;type:text;not null" json:"-"` // list of available sizes
}

func (AvatarPostgres) TableName() string {
	return viper.GetString("TableName")
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
	var avatar AvatarPostgres

	p.Gorm.Find(&avatar, "hash = ?", hash)
	if len(avatar.Hash) == 0 {
		return nil, fmt.Errorf("Cannot find avatar with hash %s", hash)
	}

	// Unmarshal the JSON object into the struct
	json.Unmarshal([]byte(avatar.Sizes), &avatar.Avatar.Sizes)

	return &avatar.Avatar, nil
}

func (p *PostgresDB) Save(a *Avatar) error {
	// Convert sizes to JSON for storage in the database
	sizes, err := json.Marshal(a.Sizes)
	if err != nil {
		return err
	}
	ap := &AvatarPostgres{
		Avatar: *a,
		Sizes:  string(sizes),
	}
	res := p.Gorm.Save(ap)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		res = p.Gorm.Create(ap)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (p *PostgresDB) Migrate() error {
	p.Gorm.AutoMigrate(&AvatarPostgres{})
	return nil
}
