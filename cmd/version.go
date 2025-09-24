package cmd

import (
	"craapi/packages/log"
	"fmt"

	"github.com/spf13/cobra"
)

func versionPrint() {
	if !v {
		log.LOGI("CRAAPI develpopment version: v0.4.0003.0005 BUG FIX")
	} else {
		fmt.Println("CRAAPI develpopment version: v0.4.0003.0005 BUG FIX")
	}
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
