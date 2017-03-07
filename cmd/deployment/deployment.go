package deployment

import "github.com/spf13/cobra"

// NewDeploymentCmd represents the deployment command
func NewDeploymentCmd() *cobra.Command {

	var DeploymentCmd = &cobra.Command{
		Use:   "deployment COMMAND",
		Short: "Manage deployments",
	}

	DeploymentCmd.AddCommand(newGetCommand())

	return DeploymentCmd
}
