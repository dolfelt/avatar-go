package data

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Application holds all the info for the app
type Application struct {
	DB    DB
	Debug bool
}

// LoadConfig loads external configuration file
func LoadConfig(path string) error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("avatar")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
		return fmt.Errorf("Unable to locate config file. (%s)\n", err)
	}

	loadDefaultSettings()

	return nil
}

func loadDefaultSettings() {
	// PostgreSQL Config
	viper.SetDefault("DBUser", "docker")
	viper.SetDefault("DBPassword", "docker")
	viper.SetDefault("DBDatabase", "avatars")
	viper.SetDefault("DBHost", "localhost")
	viper.SetDefault("DBPort", 5432)

	// DynamoDB Config
	viper.SetDefault("DynamoRegion", "us-east-1")

	viper.SetDefault("Port", 3000)
	viper.SetDefault("Debug", false)
	viper.SetDefault("TableName", "avatars")

	viper.SetDefault("Store", "postgres")
}

// DefaultSizes is the list of sizes
var DefaultSizes = map[string]int{
	"small":    128,
	"medium":   256,
	"large":    512,
	"original": 1024,
}

// DefaultSizeKeys is the list of acceptable size strings
func DefaultSizeKeys() []string {
	keys := make([]string, 0, len(DefaultSizes))
	for k := range DefaultSizes {
		keys = append(keys, k)
	}
	return keys
}

// AppError handles generic application errors
type AppError struct {
	msg string
}

func (e *AppError) Error() string {
	return e.msg
}

// MinInt is a shim for determining the minimum for integers rather than floats
func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Retry allows you to retry commands more than once
func Retry(attempts int, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return nil
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(2 * time.Second)

		log.Println("retrying...", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
