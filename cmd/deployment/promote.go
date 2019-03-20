package deployment

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentPromoteLong = `	Promote deployment on an environemt with the use of specific properties.`

	//Example command description
	deploymentPromoteExample = `	# Promote multiple deployments on an environment with specific properties. 
	liimactl.exe deployment promote --environment=Y  --fromEnvironment=B
	liimactl.exe deployment promote --environment=Y  --fromEnvironment=B --date="2018-02-01 17:00" --blacklistRuntime="Kubernetes,Kube_helm"
	liimactl.exe deployment promote --environment=Y  --fromEnvironment=B --date="2018-02-01 17:00" --blacklistAppServer="aps_bau_kube,vvn"
	liimactl.exe deployment promote --environment=Z  --fromEnvironment=I --whitelistAppServer="appServer1,appServer2" --wait --maxWaitTime=3600`

	//Flags of the command
	commandOptionsPromote client.CommandOptionsPromoteDeployments
)

//newPromoteCommand is a command to promote multiple deployments on an environment
func newPromoteCommand(cli *client.Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "promote [flags] ",
		Short:   "promote deployment",
		Long:    deploymentPromoteLong,
		Example: deploymentPromoteExample,
		Run: func(cmd *cobra.Command, args []string) {
			runPromote(cmd, cli, args)
		},
	}

	cmd.Flags().StringVarP(&commandOptionsPromote.Environment, "environment", "e", "", "Environment")
	cmd.Flags().StringVarP(&commandOptionsPromote.FromEnvironment, "fromEnvironment", "f", "", "Deploy last deployment from given environment")
	cmd.Flags().StringVarP(&commandOptionsPromote.DeploymentDate, "date", "d", "", "Deployment Date 'DD.MM.YYYY hh:mm' ")
	cmd.Flags().BoolVarP(&commandOptionsPromote.ExecuteShakedownTest, "executeShakeDownTest", "s", false, "Run Shakedowntest after the deployment")
	cmd.Flags().BoolVarP(&commandOptionsPromote.Wait, "wait", "w", false, "Wait maxWaitTime until the deployment success or failed")
	cmd.Flags().IntVarP(&commandOptionsPromote.MaxWaitTime, "maxWaitTime", "t", 600, "Max Wait time [seconds] until the deployment success or failed")
	cmd.Flags().StringSliceVarP(&commandOptionsPromote.WhitelistAppServer, "whitelistAppServer", "a", []string{}, "Whitelist with all appServer, which should be deployed, if no WhitelistAppServer is defined, the whole environment will deployed (exclusive blacklist)")
	cmd.Flags().StringSliceVarP(&commandOptionsPromote.BlacklistAppServer, "blacklistAppServer", "b", []string{}, "Blacklist with all appServer, which should not be deployed")
	cmd.Flags().StringSliceVarP(&commandOptionsPromote.BlacklistRuntime, "blacklistRuntime", "r", []string{}, "Blacklist with all runtimes, which should not be deployed")
	cmd.Flags().BoolVarP(&commandOptionsPromote.Silent, "silent", "c", false, "Silent mode, no confirmation of promote the whole environment")

	return cmd
}

//Promote a deployment on an environment and print the state of each deployment on the console
func runPromote(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Ask user for confirmation
	msg := fmt.Sprintf("Do you really want to start deployments on environment: %s", commandOptionsPromote.Environment)
	if commandOptionsPromote.Silent || AskYesNo(msg) {

		//Promote deployment
		deployments, err := client.PromoteDeployments(cli, &commandOptionsPromote)
		if err != nil {
			log.Fatal("Error Promote Deployment: ", err)
		}
		success := true

		//Sort, group be State
		sort.Sort(deployments)
		sort.Slice(deployments, func(i, j int) bool {
			switch strings.Compare(string(deployments[i].State), string(deployments[j].State)) {
			case -1:
				return true
			case 1:
				return false
			}
			return deployments[i].AppServerName < deployments[j].AppServerName
		})
		for _, deployment := range deployments {

			//Print result
			PrintDeployment(cmd, &deployment)

			//Check success
			success = success && (deployment.State == client.DeploymentStateSuccess || deployment.State == client.DeploymentStateRejected)

		}

		//Write failed, if not all deployments are successfully -> return code = 1 with log.Fatal, needed maybe for the result in a batch job
		if !success && commandOptionsPromote.Wait {
			log.Fatal("Promote failed, not all deployments are successfully")
		}

	}
}
