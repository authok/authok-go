package management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

func TestEmailProviderJSON(t *testing.T) {
	var jsonTestCases = []struct {
		name          string
		emailProvider *EmailProvider
		json          string
	}{
		{
			name:          "it can %s an empty string",
			emailProvider: &EmailProvider{},
			json:          `{}`,
		},
		{
			name: "it can %s a mandrill email provider",
			emailProvider: &EmailProvider{
				Name:               authok.String("mandrill"),
				Enabled:            authok.Bool(true),
				DefaultFromAddress: authok.String("accounts@example.com"),
				Credentials: &EmailProviderCredentialsMandrill{
					APIKey: authok.String("secretApiKey"),
				},
				Settings: &EmailProviderSettingsMandrill{
					Message: &EmailProviderSettingsMandrillMessage{
						ViewContentLink: authok.Bool(true),
					},
				},
			},
			json: `{"name":"mandrill","enabled":true,"default_from_address":"accounts@example.com","credentials":{"api_key":"secretApiKey"},"settings":{"message":{"view_content_link":true}}}`,
		},
		{
			name: "it can %s an ses email provider",
			emailProvider: &EmailProvider{
				Name:               authok.String("ses"),
				Enabled:            authok.Bool(true),
				DefaultFromAddress: authok.String("accounts@example.com"),
				Credentials: &EmailProviderCredentialsSES{
					AccessKeyID:     authok.String("accessKey"),
					SecretAccessKey: authok.String("secret"),
					Region:          authok.String("eu-west-2"),
				},
				Settings: &EmailProviderSettingsSES{
					Message: &EmailProviderSettingsSESMessage{
						ConfigurationSetName: authok.String("example"),
					},
				},
			},
			json: `{"name":"ses","enabled":true,"default_from_address":"accounts@example.com","credentials":{"accessKeyId":"accessKey","secretAccessKey":"secret","region":"eu-west-2"},"settings":{"message":{"configuration_set_name":"example"}}}`,
		},
		{
			name: "it can %s a sendgrid email provider",
			emailProvider: &EmailProvider{
				Name:               authok.String("sendgrid"),
				Enabled:            authok.Bool(true),
				DefaultFromAddress: authok.String("accounts@example.com"),
				Credentials: &EmailProviderCredentialsSendGrid{
					APIKey: authok.String("apiKey"),
				},
			},
			json: `{"name":"sendgrid","enabled":true,"default_from_address":"accounts@example.com","credentials":{"api_key":"apiKey"}}`,
		},
		{
			name: "it can %s a sparkpost email provider",
			emailProvider: &EmailProvider{
				Name:               authok.String("sparkpost"),
				Enabled:            authok.Bool(true),
				DefaultFromAddress: authok.String("accounts@example.com"),
				Credentials: &EmailProviderCredentialsSparkPost{
					APIKey: authok.String("apiKey"),
					Region: authok.String("eu"),
				},
			},
			json: `{"name":"sparkpost","enabled":true,"default_from_address":"accounts@example.com","credentials":{"api_key":"apiKey","region":"eu"}}`,
		},
		{
			name: "it can %s a mailgun email provider",
			emailProvider: &EmailProvider{
				Name:               authok.String("mailgun"),
				Enabled:            authok.Bool(true),
				DefaultFromAddress: authok.String("accounts@example.com"),
				Credentials: &EmailProviderCredentialsMailgun{
					APIKey: authok.String("apiKey"),
					Region: authok.String("eu"),
					Domain: authok.String("example.com"),
				},
			},
			json: `{"name":"mailgun","enabled":true,"default_from_address":"accounts@example.com","credentials":{"api_key":"apiKey","domain":"example.com","region":"eu"}}`,
		},
		{
			name: "it can %s an smtp email provider",
			emailProvider: &EmailProvider{
				Name:               authok.String("smtp"),
				Enabled:            authok.Bool(true),
				DefaultFromAddress: authok.String("accounts@example.com"),
				Credentials: &EmailProviderCredentialsSMTP{
					SMTPHost: authok.String("example.com"),
					SMTPPort: authok.Int(3000),
					SMTPUser: authok.String("user"),
					SMTPPass: authok.String("pass"),
				},
				Settings: &EmailProviderSettingsSMTP{
					Headers: &EmailProviderSettingsSMTPHeaders{
						XMCViewContentLink:   authok.String("true"),
						XSESConfigurationSet: authok.String("example"),
					},
				},
			},
			json: `{"name":"smtp","enabled":true,"default_from_address":"accounts@example.com","credentials":{"smtp_host":"example.com","smtp_port":3000,"smtp_user":"user","smtp_pass":"pass"},"settings":{"headers":{"X-MC-ViewContentLink":"true","X-SES-Configuration-Set":"example"}}}`,
		},
	}

	for _, testCase := range jsonTestCases {
		t.Run(fmt.Sprintf(testCase.name, "marshal"), func(t *testing.T) {
			actualJSON, err := json.Marshal(testCase.emailProvider)
			assert.NoError(t, err)
			assert.Equal(t, testCase.json, string(actualJSON))
		})
	}

	for _, testCase := range jsonTestCases {
		t.Run(fmt.Sprintf(testCase.name, "unmarshal"), func(t *testing.T) {
			var actualEmailProvider EmailProvider
			err := json.Unmarshal([]byte(testCase.json), &actualEmailProvider)
			assert.NoError(t, err)
			assert.Equal(t, testCase.emailProvider, &actualEmailProvider)
		})
	}
}

func TestEmailProviderManager_Create(t *testing.T) {
	configureHTTPTestRecordings(t)

	emailProvider := &EmailProvider{
		Name:               authok.String("smtp"),
		Enabled:            authok.Bool(true),
		DefaultFromAddress: authok.String("no-reply@example.com"),
		Credentials: &EmailProviderCredentialsSMTP{
			SMTPHost: authok.String("smtp.example.com"),
			SMTPPort: authok.Int(587),
			SMTPUser: authok.String("user"),
			SMTPPass: authok.String("pass"),
		},
	}

	err := api.EmailProvider.Create(emailProvider)
	assert.NoError(t, err)

	t.Cleanup(func() {
		cleanupEmailProvider(t)
	})
}

func TestEmailProviderManager_Read(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedEmailProvider := givenAnEmailProvider(t)

	actualEmailProvider, err := api.EmailProvider.Read()

	assert.NoError(t, err)
	assert.Equal(t, expectedEmailProvider.GetName(), actualEmailProvider.GetName())
	assert.Equal(t, expectedEmailProvider.GetEnabled(), actualEmailProvider.GetEnabled())
	assert.Equal(t, expectedEmailProvider.GetDefaultFromAddress(), actualEmailProvider.GetDefaultFromAddress())
	assert.IsType(t, &EmailProviderCredentialsSMTP{}, expectedEmailProvider.Credentials)
}

func TestEmailProviderManager_Update(t *testing.T) {
	configureHTTPTestRecordings(t)

	emailProvider := givenAnEmailProvider(t)

	emailProvider.Enabled = authok.Bool(false)
	emailProvider.DefaultFromAddress = authok.String("info@example.com")

	err := api.EmailProvider.Update(emailProvider)
	assert.NoError(t, err)

	actualEmailProvider, err := api.EmailProvider.Read()
	assert.NoError(t, err)

	assert.False(t, actualEmailProvider.GetEnabled())
	assert.Equal(t, "info@example.com", actualEmailProvider.GetDefaultFromAddress())
}

func TestEmailProviderManager_Delete(t *testing.T) {
	configureHTTPTestRecordings(t)

	givenAnEmailProvider(t)

	err := api.Email.Delete()
	assert.NoError(t, err)

	_, err = api.Email.Read()
	assert.Error(t, err)
	assert.Implements(t, (*Error)(nil), err)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status())
}

func givenAnEmailProvider(t *testing.T) *EmailProvider {
	t.Helper()

	emailProvider := &EmailProvider{
		Name:               authok.String("smtp"),
		Enabled:            authok.Bool(true),
		DefaultFromAddress: authok.String("no-reply@example.com"),
		Credentials: &EmailProviderCredentialsSMTP{
			SMTPHost: authok.String("smtp.example.com"),
			SMTPPort: authok.Int(587),
			SMTPUser: authok.String("user"),
			SMTPPass: authok.String("pass"),
		},
	}

	err := api.EmailProvider.Create(emailProvider)
	if err != nil {
		if err.(Error).Status() != http.StatusConflict {
			t.Error(err)
		}
	}

	t.Cleanup(func() {
		cleanupEmailProvider(t)
	})

	return emailProvider
}

func cleanupEmailProvider(t *testing.T) {
	t.Helper()

	err := api.EmailProvider.Delete()
	require.NoError(t, err)
}
