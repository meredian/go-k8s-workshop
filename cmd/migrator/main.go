package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gocql/gocql"
	"gitlab.k8s.gromnsk.ru/workshop/austin/migrations"
	"gitlab.k8s.gromnsk.ru/workshop/austin/pkg/config"
)

// RUN ./migrator -table=tracking
func main() {
	var tableName string
	flag.StringVar(&tableName, "table", "", "-table={YOUR_TABLE}")
	flag.Parse()

	if tableName == "" {
		log.Fatalf("the table flag must not be empty")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get a configs: %v\n", err)
	}

	hostport := fmt.Sprintf("%s:%d", cfg.Storage.Host, cfg.Storage.Port)

	fmt.Printf("run migrations ...\n")
	fmt.Printf("Keyspace: %v\n", cfg.Storage.Keyspace)
	fmt.Printf("Table: %v\n", tableName)

	clusterConfig := buildClusterConfig(hostport)

	err = createKeyspace(clusterConfig, cfg.Storage.Keyspace)
	if err != nil {
		log.Fatalf("failed to create a session: %v", err)
	}

	clusterConfig.Keyspace = cfg.Storage.Keyspace

	err = createTable(clusterConfig, tableName)
	if err != nil {
		log.Fatalf("failed to create session when creating a table %s: %v", tableName, err)
	}

	fmt.Printf("Migrations have successfully applied!\n")
}

func buildClusterConfig(hostport string) *gocql.ClusterConfig {
	return gocql.NewCluster(hostport)
}

func createKeyspace(clusterConfig *gocql.ClusterConfig, keyspaceName string) error {
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	createKeyspaceQuery := fmt.Sprintf(migrations.CreateKeyspace, keyspaceName)
	err = session.Query(createKeyspaceQuery).Exec()
	if err != nil {
		return err
	}

	return nil
}

func createTable(clusterConfig *gocql.ClusterConfig, tableName string) error {
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return err
	}
	defer session.Close()

	createTableQuery := fmt.Sprintf(migrations.CreateTable, tableName)

	err = session.Query(createTableQuery).Exec()

	if err != gocql.ErrTimeoutNoResponse && err != nil {
		return err
	}

	return nil
}
