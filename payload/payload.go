package payload

import "time"

type PushPayload struct {
	Ref          string           `json:"ref"`
	Before       string           `json:"before"`
	After        string           `json:"after"`
	CompareURL   string           `json:"compare_url"`
	Commits      []*PayloadCommit `json:"commits"`
	TotalCommits int              `json:"total_commits"`
	HeadCommit   *PayloadCommit   `json:"head_commit"`
	Repo         *Repository      `json:"repository"`
	Pusher       *User            `json:"pusher"`
	Sender       *User            `json:"sender"`
}

// Repository represents a repository
type Repository struct {
	ID            int64       `json:"id"`
	Owner         *User       `json:"owner"`
	Name          string      `json:"name"`
	FullName      string      `json:"full_name"`
	Description   string      `json:"description"`
	Empty         bool        `json:"empty"`
	Private       bool        `json:"private"`
	Fork          bool        `json:"fork"`
	Template      bool        `json:"template"`
	Parent        *Repository `json:"parent"`
	Mirror        bool        `json:"mirror"`
	Size          int         `json:"size"`
	Language      string      `json:"language"`
	LanguagesURL  string      `json:"languages_url"`
	HTMLURL       string      `json:"html_url"`
	URL           string      `json:"url"`
	Link          string      `json:"link"`
	SSHURL        string      `json:"ssh_url"`
	CloneURL      string      `json:"clone_url"`
	OriginalURL   string      `json:"original_url"`
	Website       string      `json:"website"`
	Stars         int         `json:"stars_count"`
	Forks         int         `json:"forks_count"`
	Watchers      int         `json:"watchers_count"`
	OpenIssues    int         `json:"open_issues_count"`
	OpenPulls     int         `json:"open_pr_counter"`
	Releases      int         `json:"release_counter"`
	DefaultBranch string      `json:"default_branch"`
	Archived      bool        `json:"archived"`
	// swagger:strfmt date-time
	Created time.Time `json:"created_at"`
	// swagger:strfmt date-time
	Updated                       time.Time        `json:"updated_at"`
	ArchivedAt                    time.Time        `json:"archived_at"`
	Permissions                   *Permission      `json:"permissions,omitempty"`
	HasIssues                     bool             `json:"has_issues"`
	InternalTracker               *InternalTracker `json:"internal_tracker,omitempty"`
	ExternalTracker               *ExternalTracker `json:"external_tracker,omitempty"`
	HasWiki                       bool             `json:"has_wiki"`
	ExternalWiki                  *ExternalWiki    `json:"external_wiki,omitempty"`
	HasPullRequests               bool             `json:"has_pull_requests"`
	HasProjects                   bool             `json:"has_projects"`
	ProjectsMode                  string           `json:"projects_mode"`
	HasReleases                   bool             `json:"has_releases"`
	HasPackages                   bool             `json:"has_packages"`
	HasActions                    bool             `json:"has_actions"`
	IgnoreWhitespaceConflicts     bool             `json:"ignore_whitespace_conflicts"`
	AllowMerge                    bool             `json:"allow_merge_commits"`
	AllowRebase                   bool             `json:"allow_rebase"`
	AllowRebaseMerge              bool             `json:"allow_rebase_explicit"`
	AllowSquash                   bool             `json:"allow_squash_merge"`
	AllowFastForwardOnly          bool             `json:"allow_fast_forward_only_merge"`
	AllowRebaseUpdate             bool             `json:"allow_rebase_update"`
	DefaultDeleteBranchAfterMerge bool             `json:"default_delete_branch_after_merge"`
	DefaultMergeStyle             string           `json:"default_merge_style"`
	DefaultAllowMaintainerEdit    bool             `json:"default_allow_maintainer_edit"`
	AvatarURL                     string           `json:"avatar_url"`
	Internal                      bool             `json:"internal"`
	MirrorInterval                string           `json:"mirror_interval"`
	// ObjectFormatName of the underlying git repository
	// enum: sha1,sha256
	ObjectFormatName string `json:"object_format_name"`
	// swagger:strfmt date-time
	MirrorUpdated time.Time     `json:"mirror_updated,omitempty"`
	RepoTransfer  *RepoTransfer `json:"repo_transfer"`
}
type RepoTransfer struct {
	Doer      *User   `json:"doer"`
	Recipient *User   `json:"recipient"`
	Teams     []*Team `json:"teams"`
}
type Team struct {
	ID                      int64         `json:"id"`
	Name                    string        `json:"name"`
	Description             string        `json:"description"`
	Organization            *Organization `json:"organization"`
	IncludesAllRepositories bool          `json:"includes_all_repositories"`
	// enum: none,read,write,admin,owner
	Permission string `json:"permission"`
	// example: ["repo.code","repo.issues","repo.ext_issues","repo.wiki","repo.pulls","repo.releases","repo.projects","repo.ext_wiki"]
	// Deprecated: This variable should be replaced by UnitsMap and will be dropped in later versions.
	Units []string `json:"units"`
	// example: {"repo.code":"read","repo.issues":"write","repo.ext_issues":"none","repo.wiki":"admin","repo.pulls":"owner","repo.releases":"none","repo.projects":"none","repo.ext_wiki":"none"}
	UnitsMap         map[string]string `json:"units_map"`
	CanCreateOrgRepo bool              `json:"can_create_org_repo"`
}

type Organization struct {
	ID                        int64  `json:"id"`
	Name                      string `json:"name"`
	FullName                  string `json:"full_name"`
	Email                     string `json:"email"`
	AvatarURL                 string `json:"avatar_url"`
	Description               string `json:"description"`
	Website                   string `json:"website"`
	Location                  string `json:"location"`
	Visibility                string `json:"visibility"`
	RepoAdminChangeTeamAccess bool   `json:"repo_admin_change_team_access"`
	// deprecated
	UserName string `json:"username"`
}

// PayloadCommit represents a commit
type PayloadCommit struct {
	// sha1 hash of the commit
	ID           string                     `json:"id"`
	Message      string                     `json:"message"`
	URL          string                     `json:"url"`
	Author       *PayloadUser               `json:"author"`
	Committer    *PayloadUser               `json:"committer"`
	Verification *PayloadCommitVerification `json:"verification"`
	// swagger:strfmt date-time
	Timestamp time.Time `json:"timestamp"`
	Added     []string  `json:"added"`
	Removed   []string  `json:"removed"`
	Modified  []string  `json:"modified"`
}

// PayloadCommitVerification represents the GPG verification of a commit
type PayloadCommitVerification struct {
	Verified  bool         `json:"verified"`
	Reason    string       `json:"reason"`
	Signature string       `json:"signature"`
	Signer    *PayloadUser `json:"signer"`
	Payload   string       `json:"payload"`
}

// PayloadUser represents the author or committer of a commit
type PayloadUser struct {
	// Full name of the commit author
	Name string `json:"name"`
	// swagger:strfmt email
	Email    string `json:"email"`
	UserName string `json:"username"`
}

// User represents a user
// swagger:model
type User struct {
	// the user's id
	ID int64 `json:"id"`
	// the user's username
	UserName string `json:"login"`
	// the user's authentication sign-in name.
	// default: empty
	LoginName string `json:"login_name"`
	// The ID of the user's Authentication Source
	SourceID int64 `json:"source_id"`
	// the user's full name
	FullName string `json:"full_name"`
	// swagger:strfmt email
	Email string `json:"email"`
	// URL to the user's avatar
	AvatarURL string `json:"avatar_url"`
	// URL to the user's gitea page
	HTMLURL string `json:"html_url"`
	// User locale
	Language string `json:"language"`
	// Is the user an administrator
	IsAdmin bool `json:"is_admin"`
	// swagger:strfmt date-time
	LastLogin time.Time `json:"last_login,omitempty"`
	// swagger:strfmt date-time
	Created time.Time `json:"created,omitempty"`
	// Is user restricted
	Restricted bool `json:"restricted"`
	// Is user active
	IsActive bool `json:"active"`
	// Is user login prohibited
	ProhibitLogin bool `json:"prohibit_login"`
	// the user's location
	Location string `json:"location"`
	// the user's website
	Website string `json:"website"`
	// the user's description
	Description string `json:"description"`
	// User visibility level option: public, limited, private
	Visibility string `json:"visibility"`

	// user counts
	Followers    int `json:"followers_count"`
	Following    int `json:"following_count"`
	StarredRepos int `json:"starred_repos_count"`
}

// Permission represents a set of permissions
type Permission struct {
	Admin bool `json:"admin"` // Admin indicates if the user is an administrator of the repository.
	Push  bool `json:"push"`  // Push indicates if the user can push code to the repository.
	Pull  bool `json:"pull"`  // Pull indicates if the user can pull code from the repository.
}

// InternalTracker represents settings for internal tracker
// swagger:model
type InternalTracker struct {
	// Enable time tracking (Built-in issue tracker)
	EnableTimeTracker bool `json:"enable_time_tracker"`
	// Let only contributors track time (Built-in issue tracker)
	AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
	// Enable dependencies for issues and pull requests (Built-in issue tracker)
	EnableIssueDependencies bool `json:"enable_issue_dependencies"`
}

// ExternalTracker represents settings for external tracker
// swagger:model
type ExternalTracker struct {
	// URL of external issue tracker.
	ExternalTrackerURL string `json:"external_tracker_url"`
	// External Issue Tracker URL Format. Use the placeholders {user}, {repo} and {index} for the username, repository name and issue index.
	ExternalTrackerFormat string `json:"external_tracker_format"`
	// External Issue Tracker Number Format, either `numeric`, `alphanumeric`, or `regexp`
	ExternalTrackerStyle string `json:"external_tracker_style"`
	// External Issue Tracker issue regular expression
	ExternalTrackerRegexpPattern string `json:"external_tracker_regexp_pattern"`
}

// ExternalWiki represents setting for external wiki
// swagger:model
type ExternalWiki struct {
	// URL of external wiki.
	ExternalWikiURL string `json:"external_wiki_url"`
}
