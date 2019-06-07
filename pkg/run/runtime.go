package run

import (
	"github.com/gocql/gocql"
	"github.com/jsanda/tlp-stress-go/pkg/generators"
	"github.com/jsanda/tlp-stress-go/pkg/profiles"
	"github.com/jsanda/tlp-stress-go/pkg/metrics"
	"log"
	"gopkg.in/cheggaaa/pb.v1"
	"time"
)

type RuntimeConfig struct {
	Profile               string
	CqlConfig
	Populate              int64
	Partitions            int64
	Concurrency           int64
	PartitionKeyGenerator string
	Id                    string
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

type StressCfg struct {
	Session               *gocql.Session
	Registry              *generators.Registry
	Plugin                *profiles.Plugin
	Concurrency           int64
	PartitionKeyGenerator string
	Id                    string
	Population            int64
	Metrics               *metrics.Metrics
}

func NewRuntime(cfg *RuntimeConfig) *Runtime {
	return &Runtime{RuntimeConfig: *cfg}
}

func (r *Runtime) Exec() {
	plugin, ok := profiles.GetPlugin(r.Profile)
	if !ok {
		log.Fatalf("%s is not a valid stress profile", r.Profile)
	}

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
	defer session.Close()

	createSchema(session, plugin)
	// TODO implement executeAdditionalCql as done in kotlin version

	//fieldRegistry := createFieldRegistry(plugin)

	// TODO create metrics

	stressCfg := &StressCfg{
		Session: session,
		Registry: createFieldRegistry(plugin),
		Plugin: plugin,
		Concurrency: r.Concurrency,
		PartitionKeyGenerator: r.PartitionKeyGenerator,
		Id: r.Id,
		Population: r.Populate,
		Metrics: metrics.NewMetrics(),
	}
	runner := createRunners(stressCfg)

	populateData(plugin, runner, r.Populate)
	//log.Println("Done!")
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

func createSchema(session *gocql.Session, plugin *profiles.Plugin) {
	log.Println("Creating tables...")
	for _, statement := range plugin.Instance.Schema() {
		log.Println(statement)
		if err := session.Query(statement).Exec(); err != nil {
			log.Fatalf("Failed to execute %s: ", statement, err)
		}
	}
}

func createFieldRegistry(plugin *profiles.Plugin) *generators.Registry {
	registry := generators.NewRegistry()

	for field, generator := range plugin.Instance.GetFieldGenerators() {
		registry.SetDefault(field, generator)
	}

	// TODO add support for overriding default field generators

	return registry
}

func populateData(plugin *profiles.Plugin, runner *profileRunner, populate int64) {
	if populate > 0 {
		log.Printf("Prepopulating Cassandra with %d records\n", populate)
		done := make(chan struct{})
		count := int64(0)
		bar := pb.StartNew(int(populate))

		ticker := time.NewTicker(1 * time.Second)

		go func() {
			for range ticker.C {
				// TODO hook in metrics here so we can update the progress bard with the count of the population metric
				bar.Set64(count)
			}
		}()

		go runner.Populate(populate, &count, done)

		<-done
		log.Println("Pre-populate complete")
	}
}