package management

// Rule is used as part of the authentication pipeline.
type Rule struct {
	// The rule's identifier.
	ID *string `json:"id,omitempty"`

	// The name of the rule. Can only contain alphanumeric characters, spaces
	// and '-'. Can neither start nor end with '-' or spaces.
	Name *string `json:"name,omitempty"`

	// A script that contains the rule's code.
	Script *string `json:"script,omitempty"`

	// The rule's order in relation to other rules. A rule with a lower order
	// than another rule executes first. If no order is provided it will
	// automatically be one greater than the current maximum.
	Order *int `json:"order,omitempty"`

	// Enabled should be set to true if the rule is enabled, false otherwise.
	Enabled *bool `json:"enabled,omitempty"`
}

// RuleList holds a list of Rules.
type RuleList struct {
	List
	Rules []*Rule `json:"rules"`
}

// RuleManager manages Authok Rule resources.
type RuleManager struct {
	*Management
}

func newRuleManager(m *Management) *RuleManager {
	return &RuleManager{m}
}

// Create a new rule.
//
// Note: Changing a rule's stage of execution from the default `login_success`
// can change the rule's function signature to have user omitted.
//
// See: https://authok.com/docs/api/management/v1#!/Rules/post_rules
func (m *RuleManager) Create(r *Rule, opts ...RequestOption) error {
	return m.Request("POST", m.URI("rules"), r, opts...)
}

// Retrieve rule details. Accepts a list of fields to include or exclude in the result.
//
// See: https://authok.com/docs/api/management/v1#!/Rules/get_rules_by_id
func (m *RuleManager) Read(id string, opts ...RequestOption) (r *Rule, err error) {
	err = m.Request("GET", m.URI("rules", id), &r, opts...)
	return
}

// Update an existing rule.
//
// See: https://authok.com/docs/api/management/v1#!/Rules/patch_rules_by_id
func (m *RuleManager) Update(id string, r *Rule, opts ...RequestOption) error {
	return m.Request("PATCH", m.URI("rules", id), r, opts...)
}

// Delete a rule.
//
// See: https://authok.com/docs/api/management/v1#!/Rules/delete_rules_by_id
func (m *RuleManager) Delete(id string, opts ...RequestOption) error {
	return m.Request("DELETE", m.URI("rules", id), nil, opts...)
}

// List all rules.
//
// See: https://authok.com/docs/api/management/v1#!/Rules/get_rules
func (m *RuleManager) List(opts ...RequestOption) (r *RuleList, err error) {
	err = m.Request("GET", m.URI("rules"), &r, applyListDefaults(opts))
	return
}
