package deployment

import (
	"fmt"
	"log"
	"time"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentCreateLong = fmt.Sprintf(` 
    Create deployment with the use of specific properties.`)

	//Example command description
	deploymentCreateExample = fmt.Sprintf(` 
    # Create a deplyoment with specific properties. 
	liimactl.exe deployment create --appServer=test_application --appName=ch_mobi_app1 --version="1.0.0" --appName=ch_mobi_app2 --version="1.0.1" --environment=I
	liimactl.exe deployment create --appServer=aps_bau --appName=ch_mobi_aps_bau --version="1.0.32" --environment=W --date="22.08.2017 17:00"
	liimactl.exe deployment create --appServer=generic_test --appName=ch_mobi_generic_test --version="1.0.1" --environment=U --wait`)

	//Flags of the command
	commandOptionsCreate client.CommandOptionsCreateDeployment
)

//newCreateCommand is a command to create a deployment
func newCreateCommand(cli *client.Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "create [flags] ",
		Short:   "create deployment",
		Long:    deploymentCreateLong,
		Example: deploymentCreateExample,
		Run: func(cmd *cobra.Command, args []string) {
			runCreate(cmd, cli, args)
		},
	}

	cmd.Flags().StringVarP(&commandOptionsCreate.AppServer, "appServer", "a", "", "Application Server Name")
	cmd.Flags().StringSliceVarP(&commandOptionsCreate.AppName, "appName", "n", []string{}, "Application Name")
	cmd.Flags().StringSliceVarP(&commandOptionsCreate.AppVersion, "version", "v", []string{}, "Application Version")
	cmd.Flags().StringVarP(&commandOptionsCreate.Environment, "environment", "e", "", "Environment")
	cmd.Flags().StringVarP(&commandOptionsCreate.Release, "release", "r", "", "Release")
	cmd.Flags().StringVarP(&commandOptionsCreate.DeploymentDate, "date", "d", "", "Deployment Date 'DD.MM.YYYY hh:mm' ")
	cmd.Flags().BoolVarP(&commandOptionsCreate.ExecuteShakedownTest, "executeShakeDownTest", "s", false, "Run Shakedowntest after the deplyoment")
	cmd.Flags().StringSliceVarP(&commandOptionsCreate.Key, "key", "k", []string{}, "Deploymentparameter Key")
	cmd.Flags().StringSliceVarP(&commandOptionsCreate.Value, "value", "x", []string{}, "Deploymentparameter Value")
	cmd.Flags().BoolVarP(&commandOptionsCreate.Wait, "wait", "w", false, "Wait until the deplyoment success or failed")
	cmd.Flags().StringVarP(&commandOptionsCreate.FromEnvironment, "fromEnvironment", "f", "", "Deploy last deplyoment from given environment")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runCreate(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Create deplyoment
	deplyoment, err := client.CreateDeployment(cli, &commandOptionsCreate)
	if err != nil {
		log.Fatal("Error Create Deployment: ", err)
	}

	//Print result
	if deplyoment.AppServerName != "" {
		cmd.Printf("%s ", deplyoment.AppServerName)
	}
	if deplyoment.EnvironmentName != "" {
		cmd.Printf("%s ", deplyoment.EnvironmentName)
	}
	if deplyoment.ReleaseName != "" {
		cmd.Printf("%s ", deplyoment.ReleaseName)
	}
	if deplyoment.DeploymentDate != 0 {
		cmd.Printf("%s ", time.Unix(0, deplyoment.DeploymentDate*int64(time.Millisecond)).Format("2006-01-02T15:04"))
	}
	if deplyoment.State != "" {
		cmd.Println(deplyoment.State)
	}
	for _, appsWithVersion := range deplyoment.AppsWithVersion {
		cmd.Printf("%s ", appsWithVersion.ApplicationName)
		cmd.Println(appsWithVersion.Version)
	}

	//Write error failed -> return code = 1 with log.Fatal
	if deplyoment.State == client.DeploymentStateFailed {
		log.Fatal("Deployment failed with state: ", deplyoment.State)
	}

}
