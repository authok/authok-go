package management

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

func TestRuleManager_Create(t *testing.T) {
	configureHTTPTestRecordings(t)

	rule := &Rule{
		Name:    authok.String("test-rule"),
		Script:  authok.String("function (user, context, callback) { callback(null, user, context); }"),
		Enabled: authok.Bool(false),
	}

	err := api.Rule.Create(rule)
	assert.NoError(t, err)
	assert.NotEmpty(t, rule.GetID())

	t.Cleanup(func() {
		cleanupRule(t, rule.GetID())
	})
}

func TestRuleManager_Read(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedRule := givenARule(t)

	actualRule, err := api.Rule.Read(expectedRule.GetID())

	assert.NoError(t, err)
	assert.Equal(t, expectedRule, actualRule)
}

func TestRuleManager_Update(t *testing.T) {
	configureHTTPTestRecordings(t)

	rule := givenARule(t)
	updatedRule := &Rule{
		Order:   authok.Int(5),
		Enabled: authok.Bool(true),
	}

	err := api.Rule.Update(rule.GetID(), updatedRule)
	assert.NoError(t, err)

	actualRule, err := api.Rule.Read(rule.GetID())
	assert.NoError(t, err)
	assert.Equal(t, updatedRule.GetOrder(), actualRule.GetOrder())
	assert.Equal(t, updatedRule.GetEnabled(), actualRule.GetEnabled())
}

func TestRuleManager_Delete(t *testing.T) {
	configureHTTPTestRecordings(t)

	rule := givenARule(t)

	err := api.Rule.Delete(rule.GetID())
	assert.NoError(t, err)

	actualRule, err := api.Rule.Read(rule.GetID())
	assert.Empty(t, actualRule)
	assert.Error(t, err)
	assert.Implements(t, (*Error)(nil), err)
	assert.Equal(t, http.StatusNotFound, err.(Error).Status())
}

func TestRuleManager_List(t *testing.T) {
	configureHTTPTestRecordings(t)

	rule := givenARule(t)

	ruleList, err := api.Rule.List(IncludeFields("id"))

	assert.NoError(t, err)
	assert.Contains(t, ruleList.Rules, &Rule{ID: rule.ID})
}

func givenARule(t *testing.T) *Rule {
	t.Helper()

	rule := &Rule{
		Name:    authok.String(fmt.Sprintf("test-rule%d", rand.Intn(999))),
		Script:  authok.String("function (user, context, callback) { callback(null, user, context); }"),
		Enabled: authok.Bool(false),
	}

	err := api.Rule.Create(rule)
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupRule(t, rule.GetID())
	})

	return rule
}

func cleanupRule(t *testing.T, ruleID string) {
	t.Helper()

	err := api.Rule.Delete(ruleID)
	require.NoError(t, err)
}
