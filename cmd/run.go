package cmd

import (
	"fmt"
	"github.com/jsanda/tlp-stress-go/pkg/flags"
	"github.com/jsanda/tlp-stress-go/pkg/run"
	"github.com/spf13/cobra"
)

func init() {
	runCmd.Flags().String("keyspace", "tlp_stress", "Keyspace to use")
	runCmd.Flags().String("host", "127.0.0.1", "Address of Cassandra node used for initial connection")
	runCmd.Flags().Bool("drop", false, "Drop the keyspace before starting")
	runCmd.Flags().String("username", "cassandra", "The username with which to authenticate")
	runCmd.Flags().String("password", "cassandra", "The password with which to authenticate")
	runCmd.Flags().String("profile", "BasicTimeSeries", "The stress profile to execute")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "",
	Long:  "Run a tlp-stress profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running tlp-stress...")
		exec(cmd)
	},
}

func exec(cmd *cobra.Command) {
	flags.Init(cmd)

	runtime := run.NewRuntime(&run.RuntimeConfig{
		Profile: flags.GetString("profile"),
		CqlConfig: run.CqlConfig{
			Hosts: []string {flags.GetString("host")},
			Keyspace: flags.GetString("keyspace"),
			DropKeyspace: flags.GetBool("drop"),
			Username: flags.GetString("username"),
			Password: flags.GetString("password"),
		},
	})

	runtime.Exec()
}