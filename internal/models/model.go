package models

import "github.com/google/uuid"

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

type BiMap struct {
	ForwardMap map[uuid.UUID]string
	ReverseMap map[string]uuid.UUID
}

func (b *BiMap) Put(uuid uuid.UUID, path string) {
	b.ForwardMap[uuid] = path
	b.ReverseMap[path] = uuid
}

func (b *BiMap) GetByKey(key uuid.UUID) (string, bool) {
	v, ok := b.ForwardMap[key]
	return v, ok
}

func (b *BiMap) GetByVal(val string) (uuid.UUID, bool) {
	k, ok := b.ReverseMap[val]
	return k, ok
}
