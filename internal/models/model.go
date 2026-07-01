package models

import (
	"github.com/google/uuid"
)

type Vault struct {
	Name string
	Path string
}

type Account struct {
	Username  string            `toml:"username"`
	Password  string            `toml:"password"`
	Metadata  map[string]string `toml:"metadata"`
	UpdatedAt string            `toml:"updated_at"`
}

type Config struct {
	Theme string
}

type State struct {
	CurrentVault string `yaml:"current_vault"`
}

type Registry struct {
	Vaults map[string]string `yaml:"vaults"`
}

type VaultIndex struct {
	Version  int                         `json:"version"`
	Services map[uuid.UUID]*ServiceEntry `json:"services"`
}

type ServiceEntry struct {
	Name     string                      `json:"name"`
	Accounts map[uuid.UUID]*AccountEntry `json:"accounts"`
}

type AccountEntry struct {
	Name string `json:"name"`
}

type RuntimeIndex struct {
	// Vault pane
	Vaults []*Vault
	// Service pane
	Services        []*ServiceEntry
	ServiceNameToID map[string]uuid.UUID
	ServiceIDToName map[uuid.UUID]string
	// Accounts pane
	AccountsByService map[uuid.UUID][]*AccountEntry
	AccountNameToID   map[uuid.UUID]map[string]uuid.UUID
	AccountIDToName   map[uuid.UUID]map[uuid.UUID]string
}
