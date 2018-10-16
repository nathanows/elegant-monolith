package app

import (
	"fmt"

	"github.com/imdario/mergo"
)

// The Config struct wraps the available application level config. Viper is used
// to marshal config files/env vars/flags to Config
type Config struct {
	GRPCAddr       int
	HTTPAddr       int
	DatabaseConfig DatabaseConfig
}

// DatabaseConfig is an environment agnostic config struct for DB setup
type DatabaseConfig struct {
	Username string
	Password string
	Hostname string
	Database string
	Port     int32
	Sslmode  string
}

// BuildDbConnectionStr returns a postgres compliant connection string
func (dbConfig DatabaseConfig) BuildDbConnectionStr() string {
	defaultConfig := &DatabaseConfig{Password: "", Hostname: "localhost", Database: "elegant-monolith", Port: 5432, Sslmode: "disable"}
	mergo.Merge(&dbConfig, defaultConfig)

	connectionStr := fmt.Sprintf("host=%s port=%d user=%s password='%s' dbname=%s sslmode=%s",
		dbConfig.Hostname, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Sslmode)

	return connectionStr
}
