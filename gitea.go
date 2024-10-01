package main

import (
	"encoding/json"

	"code.gitea.io/gitea/modules/structs"

	"github.com/ggdream/mcc/payload"
)

var _ GitForward = (*Gitea)(nil)

type Gitea struct {
	event string
}

func NewGitea(event string) *Gitea {
	return &Gitea{
		event: event,
	}
}

func (Gitea) Name() string   { return SourceGitea.String() }
func (Gitea) Source() Source { return SourceGitea }

// Event implements GitForward.
func (g *Gitea) Event() Event {
	switch g.event {
	case "push":
		return Push
	default:
		return Unknown
	}
}

// GetPushPayload implements GitForward.
func (g *Gitea) GetPushPayload(data []byte) (*payload.PushPayload, error) {
	var p structs.PushPayload
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}

	res := &payload.PushPayload{
		Ref:          p.Ref,
		After:        p.After,
		Before:       p.Before,
		TotalCommits: p.TotalCommits,
		CompareURL:   p.CompareURL,
		Repo: &payload.Repository{
			ID:        p.Repo.ID,
			Name:      p.Repo.Name,
			FullName:  p.Repo.FullName,
			URL:       p.Repo.HTMLURL,
			SSHURL:    p.Repo.SSHURL,
			CloneURL:  p.Repo.CloneURL,
			Private:   p.Repo.Private,
			AvatarURL: p.Repo.AvatarURL,
			Owner: &payload.User{
				ID:       p.Repo.Owner.ID,
				UserName: p.Repo.Owner.UserName,
				Email:    p.Repo.Owner.Email,
			},
			HTMLURL:       p.Repo.HTMLURL,
			Fork:          p.Repo.Fork,
			Language:      p.Repo.Language,
			DefaultBranch: p.Repo.DefaultBranch,
			Stars:         p.Repo.Stars,
			Watchers:      p.Repo.Watchers,
			Forks:         p.Repo.Forks,
			Size:          p.Repo.Size,
			Description:   p.Repo.Description,
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
			Timestamp: p.HeadCommit.Timestamp,
			Added:     p.HeadCommit.Added,
			Modified:  p.HeadCommit.Modified,
			Removed:   p.HeadCommit.Removed,
			// Verification: &payload.PayloadCommitVerification{
			// 	Payload:   p.HeadCommit.Verification.Payload,
			// 	Reason:    p.HeadCommit.Verification.Reason,
			// 	Signature: p.HeadCommit.Verification.Signature,
			// },
		},
		Pusher: &payload.User{
			ID:       p.Pusher.ID,
			UserName: p.Pusher.UserName,
		},
		Sender: &payload.User{
			ID:       p.Sender.ID,
			UserName: p.Sender.UserName,
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
				UserName: v.Author.UserName,
			},
			Committer: &payload.PayloadUser{
				Name:     v.Committer.Name,
				Email:    v.Committer.Email,
				UserName: v.Committer.UserName,
			},
			Timestamp: v.Timestamp,
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
