package config

type Config struct {
	Port        string
	Database    Database
	AuthService AuthService `mapstructure:"auth_service"`
}
