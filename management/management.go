package management

//go:generate go run gen-methods.go

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	"github.com/authok/authok-go/internal/client"
)

// Management is an Authok management client used to interact with the Authok
// Management API v2.
type Management struct {
	// Client manages Authok Client (also known as Application) resources.
	Client *ClientManager

	// ClientGrant manages Authok ClientGrant resources.
	ClientGrant *ClientGrantManager

	// ResourceServer manages Authok Resource Server (also known as API)
	// resources.
	ResourceServer *ResourceServerManager

	// Connection manages Authok Connection resources.
	Connection *ConnectionManager

	// CustomDomain manages Authok Custom Domains.
	CustomDomain *CustomDomainManager

	// Grant manages Authok Grants.
	Grant *GrantManager

	// Log reads Authok Logs.
	Log *LogManager

	// LogStream reads Authok Logs.
	LogStream *LogStreamManager

	// RoleManager manages Authok Roles.
	Role *RoleManager

	// RuleManager manages Authok Rules.
	Rule *RuleManager

	// HookManager manages Authok Hooks
	Hook *HookManager

	// RuleManager manages Authok Rule Configurations.
	RuleConfig *RuleConfigManager

	// Email manages Authok Email Providers.
	// Deprecated: Use EmailProvider instead.
	Email *EmailManager

	// EmailTemplate manages Authok Email Templates.
	EmailTemplate *EmailTemplateManager

	// User manages Authok User resources.
	User *UserManager

	// Job manages Authok jobs.
	Job *JobManager

	// Tenant manages your Authok Tenant.
	Tenant *TenantManager

	// Ticket creates verify email or change password tickets.
	Ticket *TicketManager

	// Stat is used to retrieve usage statistics.
	Stat *StatManager

	// Branding settings such as company logo or primary color.
	Branding *BrandingManager

	// Guardian manages your Authok Guardian settings
	Guardian *GuardianManager

	// Prompt manages your prompt settings.
	Prompt *PromptManager

	// Blacklist manages the authok blacklists
	Blacklist *BlacklistManager

	// SigningKey manages Authok Application Signing Keys.
	SigningKey *SigningKeyManager

	// Anomaly manages the IP blocks
	Anomaly *AnomalyManager

	// Actions manages Actions extensibility
	Action *ActionManager

	// Organization manages Authok Organizations.
	Organization *OrganizationManager

	// AttackProtection manages Authok Attack Protection.
	AttackProtection *AttackProtectionManager

	// BrandingTheme manages Authok Branding Themes.
	BrandingTheme *BrandingThemeManager

	// EmailProvider manages Authok Email Providers.
	EmailProvider *EmailProviderManager

	url              *url.URL
	basePath         string
	userAgent        string
	debug            bool
	ctx              context.Context
	tokenSource      oauth2.TokenSource
	http             *http.Client
	authokClientInfo *client.AuthokClientInfo
}

// New creates a new Authok Management client by authenticating using the
// supplied client id and secret.
func New(domain string, options ...Option) (*Management, error) {
	// Ignore the scheme if it was defined in the domain variable, then prefix
	// with https as it's the only scheme supported by the Authok API.
	if i := strings.Index(domain, "//"); i != -1 {
		domain = domain[i+2:]
	}
	domain = "https://" + domain

	u, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}

	m := &Management{
		url:              u,
		basePath:         "api/v1",
		userAgent:        client.UserAgent,
		debug:            false,
		ctx:              context.Background(),
		http:             http.DefaultClient,
		authokClientInfo: client.DefaultAuthokClientInfo,
	}

	for _, option := range options {
		option(m)
	}

	m.http = client.Wrap(
		m.http,
		m.tokenSource,
		client.WithDebug(m.debug),
		client.WithUserAgent(m.userAgent),
		client.WithRateLimit(),
		client.WithAuthokClientInfo(m.authokClientInfo),
	)

	m.Client = newClientManager(m)
	m.ClientGrant = newClientGrantManager(m)
	m.Connection = newConnectionManager(m)
	m.CustomDomain = newCustomDomainManager(m)
	m.Grant = newGrantManager(m)
	m.LogStream = newLogStreamManager(m)
	m.Log = newLogManager(m)
	m.ResourceServer = newResourceServerManager(m)
	m.Role = newRoleManager(m)
	m.Rule = newRuleManager(m)
	m.Hook = newHookManager(m)
	m.RuleConfig = newRuleConfigManager(m)
	m.EmailTemplate = newEmailTemplateManager(m)
	m.Email = newEmailManager(m)
	m.User = newUserManager(m)
	m.Job = newJobManager(m)
	m.Tenant = newTenantManager(m)
	m.Ticket = newTicketManager(m)
	m.Stat = newStatManager(m)
	m.Branding = newBrandingManager(m)
	m.Guardian = newGuardianManager(m)
	m.Prompt = newPromptManager(m)
	m.Blacklist = newBlacklistManager(m)
	m.SigningKey = newSigningKeyManager(m)
	m.Anomaly = newAnomalyManager(m)
	m.Action = newActionManager(m)
	m.Organization = newOrganizationManager(m)
	m.AttackProtection = newAttackProtectionManager(m)
	m.BrandingTheme = newBrandingThemeManager(m)
	m.EmailProvider = newEmailProviderManager(m)

	return m, nil
}
