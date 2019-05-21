package cmd

import (
	"fmt"
	"github.com/gocql/gocql"
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
	contactPoint, err := cmd.Flags().GetString("host")
	if err != nil {
		log.Fatalf("Failed to get value of host flag: %s", err)
	}

	keyspace, err := cmd.Flags().GetString("keyspace")
	if err != nil {
		log.Fatalf("Failed to get value of keyspace flag: %s", err)
	}

	dropKeyspace, err := cmd.Flags().GetBool("drop")
	if err != nil {
		log.Fatalf("Failed to get value of drop flag: %s", err)
	}

	username, password := getUsernameAndPassword(cmd)

	cluster := gocql.NewCluster(contactPoint)
	//cluster.Keyspace = keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to initialize Cassandra session: %s", err)
	}

	createKeyspace(session, keyspace, dropKeyspace)

	session.Close()
	log.Println("Done!")
}

func getUsernameAndPassword(cmd *cobra.Command) (string, string) {
	username, err := cmd.Flags().GetString("username")
	if err != nil {
		log.Fatalf("Failed to parse username flag: %s", err)
	}

	password, err := cmd.Flags().GetString("password")
	if err != nil {
		log.Fatalf("Faield to parse password flag: %s", err)
	}

	return username, password
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