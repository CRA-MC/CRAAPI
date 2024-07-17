package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionPrint() {
	fmt.Println("CRAAPI developing version: v0.3.0001")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of craapi",
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		versionPrint()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
