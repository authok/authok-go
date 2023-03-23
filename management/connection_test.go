package management

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

var connectionTestCases = []connectionTestCase{
	{
		name: "Authok Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Authok-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("authok"),
		},
		options: &ConnectionOptions{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "GoogleOAuth2 Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-GoogleOAuth2-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("google-oauth2"),
		},
		options: &ConnectionOptionsGoogleOAuth2{
			AllowedAudiences: &[]string{
				"example.com",
				"api.example.com",
			},
			Profile:  authok.Bool(true),
			Calendar: authok.Bool(true),
			Youtube:  authok.Bool(false),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "GoogleApps Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-GoogleApps-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("google-apps"),
		},
		options: &ConnectionOptionsGoogleApps{
			Domain:          authok.String("example.com"),
			TenantDomain:    authok.String("example.com"),
			BasicProfile:    authok.Bool(true),
			ExtendedProfile: authok.Bool(true),
			Groups:          authok.Bool(true),
			Admin:           authok.Bool(true),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "Email Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Email-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("email"),
		},
		options: &ConnectionOptionsEmail{
			Email: &ConnectionOptionsEmailSettings{
				Syntax:  authok.String("liquid"),
				From:    authok.String("{{application.name}} <test@example.com>"),
				Subject: authok.String("Email Login - {{application.name}}"),
				Body:    authok.String("<html><body>email contents</body></html>"),
			},
			OTP: &ConnectionOptionsOTP{
				TimeStep: authok.Int(100),
				Length:   authok.Int(4),
			},
			AuthParams: map[string]string{
				"scope": "openid profile",
			},
			BruteForceProtection: authok.Bool(true),
			DisableSignup:        authok.Bool(true),
			Name:                 authok.String("Test-Connection-Email"),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "SMS Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-SMS-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("sms"),
		},
		options: &ConnectionOptionsSMS{
			From:     authok.String("+17777777777"),
			Template: authok.String("Your verification code is { code }}"),
			Syntax:   authok.String("liquid"),
			OTP: &ConnectionOptionsOTP{
				TimeStep: authok.Int(110),
				Length:   authok.Int(5),
			},
			AuthParams: map[string]string{
				"scope": "openid profile",
			},
			BruteForceProtection: authok.Bool(true),
			DisableSignup:        authok.Bool(true),
			Name:                 authok.String("Test-Connection-SMS"),
			TwilioSID:            authok.String("abc132asdfasdf56"),
			TwilioToken:          authok.String("234127asdfsada23"),
			MessagingServiceSID:  authok.String("273248090982390423"),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "Custom SMS Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Custom-SMS-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("sms"),
		},
		options: &ConnectionOptionsSMS{
			From:     authok.String("+17777777777"),
			Template: authok.String("Your verification code is { code }}"),
			Syntax:   authok.String("liquid"),
			OTP: &ConnectionOptionsOTP{
				TimeStep: authok.Int(110),
				Length:   authok.Int(5),
			},
			BruteForceProtection: authok.Bool(true),
			DisableSignup:        authok.Bool(true),
			Name:                 authok.String("Test-Connection-Custom-SMS"),
			Provider:             authok.String("sms_gateway"),
			GatewayURL:           authok.String("https://test.com/sms-gateway"),
			GatewayAuthentication: &ConnectionGatewayAuthentication{
				Method:              authok.String("bearer"),
				Subject:             authok.String("test.us.authok.com:sms"),
				Audience:            authok.String("test.com/sms-gateway"),
				Secret:              authok.String("my-secret"),
				SecretBase64Encoded: authok.Bool(false),
			},
			ForwardRequestInfo: authok.Bool(true),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "SAML Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-SAML-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("samlp"),
		},
		options: &ConnectionOptionsSAML{
			SignInEndpoint: authok.String("https://saml.identity/provider"),
			SigningCert: authok.String(`-----BEGIN CERTIFICATE-----
MIID6TCCA1ICAQEwDQYJKoZIhvcNAQEFBQAwgYsxCzAJBgNVBAYTAlVTMRMwEQYD
VQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRQwEgYDVQQK
EwtHb29nbGUgSW5jLjEMMAoGA1UECxMDRW5nMQwwCgYDVQQDEwNhZ2wxHTAbBgkq
hkiG9w0BCQEWDmFnbEBnb29nbGUuY29tMB4XDTA5MDkwOTIyMDU0M1oXDTEwMDkw
OTIyMDU0M1owajELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAf
BgNVBAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEjMCEGA1UEAxMaZXVyb3Bh
LnNmby5jb3JwLmdvb2dsZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIK
AoICAQC6pgYt7/EibBDumASF+S0qvqdL/f+nouJw2T1Qc8GmXF/iiUcrsgzh/Fd8
pDhz/T96Qg9IyR4ztuc2MXrmPra+zAuSf5bevFReSqvpIt8Duv0HbDbcqs/XKPfB
uMDe+of7a9GCywvAZ4ZUJcp0thqD9fKTTjUWOBzHY1uNE4RitrhmJCrbBGXbJ249
bvgmb7jgdInH2PU7PT55hujvOoIsQW2osXBFRur4pF1wmVh4W4lTLD6pjfIMUcML
ICHEXEN73PDic8KS3EtNYCwoIld+tpIBjE1QOb1KOyuJBNW6Esw9ALZn7stWdYcE
qAwvv20egN2tEXqj7Q4/1ccyPZc3PQgC3FJ8Be2mtllM+80qf4dAaQ/fWvCtOrQ5
pnfe9juQvCo8Y0VGlFcrSys/MzSg9LJ/24jZVgzQved/Qupsp89wVidwIzjt+WdS
fyWfH0/v1aQLvu5cMYuW//C0W2nlYziL5blETntM8My2ybNARy3ICHxCBv2RNtPI
WQVm+E9/W5rwh2IJR4DHn2LHwUVmT/hHNTdBLl5Uhwr4Wc7JhE7AVqb14pVNz1lr
5jxsp//ncIwftb7mZQ3DF03Yna+jJhpzx8CQoeLT6aQCHyzmH68MrHHT4MALPyUs
Pomjn71GNTtDeWAXibjCgdL6iHACCF6Htbl0zGlG0OAK+bdn0QIDAQABMA0GCSqG
SIb3DQEBBQUAA4GBAOKnQDtqBV24vVqvesL5dnmyFpFPXBn3WdFfwD6DzEb21UVG
5krmJiu+ViipORJPGMkgoL6BjU21XI95VQbun5P8vvg8Z+FnFsvRFY3e1CCzAVQY
ZsUkLw2I7zI/dNlWdB8Xp7v+3w9sX5N3J/WuJ1KOO5m26kRlHQo7EzT3974g
-----END CERTIFICATE-----`),
			TenantDomain: authok.String("example.com"),
			FieldsMap: map[string]interface{}{
				"email":       "EmailAddress",
				"given_name":  "FirstName",
				"family_name": "LastName",
			},
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "AD Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-AD-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("ad"),
		},
		options: &ConnectionOptionsAD{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "ADFS Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-ADFS-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("adfs"),
		},
		options: &ConnectionOptionsADFS{
			FedMetadataXML: authok.String(`<?xml version="1.0" encoding="utf-8"?>
<EntityDescriptor entityID="https://example.com"
                  xmlns="urn:oasis:names:tc:SAML:2.0:metadata">
    <RoleDescriptor xsi:type="fed:ApplicationServiceType"
                    protocolSupportEnumeration="http://docs.oasis-open.org/wsfed/federation/200706"
                    xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                    xmlns:fed="http://docs.oasis-open.org/wsfed/federation/200706">
        <fed:TargetScopes>
            <wsa:EndpointReference xmlns:wsa="http://www.w3.org/2005/08/addressing">
                <wsa:Address>https://adfs.provider/</wsa:Address>
            </wsa:EndpointReference>
        </fed:TargetScopes>
        <fed:ApplicationServiceEndpoint>
            <wsa:EndpointReference xmlns:wsa="http://www.w3.org/2005/08/addressing">
                <wsa:Address>https://adfs.provider/wsfed</wsa:Address>
            </wsa:EndpointReference>
        </fed:ApplicationServiceEndpoint>
        <fed:PassiveRequestorEndpoint>
            <wsa:EndpointReference xmlns:wsa="http://www.w3.org/2005/08/addressing">
                <wsa:Address>https://adfs.provider/wsfed</wsa:Address>
            </wsa:EndpointReference>
        </fed:PassiveRequestorEndpoint>
    </RoleDescriptor>
    <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
                             Location="https://adfs.provider/sign_out"/>
        <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
                             Location="https://adfs.provider/sign_in"/>
    </IDPSSODescriptor>
</EntityDescriptor>
`),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "Facebook Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Facebook-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("facebook"),
		},
		options: &ConnectionOptionsFacebook{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "Apple Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Apple-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("apple"),
		},
		options: &ConnectionOptionsApple{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "LinkedIn Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-LinkedIn-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("linkedin"),
		},
		options: &ConnectionOptionsLinkedin{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "GitHub Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-GitHub-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("github"),
		},
		options: &ConnectionOptionsGitHub{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "WindowsLive Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-WindowsLive-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("windowslive"),
		},
		options: &ConnectionOptionsWindowsLive{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "SalesForce Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-SalesForce-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("salesforce"),
		},
		options: &ConnectionOptionsSalesforce{
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "OIDC Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-OIDC-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("oidc"),
		},
		options: &ConnectionOptionsOIDC{
			ClientID:              authok.String("4ef8d976-71bd-4473-a7ce-087d3f0fafd8"),
			Scope:                 authok.String("openid"),
			Issuer:                authok.String("https://example.com"),
			AuthorizationEndpoint: authok.String("https://example.com"),
			JWKSURI:               authok.String("https://example.com/jwks"),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "Okta Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Okta-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("okta"),
		},
		options: &ConnectionOptionsOkta{
			ClientID:              authok.String("4ef8d976-71bd-4473-a7ce-087d3f0fafd8"),
			ClientSecret:          authok.String("mySecret"),
			Scope:                 authok.String("openid"),
			Domain:                authok.String("domain.okta.com"),
			Issuer:                authok.String("https://example.com"),
			AuthorizationEndpoint: authok.String("https://example.com"),
			JWKSURI:               authok.String("https://example.com/jwks"),
			UpstreamParams: map[string]interface{}{
				"screen_name": map[string]interface{}{
					"alias": "login_hint",
				},
			},
		},
	},
	{
		name: "Ping Federate Connection",
		connection: Connection{
			Name:     authok.Stringf("Test-Ping-Federate-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("pingfederate"),
		},
		options: &ConnectionOptionsPingFederate{
			PingFederateBaseURL: authok.String("https://ping.example.com"),
			SigningCert: authok.String(`-----BEGIN CERTIFICATE-----
MIID6TCCA1ICAQEwDQYJKoZIhvcNAQEFBQAwgYsxCzAJBgNVBAYTAlVTMRMwEQYD
VQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRQwEgYDVQQK
EwtHb29nbGUgSW5jLjEMMAoGA1UECxMDRW5nMQwwCgYDVQQDEwNhZ2wxHTAbBgkq
hkiG9w0BCQEWDmFnbEBnb29nbGUuY29tMB4XDTA5MDkwOTIyMDU0M1oXDTEwMDkw
OTIyMDU0M1owajELMAkGA1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAf
BgNVBAoTGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEjMCEGA1UEAxMaZXVyb3Bh
LnNmby5jb3JwLmdvb2dsZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIK
AoICAQC6pgYt7/EibBDumASF+S0qvqdL/f+nouJw2T1Qc8GmXF/iiUcrsgzh/Fd8
pDhz/T96Qg9IyR4ztuc2MXrmPra+zAuSf5bevFReSqvpIt8Duv0HbDbcqs/XKPfB
uMDe+of7a9GCywvAZ4ZUJcp0thqD9fKTTjUWOBzHY1uNE4RitrhmJCrbBGXbJ249
bvgmb7jgdInH2PU7PT55hujvOoIsQW2osXBFRur4pF1wmVh4W4lTLD6pjfIMUcML
ICHEXEN73PDic8KS3EtNYCwoIld+tpIBjE1QOb1KOyuJBNW6Esw9ALZn7stWdYcE
qAwvv20egN2tEXqj7Q4/1ccyPZc3PQgC3FJ8Be2mtllM+80qf4dAaQ/fWvCtOrQ5
pnfe9juQvCo8Y0VGlFcrSys/MzSg9LJ/24jZVgzQved/Qupsp89wVidwIzjt+WdS
fyWfH0/v1aQLvu5cMYuW//C0W2nlYziL5blETntM8My2ybNARy3ICHxCBv2RNtPI
WQVm+E9/W5rwh2IJR4DHn2LHwUVmT/hHNTdBLl5Uhwr4Wc7JhE7AVqb14pVNz1lr
5jxsp//ncIwftb7mZQ3DF03Yna+jJhpzx8CQoeLT6aQCHyzmH68MrHHT4MALPyUs
Pomjn71GNTtDeWAXibjCgdL6iHACCF6Htbl0zGlG0OAK+bdn0QIDAQABMA0GCSqG
SIb3DQEBBQUAA4GBAOKnQDtqBV24vVqvesL5dnmyFpFPXBn3WdFfwD6DzEb21UVG
5krmJiu+ViipORJPGMkgoL6BjU21XI95VQbun5P8vvg8Z+FnFsvRFY3e1CCzAVQY
ZsUkLw2I7zI/dNlWdB8Xp7v+3w9sX5N3J/WuJ1KOO5m26kRlHQo7EzT3974g
-----END CERTIFICATE-----`),
			SignSAMLRequest:    authok.Bool(true),
			SignatureAlgorithm: authok.String("rsa-sha256"),
			DigestAlgorithm:    authok.String("sha256"),
		},
	},
}

type connectionTestCase struct {
	name       string
	connection Connection
	options    interface{}
}

func TestConnectionManager_Create(t *testing.T) {
	for _, testCase := range connectionTestCases {
		t.Run("It can successfully create a "+testCase.name, func(t *testing.T) {
			configureHTTPTestRecordings(t)

			expectedConnection := testCase.connection
			expectedConnection.Options = testCase.options

			err := api.Connection.Create(&expectedConnection)

			assert.NoError(t, err)
			assert.NotEmpty(t, expectedConnection.GetID())
			assert.IsType(t, testCase.options, expectedConnection.Options)

			t.Cleanup(func() {
				cleanupConnection(t, expectedConnection.GetID())
			})
		})
	}
}

func TestConnectionManager_Read(t *testing.T) {
	for _, testCase := range connectionTestCases {
		t.Run("It can successfully read a "+testCase.name, func(t *testing.T) {
			configureHTTPTestRecordings(t)

			expectedConnection := givenAConnection(t, testCase)

			actualConnection, err := api.Connection.Read(expectedConnection.GetID())

			assert.NoError(t, err)
			assert.Equal(t, expectedConnection.GetID(), actualConnection.GetID())
			assert.Equal(t, expectedConnection.GetName(), actualConnection.GetName())
			assert.Equal(t, expectedConnection.GetStrategy(), actualConnection.GetStrategy())
			assert.IsType(t, testCase.options, actualConnection.Options)

			t.Cleanup(func() {
				cleanupConnection(t, expectedConnection.GetID())
			})
		})
	}
}

func TestConnectionManager_ReadByName(t *testing.T) {
	for _, testCase := range connectionTestCases {
		t.Run("It can successfully find a "+testCase.name+" by its name", func(t *testing.T) {
			configureHTTPTestRecordings(t)

			expectedConnection := givenAConnection(t, testCase)

			actualConnection, err := api.Connection.ReadByName(expectedConnection.GetName())

			assert.NoError(t, err)
			assert.Equal(t, expectedConnection.GetID(), actualConnection.GetID())
			assert.Equal(t, expectedConnection.GetName(), actualConnection.GetName())
			assert.Equal(t, expectedConnection.GetStrategy(), actualConnection.GetStrategy())
			assert.IsType(t, testCase.options, actualConnection.Options)

			t.Cleanup(func() {
				cleanupConnection(t, expectedConnection.GetID())
			})
		})
	}

	t.Run("throw an error when connection name is empty", func(t *testing.T) {
		actualConnection, err := api.Connection.ReadByName("")

		assert.EqualError(t, err, "400 Bad Request: Name cannot be empty")
		assert.Empty(t, actualConnection)
	})
}

func TestConnectionManager_Update(t *testing.T) {
	for _, testCase := range connectionTestCases {
		t.Run("It can successfully update a "+testCase.name, func(t *testing.T) {
			if testCase.connection.GetStrategy() == "oidc" ||
				testCase.connection.GetStrategy() == "samlp" ||
				testCase.connection.GetStrategy() == "okta" ||
				testCase.connection.GetStrategy() == "adfs" ||
				testCase.connection.GetStrategy() == "pingfederate" {
				t.Skip("Skipping because we can't create an oidc, okta, samlp, adfs, or pingfederate connection with no options")
			}

			configureHTTPTestRecordings(t)

			connection := givenAConnection(t, connectionTestCase{connection: testCase.connection})

			connectionWithUpdatedOptions := &Connection{
				Options: testCase.options,
			}

			err := api.Connection.Update(connection.GetID(), connectionWithUpdatedOptions)
			assert.NoError(t, err)

			actualConnection, err := api.Connection.Read(connection.GetID())
			assert.NoError(t, err)
			assert.ObjectsAreEqualValues(testCase.options, actualConnection.Options)
		})
	}
}

func TestConnectionManager_Delete(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedConnection := givenAConnection(t, connectionTestCase{
		connection: Connection{
			Name:     authok.Stringf("Test-Authok-Connection-%d", time.Now().Unix()),
			Strategy: authok.String("authok"),
		},
	})

	err := api.Connection.Delete(expectedConnection.GetID())
	assert.NoError(t, err)

	actualConnection, err := api.Connection.Read(expectedConnection.GetID())
	assert.Nil(t, actualConnection)
	assert.Error(t, err)
	assert.Implements(t, (*Error)(nil), err)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status())
}

func TestConnectionManager_List(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedConnection := givenAConnection(t, connectionTestCase{
		connection: Connection{
			Name:     authok.Stringf("Test-Authok-Connection-List-%d", time.Now().Unix()),
			Strategy: authok.String("authok"),
		},
	})

	needle := &Connection{
		ID:                 expectedConnection.ID,
		IsDomainConnection: authok.Bool(false),
	}
	connectionList, err := api.Connection.List(IncludeFields("id"))
	assert.NoError(t, err)
	assert.Contains(t, connectionList.Connections, needle)
}

func TestConnectionOptionsScopes(t *testing.T) {
	t.Run("It can successfully set the scopes on the options of a OIDC connection", func(t *testing.T) {
		options := &ConnectionOptionsOIDC{}

		options.SetScopes(true, "foo", "bar", "baz")
		assert.Equal(t, []string{"bar", "baz", "foo"}, options.Scopes())

		options.SetScopes(false, "foo", "baz")
		assert.Equal(t, []string{"bar"}, options.Scopes())
	})

	t.Run("It can successfully set the scopes on the options of an OAuth2 connection", func(t *testing.T) {
		options := &ConnectionOptionsOAuth2{}

		options.SetScopes(true, "foo", "bar", "baz")
		assert.Equal(t, []string{"bar", "baz", "foo"}, options.Scopes())

		options.SetScopes(false, "foo", "baz")
		assert.Equal(t, []string{"bar"}, options.Scopes())
	})

	t.Run("It can successfully set the scopes on the options of an Okta connection", func(t *testing.T) {
		options := &ConnectionOptionsOkta{}

		options.SetScopes(true, "foo", "bar", "baz")
		assert.Equal(t, []string{"bar", "baz", "foo"}, options.Scopes())

		options.SetScopes(false, "foo", "baz")
		assert.Equal(t, []string{"bar"}, options.Scopes())
	})
}

func TestGoogleOauth2Connection_MarshalJSON(t *testing.T) {
	var emptySlice []string
	for connection, expected := range map[*ConnectionOptionsGoogleOAuth2]string{
		{AllowedAudiences: nil}:                                               `{}`,
		{AllowedAudiences: &emptySlice}:                                       `{"allowed_audiences":null}`,
		{AllowedAudiences: &[]string{}}:                                       `{"allowed_audiences":[]}`,
		{AllowedAudiences: &[]string{"foo", "bar"}}:                           `{"allowed_audiences":["foo","bar"]}`,
		{AllowedAudiences: &[]string{"foo", "bar"}, Email: authok.Bool(true)}: `{"email":true,"allowed_audiences":["foo","bar"]}`,
	} {
		payload, err := json.Marshal(connection)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(payload))
	}
}

func TestGoogleOauth2Connection_UnmarshalJSON(t *testing.T) {
	for expectedAsString, expected := range map[string]*ConnectionOptionsGoogleOAuth2{
		`{}`:                          {},
		`{"allowed_audiences": null}`: {},
		`{"allowed_audiences": ""}`:   {AllowedAudiences: &[]string{}},
		`{"allowed_audiences": []}`:   {AllowedAudiences: &[]string{}},
		`{"allowed_audiences": ["foo", "bar"], "scope": ["email"] }`: {AllowedAudiences: &[]string{"foo", "bar"}, Scope: []interface{}{"email"}},
	} {
		var actual *ConnectionOptionsGoogleOAuth2
		err := json.Unmarshal([]byte(expectedAsString), &actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	t.Run("Throws an unexpected type error", func(t *testing.T) {
		var actual *ConnectionOptionsGoogleOAuth2
		err := json.Unmarshal([]byte(`{"allowed_audiences": 1}`), &actual)
		assert.EqualError(t, err, "unexpected type for field allowed_audiences: float64")
	})
}

func cleanupConnection(t *testing.T, connectionID string) {
	t.Helper()

	err := api.Connection.Delete(connectionID)
	require.NoError(t, err)
}

func givenAConnection(t *testing.T, testCase connectionTestCase) *Connection {
	t.Helper()

	connection := testCase.connection
	connection.Options = testCase.options

	err := api.Connection.Create(&connection)
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupConnection(t, connection.GetID())
	})

	return &connection
}
