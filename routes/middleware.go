package routes

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CORSHeaders allows for cross-origin requests against the API
func CORSHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Cache-Control")
		c.Header("Access-Control-Allow-Methods", "GET,HEAD,POST,DELETE,OPTIONS")

		c.Next()
	}
}

// AuthRequired detects if a JWT token has been sent with the request and
// validates the token before completing the request.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Load the token from the Query String
		token := c.Query("token")

		// Load the hash from the URL
		hash := c.Param("hash")

		if len(token) == 0 {
			// Accept the token in the post body as well
			token = c.PostForm("token")
			if len(token) == 0 {
				c.JSON(400, gin.H{"msg": "no token was found"})
				return
			}
		}

		auth, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			if t.Claims["hash"] == hash {
				// TODO: Get signing key from ENV
				return viper.GetString("JwtKey"), nil
			}
			return nil, fmt.Errorf("signed hash does not match: %v", t.Claims["hash"])
		})

		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}

		if !auth.Valid {
			c.JSON(400, gin.H{"msg": "token was invalid for unknown reason"})
		}

		c.Next()
	}
}
