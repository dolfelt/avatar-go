package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Avatar",
	Long:  `All software has versions. This is Avatar's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Avatar Go v0.1 -- HEAD")
	},
}
