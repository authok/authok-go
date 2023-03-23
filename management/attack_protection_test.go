package management

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/authok/authok-go"
)

func TestAttackProtection(t *testing.T) {
	t.Run("Get breached password detection settings", func(t *testing.T) {
		configureHTTPTestRecordings(t)

		breachedPasswordDetection, err := api.AttackProtection.GetBreachedPasswordDetection()
		assert.NoError(t, err)
		assert.IsType(t, &BreachedPasswordDetection{}, breachedPasswordDetection)
	})

	t.Run("Update breached password detection settings", func(t *testing.T) {
		configureHTTPTestRecordings(t)

		// Save initial settings.
		preTestBPDSettings, err := api.AttackProtection.GetBreachedPasswordDetection()
		assert.NoError(t, err)

		expected := &BreachedPasswordDetection{
			Enabled: authok.Bool(true),
			Method:  authok.String("standard"),
			Stage: &BreachedPasswordDetectionStage{
				PreUserRegistration: &BreachedPasswordDetectionPreUserRegistration{
					Shields: &[]string{"block"},
				},
			},
		}

		err = api.AttackProtection.UpdateBreachedPasswordDetection(expected)
		assert.NoError(t, err)

		actual, err := api.AttackProtection.GetBreachedPasswordDetection()
		assert.NoError(t, err)
		assert.Equal(t, expected.GetEnabled(), actual.GetEnabled())
		assert.Equal(t, expected.GetMethod(), actual.GetMethod())
		assert.Equal(t, expected.GetStage().GetPreUserRegistration().GetShields(), actual.GetStage().GetPreUserRegistration().GetShields())

		// Restore initial settings.
		err = api.AttackProtection.UpdateBreachedPasswordDetection(preTestBPDSettings)
		assert.NoError(t, err)
	})

	t.Run("Get the brute force configuration", func(t *testing.T) {
		configureHTTPTestRecordings(t)

		bruteForceProtection, err := api.AttackProtection.GetBruteForceProtection()
		assert.NoError(t, err)
		assert.IsType(t, &BruteForceProtection{}, bruteForceProtection)
	})

	t.Run("Update the brute force configuration", func(t *testing.T) {
		configureHTTPTestRecordings(t)

		// Save initial settings.
		preTestBFPSettings, err := api.AttackProtection.GetBruteForceProtection()
		assert.NoError(t, err)

		expected := &BruteForceProtection{
			Enabled:     authok.Bool(true),
			MaxAttempts: authok.Int(10),
		}

		err = api.AttackProtection.UpdateBruteForceProtection(expected)
		assert.NoError(t, err)

		actual, err := api.AttackProtection.GetBruteForceProtection()
		assert.NoError(t, err)
		assert.Equal(t, expected.GetEnabled(), actual.GetEnabled())
		assert.Equal(t, expected.GetMaxAttempts(), actual.GetMaxAttempts())

		// Restore initial settings.
		err = api.AttackProtection.UpdateBruteForceProtection(preTestBFPSettings)
		assert.NoError(t, err)
	})

	t.Run("Get the suspicious IP throttling configuration", func(t *testing.T) {
		configureHTTPTestRecordings(t)

		suspiciousIPThrottling, err := api.AttackProtection.GetSuspiciousIPThrottling()
		assert.NoError(t, err)
		assert.IsType(t, &SuspiciousIPThrottling{}, suspiciousIPThrottling)
	})

	t.Run("Update the suspicious IP throttling configuration", func(t *testing.T) {
		configureHTTPTestRecordings(t)

		// Save initial settings.
		preTestSIPSettings, err := api.AttackProtection.GetSuspiciousIPThrottling()
		assert.NoError(t, err)

		expected := &SuspiciousIPThrottling{
			Enabled: authok.Bool(true),
			Stage: &Stage{
				PreLogin: &PreLogin{
					MaxAttempts: authok.Int(100),
					Rate:        authok.Int(864000),
				},
				PreUserRegistration: &PreUserRegistration{
					MaxAttempts: authok.Int(50),
					Rate:        authok.Int(1200),
				},
			},
		}

		err = api.AttackProtection.UpdateSuspiciousIPThrottling(expected)
		assert.NoError(t, err)

		actual, err := api.AttackProtection.GetSuspiciousIPThrottling()
		assert.NoError(t, err)
		assert.Equal(t, expected.GetEnabled(), actual.GetEnabled())
		assert.Equal(t, expected.GetStage().GetPreLogin().GetRate(), actual.GetStage().GetPreLogin().GetRate())
		assert.Equal(t, expected.GetStage().GetPreLogin().GetMaxAttempts(), actual.GetStage().GetPreLogin().GetMaxAttempts())
		assert.Equal(t, expected.GetStage().GetPreUserRegistration().GetRate(), actual.GetStage().GetPreUserRegistration().GetRate())
		assert.Equal(t, expected.GetStage().GetPreUserRegistration().GetMaxAttempts(), actual.GetStage().GetPreUserRegistration().GetMaxAttempts())

		// Restore initial settings.
		err = api.AttackProtection.UpdateSuspiciousIPThrottling(preTestSIPSettings)
		assert.NoError(t, err)
	})
}
