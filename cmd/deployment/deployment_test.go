package deployment

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/pflag"
)

//Tests the command "deployment get"
func TestNewDeploymentGetCmd(t *testing.T) {

	//Tests
	tests := []struct {
		name string   //Name of the test
		args []string //Arguments
		want string   //Wanted testresult
	}{
		{"Test1", []string{"get", "--appServer=testApp"}, "------\nTest"},
		{"Test2", []string{"get", "--appServer=testApp2", "--environment=T"}, "------\nTest"},
		{"Test3", []string{"get", "--filter=[{\"name\":\"Environment\",\"comp\":\"eq\",\"val\":\"Y\"},{\"name\":\"Application server\",\"comp\":\"eq\",\"val\":\"testApp3\"}]"}, "------\nTest"},
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
			cmd := NewDeploymentCmd(liimacli)

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
			if got := buf.String(); !strings.HasPrefix(got, tt.want) {
				t.Errorf("Commands-Output = %v, want %v", got, tt.want)
			}
		})
	}

}

//Tests the command "deployment create"
func TestNewDeploymentCreateCmd(t *testing.T) {

	//Tests
	tests := []struct {
		name string   //Name of the test
		args []string //Arguments
		want string   //Wanted testresult
	}{
		{"Test1", []string{"create", "--appServer=testApp", "--environment=T", "--appName=test1", "--version=1.1.1"}, "------\nSUCCESS\n"},
		{"Test2", []string{"promote", "--environment=Y", "--fromEnvironment=B", "--date=2018-02-01 17:00", "--silent"}, "------\nSUCCESS\n"},
		{"Test3", []string{"promote", "--environment=Y", "--fromEnvironment=B", "--date=2018-02-01 17:00", "-c", "--blacklistAppServer=Test"}, ""},
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
			cmd := NewDeploymentCmd(liimacli)

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
			if got := buf.String(); !strings.HasPrefix(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}

}

// initConfig reads in config file and ENV variables if set.
func initConfig(flags *pflag.FlagSet) (*client.Config, error) {

	var config client.Config
	return &config, nil
}
