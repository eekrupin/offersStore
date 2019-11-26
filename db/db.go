package db

import (
	"github.com/gocql/gocql"
	"time"
)

type Config struct {
	Hosts           []string
	Port            int
	Password        string
	User            string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	Keyspace        string
	Consistency     gocql.Consistency
	Cluster_timeout int
}

var DB *gocql.ClusterConfig

func Open(c *Config) (dbConnection *gocql.ClusterConfig, err error) {

	cluster := gocql.NewCluster(c.Hosts...)
	cluster.Keyspace = c.Keyspace
	cluster.Consistency = c.Consistency
	cluster.Timeout = time.Duration(time.Duration(c.Cluster_timeout) * time.Second)
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	session.Close()
	return cluster, nil
}
