package deployment

import (
	"fmt"
	"strings"
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
	DeploymentCmd.AddCommand(newGetFilteredDeploymentsCommand(cli))
	DeploymentCmd.AddCommand(newCreateCommand(cli))
	DeploymentCmd.AddCommand(newPromoteCommand(cli))

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

//AskYesNo want's a user confirmation yes from the console
func AskYesNo(message string) bool {
	var s string

	fmt.Printf("%s [y/n]: ", message)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}
