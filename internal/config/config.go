package config

// Config holds all application configuration.
type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SslMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	Payment struct {
		APIURL      string `mapstructure:"payment_api_url"`
		APIKey      string `mapstructure:"payment_api_key"`
		RefreshTime int    `mapstructure:"payment_refresh_time"`
	} `mapstructure:"Payment"`

	Server struct {
		Port     int `mapstructure:"port"`
		HTTPPort int `mapstructure:"httpPort"`
	} `mapstructure:"server"`
}