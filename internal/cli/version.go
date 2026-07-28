package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Freelunch-AI/freelunch-ide/internal/buildinfo"
)

func newVersionCommand() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the freelunch version and build metadata",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			info := buildinfo.Get()
			out := cmd.OutOrStdout()

			if asJSON {
				enc := json.NewEncoder(out)
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}

			_, err := fmt.Fprintf(out,
				"freelunch %s\n  commit:  %s\n  built:   %s\n  go:      %s\n  platform: %s\n",
				info.Version, info.Commit, info.Date, info.GoVersion, info.Platform)
			return err
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "emit build metadata as JSON")

	return cmd
}
