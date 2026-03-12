package model

type Account struct {
	Username  string
	Password  string
	Metadata  map[string]string
	UpdatedAt string
}

type Config struct {
	CurrentVault string            `yaml:"current_vault"`
	Vaults       map[string]string `yaml:"vaults"`
}
