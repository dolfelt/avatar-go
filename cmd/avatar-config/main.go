package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	app "github.com/dolfelt/avatar-go"
	_ "github.com/lib/pq"
)

const (
	colorRed   = "\033[0;31m"
	colorGreen = "\033[0;32m"
	colorGray  = "\033[0;37m"
	colorClear = "\033[0m"
)

const (
	sqlCreateTable = `CREATE TABLE images (
      hash         varchar(40) NOT NULL,
      type         char(4) NOT NULL,
      sizes        json NOT NULL,
      updated_at   timestamp DEFAULT current_timestamp,
      created_at   timestamp DEFAULT current_timestamp,
      CONSTRAINT hash UNIQUE(hash)
    );`
)

func setupTables(db *sql.DB) {

	// Detect if the table exists
	// Create it if it doesn't
	_, err := db.Exec("SELECT hash FROM images LIMIT 1")
	if err != nil {
		fmt.Println(colorGray+"Cannot find table \"images\". Creating...", colorClear)

		res, tblErr := db.Exec(sqlCreateTable)
		if tblErr != nil {
			fmt.Println(colorRed+"Problem creating tables:", tblErr, colorClear)
			os.Exit(1)
		}
		fmt.Println(colorGreen+"Table \"images\" created.", res, colorClear)
	}

	fmt.Println(colorGreen+"Setup Complete!", colorClear)

}

func setDefaultAvatar(db *sql.DB, path string, config app.Configuration, debug bool) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(colorRed+"Cannot open file: ", err, colorClear)
		return
	}

	ext, err := app.GetFileExt(file)
	if err != nil {
		fmt.Println(colorRed+"Cannot read file properly: ", err, colorClear)
		return
	}

	application := &app.Application{DB: db, Config: config, Debug: debug}

	sha := sha1.New()
	sha.Write([]byte("default"))
	hash := hex.EncodeToString(sha.Sum(nil))

	newAvatar := app.Avatar{
		Hash:      hash,
		Type:      ext,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Println("Uploading default avatar...")

	data, err := app.ProcessImageUpload(application, newAvatar, file)
	if err != nil {
		fmt.Println(colorRed+"Error uploading: ", err, colorClear)
		return
	}
	if len(data) == 0 {
		fmt.Println(colorRed+"No images uploaded. File must be at least 128px. ", colorClear)
		return
	}

	for size := range data {
		newAvatar.Sizes = append(newAvatar.Sizes, size)
	}

	config.DefaultAvatar = newAvatar
	saveConfig(config)

	fmt.Println(colorGreen+"Upload Complete!", colorClear)
}

func saveConfig(config app.Configuration) {

	j, jerr := json.MarshalIndent(config, "", "  ")
	if jerr != nil {
		fmt.Println(colorRed+"Error Parsing JSON: ", jerr, colorClear)
		return
	}

	file, err := os.Create("config.json")
	if err != nil {
		fmt.Println(colorRed+"Error Opening config.json: ", err, colorClear)
		return
	}

	_, werr := file.Write(j)
	if werr != nil {
		fmt.Println(colorRed+"Error Writing Config: ", werr.Error(), colorClear)
		return
	}

	fmt.Println("Configuration file created! [at ./config.json]")
}

func main() {

	debug := flag.Bool("debug", false, "Disables security and increases logging")
	flag.Parse()

	args := flag.Args()

	filename := "config.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = "config.example.json"
	}
	config := app.LoadConfig(filename)

	defer saveConfig(config)

	db, err := config.DB.Connect()

	if err != nil {
		fmt.Println(colorRed+"Error loading connection: ", err, colorClear)
		return
	}

	// Check if the Database connection is valid
	// If it isn't, make sure to throw an error
	err = db.Ping()
	if err != nil {

		if !strings.Contains(err.Error(), config.DB.Database) {
			fmt.Println(colorRed+"Please check your connection configuration in config.json.", err, colorClear)
		} else {
			fmt.Println(colorRed+"Cannot connect to database \""+config.DB.Database+"\". Please create it.", err, colorClear)
		}
		return
	}

	var command string
	if len(args) == 0 {
		command = "install"
	} else {
		command = args[0]
	}

	if command == "install" {
		setupTables(db)
		return
	}

	if command == "setdefault" {
		if len(args) < 2 {
			log.Fatal(colorRed+"Please include a file path to upload.", colorClear)
		}
		setDefaultAvatar(db, args[1], config, *debug)
		return
	}
}
