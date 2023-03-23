package management

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

func TestEmailManager_Create(t *testing.T) {
	configureHTTPTestRecordings(t)

	emailProvider := &Email{
		Name:               authok.String("smtp"),
		Enabled:            authok.Bool(true),
		DefaultFromAddress: authok.String("no-reply@example.com"),
		Credentials: &EmailCredentials{
			SMTPHost: authok.String("smtp.example.com"),
			SMTPPort: authok.Int(587),
			SMTPUser: authok.String("user"),
			SMTPPass: authok.String("pass"),
		},
	}

	err := api.Email.Create(emailProvider)
	assert.NoError(t, err)

	t.Cleanup(func() {
		cleanupEmail(t)
	})
}

func TestEmailManager_Read(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedEmailProvider := givenAnEmail(t)

	actualEmailProvider, err := api.Email.Read()

	assert.NoError(t, err)
	assert.Equal(t, expectedEmailProvider.GetName(), actualEmailProvider.GetName())
	assert.Equal(t, expectedEmailProvider.GetEnabled(), actualEmailProvider.GetEnabled())
	assert.Equal(t, expectedEmailProvider.GetDefaultFromAddress(), actualEmailProvider.GetDefaultFromAddress())
	assert.Equal(
		t,
		expectedEmailProvider.GetCredentials().GetSMTPUser(),
		actualEmailProvider.GetCredentials().GetSMTPUser(),
	)
	assert.Equal(
		t,
		"",
		actualEmailProvider.GetCredentials().GetSMTPPass(),
	) // Passwords are not returned from the Authok API.
}

func TestEmailManager_Update(t *testing.T) {
	configureHTTPTestRecordings(t)

	emailProvider := givenAnEmail(t)

	emailProvider.Enabled = authok.Bool(false)
	emailProvider.DefaultFromAddress = authok.String("info@example.com")

	err := api.Email.Update(emailProvider)
	assert.NoError(t, err)

	actualEmailProvider, err := api.Email.Read()
	assert.NoError(t, err)

	assert.False(t, actualEmailProvider.GetEnabled())
	assert.Equal(t, "info@example.com", actualEmailProvider.GetDefaultFromAddress())
}

func TestEmailManager_Delete(t *testing.T) {
	configureHTTPTestRecordings(t)

	givenAnEmail(t)

	err := api.Email.Delete()
	assert.NoError(t, err)

	_, err = api.Email.Read()
	assert.Error(t, err)
	assert.Implements(t, (*Error)(nil), err)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status())
}

func givenAnEmail(t *testing.T) *Email {
	t.Helper()

	emailProvider := &Email{
		Name:               authok.String("smtp"),
		Enabled:            authok.Bool(true),
		DefaultFromAddress: authok.String("no-reply@example.com"),
		Credentials: &EmailCredentials{
			SMTPHost: authok.String("smtp.example.com"),
			SMTPPort: authok.Int(587),
			SMTPUser: authok.String("user"),
			SMTPPass: authok.String("pass"),
		},
	}

	err := api.Email.Create(emailProvider)
	if err != nil {
		if err.(Error).Status() != http.StatusConflict {
			t.Error(err)
		}
	}

	t.Cleanup(func() {
		cleanupEmail(t)
	})

	return emailProvider
}

func cleanupEmail(t *testing.T) {
	t.Helper()

	err := api.Email.Delete()
	require.NoError(t, err)
}
