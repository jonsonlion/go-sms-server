package uuid

import (
	"github.com/sluu99/uuid"
	"strings"
)

func GetUUID() string {
	id := uuid.Rand()
	return id.Hex()
}

func GetUUIDWithoutLine() string {
	id := uuid.Rand()
	return strings.Replace(id.Hex(), "-", "", -1)
}
