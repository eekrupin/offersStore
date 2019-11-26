package config

import (
	"fmt"
	"github.com/eekrupin/offersStore/db"
	"github.com/eekrupin/offersStore/modules"
	"github.com/eekrupin/offersStore/services/loggerService"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type HTTPServerConfig struct {
	Host         string
	InternalPort uint16
}

type AppConfig struct {
	HTTPServer  *HTTPServerConfig
	DBConfig    *db.Config
	Debug       bool
	ENV         string
	MAX_WORKERS int
	BatchSize   int
}

var Config AppConfig

const (
	KeyResponse = "responseData"
)

func getDefaultEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

func setConfigFromEnv(config *AppConfig) {
	DB_HOST := os.Getenv("DB_HOST")
	if DB_HOST == "" {
		DB_HOST = "localhost"
	}
	Config.DBConfig.Hosts = strings.Split(DB_HOST, ",")

	DBPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		DBPort = 9040
	}

	DB_USER := os.Getenv("DB_USER")
	if DB_USER == "" {
		DB_USER = "cassandra"
	}

	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	if DB_PASSWORD == "" {
		DB_PASSWORD = "cassandra"
	}

	Keyspace := os.Getenv("Keyspace")
	if Keyspace == "" {
		Keyspace = "offers"
	}

	Consistency := os.Getenv("Consistency")
	if Consistency == "" {
		Consistency = "One"
	}

	MAX_WORKERS, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))
	if err != nil {
		MAX_WORKERS = 1
	}

	CLUSTER_TIMEOUT, err := strconv.Atoi(os.Getenv("CLUSTER_TIMEOUT"))
	if err != nil {
		CLUSTER_TIMEOUT = 120
	}

	BATCH_SIZE, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		BATCH_SIZE = 250
	}

	Config.DBConfig.Port = DBPort
	Config.DBConfig.User = DB_USER
	Config.DBConfig.Password = DB_PASSWORD
	Config.DBConfig.Keyspace = Keyspace
	Config.DBConfig.Consistency = gocql.ParseConsistency(Consistency)
	Config.DBConfig.Cluster_timeout = CLUSTER_TIMEOUT
	Config.MAX_WORKERS = MAX_WORKERS
	Config.BatchSize = BATCH_SIZE
}

func init() {
	var err error

	Config = AppConfig{
		HTTPServer: &HTTPServerConfig{},
		DBConfig:   &db.Config{},
	}

	err = godotenv.Load()
	if err != nil {
		loggerService.GetMainLogger().Warn(nil, err.Error())
	}

	setConfigFromEnv(&Config)

	db.DB, err = db.Open(Config.DBConfig)
	if err != nil {
		panic(err)
	}

	modules.MaxWorkers = Config.MAX_WORKERS
	modules.BatchSize = Config.BatchSize

	Config.Debug = getDefaultEnv("IS_DEBUG", "0") == "1"
	if Config.Debug {
		loggerService.GetMainLogger().Info(nil, "Environment variables")
		for _, environ := range os.Environ() {
			loggerService.GetMainLogger().Info(nil, environ)
		}
	}

	Config.ENV = getDefaultEnv("ENV", "dev")
	loggerService.GetMainLogger().Info(nil, fmt.Sprintf("ENV: %s", Config.ENV))

	Config.HTTPServer.Host = getDefaultEnv("HTTP_SERVER_HOST", "")
	httpInternalServerPort, err := strconv.ParseInt(getDefaultEnv("HTTP_INTERNAL_SERVER_PORT", "80"), 10, 32)
	if err == nil {
		Config.HTTPServer.InternalPort = uint16(httpInternalServerPort)
	} else {
		Config.HTTPServer.InternalPort = uint16(80)
	}

}
