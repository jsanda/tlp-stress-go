package cmd

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/flags"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	runCmd.Flags().String("keyspace", "tlp_stress", "Keyspace to use")
	runCmd.Flags().String("host", "127.0.0.1", "Address of Cassandra node used for initial connection")
	runCmd.Flags().Bool("drop", false, "Drop the keyspace before starting")
	runCmd.Flags().String("username", "cassandra", "The username with which to authenticate")
	runCmd.Flags().String("password", "cassandra", "The password with which to authenticate")

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

	contactPoint := flags.GetString("host")
	keyspace := flags.GetString("keyspace")
	dropKeyspace := flags.GetBool("drop")
	username := flags.GetString("username")
	password := flags.GetString("password")

	cluster := gocql.NewCluster(contactPoint)
	//cluster.Keyspace = keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to initialize Cassandra session: %s", err)
	}

	createKeyspace(session, keyspace, dropKeyspace)
	session.Close()

	cluster.Keyspace = keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to initialize Cassandra session: %s", err)
	}

	session.Close()
	log.Println("Done!")
}

func createKeyspace(session *gocql.Session, keyspace string, dropKeyspace bool) {
	if (dropKeyspace) {
		log.Printf("Dropping keyspace %s", keyspace)
		if err:= session.Query("DROP KEYSPACE IF EXISTS " + keyspace).Exec(); err != nil {
			log.Fatalf("Failed to drop keyspace %s: %s", keyspace, err)
		}
	}

	query :=
		`CREATE KEYSPACE IF NOT EXISTS ` + keyspace +
        ` WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`
	if err := session.Query(query, keyspace).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace %s: %s", keyspace, err)
	}
}