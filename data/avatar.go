package data

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Sizes is the available sizes for an avatar object.
type Sizes []string

// Avatar stores the data for each object
type Avatar struct {
	Hash      string    `gorm:"type:varchar(40);not null;primary_key" json:"hash"` // hash identifier of the object
	Type      string    `gorm:"type:char(4);not null" json:"type"`                 // file extension of the avatar
	SizesRaw  string    `gorm:"column:sizes;type:text;not null" json:"-"`          // list of available sizes
	Sizes     Sizes     `gorm:"-" sql:"-" json:"sizes"`                            // list of available sizes
	CreatedAt time.Time `json:"createdAt"`                                         // when the avatar was first created
	UpdatedAt time.Time `json:"updatedAt"`                                         // last update of the avatar
}

func (Avatar) TableName() string {
	return viper.GetString("TableName")
}

// FindAvatar searches the database for an avatar object
// given a hash string
func FindAvatar(db *DB, hash string) *Avatar {
	var avatar Avatar

	db.Find(&avatar, "hash = ?", hash)
	if len(avatar.Hash) == 0 {
		log.Printf("Cannot find avatar with hash %s", hash)
		return nil
	}

	// Parse JSON sizes array
	json.Unmarshal([]byte(avatar.SizesRaw), &avatar.Sizes)

	return &avatar
}

// Save the avatar to the database
func (a Avatar) Save(db *DB) error {

	// Convert sizes to JSON for storage in the database
	sizes, err := json.Marshal(a.Sizes)
	if err != nil {
		return err
	}
	a.SizesRaw = string(sizes)
	res := db.Save(&a)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		res = db.Create(&a)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

// DefaultAvatar returns the default placeholder object for no avatar
func DefaultAvatar(app *Application) *Avatar {
	return &Avatar{
		Hash:  "7505d64a54e061b7acd54ccd58b49dc43500b635",
		Type:  "jpg",
		Sizes: DefaultSizeKeys(),
	}
}

// GetURL gets the full URL of the avatar object for a given size, including
// the S3 bucket path.
func (a Avatar) GetURL(size string, bucket string) string {
	return "//s3.amazonaws.com/" + bucket + "/" + a.GetPath(size)
}

// GetPath returns the path to the file object for a given size.
func (a Avatar) GetPath(size string) string {
	// Provides segmentation to prevent any single directory from becoming
	// too large to be easily navigated.
	file := a.GetFilename(size)
	return file[:1] + "/" + file[1:3] + "/" + file
}

// GetFilename generates the file name of the object for a given size.
func (a Avatar) GetFilename(size string) string {
	return a.Hash + "." + size + "." + a.Type
}

// BestSize determines the best size for the avatar, using the requested size
// as a reference.
func (a Avatar) BestSize(size string) string {
	for _, s := range a.Sizes {
		if s == size {
			return size
		}
	}
	return a.Sizes[len(a.Sizes)-1]
}

// ConvertToSize converts a pixel int to the corresponding size string.
func ConvertToSize(pixels int) string {
	if pixels <= DefaultSizes["small"] {
		return "small"
	}
	for name, size := range DefaultSizes {
		if size >= pixels {
			return name
		}
	}
	return "medium"
}

// ValidAvatarSize determines if the string is a valid identifier
func ValidAvatarSize(size string) bool {
	for s := range DefaultSizes {
		if s == size {
			return true
		}
	}
	return false
}

// CheckAvatarSize determines if the size is valid and returns the closest identifier
// if it is not valid
func CheckAvatarSize(size string) string {
	if len(size) > 0 {
		if pixels, err := strconv.Atoi(size); err == nil {
			return ConvertToSize(pixels)
		} else if ValidAvatarSize(size) {
			return size
		}
	}

	return "medium"
}
