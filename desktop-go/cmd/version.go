package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"xcpng-center/internal/buildinfo"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.json {
				return printJSON(map[string]string{"version": buildinfo.Version, "commit": buildinfo.Commit, "go": runtime.Version()})
			}
			fmt.Printf("xcp-ng-center-go %s (%s, %s)\n", buildinfo.Version, runtime.Version(), buildinfo.Commit)
			return nil
		},
	})
}
