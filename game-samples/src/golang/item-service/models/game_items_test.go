package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUUID(t *testing.T) {
	id := generateUUID()
	_, err := uuid.Parse(id)

	assert.Nil(t, err)
}
