package hostname

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/liimaorg/liimactl/cmd"
	"github.com/liimaorg/liimactl/cmd/util"
	"github.com/spf13/cobra"
)

//CommandOptions used for the command options (flags)
type CommandOptions struct {
	AppServer    []string `json:"appServer"`
	Runtime      []string `json:"runtime"`
	Environment  []string `json:"environment"`
	Host         []string `json:"host"`
	Node         []string `json:"node"`
	DisableMerge bool     `json:"disableMerge"`
}

var (
	//Long command description
	hostnameLong = fmt.Sprintf(` 
    Get a hostname with the use of specific filter.`)

	//Example command description
	hostnameExample = fmt.Sprintf(` 
    # Get a hostname with specific filter. 
    liimactl.exe hostname get --appServer=test_application --envrionment=I`)

	//Flags of the command
	commandOptions CommandOptions
)

//newGetCommand is a command to get hostnames
func newGetCommand(cli *cmd.Cli) *cobra.Command {

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
func runGet(cmd *cobra.Command, cli *cmd.Cli, args []string) {

	//Build URL
	url := fmt.Sprintf("resources/./hostNames?")
	url += util.BuildCommandURL(&commandOptions)

	//Call rest client
	hostnames := Hostnames{}
	if err := cli.Client.DoRequest(http.MethodGet, url, nil, &hostnames); err != nil {
		log.Fatal("Error rest call: ", err)
	}

	//Print result
	sort.Sort(hostnames)
	for _, hostname := range hostnames {
		fmt.Printf("%s ", hostname.AppServer)
		fmt.Printf("%s ", hostname.Environment)
		fmt.Printf("%s ", hostname.Host)
		fmt.Printf("%s ", hostname.Runtime)
		fmt.Printf("%s ", hostname.Node)
		fmt.Printf("%s ", hostname.NodeRelease)
		fmt.Println(hostname.Domain)
	}

}
