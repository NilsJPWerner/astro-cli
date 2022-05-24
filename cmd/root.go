package cmd

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	astro "github.com/astronomer/astro-cli/astro-client"
	cloudCmd "github.com/astronomer/astro-cli/cmd/cloud"
	softwareCmd "github.com/astronomer/astro-cli/cmd/software"
	"github.com/astronomer/astro-cli/config"
	"github.com/astronomer/astro-cli/context"
	"github.com/astronomer/astro-cli/houston"
	"github.com/astronomer/astro-cli/pkg/httputil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	houstonClient houston.ClientInterface
	verboseLevel  string
)

// NewRootCmd adds all of the primary commands for the cli
func NewRootCmd() *cobra.Command {
	httpClient := httputil.NewHTTPClient()
	// configure http transport
	dialTimeout := config.CFG.HoustonDialTimeout.GetInt()
	// #nosec
	httpClient.HTTPClient.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(dialTimeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(dialTimeout) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.CFG.HoustonSkipVerifyTLS.GetBool()},
	}
	houstonClient = houston.NewClient(httpClient)

	astroClient := astro.NewAstroClient(httputil.NewHTTPClient())
	rootCmd := &cobra.Command{
		Use:   "astro",
		Short: "Run Apache Airflow locally and interact with Astronomer",
		Long:  "Welcome to the Astro CLI. Astro is the modern command line interface for data orchestration. You can use it for Astro, Astronomer Software, or local development.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if context.IsCloudContext() {
				return cloudCmd.Setup(cmd, args, astroClient)
			}
			// Software PersistentPreRunE component
			// setting up log verbosity and dumping debug logs collected during CLI-initialization
			if err := softwareCmd.SetUpLogs(os.Stdout, verboseLevel); err != nil {
				return err
			}
			softwareCmd.PrintDebugLogs()
			return nil
		},
	}

	rootCmd.AddCommand(
		newLoginCommand(astroClient, os.Stdout),
		newLogoutCommand(os.Stdout),
		newVersionCommand(),
		newDevRootCmd(),
		newContextCmd(os.Stdout),
		newConfigRootCmd(os.Stdout),
		newAuthCommand(),
	)

	if context.IsCloudContext() { // Include all the commands to be exposed for cloud users
		rootCmd.AddCommand(
			cloudCmd.AddCmds(astroClient, os.Stdout)...,
		)
	} else { // Include all the commands to be exposed for software users
		rootCmd.AddCommand(
			softwareCmd.AddCmds(houstonClient, os.Stdout)...,
		)
		rootCmd.PersistentFlags().StringVarP(&verboseLevel, "verbosity", "", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")
	}
	return rootCmd
}
