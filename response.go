package main

import (
	"encoding/json"
	"fmt"

	"github.com/gocraft/web"
)

// Response holds the API response
type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

// WriteResponse writes the JSON API response to the ResponseWriter with
// the appropriate headers
func WriteResponse(rw web.ResponseWriter, data Response, code ...int) {
	rw.Header().Set("Content-Type", "application/json")
	if len(code) > 0 {
		rw.WriteHeader(code[0])
	}
	fmt.Fprint(rw, data)
}

// WriteDocs Writes the API Documentation for this avatar API.
func WriteDocs(rw web.ResponseWriter) {
	docs := Response{
		"links": Response{
			"avatar.exists": Response{
				"type":   "endpoint",
				"href":   "/:hash",
				"method": "HEAD"},
			"avatar.read": Response{
				"type":     "endpoint",
				"href":     "/:hash/:backup/:size",
				"method":   "GET",
				"optional": []string{":size", ":backup"}},
			"avatar.write": Response{
				"type":   "endpoint",
				"href":   "/:hash",
				"method": "POST"},
			"avatar.delete": Response{
				"type":   "endpoint",
				"href":   "/:hash",
				"method": "DELETE"}},
		"meta": Response{
			"parameters": Response{
				":hash": Response{
					"desc": "sha1 hash of the prefixed user id"},
				":backup": Response{
					"desc": "another sha1 hash to use if the given :hash does not exist"},
				":size": Response{
					"desc":    "one of the possible sizes",
					"note":    "if the requested size is not available, the next largest size will be used",
					"choices": DefaultSizeKeys(),
					"default": "medium",
					"sizes":   DefaultSizes}}}}

	WriteResponse(rw, docs)
}
