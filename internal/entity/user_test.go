package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "john@doe.com", "abc123")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@doe.com", user.Email)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
}

/* This is kinda of useless 🤷‍♂️ */
func TestNewUserWhenInvalidPassword(t *testing.T) {
	password := make([]byte, 73)
	user, err := NewUser("John Doe", "john@doe.com", string(password))
	assert.NotNil(t, err)
	assert.Nil(t, user)
}

func TestUserValidatePassword(t *testing.T) {
	password := "abc123"
	user, err := NewUser("John Doe", "john@doe.com", password)
	assert.Nil(t, err)
	assert.True(
		t,
		user.ValidatePassword(password),
		"Correct password should return true",
	)
	assert.False(
		t,
		user.ValidatePassword(password+"*"),
		"Incorrect password should return false",
	)
	assert.NotEqual(
		t,
		password,
		user.Password,
		"Password is not stored as plaintext",
	)
}
