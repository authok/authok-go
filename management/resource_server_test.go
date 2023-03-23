package management

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

func TestResourceServer_Create(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedResourceServer := &ResourceServer{
		Name:                authok.Stringf("Test Resource Server (%s)", time.Now().Format(time.StampMilli)),
		Identifier:          authok.String("https://api.example.com/"),
		SigningAlgorithm:    authok.String("HS256"),
		TokenLifetime:       authok.Int(7200),
		TokenLifetimeForWeb: authok.Int(3600),
		Scopes: &[]ResourceServerScope{
			{
				Value:       authok.String("create:resource"),
				Description: authok.String("Create Resource"),
			},
		},
	}

	err := api.ResourceServer.Create(expectedResourceServer)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedResourceServer.GetID())

	t.Cleanup(func() {
		cleanupResourceServer(t, expectedResourceServer.GetID())
	})
}

func TestResourceServer_Read(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedResourceServer := givenAResourceServer(t)

	actualResourceServer, err := api.ResourceServer.Read(expectedResourceServer.GetID())

	assert.NoError(t, err)
	assert.Equal(t, expectedResourceServer.GetName(), actualResourceServer.GetName())
}

func TestResourceServer_Update(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedResourceServer := givenAResourceServer(t)

	resourceServerID := expectedResourceServer.GetID()

	expectedResourceServer.ID = nil         // Read-Only: Additional properties not allowed.
	expectedResourceServer.Identifier = nil // Read-Only: Additional properties not allowed.
	expectedResourceServer.SigningSecret = nil

	expectedResourceServer.AllowOfflineAccess = authok.Bool(true)
	expectedResourceServer.SigningAlgorithm = authok.String("RS256")
	expectedResourceServer.SkipConsentForVerifiableFirstPartyClients = authok.Bool(true)
	expectedResourceServer.TokenLifetime = authok.Int(7200)
	expectedResourceServer.TokenLifetimeForWeb = authok.Int(5400)
	scopes := append(expectedResourceServer.GetScopes(), ResourceServerScope{
		Value:       authok.String("update:resource"),
		Description: authok.String("Update Resource"),
	})
	expectedResourceServer.Scopes = &scopes

	err := api.ResourceServer.Update(resourceServerID, expectedResourceServer)

	assert.NoError(t, err)
	assert.Equal(t, expectedResourceServer.GetAllowOfflineAccess(), true)
	assert.Equal(t, expectedResourceServer.GetSigningAlgorithm(), "RS256")
	assert.Equal(t, expectedResourceServer.GetSkipConsentForVerifiableFirstPartyClients(), true)
	assert.Equal(t, expectedResourceServer.GetTokenLifetime(), 7200)
	assert.Equal(t, expectedResourceServer.GetTokenLifetimeForWeb(), 5400)
	assert.Equal(t, len(expectedResourceServer.GetScopes()), 2)
}

func TestResourceServer_Delete(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedResourceServer := givenAResourceServer(t)

	err := api.ResourceServer.Delete(expectedResourceServer.GetID())
	assert.NoError(t, err)

	actualResourceServer, err := api.ResourceServer.Read(expectedResourceServer.GetID())
	assert.Empty(t, actualResourceServer)
	assert.Error(t, err)
	assert.Implements(t, (*Error)(nil), err)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status())
}

func TestResourceServer_List(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedResourceServer := givenAResourceServer(t)

	resourceServerList, err := api.ResourceServer.List(IncludeFields("id"))

	assert.NoError(t, err)
	assert.Contains(t, resourceServerList.ResourceServers, &ResourceServer{ID: expectedResourceServer.ID})
}

func givenAResourceServer(t *testing.T) *ResourceServer {
	t.Helper()

	resourceServer := &ResourceServer{
		Name:                authok.Stringf("Test Resource Server (%s)", time.Now().Format(time.StampMilli)),
		Identifier:          authok.String("https://api.example.com/"),
		SigningAlgorithm:    authok.String("HS256"),
		TokenLifetime:       authok.Int(7200),
		TokenLifetimeForWeb: authok.Int(3600),
		Scopes: &[]ResourceServerScope{
			{
				Value:       authok.String("create:resource"),
				Description: authok.String("Create Resource"),
			},
		},
	}

	err := api.ResourceServer.Create(resourceServer)
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupResourceServer(t, resourceServer.GetID())
	})

	return resourceServer
}

func cleanupResourceServer(t *testing.T, resourceServerID string) {
	t.Helper()

	err := api.ResourceServer.Delete(resourceServerID)
	require.NoError(t, err)
}
