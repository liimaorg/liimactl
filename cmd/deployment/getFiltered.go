package deployment

import (
	"fmt"
	"log"
	"sort"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	deploymentGetFilteredLong = fmt.Sprintf(` 
    Get deployment with the use of given filters.`)

	//Example command description
	deploymentGetFilteredExample = fmt.Sprintf(` 
    # Get a deplyoment with specific filters. 
	liimactl.exe deployment getFiltered --filter='[{"name":"Environment","comp":"eq","val":"Y"},{"name":"Application server","comp":"eq","val":"aps_bau_kube"}]'
	liimactl.exe deployment getFiltered --filter='[{"name":"Environment","comp":"eq","val":"Y"},{"name":"Latest deployment job for App Server and Env","comp":"eq","val":"true"}]'`)

	//Flags of the command
	commandOptionsGetFilteredDeployments client.CommandOptionsGetFilteredDeployments
)

//newGetFilteredDeploymentsCommand is a command to get deployments filter
func newGetFilteredDeploymentsCommand(cli *client.Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "getFiltered [flags] ",
		Short:   "Get deployments filtered",
		Long:    deploymentGetFilteredLong,
		Example: deploymentGetFilteredExample,
		Run: func(cmd *cobra.Command, args []string) {
			runGetFilteredDeplyoments(cmd, cli, args)
		},
	}

	cmd.Flags().StringVarP(&commandOptionsGetFilteredDeployments.Filter, "filter", "f", "", "Filter")

	return cmd
}

//Get the deployments properties given by the arguments (see type Deplyoments) and print it on the console
func runGetFilteredDeplyoments(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get deployments
	deployments, err := client.GetFilteredDeployments(cli, &commandOptionsGetFilteredDeployments)
	if err != nil {
		log.Fatal("Error on getting the filtered deployments: ", err)
	}

	const createdFormat = "2017-07-06 21:00" //"Jan 2, 2006 at 3:04pm (MST)"

	//Print result
	sort.Sort(deployments)
	for _, deployment := range deployments {
		PrintDeployment(cmd, deployment)
	}

}
