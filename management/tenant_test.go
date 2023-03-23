package management

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

func TestTenantManager(t *testing.T) {
	configureHTTPTestRecordings(t)

	initialSettings, err := api.Tenant.Read()
	assert.NoError(t, err)

	t.Cleanup(func() {
		initialSettings.SandboxVersionAvailable = nil
		initialSettings.UniversalLogin = nil
		initialSettings.Flags = nil
		err := api.Tenant.Update(initialSettings)
		require.NoError(t, err)
	})

	newTenantSettings := &Tenant{
		FriendlyName:          authok.String("My Example Tenant"),
		SupportURL:            authok.String("https://support.example.com"),
		SupportEmail:          authok.String("support@example.com"),
		DefaultRedirectionURI: authok.String("https://example.com/login"),
		SessionLifetime:       authok.Float64(1080),
		IdleSessionLifetime:   authok.Float64(720.2), // will be rounded off
		SessionCookie: &TenantSessionCookie{
			Mode: authok.String("non-persistent"),
		},
		AllowedLogoutURLs:       &[]string{"https://app.com/logout", "http://localhost/logout"},
		EnabledLocales:          &[]string{"fr", "en", "es"},
		SandboxVersionAvailable: nil,
	}
	err = api.Tenant.Update(newTenantSettings)
	assert.NoError(t, err)

	actualTenantSettings, err := api.Tenant.Read()
	assert.NoError(t, err)
	assert.Equal(t, newTenantSettings.GetFriendlyName(), actualTenantSettings.GetFriendlyName())
	assert.Equal(t, newTenantSettings.GetIdleSessionLifetime(), actualTenantSettings.GetIdleSessionLifetime())
	assert.Equal(t, newTenantSettings.GetIdleSessionLifetime(), 720.0) // it got rounded off
	assert.Equal(t, newTenantSettings.GetSessionLifetime(), actualTenantSettings.GetSessionLifetime())
	assert.Equal(t, newTenantSettings.GetSupportEmail(), actualTenantSettings.GetSupportEmail())
	assert.Equal(t, newTenantSettings.GetSupportURL(), actualTenantSettings.GetSupportURL())
	assert.Equal(t, newTenantSettings.SessionCookie.GetMode(), actualTenantSettings.SessionCookie.GetMode())
	assert.Equal(t, newTenantSettings.GetAllowedLogoutURLs(), actualTenantSettings.GetAllowedLogoutURLs())
	assert.Equal(t, newTenantSettings.GetEnabledLocales(), actualTenantSettings.GetEnabledLocales())
	assert.Equal(t, newTenantSettings.GetSandboxVersion(), actualTenantSettings.GetSandboxVersion())
}

func TestTenant_MarshalJSON(t *testing.T) {
	for tenant, expected := range map[*Tenant]string{
		{}:                                          `{}`,
		{SessionLifetime: authok.Float64(1.2)}:      `{"session_lifetime":1}`,
		{SessionLifetime: authok.Float64(1.19)}:     `{"session_lifetime":1}`,
		{SessionLifetime: authok.Float64(1)}:        `{"session_lifetime":1}`,
		{SessionLifetime: authok.Float64(720)}:      `{"session_lifetime":720}`,
		{IdleSessionLifetime: authok.Float64(1)}:    `{"idle_session_lifetime":1}`,
		{IdleSessionLifetime: authok.Float64(1.2)}:  `{"idle_session_lifetime":1}`,
		{SessionLifetime: authok.Float64(0.25)}:     `{"session_lifetime_in_minutes":15}`,
		{SessionLifetime: authok.Float64(0.5)}:      `{"session_lifetime_in_minutes":30}`,
		{SessionLifetime: authok.Float64(0.99)}:     `{"session_lifetime_in_minutes":59}`,
		{IdleSessionLifetime: authok.Float64(0.25)}: `{"idle_session_lifetime_in_minutes":15}`,
		{AllowedLogoutURLs: nil}:                    `{}`,
		{AllowedLogoutURLs: &[]string{}}:            `{"allowed_logout_urls":[]}`,
	} {
		payload, err := json.Marshal(tenant)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(payload))
	}
}

func TestTenantUniversalLoginColors_MarshalJSON(t *testing.T) {
	for tenantUniversalLoginColors, expected := range map[*TenantUniversalLoginColors]string{
		{}: `{}`,
		{
			PageBackground: authok.String("#ffffff"),
		}: `{"page_background":"#ffffff"}`,
		{
			PageBackgroundGradient: &BrandingPageBackgroundGradient{
				Type:        authok.String("linear-gradient"),
				Start:       authok.String("#ffffff"),
				End:         authok.String("#000000"),
				AngleDegree: authok.Int(3),
			},
		}: `{"page_background":{"type":"linear-gradient","start":"#ffffff","end":"#000000","angle_deg":3}}`,
	} {
		payload, err := json.Marshal(tenantUniversalLoginColors)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(payload))
	}
}

func TestTenantUniversalLoginColors_UnmarshalJSON(t *testing.T) {
	for jsonBody, expected := range map[string]*TenantUniversalLoginColors{
		`{}`: {},
		`{"page_background":"#ffffff"}`: {
			PageBackground: authok.String("#ffffff"),
		},
		`{"page_background":{"type":"linear-gradient","start":"#ffffff","end":"#000000","angle_deg":3}}`: {
			PageBackgroundGradient: &BrandingPageBackgroundGradient{
				Type:        authok.String("linear-gradient"),
				Start:       authok.String("#ffffff"),
				End:         authok.String("#000000"),
				AngleDegree: authok.Int(3),
			},
		},
	} {
		var actual TenantUniversalLoginColors
		err := json.Unmarshal([]byte(jsonBody), &actual)
		assert.NoError(t, err)
		assert.Equal(t, &actual, expected)
	}
}
