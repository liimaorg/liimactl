package deployment

import (
	"github.com/liimaorg/liimactl/cmd"
	"github.com/spf13/cobra"
)

// NewDeploymentCmd represents the deployment command
func NewDeploymentCmd(cli *cmd.Cli) *cobra.Command {

	var DeploymentCmd = &cobra.Command{
		Use:   "deployment COMMAND",
		Short: "Manage deployments",
	}

	DeploymentCmd.AddCommand(newGetCommand(cli))

	return DeploymentCmd
}
