/*
Handles all the rounting and web requests coming in
*/
package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gocraft/web"
)

// Middlware custom functions

// APIHeaders injects the headers for accessing the API cross-origin
func (ctx *Application) APIHeaders(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Cache-Control")
	rw.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,POST,DELETE,OPTIONS")

	next(rw, req)
}

// AuthRequired detects if a JWT token has been sent with the request and
// validates the token before completing the request.
func (ctx *Application) AuthRequired(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {

	// Load the token from the Query String
	token := req.URL.Query().Get("token")

	// Load the hash from the URL
	hash := req.PathParams["hash"]

	if len(token) == 0 {
		// Accept the token in the post body as well
		token = req.FormValue("token")
		if len(token) == 0 {
			WriteResponse(rw, Response{"msg": "no token was found"}, 400)
			return
		}
	}

	auth, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		if t.Claims["hash"] == hash {
			return ctx.Config.SigningKey(), nil
		}
		return nil, fmt.Errorf("signed hash does not match: %v", t.Claims["hash"])
	})

	if err != nil {
		WriteResponse(rw, Response{"msg": err.Error()}, 400)
		return
	}

	if !auth.Valid {
		WriteResponse(rw, Response{"msg": "token was invalid for unknown reason"}, 400)
	}

	next(rw, req)
}

// TODO: Added authentication middleware function for certain requests

// GetRouter returns the router with all the routes defined
func GetRouter(context Application) *web.Router {
	router := web.New(context)
	router.Middleware(web.LoggerMiddleware).
		Middleware(web.ShowErrorsMiddleware).
		Middleware(context.APIHeaders)

	// Get endpoints for displaying the avatar
	router.Get("/:hash:[0-9a-f]{40}", context.read).
		Get("/:hash:[0-9a-f]{40}/:size:[a-z0-9]{2,10}", context.read).
		Get("/:hash:[0-9a-f]{40}/:backup:[0-9a-f]{40}", context.read).
		Get("/:hash:[0-9a-f]{40}/:backup:[0-9a-f]{40}/:size", context.read).
		Get("/", context.index)

	// Head endpoint for determining if the avatar exists
	router.Head("/:hash:[0-9a-f]{40}", context.exists)

	authRouter := router.Subrouter(context, "")

	if !context.Debug {
		authRouter.Middleware(context.AuthRequired)
	}

	// Options endpoint for available methods
	authRouter.Options("/:hash:[0-9a-f]{40}", context.options).
		Post("/:hash:[0-9a-f]{40}", context.write).
		Delete("/:hash:[0-9a-f]{40}", context.delete)

	return router
}
