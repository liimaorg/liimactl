package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"log"

	"github.com/liimaorg/liimactl/client"
	"github.com/liimaorg/liimactl/cmd/deployment"
	"github.com/liimaorg/liimactl/cmd/hostname"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfgFile string

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd := newRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// newRootCmd represents the base command when called without any subcommands
func newRootCmd() *cobra.Command {
	var flags *pflag.FlagSet
	liimacli := &client.Cli{}

	var rootCmd = &cobra.Command{
		Use:   "liimactl",
		Short: "Comandline tool for liima",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			config, err := initConfig(flags)
			if err != nil {
				return err
			}
			liimacli.Client, err = client.NewClient(config)
			return err
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.liimactl/config.yaml)")
	rootCmd.PersistentFlags().String("host", "", "liima host")
	flags = rootCmd.Flags()
	rootCmd.AddCommand(deployment.NewDeploymentCmd(liimacli))
	rootCmd.AddCommand(hostname.NewHostnameCmd(liimacli))

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(flags *pflag.FlagSet) (*client.Config, error) {

	var config client.Config
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	//Get path of executable
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(exPath)
	viper.AddConfigPath("$HOME/.liimactl")
	viper.SetEnvPrefix("LIIMA")
	viper.BindEnv("HOST")
	viper.BindPFlag("host", flags.Lookup("host"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into config, %v", err)
	}

	return &config, nil
}
