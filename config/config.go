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
	MongoURI          string
	MongoDatabaseName string

	RedisURI      string
	RedisPassword string

	JWTSecretKey string

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
	MongoURI = os.Getenv("MONGO_URI")
	MongoDatabaseName = os.Getenv("MONGO_DATABASE_NAME")

	RedisURI = os.Getenv("REDIS_URI")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
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
	mongoURI := viper.GetString("mongoURI")
	mongoDatabaseName := viper.GetString("mongoDatabaseName")

	redisURI := viper.GetString("redisURI")
	redisPassword := viper.GetString("redisPassword")

	jwtSecretKey := viper.GetString("jwtSecretKey")

	// set the host OS env vars
	os.Setenv("MONGO_URI", mongoURI)
	os.Setenv("MONGO_DATABASE_NAME", mongoDatabaseName)

	os.Setenv("REDIS_URI", redisURI)
	os.Setenv("REDIS_PASSWORD", redisPassword)

	os.Setenv("JWT_SECRET_KEY", jwtSecretKey)
}
