package hostname

import (
	"github.com/liimaorg/liimactl/cmd"
	"github.com/spf13/cobra"
)

//Hostnames type
type Hostnames []struct {
	Host             string `json:"host"`
	AppServer        string `json:"appServer"`
	AppServerRelease string `json:"appServerRelease"`
	Runtime          string `json:"runtime"`
	Node             string `json:"node"`
	NodeRelease      string `json:"nodeRelease"`
	Environment      string `json:"environment"`
	Domain           string `json:"domain"`
	DefinedOnNode    bool   `json:"definedOnNode"`
}

//sort.Interface
func (slice Hostnames) Len() int {
	return len(slice)
}

func (slice Hostnames) Less(i, j int) bool {
	return slice[i].Domain < slice[j].Domain
}

func (slice Hostnames) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//NewHostnameCmd is a command to manage hostnames
func NewHostnameCmd(cli *cmd.Cli) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "hostname COMMAND",
		Short: "Manage hostnames",
	}

	cmd.AddCommand(newGetCommand(cli))

	return cmd
}
