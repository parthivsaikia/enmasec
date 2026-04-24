package service

import (
	"github.com/google/uuid"
	"github.com/parthivsaikia/enmasec/internal/models"
	"github.com/parthivsaikia/enmasec/internal/store"
)

var LocationUUIDMap = models.BiMap{
	ForwardMap: map[uuid.UUID]string{},
	ReverseMap: map[string]uuid.UUID{},
}

func InitIndexMap() {
	biMapData, err := store.DecryptIndexFile(vault, password)
}
