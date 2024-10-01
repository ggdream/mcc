package git

import (
	"encoding/json"

	"github.com/go-playground/webhooks/v6/gitlab"

	"github.com/ggdream/mcc/payload"
)

var _ Git = (*GitLab)(nil)

type GitLab struct {
	event string
}

func NewGitlab(event string) *GitLab {
	return &GitLab{
		event: event,
	}
}

// Event implements Git.
func (g *GitLab) Event() Event {
	switch g.event {
	case string(gitlab.PushEvents):
		return Push
	default:
		return Unknown
	}
}

// GetPushPayload implements Git.
func (g *GitLab) GetPushPayload(data []byte) (*payload.PushPayload, error) {
	var p gitlab.PushEventPayload
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}

	res := &payload.PushPayload{
		Ref:          p.Ref,
		After:        p.After,
		Before:       p.Before,
		TotalCommits: int(p.TotalCommitsCount),
		Repo: &payload.Repository{
			ID:          p.Project.ID,
			Name:        p.Project.Name,
			FullName:    p.Project.PathWithNamespace,
			URL:         p.Project.URL,
			SSHURL:      p.Project.SSHURL,
			CloneURL:    p.Project.HTTPURL,
			HTMLURL:     p.Project.WebURL,
			Description: p.Project.Description,
		},
	}
	headCommit := p.Commits[0]
	res.HeadCommit = &payload.PayloadCommit{
		ID:      headCommit.ID,
		Message: headCommit.Message,
		URL:     headCommit.URL,
		Author: &payload.PayloadUser{
			Name:  headCommit.Author.Name,
			Email: headCommit.Author.Email,
		},
		Added:     headCommit.Added,
		Modified:  headCommit.Modified,
		Removed:   headCommit.Removed,
		Timestamp: headCommit.Timestamp.Time,
	}
	for _, v := range p.Commits {
		res.Commits = append(res.Commits, &payload.PayloadCommit{
			ID:      v.ID,
			Message: v.Message,
			URL:     v.URL,
			Author: &payload.PayloadUser{
				Name:  v.Author.Name,
				Email: v.Author.Email,
			},
			Added:     v.Added,
			Modified:  v.Modified,
			Removed:   v.Removed,
			Timestamp: v.Timestamp.Time,
		})
	}

	return res, nil
}

// Name implements Git.
func (g *GitLab) Name() string {
	return SourceGitlab.String()
}

// Source implements Git.
func (g *GitLab) Source() Source {
	return SourceGitlab
}
