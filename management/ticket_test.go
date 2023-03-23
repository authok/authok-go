package management

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/authok/authok-go"
)

func TestTicketManager_VerifyEmail(t *testing.T) {
	configureHTTPTestRecordings(t)

	user := givenAUser(t)
	ticket := &Ticket{
		ResultURL: authok.String("https://example.com/verify-email"),
		UserID:    user.ID,
		TTLSec:    authok.Int(3600),
	}

	err := api.Ticket.VerifyEmail(ticket)
	assert.NoError(t, err)
}

func TestTicketManager_ChangePassword(t *testing.T) {
	configureHTTPTestRecordings(t)

	user := givenAUser(t)
	ticket := &Ticket{
		ResultURL:              authok.String("https://example.com/change-password"),
		UserID:                 user.ID,
		TTLSec:                 authok.Int(3600),
		MarkEmailAsVerified:    authok.Bool(true),
		IncludeEmailInRedirect: authok.Bool(true),
	}

	err := api.Ticket.ChangePassword(ticket)
	assert.NoError(t, err)
}
