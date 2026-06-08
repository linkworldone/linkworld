package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualNumberGenerator_Generate(t *testing.T) {
	gen := NewVirtualNumberGenerator()

	num, pwd, err := gen.Generate("US")
	assert.NoError(t, err)
	assert.NotEmpty(t, num)
	assert.Len(t, pwd, 8)
	assert.Contains(t, num, "+1")
}

func TestVirtualNumberGenerator_UnsupportedCountry(t *testing.T) {
	gen := NewVirtualNumberGenerator()

	num, _, err := gen.Generate("XX")
	assert.Error(t, err)
	assert.Empty(t, num)
}

func TestServiceError(t *testing.T) {
	err := ErrAlreadyRegistered
	assert.Equal(t, "User already registered", err.Error())
}

func TestErrServiceAlreadyActive(t *testing.T) {
	err := ErrServiceAlreadyActive
	assert.Equal(t, "Service already active", err.Error())
}