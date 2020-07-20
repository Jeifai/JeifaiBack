package cmd

import (
    "github.com/spf13/cobra"
    matcher "github.com/Jeifai/JeifaiBack/cmd/matcher"
)

// matchCmd represents the match command
var matchCmd = &cobra.Command{
	Use:   "match",
	Short: "Run the matcher",
    Long: `Run the matcher for specific targets or for all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
        matcher.TempMatch(company)
	},
}

func init() {
	rootCmd.AddCommand(matchCmd)
	scrapeCmd.Flags().StringVarP(&company, "matchcompany", "m", "", "Specify a company or all of them")
}
