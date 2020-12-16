package awscmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var rootCommand = &cobra.Command{
	Use: "awstools",
	Short: "Amazon cloud operations",
	Long: "Used for running operation against AWS cloud",

	RunE: func(command *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCommand.PersistentFlags().StringP("profile", "p", "default", "Your AWS profile in config file")
	rootCommand.PersistentFlags().StringP("awskey", "k", "", "AWS Key ID")
	rootCommand.PersistentFlags().StringP("awssecret", "s", "", "AWS Secret")
	rootCommand.AddCommand(StsTokenCommand)
}

// Execute Used for execute root command
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Printf("Error execute command %v\n", err)
		os.Exit(1)
	}
}