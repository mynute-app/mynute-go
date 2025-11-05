package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	t.Run("should return Resend adapter when provider is resend", func(t *testing.T) {
		os.Setenv("RESEND_API_KEY", "test-key")
		os.Setenv("RESEND_DEFAULT_FROM", "test@example.com")
		defer os.Unsetenv("RESEND_API_KEY")
		defer os.Unsetenv("RESEND_DEFAULT_FROM")

		provider, err := NewProvider(&ProviderOpts{Provider: "resend"})

		assert.NoError(t, err)
		assert.NotNil(t, provider)
		_, ok := provider.(*ResendAdapter)
		assert.True(t, ok)
	})

	t.Run("should return MailHog adapter when provider is mailhog", func(t *testing.T) {
		provider, err := NewProvider(&ProviderOpts{Provider: "mailhog"})

		assert.NoError(t, err)
		assert.NotNil(t, provider)
		_, ok := provider.(*MailHogAdapter)
		assert.True(t, ok)
	})

	t.Run("should return error for unsupported provider", func(t *testing.T) {
		provider, err := NewProvider(&ProviderOpts{Provider: "unsupported"})

		assert.Error(t, err)
		assert.Nil(t, provider)
		assert.Equal(t, "email provider (unsupported) not implemented", err.Error())
	})
}
