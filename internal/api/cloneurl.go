package api

import (
	"strings"

	"github.com/superplanehq/superplane/supergit/internal/storage"
)

func cloneURL(publicURL, repoID string) string {
	publicURL = strings.TrimRight(strings.TrimSpace(publicURL), "/")
	repoID = strings.Trim(strings.TrimSpace(repoID), "/")
	if publicURL == "" || repoID == "" {
		return ""
	}

	return publicURL + "/" + repoID + ".git"
}

func (s *Server) withCloneURL(repo storage.Repository) storage.Repository {
	repo.CloneURL = cloneURL(s.config.PublicURL, repo.ID)
	return repo
}
