package deployment

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get deployments",
		Run:   runGet,
	}

	return getCmd
}

func runGet(cmd *cobra.Command, args []string) {
	// TODO: Work your own magic here
	fmt.Println("get called")
}
