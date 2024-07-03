package config

import (
	"github.com/spf13/viper"
	log "k8s.io/klog/v2"
)

type Config struct {
	Name string `mapstructure:"POWERSCALE_NAME"`
	// URL address of the Powerscale API endpoint.
	// Example: `https://isilon.example.com:8080`
	ApiEndpoint string `mapstructure:"POWERSCALE_API_ENDPOINT"`
	ApiUsername string `mapstructure:"POWERSCALE_API_USERNAME"`
	ApiPassword string `mapstructure:"POWERSCALE_API_PASSWORD"`
	// URL address of the S3 Powerscale endpoint.
	// Example: `https://data.nas.example.com:9021`
	S3Endpoint string `mapstructure:"POWERSCALE_S3_ENDPOINT"`
	S3Region   string `mapstructure:"POWERSCALE_S3_REGION"`
	Zone       string `mapstructure:"POWERSCALE_ZONE"`
	BasePath   string `mapstructure:"POWERSCALE_BASE_PATH"`
	// TLS options
	TlsInsecureSkipVerify bool   `mapstructure:"POWERSCALE_TLS_INSECURE_SKIP_VERIFY"`
	TlsClientCert         string `mapstructure:"POWERSCALE_TLS_CLIENT_CERT"`
	TlsClientKey          string `mapstructure:"POWERSCALE_TLS_CLIENT_KEY"`
	TlsCacert             string `mapstructure:"POWERSCALE_TLS_CACERT"`
}

func New() *Config {

	// Required env
	viper.BindEnv("POWERSCALE_NAME")
	viper.BindEnv("POWERSCALE_API_ENDPOINT")
	viper.BindEnv("POWERSCALE_API_USERNAME")
	viper.BindEnv("POWERSCALE_API_PASSWORD")
	viper.BindEnv("POWERSCALE_S3_ENDPOINT")
	viper.BindEnv("POWERSCALE_S3_REGION")
	viper.BindEnv("POWERSCALE_ZONE")
	viper.BindEnv("POWERSCALE_BASE_PATH")

	// Optional
	viper.SetDefault("POWERSCALE_TLS_INSECURE_SKIP_VERIFY", false)
	viper.SetDefault("POWERSCALE_TLS_CLIENT_CERT", "")
	viper.SetDefault("POWERSCALE_TLS_CLIENT_KEY", "")
	viper.SetDefault("POWERSCALE_TLS_CACERT", "")

	var cfg Config

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
	log.Infof("Loaded config: %+v", cfg)

	return &cfg
}
