package deployment

import (
	"time"

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
	DeploymentCmd.AddCommand(newGetFilterCommand(cli))
	DeploymentCmd.AddCommand(newCreateCommand(cli))
	DeploymentCmd.AddCommand(newInstallCommand(cli))

	return DeploymentCmd
}

//PrintDeployment prints out the properties of a DeploymentResponse
func PrintDeployment(cmd *cobra.Command, deployment client.DeploymentResponse) {

	//Print result
	cmd.Println("------")

	if deployment.AppServerName != "" {
		cmd.Printf("%s ", deployment.AppServerName)
	}
	if deployment.EnvironmentName != "" {
		cmd.Printf("%s ", deployment.EnvironmentName)
	}
	if deployment.ReleaseName != "" {
		cmd.Printf("%s ", deployment.ReleaseName)
	}
	if deployment.DeploymentDate != 0 {
		cmd.Printf("%s ", time.Unix(0, deployment.DeploymentDate*int64(time.Millisecond)).Format("2006-01-02T15:04"))
	}
	if deployment.State != "" {
		cmd.Println(deployment.State)
	}
	for _, appsWithVersion := range deployment.AppsWithVersion {
		cmd.Printf("%s ", appsWithVersion.ApplicationName)
		cmd.Println(appsWithVersion.Version)
	}

}
