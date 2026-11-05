package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"xcpng-center/internal/app"
	"xcpng-center/internal/config"
	"xcpng-center/internal/xapi"
)

func init() {
	server := &cobra.Command{Use: "server", Short: "The server list"}
	server.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List the configured servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := loadSettings()
			if err != nil {
				return err
			}
			if flags.json {
				return printJSON(s.Servers)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
			fmt.Fprintln(w, "HOST\tPORT\tUSER\tNAME")
			for _, e := range s.Servers {
				fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", e.Host, e.Port, e.User, e.Name)
			}
			return w.Flush()
		},
	})
	server.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "Connect to --host and add it to the server list (its certificate is trusted on the way)",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := loadSettings()
			if err != nil {
				return err
			}
			conn, err := connect(cmd.Context(), s)
			if err != nil {
				return err
			}
			defer conn.Close()
			if err := s.AddServer(config.Server{Host: conn.Host, Port: conn.Port, User: conn.User}); err != nil {
				return err
			}
			fmt.Printf("added %s\n", xapi.HostKey(conn.Host, conn.Port))
			return nil
		},
	})
	rootCmd.AddCommand(server)

	for _, kind := range app.Kinds {
		kind := kind
		name := kind[:len(kind)-1] // vms -> vm
		group := &cobra.Command{Use: name, Short: "Commands for " + kind}
		group.AddCommand(&cobra.Command{
			Use:   "list",
			Short: "List the " + kind + " of --host",
			RunE: func(cmd *cobra.Command, args []string) error {
				s, err := loadSettings()
				if err != nil {
					return err
				}
				conn, err := connect(cmd.Context(), s)
				if err != nil {
					return err
				}
				defer conn.Close()
				objs, err := app.ListObjects(conn, kind)
				if err != nil {
					return err
				}
				for i := range objs { // server data never echoes the password to the terminal
					objs[i].Name, objs[i].Description, objs[i].Host = redact(objs[i].Name), redact(objs[i].Description), redact(objs[i].Host)
				}
				if flags.json {
					return printJSON(objs)
				}
				w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
				fmt.Fprintln(w, "NAME\tUUID\tPOWER\tHOST")
				for _, o := range objs {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", o.Name, o.UUID, o.PowerState, o.Host)
				}
				return w.Flush()
			},
		})
		rootCmd.AddCommand(group)
	}
}
