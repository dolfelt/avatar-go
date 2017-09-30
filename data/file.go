package data

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

// FileData describes the file
type FileData map[string]string

// GetFileExt gets the formatted extension of the supported image file
func GetFileExt(file io.ReadSeeker) (string, error) {
	buff := make([]byte, 512) // 512 bytes because --> http://golang.org/pkg/net/http/#DetectContentType
	_, err := file.Read(buff)

	defer file.Seek(0, 0)

	if err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	switch filetype {
	case "image/jpeg", "image/jpg":
		return "jpg", nil

	case "image/gif":
		return "gif", nil

	case "image/png":
		return "png", nil
	default:
		return "", &AppError{"not a supported image file"}
	}
}

// GetUploadedFile returns the file that was attempted to be uploaded
func GetUploadedFile(c *gin.Context) (io.ReadSeeker, string, error) {
	var file io.ReadSeeker
	file, _, err := c.Request.FormFile("avatar")
	if err != nil {
		b, errb := c.GetRawData()
		if errb != nil {
			return nil, "", err
		}
		file = bytes.NewReader(b)
	}

	ext, err := GetFileExt(file)
	if err != nil {
		return nil, "", err
	}

	return file, ext, nil
}

// ProcessImageUpload processes uploaded images into the appropriate size
func ProcessImageUpload(app *Application, avatar Avatar, file io.ReadSeeker) (FileData, error) {
	files := make(FileData, 0)

	if img, _, err := image.Decode(file); err == nil {
		file.Seek(0, 0)
		config, _, _ := image.DecodeConfig(file)

		// Find the max square size we can make the avatar
		maxSize := MinInt(config.Width, config.Height)

		// Loop through all the sizes and create the avatars
		for size, pixels := range DefaultSizes {
			if maxSize < pixels && size != "small" {
				if app.Debug {
					log.Println("Skipping size:", size)
				}
				continue
			}
			data := imaging.Thumbnail(img, pixels, pixels, imaging.CatmullRom)

			buf := new(bytes.Buffer)
			switch avatar.Type {
			case "jpg":
				jpeg.Encode(buf, data, &jpeg.Options{Quality: 80})
				break
			case "png":
				png.Encode(buf, data)
				break
			case "gif":
				gif.Encode(buf, data, &gif.Options{NumColors: 256})
				break
			}

			if path, errs := uploadImageS3(app, avatar, buf.Bytes(), size); errs == nil {
				files[size] = path
			} else if app.Debug {
				log.Println("Error uploading", size, errs)
			}
		}
	} else {
		return nil, err
	}

	return files, nil
}

func getS3Bucket() *s3.Bucket {
	auth := aws.Auth{
		AccessKey: viper.GetString("AWSKey"),
		SecretKey: viper.GetString("AWSSecret"),
	}
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket(viper.GetString("AWSBucket"))

	return bucket
}

func uploadImageS3(app *Application, avatar Avatar, data []byte, size string) (string, error) {

	bucket := getS3Bucket()

	err := bucket.Put(avatar.GetPath(size), data, "image/"+avatar.Type, s3.PublicRead)
	if err != nil {
		return "", err
	}

	if app.Debug {
		log.Println("Uploaded size ", size, "to", avatar.GetPath(size))
	}

	return avatar.GetURL(size, viper.GetString("AWSBucket")), nil
}

// ClearAvatarFiles removes all unneeded files from S3
func ClearAvatarFiles(avatar Avatar) error {
	bucket := getS3Bucket()

	sizes := avatar.Sizes

	for _, size := range sizes {
		path := avatar.GetPath(size)
		err := bucket.Del(path)
		if err != nil {
			return err
		}
	}

	return nil
}
