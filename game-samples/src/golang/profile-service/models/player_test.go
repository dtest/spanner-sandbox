package models

import (
	"testing"
)

// var validate *validator.Validate

func TestValidEmails(t *testing.T) {
	var tests = []string{"good@gmail.com", "good@somedomain.net", "good.email@somedomain.org"}

	for _, testEmail := range tests {
		var player = Player{Email: testEmail, Password: "testpassword", Player_name: "Test Player"}

		t.Logf("Testing '%s'", testEmail)

		err := player.Validate()

		if err != nil {
			t.Error(
				"For", testEmail,
				"expected", nil,
				"got", err,
			)
		}
	}
}

func TestInvalidEmails(t *testing.T) {
	var tests = []string{"bademail", "bad@gmail"}

	for _, testEmail := range tests {
		var player = Player{Email: testEmail, Password: "testpassword", Player_name: "Test Player"}

		t.Logf("Testing '%s'", testEmail)

		err := player.Validate()

		if err == nil {
			t.Error(
				"For", testEmail,
				"expected", "Invalid email",
				"got", err,
			)
		}
	}
}

func TestHashPassword(t *testing.T) {
	var tests = []string{"mypass", "somepass", "som1pass"}

	for _, pass := range tests {
		hash, err := hashPassword(pass)

		if err != nil {
			t.Error(
				"For", pass,
				"expected", "no error",
				"got", err,
			)
		}

		err = validatePassword(pass, hash)

		if err != nil {
			t.Error(
				"For", pass,
				"expected", "no error",
				"got", err,
			)
		}
	}

}
