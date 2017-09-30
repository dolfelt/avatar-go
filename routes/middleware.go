package routes

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func getCORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cache-Control"}
	return corsConfig
}

// AuthRequired detects if a JWT token has been sent with the request and
// validates the token before completing the request.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Load the token from the Query String
		token := c.Query("token")
		if len(token) == 0 {
			// Accept the token in the post body as well
			token = c.PostForm("token")
		}
		if len(token) == 0 {
			// Accept the token in the Authorization header as well
			token = c.Request.Header.Get("Authorization")
		}
		if len(token) == 0 {
			c.JSON(400, gin.H{"msg": "no token was found"})
			c.AbortWithStatus(400)
			return
		}

		// Load the hash from the URL
		hash := c.Param("hash")

		auth, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			if t.Claims["hash"] == hash {
				return []byte(viper.GetString("JwtKey")), nil
			}
			return nil, fmt.Errorf("signed hash does not match: %v", t.Claims["hash"])
		})

		if err != nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Problem with token: %s", err.Error())})
			c.Abort()
			return
		}

		if !auth.Valid {
			c.JSON(400, gin.H{"error": "token was invalid for unknown reason"})
			c.Abort()
		}

		c.Next()
	}
}
