package routes

import (
	"net/http"
	"time"

	"github.com/dolfelt/avatar-go/data"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func read(app *data.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		sizeOrBackup := c.Param("size_or_backup")
		size := c.Param("size")

		var backup string
		if len(sizeOrBackup) == 40 {
			backup = sizeOrBackup
		} else if len(size) == 0 {
			size = sizeOrBackup
		}
		size = data.CheckAvatarSize(size)

		avatar := data.FindAvatar(app.DB, hash)

		if avatar == nil {
			// Check for backup avatar
			if len(backup) > 0 {
				avatar = data.FindAvatar(app.DB, backup)
			}
		}

		// Do default fallback to something
		if avatar == nil {
			avatar = data.DefaultAvatar(app)
			if len(avatar.Hash) == 0 || len(avatar.Sizes) == 0 {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		size = avatar.BestSize(size)
		path := avatar.GetURL(size, viper.GetString("AWSBucket"))

		c.Header("Location", path)
		c.Header("Last-Modified", avatar.UpdatedAt.Format(time.RFC822))
		c.Status(302)
	}
}

func exists(app *data.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		avatar := data.FindAvatar(app.DB, hash)
		if avatar == nil {
			// http/net package does not support a response body for HEAD requests. :(
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusNoContent)
	}
}

func options(app *data.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		avatar := data.FindAvatar(app.DB, hash)
		if avatar == nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": "no matching avatar found", "hash": hash})
			return
		}
		methods := []string{"HEAD", "GET", "POST", "DELETE"}
		c.JSON(200, gin.H{"methods": methods, "hash": hash})
	}
}
