package main

import (
	"fmt"
	"os"

	"github.com/dolfelt/avatar-go/cmd"
	_ "github.com/lib/pq"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
