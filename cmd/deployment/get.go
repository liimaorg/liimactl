package deployment

import (
	"fmt"

	"github.com/liimaorg/liimactl/cmd"
	"github.com/spf13/cobra"
)

func newGetCommand(cli *cmd.Cli) *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get deployments",
		Run: func(cmd *cobra.Command, args []string) {
			runGet(cli, args)
		},
	}

	return getCmd
}

func runGet(cli *cmd.Cli, args []string) {
	// TODO: Work your own magic here
	fmt.Printf("get called, client: %+v", cli.Client)
}
