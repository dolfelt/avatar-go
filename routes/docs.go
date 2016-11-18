package routes

import (
	"github.com/dolfelt/avatar-go/data"
	"github.com/gin-gonic/gin"
)

func index() gin.HandlerFunc {
	return func(c *gin.Context) {
		docs := gin.H{
			"links": gin.H{
				"avatar.exists": gin.H{
					"type":   "endpoint",
					"href":   "/:hash",
					"method": "HEAD",
				},
				"avatar.read": gin.H{
					"type":     "endpoint",
					"href":     "/:hash/:backup/:size",
					"method":   "GET",
					"optional": []string{":size", ":backup"},
				},
				"avatar.write": gin.H{
					"type":   "endpoint",
					"href":   "/:hash",
					"method": "POST",
				},
				"avatar.delete": gin.H{
					"type":   "endpoint",
					"href":   "/:hash",
					"method": "DELETE",
				},
			},
			"meta": gin.H{
				"parameters": gin.H{
					":hash": gin.H{
						"desc": "sha1 hash of the prefixed user id",
					},
					":backup": gin.H{
						"desc": "another sha1 hash to use if the given :hash does not exist",
					},
					":size": gin.H{
						"desc":    "one of the possible sizes",
						"note":    "if the requested size is not available, the next largest size will be used",
						"choices": data.DefaultSizeKeys(),
						"default": "medium",
						"sizes":   data.DefaultSizes,
					},
				},
			},
		}

		c.JSON(200, docs)
	}
}
