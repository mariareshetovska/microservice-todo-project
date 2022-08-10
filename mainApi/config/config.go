package config

type ConfigPostgres struct {
	DB struct {
		ConnUrl       string `yaml:"conn_url"`
		MigrationPath string `yaml:"migration_path"`
	} `yaml:"DB"`
	Port string `yaml:"Port"`
}
