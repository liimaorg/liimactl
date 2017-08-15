package hostname

import (
	"bytes"
	"testing"

	"fmt"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/pflag"
)

//Tests the command "hostname get"
func TestNewHostnameCmd(t *testing.T) {

	//Tests
	tests := []struct {
		name string   //Name of the test
		args []string //Arguments
		want string   //Wanted testresult
	}{
		{"Test1", []string{"get", "--appServer=testApp"}, "testApp "},
		{"Test2", []string{"get", "--appServer=testApp2", "--environment=T"}, "testApp2 T "},
	}

	//Init config
	var flags *pflag.FlagSet
	liimacli := &client.Cli{}
	config, err := initConfig(flags)
	if err != nil {
		fmt.Println(err)
	}

	//Create mock client
	liimacli.Client, err = client.NewMockClient(config)

	//Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//Create command
			cmd := NewHostnameCmd(liimacli)

			//Set commands output to buffer
			buf := new(bytes.Buffer)
			cmd.SetOutput(buf)

			//Set test arguments
			cmd.SetArgs(tt.args)

			//Execute command
			err = cmd.Execute()
			if err != nil {
				t.Errorf("Execute() failed with %v", err)
			}
			//Check result
			if got := buf.String(); got != tt.want {
				t.Errorf("Commands-Output = %v, want %v", got, tt.want)
			}

		})
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig(flags *pflag.FlagSet) (*client.Config, error) {

	var config client.Config
	return &config, nil
}
