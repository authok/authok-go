package management

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authok/authok-go"
)

var logStreamTestCases = []logStreamTestCase{
	{
		name: "AmazonEventBridge LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeAmazonEventBridge),
			Sink: &LogStreamSinkAmazonEventBridge{
				AccountID: authok.String("999999999999"),
				Region:    authok.String("us-west-2"),
			},
		},
	},
	// { // This test requires an active subscription.
	// 	name: "AzureEventGrid LogStream",
	// 	logStream: LogStream{
	// 		Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
	// 		Type: authok.String(LogStreamTypeAzureEventGrid),
	// 		Sink: &LogStreamSinkAzureEventGrid{
	// 			SubscriptionID: authok.String("b69a6835-57c7-4d53-b0d5-1c6ae580b6d5"),
	// 			Region:         authok.String("northeurope"),
	// 			ResourceGroup:  authok.String("azure-logs-rg"),
	// 		},
	// 	},
	// },
	{
		name: "HTTP LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeHTTP),
			Sink: &LogStreamSinkHTTP{
				Endpoint:      authok.String("https://example.com/logs"),
				Authorization: authok.String("Bearer f2368bbe77074527a37be2fdd5b92bad"),
				ContentFormat: authok.String("JSONLINES"),
				ContentType:   authok.String("application/json"),
			},
		},
	},
	{
		name: "DataDog LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeDatadog),
			Sink: &LogStreamSinkDatadog{
				APIKey: authok.String("121233123455"),
				Region: authok.String("us"),
			},
		},
	},
	{
		name: "Segment LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeSegment),
			Sink: &LogStreamSinkSegment{
				WriteKey: authok.String("121233123455"),
			},
		},
	},
	{
		name: "Splunk LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeSplunk),
			Sink: &LogStreamSinkSplunk{
				Domain: authok.String("demo.splunk.com"),
				Port:   authok.String("8080"),
				Secure: authok.Bool(true),
				Token:  authok.String("12a34ab5-c6d7-8901-23ef-456b7c89d0c1"),
			},
		},
	},
	{
		name: "Sumo LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeSumo),
			Sink: &LogStreamSinkSumo{
				SourceAddress: authok.String("https://example.com"),
			},
		},
	},
	{
		name: "Mixpanel LogStream",
		logStream: LogStream{
			Name: authok.Stringf("Test-LogStream-%d", time.Now().Unix()),
			Type: authok.String(LogStreamTypeMixpanel),
			Sink: &LogStreamSinkMixpanel{
				Region:                 authok.String("us"),
				ProjectID:              authok.String("123456789"),
				ServiceAccountUsername: authok.String("fake-account.123abc.mp-service-account"),
				ServiceAccountPassword: authok.String("8iwyKSzwV2brfakepassGGKhsZ3INozo"),
			},
		},
	},
}

type logStreamTestCase struct {
	name      string
	logStream LogStream
}

func TestLogStreamManager_Create(t *testing.T) {
	for _, testCase := range logStreamTestCases {
		t.Run("It can successfully create a "+testCase.name, func(t *testing.T) {
			configureHTTPTestRecordings(t)

			expectedLogStream := testCase.logStream

			err := api.LogStream.Create(&expectedLogStream)
			assert.NoError(t, err)
			assert.NotEmpty(t, expectedLogStream.GetID())

			t.Cleanup(func() {
				cleanupLogStream(t, expectedLogStream.GetID())
			})
		})
	}
}

func TestLogStreamManager_Read(t *testing.T) {
	for _, testCase := range logStreamTestCases {
		t.Run("It can successfully read a "+testCase.name, func(t *testing.T) {
			configureHTTPTestRecordings(t)

			expectedLogStream := givenALogStream(t, testCase)

			actualLogStream, err := api.LogStream.Read(expectedLogStream.GetID())

			assert.NoError(t, err)
			assert.Equal(t, expectedLogStream, actualLogStream)
		})
	}
}

func TestLogStreamManager_Update(t *testing.T) {
	for _, testCase := range logStreamTestCases {
		t.Run("It can successfully update a "+testCase.name, func(t *testing.T) {
			configureHTTPTestRecordings(t)

			logStream := givenALogStream(t, testCase)
			updatedLogStream := &LogStream{
				Filters: &[]map[string]string{
					{
						"type": "category",
						"name": "auth.login.fail",
					},
				},
			}

			err := api.LogStream.Update(logStream.GetID(), updatedLogStream)
			assert.NoError(t, err)

			actualLogStream, err := api.LogStream.Read(logStream.GetID())
			assert.NoError(t, err)
			assert.Equal(t, updatedLogStream.Filters, actualLogStream.Filters)
		})
	}
}

func TestLogStreamManager_Delete(t *testing.T) {
	for _, testCase := range logStreamTestCases {
		t.Run("It can successfully delete a "+testCase.name, func(t *testing.T) {
			configureHTTPTestRecordings(t)

			logStream := givenALogStream(t, testCase)

			err := api.LogStream.Delete(logStream.GetID())
			assert.NoError(t, err)

			actualLogStream, err := api.LogStream.Read(logStream.GetID())
			assert.Nil(t, actualLogStream)
			assert.Error(t, err)
			assert.Implements(t, (*Error)(nil), err)
			assert.Equal(t, http.StatusNotFound, err.(Error).Status())
		})
	}
}

func TestLogStreamManager_List(t *testing.T) {
	configureHTTPTestRecordings(t)

	// There are no params we can add here, unfortunately.
	logStreamList, err := api.LogStream.List()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(logStreamList), 0)
}

func givenALogStream(t *testing.T, testCase logStreamTestCase) *LogStream {
	t.Helper()

	logStream := testCase.logStream

	err := api.LogStream.Create(&logStream)
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanupLogStream(t, logStream.GetID())
	})

	return &logStream
}

func cleanupLogStream(t *testing.T, logStreamID string) {
	t.Helper()

	err := api.LogStream.Delete(logStreamID)
	require.NoError(t, err)
}
