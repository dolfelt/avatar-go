package cmd

import (
	"fmt"
	"log"

	"github.com/dolfelt/avatar-go/data"
	"github.com/dolfelt/avatar-go/routes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serveCmd.Flags().BoolP("debug", "d", false, "disables security and increases logging")
	serveCmd.Flags().StringP("port", "p", "3000", "choose a custom port")
	serveCmd.Flags().StringP("addr", "a", "", "address to bind this service to")

	viper.BindPFlag("Port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("Debug", serveCmd.Flags().Lookup("debug"))
	viper.BindPFlag("IPAddress", serveCmd.Flags().Lookup("addr"))

	RootCmd.AddCommand(serveCmd)
}

func serveLoadConfig() *data.Application {
	configErr := data.LoadConfig("config")
	if configErr != nil {
		log.Println(configErr)
	}
	db := &data.PostgresDB{}
	err := data.Retry(10, func() error {
		var err error
		err = db.Connect()
		return err
	})

	if err != nil {
		log.Fatalln("Please make sure Postgres is installed and configured.", err)
	}

	return &data.Application{
		DB:    db,
		Debug: viper.GetBool("Debug"),
	}
}

func serveRun(cmd *cobra.Command, args []string) {
	app := serveLoadConfig()

	if app.Debug {
		fmt.Println("Debugging mode enabled.")
	}

	app.DB.Migrate()

	router := routes.Register(app)

	router.Run(viper.GetString("IPAddress") + ":" + viper.GetString("Port"))
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start up the avatar service",
	Long:  `Runs Avatar Go and connects to the PostgreSQL database`,
	Run:   serveRun,
}

// _, err = db.Exec("SELECT hash FROM images LIMIT 1")
// if err != nil {
//   log.Fatalln("Please make sure Postgres is configured.", err)
// }

// fmt.Println("Running server on " + *host + ":" + *port)
