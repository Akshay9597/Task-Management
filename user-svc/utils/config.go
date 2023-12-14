package utils

import (
	"github.com/spf13/viper"
)

const defaultDbSSLMode = "disable"

type(
	DBConfig struct {
		Host string
		Port string
		Name string
		Username string
		Password string
		SSLMode string
	}
)

func Init() (DBConfig, error){
	if err:= setUpViper(); err != nil {
		return DBConfig{}, err
	}
	return setConfig(), nil
}

func setUpViper() error {
	populateDefaults()

	if err := parseDbEnvVariables(); err != nil {
		return err
	}

	return parsePasswordEnvVariables()

}

func parseDbEnvVariables() error {
	viper.SetEnvPrefix("db")
	if err := viper.BindEnv("host"); err != nil {
		return err
	}

	if err := viper.BindEnv("port"); err != nil {
		return err
	}

	if err := viper.BindEnv("name"); err != nil {
		return err
	}

	if err := viper.BindEnv("user"); err != nil {
		return err
	}

	if err := viper.BindEnv("pass"); err != nil {
		return err
	}

	return viper.BindEnv("sslmode")
}

func populateDefaults() {
	viper.SetDefault("db.sslmode", defaultDbSSLMode)
}

func setConfig() DBConfig {
	return DBConfig{
		Host:     viper.GetString("host"),
		Port:     viper.GetString("port"),
		Name:     viper.GetString("name"),
		Username: viper.GetString("user"),
		Password: viper.GetString("pass"),
		SSLMode:  viper.GetString("sslmode"),
	}
}

func parsePasswordEnvVariables() error {
	viper.SetEnvPrefix("password")
	return viper.BindEnv("salt")
}