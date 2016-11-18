package routes

/*
Handles all the rounting and web requests coming in
*/

import (
	"github.com/dolfelt/avatar-go/data"
	"github.com/gin-gonic/gin"
)

// Register creates and registers all the routes to the router
func Register(app *data.Application) *gin.Engine {
	router := gin.New()

	return configRouter(router, app)
}

func configRouter(router *gin.Engine, app *data.Application) *gin.Engine {
	// GetRouter returns the router with all the routes defined
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSHeaders())

	// Get endpoints for displaying the avatar
	router.GET("/:hash", read(app))
	router.GET("/:hash/:size_or_backup", read(app))
	router.GET("/:hash/:size_or_backup/:size", read(app))
	// router.GET("/:hash/:backup/:size", read(app))
	// router.GET("/:hash:[0-9a-f]{40}/:backup:[0-9a-f]{40}/:size", read(app))
	router.GET("/", index())

	// Head endpoint for determining if the avatar exists
	router.HEAD("/:hash", exists(app))

	authRouter := router.Group("/")
	if !app.Debug {
		authRouter.Use(AuthRequired())
	}

	// Options endpoint for available methods
	authRouter.OPTIONS("/:hash", options(app))
	authRouter.POST("/:hash", write(app))
	authRouter.DELETE("/:hash", delete(app))

	return router
}
