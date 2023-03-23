package management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"

	"github.com/authok/authok-go"
)

const (
	recordingsDIR    = "./../test/data/recordings/"
	recordingsDomain = "authok-go-dev.eu.authok.com"
)

func configureHTTPTestRecordings(t *testing.T) {
	t.Helper()

	if !httpRecordingsEnabled {
		return
	}

	initialTransport := api.http.Transport

	recorderTransport, err := recorder.NewWithOptions(
		&recorder.Options{
			CassetteName:       recordingsDIR + t.Name(),
			Mode:               recorder.ModeRecordOnce,
			RealTransport:      api.http.Transport,
			SkipRequestLatency: true,
		},
	)
	require.NoError(t, err)

	removeSensitiveDataFromRecordings(t, recorderTransport)

	api.http.Transport = recorderTransport

	t.Cleanup(func() {
		err := recorderTransport.Stop()
		require.NoError(t, err)
		api.http.Transport = initialTransport
	})
}

func removeSensitiveDataFromRecordings(t *testing.T, recorderTransport *recorder.Recorder) {
	recorderTransport.AddHook(
		func(i *cassette.Interaction) error {
			skip429Response(i)
			redactHeaders(i)
			redactSensitiveDataInSigningKey(t, i)
			redactSensitiveDataInClient(t, i)
			redactSensitiveDataInResourceServer(t, i)
			redactSensitiveDataInLogSession(t, i)
			redactSensitiveDataInLogs(t, i)

			// Redact domain should always be ran last
			redactDomain(i, domain)

			return nil
		},
		recorder.BeforeSaveHook,
	)
}

func skip429Response(i *cassette.Interaction) {
	if i.Response.Code == http.StatusTooManyRequests {
		i.DiscardOnSave = true
	}
}

func redactHeaders(i *cassette.Interaction) {
	allowedHeaders := map[string]bool{
		"Content-Type": true,
		"User-Agent":   true,
	}

	for header := range i.Request.Headers {
		if _, ok := allowedHeaders[header]; !ok {
			delete(i.Request.Headers, header)
		}
	}
	for header := range i.Response.Headers {
		if _, ok := allowedHeaders[header]; !ok {
			delete(i.Response.Headers, header)
		}
	}
}

func redactDomain(i *cassette.Interaction, domain string) {
	i.Request.Host = strings.ReplaceAll(i.Request.Host, domain, recordingsDomain)
	i.Request.URL = strings.ReplaceAll(i.Request.URL, domain, recordingsDomain)

	domainParts := strings.Split(domain, ".")

	i.Response.Body = strings.ReplaceAll(i.Response.Body, domainParts[0], recordingsDomain)
	i.Request.Body = strings.ReplaceAll(i.Request.Body, domainParts[0], recordingsDomain)
}

func redactSensitiveDataInSigningKey(t *testing.T, i *cassette.Interaction) {
	signingKey := &SigningKey{
		KID:         authok.String("111111111111111111111"),
		Cert:        authok.String("-----BEGIN CERTIFICATE-----\\r\\n[REDACTED]\\r\\n-----END CERTIFICATE-----"),
		PKCS7:       authok.String("-----BEGIN PKCS7-----\\r\\n[REDACTED]\\r\\n-----END PKCS7-----"),
		Current:     authok.Bool(true),
		Next:        authok.Bool(false),
		Previous:    authok.Bool(true),
		Fingerprint: authok.String("[REDACTED]"),
		Thumbprint:  authok.String("[REDACTED]"),
		Revoked:     authok.Bool(false),
	}
	previousSigningKey := &SigningKey{
		KID:         authok.String("222222222222222222222"),
		Cert:        authok.String("-----BEGIN CERTIFICATE-----\\r\\n[REDACTED]\\r\\n-----END CERTIFICATE-----"),
		PKCS7:       authok.String("-----BEGIN PKCS7-----\\r\\n[REDACTED]\\r\\n-----END PKCS7-----"),
		Current:     authok.Bool(false),
		Next:        authok.Bool(true),
		Previous:    authok.Bool(true),
		Fingerprint: authok.String("[REDACTED]"),
		Thumbprint:  authok.String("[REDACTED]"),
		Revoked:     authok.Bool(false),
	}

	if i.Request.URL == "https://"+domain+"/api/v1/keys/signing" && i.Request.Method == http.MethodGet {
		signingKeyBody, err := json.Marshal(signingKey)
		require.NoError(t, err)
		previousSigningKeyBody, err := json.Marshal(previousSigningKey)
		require.NoError(t, err)

		i.Response.Body = fmt.Sprintf(`[%s,%s]`, signingKeyBody, previousSigningKeyBody)
		return
	}

	isSigningKeysURL := strings.Contains(i.Request.URL, "https://"+domain+"/api/v1/keys/signing")
	if isSigningKeysURL && i.Request.Method == http.MethodGet {
		i.Request.URL = "https://" + domain + "/api/v1/keys/signing/111111111111111111111"

		signingKeyBody, err := json.Marshal(signingKey)
		require.NoError(t, err)

		i.Response.Body = string(signingKeyBody)
		return
	}

	if isSigningKeysURL && strings.Contains(i.Request.URL, "/revoke") && i.Request.Method == http.MethodPut {
		i.Request.URL = "https://" + domain + "/api/v1/keys/signing/111111111111111111111/revoke"

		signingKey.RevokedAt = authok.Time(time.Now())
		signingKeyBody, err := json.Marshal(signingKey)
		require.NoError(t, err)

		i.Response.Body = string(signingKeyBody)
		return
	}

	if isSigningKeysURL && strings.Contains(i.Request.URL, "/rotate") && i.Request.Method == http.MethodPost {
		signingKeyBody, err := json.Marshal(signingKey)
		require.NoError(t, err)

		i.Response.Body = string(signingKeyBody)
		return
	}
}

func redactSensitiveDataInClient(t *testing.T, i *cassette.Interaction) {
	isClientURL := strings.Contains(i.Request.URL, "https://"+domain+"/api/v1/clients")
	create := isClientURL && i.Request.Method == http.MethodPost
	read := isClientURL && i.Request.Method == http.MethodGet
	update := isClientURL && i.Request.Method == http.MethodPatch
	rotateSecret := isClientURL && strings.Contains(i.Request.URL, "/rotate-secret")
	list := isClientURL && strings.Contains(i.Request.URL, "include_totals")

	if (create || read || update || rotateSecret) && !list {
		var client Client
		err := json.Unmarshal([]byte(i.Response.Body), &client)
		require.NoError(t, err)

		redacted := "[REDACTED]"
		if rotateSecret {
			redacted = "[ROTATED-REDACTED]"
		}

		client.SigningKeys = []map[string]string{
			{"cert": redacted},
		}
		client.ClientSecret = &redacted

		clientBody, err := json.Marshal(client)
		require.NoError(t, err)

		i.Response.Body = string(clientBody)
	}
}

func redactSensitiveDataInResourceServer(t *testing.T, i *cassette.Interaction) {
	isResourceServerURL := strings.Contains(i.Request.URL, "https://"+domain+"/api/v1/resource-servers")
	create := isResourceServerURL && i.Request.Method == http.MethodPost
	read := isResourceServerURL && i.Request.Method == http.MethodGet
	update := isResourceServerURL && i.Request.Method == http.MethodPatch
	list := isResourceServerURL && strings.Contains(i.Request.URL, "include_totals")

	if (create || read || update) && !list {
		var rs ResourceServer
		err := json.Unmarshal([]byte(i.Response.Body), &rs)
		require.NoError(t, err)

		rs.SigningSecret = nil

		rsBody, err := json.Marshal(rs)
		require.NoError(t, err)

		i.Response.Body = string(rsBody)
	}
}

func redactSensitiveDataInLogSession(t *testing.T, i *cassette.Interaction) {
	isLogSessionURL := strings.Contains(i.Request.URL, "https://"+domain+"/api/v1/actions/log-sessions")
	if isLogSessionURL {
		var logSession ActionLogSession
		err := json.Unmarshal([]byte(i.Response.Body), &logSession)
		require.NoError(t, err)

		replacedURL := "https://" + domain + "/api/v1/actions/log-sessions/tail?token=tkn_123"
		logSession.URL = &replacedURL

		logSessionBody, err := json.Marshal(logSession)
		require.NoError(t, err)

		i.Response.Body = string(logSessionBody)
	}
}

func redactSensitiveDataInLogs(t *testing.T, i *cassette.Interaction) {
	isMultipleLogsURL := strings.Contains(i.Request.URL, "https://"+domain+"/api/v1/logs?")
	if isMultipleLogsURL {
		var logs []Log
		err := json.Unmarshal([]byte(i.Response.Body), &logs)
		require.NoError(t, err)

		for i, log := range logs {
			log.IP = authok.String("[REDACTED]")
			logs[i] = log
		}

		logsBody, err := json.Marshal(logs)
		require.NoError(t, err)

		i.Response.Body = string(logsBody)
	}

	isSingleLogURL := strings.Contains(i.Request.URL, "https://"+domain+"/api/v1/logs/")
	if isSingleLogURL {
		var log Log
		err := json.Unmarshal([]byte(i.Response.Body), &log)
		require.NoError(t, err)

		log.IP = authok.String("[REDACTED]")

		logsBody, err := json.Marshal(log)
		require.NoError(t, err)

		i.Response.Body = string(logsBody)
	}
}