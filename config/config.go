package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// AppEnvironment : string wrapper for environment name
type AppEnvironment string

func (e AppEnvironment) IsLocal() bool {
	return e == AppEnvironmentLocal
}

// slice of all app environments except the `local`` env
var (
	// MongoDB
	MongoURI          string
	MongoDatabaseName string

	// Redis
	RedisURI      string
	RedisPassword string

	// Token
	JWTSecretKey string

	// Server
	ServerPort string

	AppEnvironemnts = []AppEnvironment{
		AppEnvironmentStaging,
		AppEnvironmentSandbox,
		AppEnvironmentProduction,
	}
)

const (
	// AppEnvironmentLocal : local env
	AppEnvironmentLocal = AppEnvironment("local")

	// AppEnvironmentLocal : staging env
	AppEnvironmentStaging = AppEnvironment("staging")

	// AppEnvironmentLocal : sandbox env
	AppEnvironmentSandbox = AppEnvironment("sandbox")

	// AppEnvironmentLocal : production env
	AppEnvironmentProduction = AppEnvironment("production")

	// ConfigFileName : localConfig.yaml
	ConfigFileName string = "localConfig"

	// ConfigFileType : yaml
	ConfigFileType string = "yaml"
)

func init() {
	SetConfigFromViper()
}

func SetConfigFromViper() {
	currentHostEnvironment := getCurrentHostEnvironment()
	log.Printf("🚀 Current Host Environment: %s\n", currentHostEnvironment)

	// if env is local, we set the env variables using the config file
	if currentHostEnvironment.IsLocal() {
		setEnvironmentVarsFromConfig()
	}

	// fetch the env vars and store in variables
	// MongoDB
	MongoURI = os.Getenv("MONGO_URI")
	MongoDatabaseName = os.Getenv("MONGO_DATABASE_NAME")

	// Redis
	RedisURI = os.Getenv("REDIS_URI")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	// Token
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

	// Server
	ServerPort = os.Getenv("SERVER_PORT")
}

func getCurrentHostEnvironment() AppEnvironment {
	currentHostEnvironment := os.Getenv("ENVIRONMENT")
	for _, env := range AppEnvironemnts {
		if env == AppEnvironment(currentHostEnvironment) {
			return env
		}
	}
	// if env not found return `local`` env
	return AppEnvironmentLocal
}

func setEnvironmentVarsFromConfig() {
	baseProjectPath, _ := os.Getwd()

	// add the path of the config file
	viper.AddConfigPath(baseProjectPath + "/config/")
	// set the config file name
	viper.SetConfigName(ConfigFileName)
	// set the config file type
	viper.SetConfigType(ConfigFileType)

	viper.AutomaticEnv()

	// read the env vars from the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("unable to read env vars from config file")
	}

	// get the env vars from viper
	mongoURI := viper.GetString("MONGO_URI")
	mongoDatabaseName := viper.GetString("MONGO_DB_NAME")
	redisURI := viper.GetString("REDIS_URI")
	redisPassword := viper.GetString("REDIS_PASSWORD")
	jwtSecretKey := viper.GetString("JWT_SECRET_KEY")
	serverPort := viper.GetString("SERVER_PORT")

	// set the host OS env vars
	os.Setenv("MONGO_URI", mongoURI)
	os.Setenv("MONGO_DATABASE_NAME", mongoDatabaseName)
	os.Setenv("REDIS_URI", redisURI)
	os.Setenv("REDIS_PASSWORD", redisPassword)
	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)
	os.Setenv("SERVER_PORT", serverPort)
}
