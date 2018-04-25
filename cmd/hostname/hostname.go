package hostname

import (
	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/cobra"
)

//NewHostnameCmd is a command to manage hostnames
func NewHostnameCmd(cli *client.Cli) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "hostname COMMAND",
		Short: "Manage hostnames",
	}

	cmd.AddCommand(newGetCommand(cli))

	return cmd
}
