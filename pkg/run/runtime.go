package run

import (
	"github.com/gocql/gocql"
	"log"
)

type RuntimeConfig struct {
	Profile string
	CqlConfig
}

type CqlConfig struct {
	Hosts        []string
	Keyspace     string
	DropKeyspace bool
	Username     string
	Password     string
}

type Runtime struct {
	RuntimeConfig
}

func NewRuntime(cfg *RuntimeConfig) *Runtime {
	return &Runtime{RuntimeConfig: *cfg}
}

func (r *Runtime) Exec() {
	cluster := gocql.NewCluster(r.CqlConfig.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: r.CqlConfig.Username, Password: r.CqlConfig.Password}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to initialize Cassandra session: %s", err)
	}

	r.createKeyspace(session)
	session.Close()

	cluster.Keyspace = r.CqlConfig.Keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to initialize Cassandra session: %s", err)
	}

	session.Close()
	log.Println("Done!")
}

func (r *Runtime) createKeyspace(session *gocql.Session, ) {
	if (r.CqlConfig.DropKeyspace) {
		log.Printf("Dropping keyspace %s", r.CqlConfig.Keyspace)
		if err:= session.Query("DROP KEYSPACE IF EXISTS " + r.CqlConfig.Keyspace).Exec(); err != nil {
			log.Fatalf("Failed to drop keyspace %s: %s", r.CqlConfig.Keyspace, err)
		}
	}

	query :=
		`CREATE KEYSPACE IF NOT EXISTS ` + r.CqlConfig.Keyspace +
			` WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`
	if err := session.Query(query, r.CqlConfig.Keyspace).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace %s: %s", r.CqlConfig.Keyspace, err)
	}
}