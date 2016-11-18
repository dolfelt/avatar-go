package cmd

import "github.com/spf13/cobra"

// RootCmd is the entry point for all other commands
var RootCmd = &cobra.Command{
	Use:   "avatar",
	Short: "Avatar is a simple service for serving and managing avatars",
	Long: `A simple avatar service built with love by dolfelt and friends in Go.
         Complete documentation is available at http://github.com/dolfelt/avatar-go`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
