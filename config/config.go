package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Configurations represents the app configurations including Env varaibles
type Configurations struct {
	Port            int
	HashSecretKey   string
	DatabaseURI     string
	CSRFKey         string
	EmailAPIKey     string
	IsProductionEnv bool
}

// AppConfig is the configurations for the app you should load the config
// package and this varaibale will be populated for you with all the
// configs and environment vars
var AppConfig Configurations

func init() {
	loadEnvFile()
	configs, err := newConfigurations()
	if err != nil {
		panic(err)
	}
	AppConfig = *configs
}

// loadEnvFile is responsible for loading the env variable file
func loadEnvFile() error {
	return godotenv.Load()
}

// NewConfigurations is the constructor for creating Configurations struct
// that you can use across the whole app
func newConfigurations() (*Configurations, error) {
	port, err := intEnvVariable("PORT")
	if err != nil {
		return nil, err
	}
	hashSecretKey, err := stringEnvVariable("HASH_SECRET_KEY")
	if err != nil {
		return nil, err
	}
	databaseURI, err := stringEnvVariable("DATABASE_URI")
	if err != nil {
		return nil, err
	}
	csrfKey, err := stringEnvVariable("CSRF_KEY")
	if err != nil {
		return nil, err
	}
	emailAPIKey, err := stringEnvVariable("EMAIL_API_KEY")
	if err != nil {
		return nil, err
	}
	isProductionEnv, err := boolEnvVariable("IS_PRODUCTION_ENV")
	if err != nil {
		return nil, err
	}

	return &Configurations{
		Port:            port,
		HashSecretKey:   hashSecretKey,
		DatabaseURI:     databaseURI,
		CSRFKey:         csrfKey,
		EmailAPIKey:     emailAPIKey,
		IsProductionEnv: isProductionEnv,
	}, nil
}

func stringEnvVariable(name string) (string, error) {
	val, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("env variable %v not found", name)
	}
	if val == "" {
		return "", fmt.Errorf("env variable %v can not be empty", name)
	}
	return val, nil
}

func intEnvVariable(name string) (int, error) {
	strVal, err := stringEnvVariable(name)
	if err != nil {
		return 0, err
	}

	// convert the value to int
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, err
	}
	return intVal, nil
}

func boolEnvVariable(name string) (bool, error) {
	strVal, err := stringEnvVariable(name)
	if err != nil {
		return false, err
	}

	// convert the value to int
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		return false, err
	}
	return boolVal, nil
}
