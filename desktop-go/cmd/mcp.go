package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"xcpng-center/internal/mcp"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "mcp",
		Short: "Serve the app's operations to AI agents over MCP (stdio); read-only in this version",
		Long: "Serves MCP on stdin/stdout. The servers come from the server list (server add); the password from " +
			"XCP_PASS (the agent never sees it). Certificates must already be trusted, or given in --trust-fingerprint.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := loadSettings()
			if err != nil {
				return err
			}
			srv := mcp.New(s, flags.trustFingerprint, mcp.Credentials{User: flags.user, Password: os.Getenv("XCP_PASS")})
			return srv.Run(cmd.Context())
		},
	})
}
