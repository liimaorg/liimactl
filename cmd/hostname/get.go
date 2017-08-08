package hostname

import (
	"fmt"
	"sort"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	hostnameLong = fmt.Sprintf(` 
    Get a hostname with the use of specific filter.`)

	//Example command description
	hostnameExample = fmt.Sprintf(` 
    # Get a hostname with specific filter. 
    liimactl.exe hostname get --appServer=test_application --envrionment=I`)

	//Flags of the command
	commandOptions client.CommandOptions
)

//newGetCommand is a command to get hostnames
func newGetCommand(cli *client.Cli) *cobra.Command {

	var cmd = &cobra.Command{
		Use:     "get [flags] ",
		Short:   "Get hostnames",
		Long:    hostnameLong,
		Example: hostnameExample,
		Run: func(cmd *cobra.Command, args []string) {
			runGet(cmd, cli, args)
		},
	}

	cmd.Flags().StringSliceVarP(&commandOptions.AppServer, "appServer", "", []string{}, "Application server name")
	cmd.Flags().StringSliceVarP(&commandOptions.Runtime, "runtime", "", []string{}, "Runtime name")
	cmd.Flags().StringSliceVarP(&commandOptions.Environment, "environment", "", []string{}, "Environement name")
	cmd.Flags().StringSliceVarP(&commandOptions.Host, "host", "", []string{}, "Host name")
	cmd.Flags().StringSliceVarP(&commandOptions.Node, "node", "", []string{}, "Node name")
	cmd.Flags().BoolVarP(&commandOptions.DisableMerge, "disableMerge", "", false, "Merge releases")

	return cmd
}

//Get the hostnames properties given by the arguments (see type Hostnames) and print it on the console
func runGet(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get hostnames
	hostnames := client.GetHostname(cli, &commandOptions)

	//Print result
	sort.Sort(hostnames)
	for _, hostname := range hostnames {
		cmd.Printf("%s ", hostname.AppServer)
		cmd.Printf("%s ", hostname.Environment)
		cmd.Printf("%s ", hostname.Host)
		cmd.Printf("%s ", hostname.Runtime)
		cmd.Printf("%s ", hostname.Node)
		cmd.Printf("%s ", hostname.NodeRelease)
		cmd.Println(hostname.Domain)
	}

}
