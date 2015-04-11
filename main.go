package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gocraft/web"
)

// Application holds all the context of the app to be passed around
type Application struct {
	DB     *sql.DB
	Config Configuration
	Debug  bool
}

func (ctx *Application) index(rw web.ResponseWriter, req *web.Request) {
	WriteDocs(rw)
}

func (ctx *Application) read(rw web.ResponseWriter, req *web.Request) {
	hash := req.PathParams["hash"]
	backup := req.PathParams["backup"]
	size := CheckAvatarSize(req.PathParams["size"])

	avatar := FindAvatar(ctx.DB, hash)

	if avatar == nil {
		// Check for backup avatar
		if len(backup) > 0 {
			avatar = FindAvatar(ctx.DB, backup)
		}
	}

	// Do default fallback to something
	if avatar == nil {
		avatar = DefaultAvatar(ctx)
		if len(avatar.Hash) == 0 || len(avatar.Sizes) == 0 {
			rw.WriteHeader(404)
			return
		}
	}

	size = avatar.BestSize(size)
	path := avatar.GetURL(size, ctx.Config.AWS.Bucket)

	rw.Header().Set("Location", path)
	rw.Header().Set("Last-Modified", avatar.UpdatedAt.Format(time.RFC822))
	rw.WriteHeader(302)
}

func (ctx *Application) exists(rw web.ResponseWriter, req *web.Request) {
	hash := req.PathParams["hash"]
	avatar := FindAvatar(ctx.DB, hash)
	if avatar == nil {
		// http/net package does not support a response body for HEAD requests. :(
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (ctx *Application) options(rw web.ResponseWriter, req *web.Request) {
	hash := req.PathParams["hash"]
	avatar := FindAvatar(ctx.DB, hash)
	if avatar == nil {
		WriteResponse(rw, Response{"msg": "no matching avatar found", "hash": hash}, http.StatusNotFound)
		return
	}
	methods := []string{"HEAD", "GET", "POST", "DELETE"}
	WriteResponse(rw, Response{"methods": methods, "hash": req.PathParams["hash"]})
}

func (ctx *Application) write(rw web.ResponseWriter, req *web.Request) {
	// TODO: Allow upload to S3
	file, ext, err := GetUploadedFile(req)
	if err != nil {
		WriteResponse(rw, Response{"msg": "please include an `avatar` file"}, http.StatusBadRequest)
		return
	}

	newAvatar := Avatar{
		Hash: req.PathParams["hash"],
		Type: ext,
	}
	data, err := ProcessImageUpload(ctx, newAvatar, *file)

	for size := range data {
		newAvatar.Sizes = append(newAvatar.Sizes, size)
	}

	newAvatar.Save(ctx.DB)

	if err != nil {
		WriteResponse(rw, Response{"msg": err.Error()}, http.StatusBadRequest)
		return
	}

	// TODO: Response with the appropriate JSON body
}

func (ctx *Application) delete(rw web.ResponseWriter, req *web.Request) {
	// TODO: Delete avatar from S3
}

func main() {

	// maximize CPU usage for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := flag.String("port", "3000", "Port number for process")
	host := flag.String("host", "localhost", "Host for binding")
	debug := flag.Bool("debug", false, "Disables security and increases logging")
	flag.Parse()

	config := LoadConfig("config.json")

	db, err := config.DB.Connect()
	if err != nil {
		log.Fatalln("Please make sure Postgres is installed and configured.", err)
	}

	_, err = db.Exec("SELECT hash FROM images LIMIT 1")
	if err != nil {
		log.Fatalln("Please make sure Postgres is configured.", err)
	}

	router := GetRouter(Application{DB: db, Config: config, Debug: *debug})

	if *debug {
		fmt.Println("Debugging mode enabled.")
	}

	fmt.Println("Running server on " + *host + ":" + *port)
	http.ListenAndServe(*host+":"+*port, router)
}
