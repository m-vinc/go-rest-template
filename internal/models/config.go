package models

type ConfigApplication struct {
	MasterKey   string `mapstructure:"master_key" validate:"required"`
	Environment string `mapstructure:"env" validate:"required"`
	Bind        string `mapstructure:"bind" validate:"required"`
}

type ConfigPostgres struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`

	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type Config struct {
	Application *ConfigApplication `mapstructure:"application" validate:"required,dive"`
	Postgres    *ConfigPostgres    `mapstructure:"postgres" validate:"required,dive"`
}
