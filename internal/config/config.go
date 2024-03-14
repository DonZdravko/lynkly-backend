package config

import "os"

// MongoConfig Config holds configuration settings for the application.
type MongoConfig struct {
	MongoDBURI string
}

// NewMongoConfig initializes a new Config instance with default values.
func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		MongoDBURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
	}
}

type ServerConfig struct {
	Port string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: getEnv("PORT", "8080"),
	}
}

// getEnv retrieves the value of the specified environment variable,
// or returns the default value if the environment variable is not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
