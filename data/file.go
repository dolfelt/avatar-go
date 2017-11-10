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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

			if path, errs := uploadImageS3(app, avatar, buf, size); errs == nil {
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

func getBucketName() string {
	return viper.GetString("AWSBucket")
}

func getS3() *s3.S3 {
	awsConfig := &aws.Config{
		Region: aws.String(viper.GetString("AwsBucketRegion")),
	}
	if viper.GetString("AwsKey") != "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(
			viper.GetString("AwsKey"),
			viper.GetString("AwsSecret"),
			"",
		)
	}
	return s3.New(session.Must(session.NewSession()), awsConfig)
}

func uploadImageS3(app *Application, avatar Avatar, data io.Reader, size string) (string, error) {
	up := s3manager.NewUploaderWithClient(getS3())

	upParams := &s3manager.UploadInput{
		Bucket:      aws.String(getBucketName()),
		Key:         aws.String(avatar.GetPath(size)),
		Body:        data,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/" + avatar.Type),
	}
	result, err := up.Upload(upParams)
	if err != nil {
		return "", err
	}

	if app.Debug {
		log.Printf("%#v\n", result)
		log.Println("Uploaded size ", size, "to", avatar.GetPath(size))
	}

	return avatar.GetURL(size, viper.GetString("AWSBucket")), nil
}

// ClearAvatarFiles removes all unneeded files from S3
func ClearAvatarFiles(avatar Avatar) error {
	client := getS3()

	sizes := avatar.Sizes

	for _, size := range sizes {
		path := avatar.GetPath(size)
		_, err := client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(getBucketName()),
			Key:    aws.String(path),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
