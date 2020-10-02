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
		test, err := command.Flags().GetString("test")
		if err != nil {
			return err
		} 
			
		fmt.Printf("Value of test %s\n", test)
		return nil
	},
}

func init() {
	rootCommand.PersistentFlags().StringP("test", "t", "AAAA", "Used only for test")
}

// Execute Used for execute root command
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Printf("Error execute command %v\n", err)
		os.Exit(1)
	}
}