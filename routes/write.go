package routes

import (
	"net/http"
	"time"

	"github.com/dolfelt/avatar-go/data"
	"github.com/gin-gonic/gin"
)

func write(app *data.Application) gin.HandlerFunc {

	return func(c *gin.Context) {
		// TODO: Allow upload to S3
		file, ext, err := data.GetUploadedFile(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "please include an `avatar` file"})
			return
		}

		hash := c.Param("hash")

		oldAvatar := data.FindAvatar(app.DB, hash)
		if oldAvatar != nil {
			data.ClearAvatarFiles(*oldAvatar)
		}

		newAvatar := data.Avatar{
			Hash: hash,
			Type: ext,
		}
		data, err := data.ProcessImageUpload(app, newAvatar, file)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for size := range data {
			newAvatar.Sizes = append(newAvatar.Sizes, size)
		}

		err = newAvatar.Save(app.DB)

		now := time.Now()
		newAvatar.UpdatedAt = now
		newAvatar.CreatedAt = now

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"data":  newAvatar,
			"error": nil,
		})
	}
}

func delete(app *data.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		oldAvatar := data.FindAvatar(app.DB, c.Param("hash"))
		if oldAvatar != nil {
			err := data.ClearAvatarFiles(*oldAvatar)
			if err != nil {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": err.Error()})
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.AbortWithStatus(http.StatusNotFound)
	}
}
