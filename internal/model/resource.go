package model

import "time"

// Resource is the canonical normalized representation of cloud infrastructure.
type Resource struct {
	ID         string            `json:"id"`
	Provider   string            `json:"provider"`
	Type       string            `json:"type"`
	CloudID    string            `json:"cloud_id"`
	Name       string            `json:"name"`
	Attributes map[string]any    `json:"attributes"`
	Tags       map[string]string `json:"tags"`
	Region     string            `json:"region"`
	Source     string            `json:"source"`
	Module     string            `json:"module,omitempty"`
}

// DriftKind categorizes a drift finding.
type DriftKind string

const (
	DriftMissingInCloud  DriftKind = "missing_in_cloud"
	DriftExtraInCloud    DriftKind = "extra_in_cloud"
	DriftAttributeChange DriftKind = "attribute_changed"
	DriftTagsChanged     DriftKind = "tags_changed"
)

// Severity levels for drift findings.
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// DriftFinding represents a single detected drift item.
type DriftFinding struct {
	Kind         DriftKind `json:"kind"`
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	ResourceName string    `json:"resource_name,omitempty"`
	Field        string    `json:"field,omitempty"`
	Expected     any       `json:"expected,omitempty"`
	Actual       any       `json:"actual,omitempty"`
	Severity     string    `json:"severity"`
	RiskScore    int       `json:"risk_score"`
	Remediation  string    `json:"remediation,omitempty"`
}

// DriftSummary aggregates finding counts.
type DriftSummary struct {
	TotalResources   int `json:"total_resources"`
	MissingInCloud   int `json:"missing_in_cloud"`
	ExtraInCloud     int `json:"extra_in_cloud"`
	AttributeChanges int `json:"attribute_changes"`
	TagChanges       int `json:"tag_changes"`
	TotalFindings    int `json:"total_findings"`
	TotalRiskScore   int `json:"total_risk_score"`
}

// DriftReport is the output of a drift scan.
type DriftReport struct {
	ScanID      string         `json:"scan_id"`
	WorkspaceID string         `json:"workspace_id,omitempty"`
	Workspace   string         `json:"workspace,omitempty"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt time.Time      `json:"completed_at"`
	Status      string         `json:"status"`
	Summary     DriftSummary   `json:"summary"`
	Findings    []DriftFinding `json:"findings"`
	Errors      []string       `json:"errors,omitempty"`
}

// ScanStatus values.
const (
	ScanStatusRunning   = "running"
	ScanStatusCompleted = "completed"
	ScanStatusFailed    = "failed"
)

// AuthConfig holds provider authentication details (e.g. OIDC, assume role).
type AuthConfig struct {
	RoleARN              string `json:"role_arn,omitempty" yaml:"role_arn,omitempty"`
	WebIdentityTokenFile string `json:"web_identity_token_file,omitempty" yaml:"web_identity_token_file,omitempty"`
	Profile              string `json:"profile,omitempty" yaml:"profile,omitempty"`
}

// Workspace represents a configured drift detection target.
type Workspace struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Provider string          `json:"provider"`
	Auth     *AuthConfig     `json:"auth,omitempty" yaml:"auth,omitempty"`
	State    StateConfig     `json:"state"`
	Regions  []string        `json:"regions"`
	Compare  CompareConfig   `json:"compare"`
	Schedule *ScheduleConfig `json:"schedule,omitempty"`
}

// StateConfig describes where Terraform state is stored.
type StateConfig struct {
	Backend     string            `json:"backend" yaml:"backend"`
	Path        string            `json:"path,omitempty" yaml:"path,omitempty"`
	Bucket      string            `json:"bucket,omitempty" yaml:"bucket,omitempty"`
	Key         string            `json:"key,omitempty" yaml:"key,omitempty"`
	Region      string            `json:"region,omitempty" yaml:"region,omitempty"`
	WorkspaceID string            `json:"workspace_id,omitempty" yaml:"workspace_id,omitempty"` // For TFC
	Token       string            `json:"token,omitempty" yaml:"token,omitempty"`               // For TFC
	Extra       map[string]string `json:"extra,omitempty" yaml:"extra,omitempty"`
}

// CompareConfig controls drift comparison behavior.
type CompareConfig struct {
	IgnoreTags       []string `json:"ignore_tags" yaml:"ignore_tags"`
	IgnoreAttributes []string `json:"ignore_attributes" yaml:"ignore_attributes"`
}

// ScheduleConfig defines a cron schedule for a workspace.
type ScheduleConfig struct {
	Cron string `json:"cron" yaml:"cron"`
}

// ResourceSelector targets specific resources for cloud fetch.
type ResourceSelector struct {
	Type    string
	Region  string
	CloudID string
}
