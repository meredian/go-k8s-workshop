package storage

import (
	"time"

	"github.com/gocql/gocql"
)

func InitStorage(hostport, keyspace string) (*gocql.Session, error) {
	clusterConfig := gocql.NewCluster(hostport)
	clusterConfig.Keyspace = keyspace
	clusterConfig.ConnectTimeout = 10 * time.Second
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "iD5hBMkSfDTF",
	}

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
