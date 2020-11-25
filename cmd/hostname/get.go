package hostname

import (
	"log"
	"sort"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

var (
	//Long command description
	hostnameLong = `	Get a hostname with the use of specific filter.`

	//Example command description
	hostnameExample = `	# Get a hostname with specific filter. 
	liimactl hostname get --appServer=test_application --environment=I`

	//Flags of the command
	commandOptions client.CommandOptionsHostName
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

	cmd.Flags().StringSliceVarP(&commandOptions.AppServer, "appServer", "a", []string{}, "Application server name")
	cmd.Flags().StringSliceVarP(&commandOptions.Runtime, "runtime", "r", []string{}, "Runtime name")
	cmd.Flags().StringSliceVarP(&commandOptions.Environment, "environment", "e", []string{}, "Environement name")
	cmd.Flags().StringSliceVarP(&commandOptions.Host, "host", "s", []string{}, "Host name")
	cmd.Flags().StringSliceVarP(&commandOptions.Node, "node", "n", []string{}, "Node name")
	cmd.Flags().BoolVarP(&commandOptions.DisableMerge, "disableMerge", "d", false, "Merge releases")

	return cmd
}

//Get the hostnames properties given by the arguments (see type Hostnames) and print it on the console
func runGet(cmd *cobra.Command, cli *client.Cli, args []string) {

	//Get hostnames
	hostnames, err := client.GetHostname(cli, &commandOptions)
	if err != nil {
		log.Fatalf("Couldn't get hostnames: %v", err)
	}

	//Print result
	sort.Sort(hostnames)

	//Print result
	for _, hostname := range hostnames {

		if hostname.Host != "" {
			cmd.Printf("%-34s ", hostname.Host)
		}
		if hostname.Environment != "" {
			cmd.Printf("%-8s ", hostname.Environment)
		}
		if hostname.AppServer != "" {
			cmd.Printf("%-30s ", hostname.AppServer)
		}
		if hostname.AppServerRelease != "" {
			cmd.Printf("%-10s ", hostname.AppServerRelease)
		}
		if hostname.Runtime != "" {
			cmd.Printf("%-10s ", hostname.Runtime)
		}
		if hostname.Node != "" {
			cmd.Printf("%-30s ", hostname.Node)
		}
		if hostname.NodeRelease != "" {
			cmd.Printf("%-10s ", hostname.NodeRelease)
		}
		if hostname.Domain != "" {
			cmd.Println(hostname.Domain)
		}
	}
}
