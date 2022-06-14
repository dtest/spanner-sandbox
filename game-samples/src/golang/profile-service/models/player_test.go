package models

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidEmails(t *testing.T) {
	var tests = []string{"good@gmail.com", "good@somedomain.net", "good.email@somedomain.org"}

	for _, testEmail := range tests {
		var player = Player{Email: testEmail, Password: "testpassword", Player_name: "Test Player"}

		t.Logf("Testing '%s'", testEmail)

		err := player.Validate()

		assert.Nil(t, err)
	}
}

func TestInvalidEmails(t *testing.T) {
	var tests = []string{"bademail", "bad@gmail"}

	for _, testEmail := range tests {
		var player = Player{Email: testEmail, Password: "testpassword", Player_name: "Test Player"}

		t.Logf("Testing '%s'", testEmail)

		err := player.Validate()

		assert.NotNil(t, err)
	}
}

func TestHashPassword(t *testing.T) {
	var tests = []string{"mypass", "somepass", "som1pass"}

	for _, pass := range tests {
		hash, err := hashPassword(pass)

		assert.Nil(t, err)

		err = validatePassword(pass, hash)

		assert.Nil(t, err)

	}

}

func TestGenerateUUID(t *testing.T) {
	id := generateUUID()
	_, err := uuid.Parse(id)

	assert.Nil(t, err)

}
