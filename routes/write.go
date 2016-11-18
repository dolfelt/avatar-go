package routes

import (
	"net/http"

	"github.com/dolfelt/avatar-go/data"
	"github.com/gin-gonic/gin"
)

func write(app *data.Application) gin.HandlerFunc {

	return func(c *gin.Context) {
		// TODO: Allow upload to S3
		file, ext, err := data.GetUploadedFile(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "please include an `avatar` file"})
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
		data, err := data.ProcessImageUpload(app, newAvatar, *file)

		for size := range data {
			newAvatar.Sizes = append(newAvatar.Sizes, size)
		}

		newAvatar.Save(app.DB)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		// TODO: Response with the appropriate JSON body
	}
}

func delete(app *data.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		oldAvatar := data.FindAvatar(app.DB, c.Param("hash"))
		if oldAvatar != nil {
			err := data.ClearAvatarFiles(*oldAvatar)
			if err != nil {
				c.JSON(http.StatusGatewayTimeout, gin.H{"msg": err.Error()})
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.AbortWithStatus(http.StatusNotFound)
	}
}
