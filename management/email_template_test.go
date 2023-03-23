package management

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

func TestEmailTemplateManager_Create(t *testing.T) {
	configureHTTPTestRecordings(t)

	template := &EmailTemplate{
		Template:               authok.String("verify_email"),
		Body:                   authok.String("<html><body><h1>Verify your email</h1></body></html>"),
		From:                   authok.String("me@example.com"),
		ResultURL:              authok.String("https://www.example.com/verify-email"),
		Subject:                authok.String("Verify your email"),
		Syntax:                 authok.String("liquid"),
		Enabled:                authok.Bool(true),
		IncludeEmailInRedirect: authok.Bool(true),
	}

	err := api.EmailTemplate.Create(template)
	if err != nil {
		if err, ok := err.(Error); ok && err.Status() != http.StatusConflict {
			t.Error(err)
		}
	}

	t.Cleanup(func() {
		cleanupEmailTemplate(t, template.GetTemplate())
	})
}

func TestEmailTemplateManager_Read(t *testing.T) {
	configureHTTPTestRecordings(t)

	expectedTemplate := givenAnEmailTemplate(t)

	actualTemplate, err := api.EmailTemplate.Read(expectedTemplate.GetTemplate())

	assert.NoError(t, err)
	assert.ObjectsAreEqual(expectedTemplate, actualTemplate)
}

func TestEmailTemplateManager_Update(t *testing.T) {
	configureHTTPTestRecordings(t)

	template := givenAnEmailTemplate(t)

	expectedBody := "<html><body><h1>Let's get you verified!</h1></body></html>"
	expectedIncludeEmailInRedirect := false
	err := api.EmailTemplate.Update(
		template.GetTemplate(),
		&EmailTemplate{
			Body:                   &expectedBody,
			IncludeEmailInRedirect: &expectedIncludeEmailInRedirect,
		},
	)
	assert.NoError(t, err)

	actualTemplate, err := api.EmailTemplate.Read(template.GetTemplate())
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualTemplate.GetBody())
	assert.Equal(t, expectedIncludeEmailInRedirect, actualTemplate.GetIncludeEmailInRedirect())
}

func TestEmailTemplateManager_Replace(t *testing.T) {
	configureHTTPTestRecordings(t)

	givenAnEmailProvider(t)
	template := givenAnEmailTemplate(t)

	template.Template = authok.String("verify_email")
	template.Subject = authok.String("Let's get you verified!")
	template.Body = authok.String("<html><body><h1>Let's get you verified!</h1></body></html>")
	template.From = authok.String("someone@example.com")
	template.IncludeEmailInRedirect = authok.Bool(true)

	err := api.EmailTemplate.Replace(template.GetTemplate(), template)
	assert.NoError(t, err)

	actualTemplate, err := api.EmailTemplate.Read(template.GetTemplate())
	assert.NoError(t, err)

	assert.Equal(t, actualTemplate.GetBody(), template.GetBody())
	assert.Equal(t, actualTemplate.GetSubject(), template.GetSubject())
	assert.Equal(t, actualTemplate.GetFrom(), template.GetFrom())
	assert.Equal(t, actualTemplate.GetTemplate(), template.GetTemplate())
	assert.Equal(t, actualTemplate.GetIncludeEmailInRedirect(), template.GetIncludeEmailInRedirect())
}

func givenAnEmailTemplate(t *testing.T) *EmailTemplate {
	t.Helper()

	template := &EmailTemplate{
		Template:               authok.String("verify_email"),
		Body:                   authok.String("<html><body><h1>Verify your email</h1></body></html>"),
		From:                   authok.String("me@example.com"),
		ResultURL:              authok.String("https://www.example.com/verify-email"),
		Subject:                authok.String("Verify your email"),
		Syntax:                 authok.String("liquid"),
		Enabled:                authok.Bool(true),
		IncludeEmailInRedirect: authok.Bool(true),
	}

	err := api.EmailTemplate.Create(template)
	if err != nil {
		if err, ok := err.(Error); ok && err.Status() != http.StatusConflict {
			t.Error(err)
		}
	}

	t.Cleanup(func() {
		cleanupEmailTemplate(t, template.GetTemplate())
	})

	return template
}

func cleanupEmailTemplate(t *testing.T, templateName string) {
	t.Helper()

	err := api.EmailTemplate.Update(templateName, &EmailTemplate{Enabled: authok.Bool(false)})
	require.NoError(t, err)
}
