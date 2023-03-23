package management

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"

	"github.com/authok/authok-go/internal/client"
)

var (
	domain                = os.Getenv("AUTHOK_DOMAIN")
	clientID              = os.Getenv("AUTHOK_CLIENT_ID")
	clientSecret          = os.Getenv("AUTHOK_CLIENT_SECRET")
	debug                 = os.Getenv("AUTHOK_DEBUG")
	httpRecordings        = os.Getenv("AUTHOK_HTTP_RECORDINGS")
	httpRecordingsEnabled = false
	api                   = &Management{}
)

func envVarEnabled(envVar string) bool {
	return envVar == "true" || envVar == "1" || envVar == "on"
}

func TestMain(m *testing.M) {
	httpRecordingsEnabled = envVarEnabled(httpRecordings)
	initializeTestClient()

	code := m.Run()
	os.Exit(code)
}

func initializeTestClient() {
	var err error

	api, err = New(
		domain,
		WithClientCredentials(clientID, clientSecret),
		WithDebug(envVarEnabled(debug)),
	)
	if err != nil {
		log.Fatal("failed to initialize the api client")
	}
}

func TestNew(t *testing.T) {
	for _, domain := range []string{
		"example.com ",
		" example.com",
		" example.com ",
		"%2Fexample.com",
		" a.b.c.example.com",
	} {
		_, err := New(domain)
		assert.Errorf(t, err, "expected New to fail with domain %q", domain)
	}
}

func TestOptionFields(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	IncludeFields("foo", "bar").apply(r)

	v := r.URL.Query()

	fields := v.Get("fields")
	assert.Equal(t, "foo,bar", fields)

	includeFields := v.Get("include_fields")
	assert.Equal(t, "true", includeFields)

	ExcludeFields("foo", "bar").apply(r)

	includeFields = v.Get("include_fields")
	assert.Equal(t, "true", includeFields)
}

func TestOptionPage(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	Page(3).apply(r)
	PerPage(10).apply(r)

	v := r.URL.Query()

	page := v.Get("page")
	assert.Equal(t, "3", page)

	perPage := v.Get("page_size")
	assert.Equal(t, "10", perPage)
}

func TestOptionTotals(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	IncludeTotals(true).apply(r)

	v := r.URL.Query()

	includeTotals := v.Get("include_totals")
	assert.Equal(t, "true", includeTotals)
}

func TestOptionFrom(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	From("abc").apply(r)
	Take(10).apply(r)

	v := r.URL.Query()

	from := v.Get("from")
	assert.Equal(t, "abc", from)

	take := v.Get("take")
	assert.Equal(t, "10", take)
}

func TestOptionParameter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	Parameter("foo", "123").apply(r)
	Parameter("bar", "xyz").apply(r)

	v := r.URL.Query()

	foo := v.Get("foo")
	assert.Equal(t, "123", foo)

	bar := v.Get("bar")
	assert.Equal(t, "xyz", bar)
}

func TestOptionDefaults(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	applyListDefaults([]RequestOption{
		PerPage(20),          // This should be persisted (default is 50).
		IncludeTotals(false), // This should be persisted (default is true).
	}).apply(r)

	v := r.URL.Query()

	perPage := v.Get("page_size")
	assert.Equal(t, "20", perPage)

	includeTotals := v.Get("include_totals")
	assert.Equal(t, "false", includeTotals)
}

func TestStringify(t *testing.T) {
	expected := `{
  "foo": "bar"
}`

	v := struct {
		Foo string `json:"foo"`
	}{
		"bar",
	}

	s := Stringify(v)
	assert.Equal(t, expected, s)
}

func TestRequestOptionContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the request.

	err := api.Request("GET", "/", nil, Context(ctx))
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected err to be context.Canceled, got %v", err)
	}
}

func TestRequestOptionContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(50 * time.Millisecond) // Delay until the deadline is exceeded.

	err := api.Request("GET", "/", nil, Context(ctx))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected err to be context.DeadlineExceeded, got %v", err)
	}
}

func TestNew_WithInsecure(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/users/123":
			w.Write([]byte(`{"user_id":"123"}`))
		default:
			http.NotFound(w, r)
		}
	})
	s := httptest.NewServer(h)

	m, err := New(s.URL, WithInsecure())
	assert.NoError(t, err)

	u, err := m.User.Read("123")
	assert.NoError(t, err)
	assert.Equal(t, "123", u.GetID())
}

func TestManagement_URI(t *testing.T) {
	var testCases = []struct {
		name     string
		given    []string
		expected string
	}{
		{
			name:     "encodes regular user_id",
			given:    []string{"users", "1234"},
			expected: "https://" + domain + "/api/v1/users/1234",
		},
		{
			name:     "encodes a user_id with a space",
			given:    []string{"users", "123 4"},
			expected: "https://" + domain + "/api/v1/users/123%204",
		},
		{
			name:     "encodes a user_id with a |",
			given:    []string{"users", "authok|12345678"},
			expected: "https://" + domain + "/api/v1/users/authok%7C12345678",
		},
		{
			name:     "encodes a user_id with a | and /",
			given:    []string{"users", "authok|1234/5678"},
			expected: "https://" + domain + "/api/v1/users/authok%7C1234%2F5678",
		},
		{
			name:     "encodes a user_id with a /",
			given:    []string{"users", "anotherUserId/secret"},
			expected: "https://" + domain + "/api/v1/users/anotherUserId%2Fsecret",
		},
		{
			name:     "encodes a user_id with a percentage",
			given:    []string{"users", "anotherUserId/secret%23"},
			expected: "https://" + domain + "/api/v1/users/anotherUserId%2Fsecret%2523",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := api.URI(testCase.given...)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestAuthokClient(t *testing.T) {
	t.Run("Defaults to the default data", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authok-Client")
			authokClientDecoded, err := base64.StdEncoding.DecodeString(header)
			assert.NoError(t, err)

			var authokClient client.AuthokClientInfo
			err = json.Unmarshal(authokClientDecoded, &authokClient)

			assert.NoError(t, err)
			assert.Equal(t, "authok-go", authokClient.Name)
			assert.Equal(t, "latest", authokClient.Version)
			assert.Equal(t, runtime.Version(), authokClient.Env["go"])
		})
		s := httptest.NewServer(h)

		m, err := New(
			s.URL,
			WithInsecure(),
		)
		assert.NoError(t, err)

		_, err = m.User.Read("123")

		assert.NoError(t, err)
	})

	t.Run("Allows passing custom AuthokClientInfo", func(t *testing.T) {
		customClient := client.AuthokClientInfo{Name: "test-client", Version: "1.0.0"}

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authok-Client")
			assert.Equal(t, "eyJuYW1lIjoidGVzdC1jbGllbnQiLCJ2ZXJzaW9uIjoiMS4wLjAifQ==", header)
		})
		s := httptest.NewServer(h)

		m, err := New(
			s.URL,
			WithInsecure(),
			WithAuthokClientInfo(customClient),
		)
		assert.NoError(t, err)

		_, err = m.User.Read("123")

		assert.NoError(t, err)
	})

	t.Run("Allows disabling AuthokClientInfo", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawHeader := r.Header.Get("Authok-Client")
			assert.Empty(t, rawHeader)
		})
		s := httptest.NewServer(h)

		m, err := New(
			s.URL,
			WithInsecure(),
			WithNoAuthokClientInfo(),
		)
		assert.NoError(t, err)
		_, err = m.User.Read("123")
		assert.NoError(t, err)
	})
}
