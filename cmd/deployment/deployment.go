package deployment

import (
	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

// NewDeploymentCmd represents the deployment command
func NewDeploymentCmd(cli *client.Cli) *cobra.Command {

	var DeploymentCmd = &cobra.Command{
		Use:   "deployment COMMAND",
		Short: "Manage deployments",
	}

	DeploymentCmd.AddCommand(newGetCommand(cli))
	DeploymentCmd.AddCommand(newCreateCommand(cli))

	return DeploymentCmd
}
