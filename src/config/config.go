package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"strings"
)

var Config appConfig

type appConfig struct {
	DB                     *gorm.DB
	DBErr                  error
	ServerPort             int    `mapstructure:"server_port"`
	DSN                    string `mapstructure:"dsn"`
	AuthServerUrl          string `mapstructure:"auth_server_url"`
	FamilyServerUrl        string `mapstructure:"family_server_url"`
	AuthServerLoginPath    string `mapstructure:"auth_server_login_path"`
	AuthServerUsername     string `mapstructure:"auth_server_username"`
	AuthServerPassword     string `mapstructure:"auth_server_password"`
	AuthServerJwtPublicKey string `mapstructure:"auth_server_jwt_public"`
}

func (a *appConfig) GetAuthServerUrl() string {
	return a.AuthServerUrl
}

func (a *appConfig) GetAuthServerLoginPath() string {
	return a.AuthServerLoginPath
}

func (a *appConfig) GetAuthServerUsername() string {
	return a.AuthServerUsername
}

func (a *appConfig) GetAuthServerPassword() string {
	return a.AuthServerPassword
}

func (a *appConfig) GetAuthServerJwtPublicKey() string {
	return strings.ReplaceAll(a.AuthServerJwtPublicKey, "\\n", "\n")
}

func (a *appConfig) GetFamilyServerUrl() string {
	return a.FamilyServerUrl
}

func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("server")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("measurement_service")
	v.AutomaticEnv()

	v.SetDefault("server_port", 8080)

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}
	return v.Unmarshal(&Config)
}
