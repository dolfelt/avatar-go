package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

// Sizes is the available sizes for an avatar object.
type Sizes []string

// Avatar stores the data for each object
type Avatar struct {
	Hash      string    // hash identifier of the object
	Type      string    // file extension of the avatar
	Sizes     Sizes     // list of available sizes
	CreatedAt time.Time // when the avatar was first created
	UpdatedAt time.Time // last update of the avatar
}

// FindAvatar searches the database for an avatar object
// given a hash string
func FindAvatar(db *sql.DB, hash string) *Avatar {
	var imgType string
	var sizeString string
	var created time.Time
	var updated time.Time

	row := db.QueryRow("SELECT type, sizes, created_at, updated_at FROM images WHERE hash = $1", hash)
	err := row.Scan(&imgType, &sizeString, &created, &updated)
	if err != nil {
		log.Print(err)
		return nil
	}

	// Parse JSON sizes array
	sizes := make(Sizes, 0)
	json.Unmarshal([]byte(sizeString), &sizes)

	return &Avatar{Hash: hash, Type: imgType, Sizes: sizes, CreatedAt: created, UpdatedAt: updated}
}

// Save the avatar to the database
func (a Avatar) Save(db *sql.DB) error {

	// Convert sizes to JSON for storage in the database
	sizes, err := json.Marshal(a.Sizes)
	if err != nil {
		return err
	}

	res, err := db.Exec("UPDATE images SET sizes = $1, type = $2, updated_at = $3 WHERE hash = $4",
		string(sizes),
		a.Type,
		time.Now(),
		a.Hash,
	)
	if err != nil {
		return err
	}

	if num, _ := res.RowsAffected(); num == 0 {
		_, err := db.Exec("INSERT INTO images (sizes, type, updated_at, hash) VALUES ($1, $2, $3, $4)",
			string(sizes),
			a.Type,
			time.Now(),
			a.Hash,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// DefaultAvatar returns the default placeholder object for no avatar
func DefaultAvatar(ctx *Application) *Avatar {
	return &ctx.Config.DefaultAvatar
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

// UploadAvatar resizes and uploads the avatar to S3
func UploadAvatar() {

}
