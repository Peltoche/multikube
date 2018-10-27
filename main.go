package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// Set the flag variables.
	var matchsOpt []string

	// Define the cli.
	rootCmd := cobra.Command{
		Use:   "mkube",
		Short: "Mkube runs kubernetes commands on several contexts at the same time",

		// Run only parse the options and trigger the Execute function which
		// contains all the stuff. It allows to isolate the cli logic.
		Run: func(cmd *cobra.Command, args []string) {
			err := Execute(&Cmd{
				Command: args,
				Matchs:  matchsOpt,
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	// Link the flag variables to the cli.
	rootCmd.PersistentFlags().StringArrayVarP(&matchsOpt, "match", "m", []string{}, "bar")

	// Run the cli.
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
