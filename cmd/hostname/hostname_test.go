package hostname

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"

	"github.com/liimaorg/liimactl/client"
	"github.com/spf13/pflag"
)

func TestX(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")

	}))
	defer ts.Close()

	fmt.Println(ts.URL)

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)

}

func TestNewHostnameCmd(t *testing.T) {

	var flags *pflag.FlagSet
	liimacli := &client.Cli{}
	config, err := initConfig(flags)
	if err != nil {
		fmt.Println(err)
	}

	liimacli.Client, err = client.NewMockClient(config)

	fmt.Println(config.Host)

	cmd := NewHostnameCmd(liimacli)

	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{"get", "--appServer=testApp"})

	err = cmd.Execute()
	if err != nil {
		t.Errorf("Execute() failed with %v", err)
	}

	fmt.Println("----")
	fmt.Println(buf.String())
	fmt.Println("----")
}

// initConfig reads in config file and ENV variables if set.
func initConfig(flags *pflag.FlagSet) (*client.Config, error) {

	var config client.Config
	return &config, nil
}
