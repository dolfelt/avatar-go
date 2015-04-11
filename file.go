package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gocraft/web"
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
)

type FileData map[string]string

// GetFileExt gets the formatted extension of the supported image file
func GetFileExt(file multipart.File) (string, error) {
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
		return "jpg", nil

	case "image/png":
		return "jpg", nil
	default:
		return "", &AppError{"not a supported image file"}
	}
}

// GetUploadedFile returns the file that was attempted to be uploaded
func GetUploadedFile(req *web.Request) (*multipart.File, string, error) {

	file, _, err := req.FormFile("avatar")
	if err != nil {
		return nil, "", err
	}

	ext, err := GetFileExt(file)
	if err != nil {
		return nil, "", err
	}

	return &file, ext, nil
}

// ProcessImageUpload processes uploaded images into the appropriate size
func ProcessImageUpload(ctx *Application, avatar Avatar, file multipart.File) (FileData, error) {
	files := make(FileData, 0)

	defer file.Close()

	if img, _, err := image.Decode(file); err == nil {
		file.Seek(0, 0)
		config, _, _ := image.DecodeConfig(file)

		// Find the max square size we can make the avatar
		maxSize := MinInt(config.Width, config.Height)

		// Loop through all the sizes and create the avatars
		for size, pixels := range DefaultSizes {
			if maxSize < pixels {
				if ctx.Debug {
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

			if path, err := uploadImageS3(ctx, avatar, buf.Bytes(), size); err == nil {
				files[size] = path
			} else if ctx.Debug {
				log.Println("Error uploading", size, err)
			}
		}
	} else {
		return nil, err
	}

	return files, nil
}

func uploadImageS3(ctx *Application, avatar Avatar, data []byte, size string) (string, error) {
	auth := aws.Auth{
		AccessKey: ctx.Config.AWS.AccessKey,
		SecretKey: ctx.Config.AWS.AccessSecret,
	}
	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket(ctx.Config.AWS.Bucket)

	err := bucket.Put(avatar.GetPath(size), data, "image/"+avatar.Type, s3.PublicRead)
	if err != nil {
		return "", err
	}

	if ctx.Debug {
		log.Println("Uploaded size ", size, "to", avatar.GetPath(size))
	}

	return avatar.GetURL(size, ctx.Config.AWS.Bucket), nil
}
