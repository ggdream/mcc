package git

import (
	"encoding/json"

	"github.com/go-playground/webhooks/v6/github"

	"github.com/ggdream/mcc/payload"
)

var _ Git = (*GitHub)(nil)

type GitHub struct {
	event string
}

func NewGithub(event string) *GitHub {
	return &GitHub{
		event: event,
	}
}

// Event implements Git.
func (g *GitHub) Event() Event {
	switch g.event {
	case "push":
		return Push
	default:
		return Unknown
	}
}

// GetPushPayload implements Git.
func (g *GitHub) GetPushPayload(data []byte) (*payload.PushPayload, error) {
	var p github.PushPayload
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}

	res := &payload.PushPayload{
		Ref:    p.Ref,
		After:  p.After,
		Before: p.Before,
		Repo: &payload.Repository{
			ID:        p.Repository.ID,
			Name:      p.Repository.Name,
			FullName:  p.Repository.FullName,
			URL:       p.Repository.HTMLURL,
			SSHURL:    p.Repository.SSHURL,
			CloneURL:  p.Repository.CloneURL,
			Private:   p.Repository.Private,
			Owner: &payload.User{
				ID:       p.Repository.Owner.ID,
			},
			HTMLURL:       p.Repository.HTMLURL,
			Fork:          p.Repository.Fork,
			DefaultBranch: p.Repository.DefaultBranch,
			Stars:         int(p.Repository.Stargazers),
			Watchers:      int(p.Repository.Watchers),
			Forks:         int(p.Repository.Forks),
			Size:          int(p.Repository.Size),
			Description:   p.Repository.Description,
		},
		HeadCommit: &payload.PayloadCommit{
			ID:      p.HeadCommit.ID,
			Message: p.HeadCommit.Message,
			URL:     p.HeadCommit.URL,
			Author: &payload.PayloadUser{
				Name:  p.HeadCommit.Author.Name,
				Email: p.HeadCommit.Author.Email,
			},
			Committer: &payload.PayloadUser{
				Name: p.HeadCommit.Committer.Name,
			},
			Added:    p.HeadCommit.Added,
			Modified: p.HeadCommit.Modified,
			Removed:  p.HeadCommit.Removed,
			// Verification: &payload.PayloadCommitVerification{
			// 	Payload:   p.HeadCommit.Verification.Payload,
			// 	Reason:    p.HeadCommit.Verification.Reason,
			// 	Signature: p.HeadCommit.Verification.Signature,
			// },
		},
		Pusher: &payload.User{
			Email:    p.Pusher.Email,
			UserName: p.Pusher.Name,
		},
		Sender: &payload.User{
			ID:       p.Sender.ID,
		},
	}
	for _, v := range p.Commits {
		res.Commits = append(res.Commits, &payload.PayloadCommit{
			ID:      v.ID,
			Message: v.Message,
			URL:     v.URL,
			Author: &payload.PayloadUser{
				Name:     v.Author.Name,
				Email:    v.Author.Email,
				UserName: v.Author.Username,
			},
			Committer: &payload.PayloadUser{
				Name:     v.Committer.Name,
				Email:    v.Committer.Email,
				UserName: v.Committer.Username,
			},
			Added:     v.Added,
			Modified:  v.Modified,
			Removed:   v.Removed,
			// Verification: &payload.PayloadCommitVerification{
			// 	Payload:   v.Verification.Payload,
			// 	Reason:    v.Verification.Reason,
			// 	Signature: v.Verification.Signature,
			// 	Signer: &payload.PayloadUser{
			// 		Name:     v.Verification.Signer.Name,
			// 		Email:    v.Verification.Signer.Email,
			// 		UserName: v.Verification.Signer.UserName,
			// 	},
			// 	Verified: v.Verification.Verified,
			// },
		})
	}

	return res, nil
}

// Name implements Git.
func (g *GitHub) Name() string {
	return SourceGithub.String()
}

// Source implements Git.
func (g *GitHub) Source() Source {
	return SourceGithub
}
