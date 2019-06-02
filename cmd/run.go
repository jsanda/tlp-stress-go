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
	runCmd.Flags().String("duration", "", "Duration of the stress test. Expressed in format 1d 3h 15m")
	runCmd.Flags().Int("iterations", 0, "Number of operations to run")
	runCmd.Flags().Int64("populate", 0, "Pre-population the DB with N rows before starting load test")
	runCmd.Flags().Int64("partitions", 1000000, "Max value of integer component of first partition key")
	runCmd.Flags().Int64("concurrency", 100, "Concurrent queries allowed.  Increase for larger clusters")
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
		Populate: flags.GetInt64("populate"),
		Partitions: flags.GetInt64("partitions"),
		Concurrency: flags.GetInt64("concurrency"),
	})

	runtime.Exec()
}