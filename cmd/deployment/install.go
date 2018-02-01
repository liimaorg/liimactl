package deployment

import (
	"fmt"
	"log"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentInstallLong = fmt.Sprintf(` 
    Install deployment on an environemt with the use of specific properties.`)

	//Example command description
	deploymentInstallExample = fmt.Sprintf(` 
    # Install multiple deplyoments on an environment with specific properties. 
	liimactl.exe deployment install --environment=I  --fromEnvironment=V
	liimactl.exe deployment install --environment=I  --fromEnvironment=V --date="2018-02-01 17:00"
	liimactl.exe deployment install --environment=I  --fromEnvironment=V --date="2018-02-01 17:00" --wait`)

	//Flags of the command
	commandOptionsInstall client.CommandOptionsInstallDeployment
)

//newInstallCommand is a command to install multiple deployments on an environment
func newInstallCommand(cli *client.Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "install [flags] ",
		Short:   "install deployment",
		Long:    deploymentInstallLong,
		Example: deploymentInstallExample,
		Run: func(cmd *cobra.Command, args []string) {
			runInstall(cmd, cli, args)
		},
	}

	cmd.Flags().StringVarP(&commandOptionsInstall.Environment, "environment", "e", "", "Environment")
	cmd.Flags().StringVarP(&commandOptionsInstall.FromEnvironment, "fromEnvironment", "f", "", "Deploy last deplyoment from given environment")
	cmd.Flags().StringVarP(&commandOptionsInstall.DeploymentDate, "date", "d", "", "Deployment Date 'DD.MM.YYYY hh:mm' ")
	cmd.Flags().BoolVarP(&commandOptionsInstall.ExecuteShakedownTest, "executeShakeDownTest", "s", false, "Run Shakedowntest after the deplyoment")
	cmd.Flags().BoolVarP(&commandOptionsInstall.Wait, "wait", "w", false, "Wait maxWaitTime until the deplyoment success or failed")
	cmd.Flags().IntVarP(&commandOptionsInstall.MaxWaitTime, "maxWaitTime", "t", 600, "Max Wait time [seconds] until the deplyoment success or failed")
	cmd.Flags().StringVarP(&commandOptionsInstall.Blacklist, "blacklist", "b", "", "Blacklist with all appServer, which should not be deployed [filepath]")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runInstall(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Install deplyoment
	deplyoments, err := client.InstallDeployment(cli, &commandOptionsInstall)
	if err != nil {
		log.Fatal("Error Install Deployment: ", err)
	}

	for _, deployment := range deplyoments {

		//Print result
		PrintDeployment(cmd, deployment)

		//Write error failed -> return code = 1 with log.Fatal
		if deployment.State == client.DeploymentStateFailed {
			log.Fatal("Deployment failed with state: ", deployment.State)
		}

	}

}
