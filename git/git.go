package git

import "github.com/ggdream/mcc/payload"

type Git interface {
	Name() string
	Source() Source
	Event() Event
	GitEvent
}

type GitEvent interface {
	GetPushPayload(data []byte) (*payload.PushPayload, error)
}

type Source string

func (s Source) String() string {
	return string(s)
}

const (
	SourceGitea  Source = "gitea"
	SourceGithub Source = "github"
	SourceGitlab Source = "gitlab"
)

type Event string

func (e Event) String() string {
	return string(e)
}

const (
	Unknown                   Event = "unknown"
	Create                    Event = "create"
	Delete                    Event = "delete"
	Fork                      Event = "fork"
	Push                      Event = "push"
	Issues                    Event = "issues"
	IssueAssign               Event = "issue_assign"
	IssueLabel                Event = "issue_label"
	IssueMilestone            Event = "issue_milestone"
	IssueComment              Event = "issue_comment"
	PullRequest               Event = "pull_request"
	PullRequestAssign         Event = "pull_request_assign"
	PullRequestLabel          Event = "pull_request_label"
	PullRequestMilestone      Event = "pull_request_milestone"
	PullRequestComment        Event = "pull_request_comment"
	PullRequestReviewApproved Event = "pull_request_review_approved"
	PullRequestReviewRejected Event = "pull_request_review_rejected"
	PullRequestReviewComment  Event = "pull_request_review_comment"
	PullRequestSync           Event = "pull_request_sync"
	PullRequestReviewRequest  Event = "pull_request_review_request"
	Wiki                      Event = "wiki"
	Repository                Event = "repository"
	Release                   Event = "release"
	Package                   Event = "package"
	Schedule                  Event = "schedule"
)
