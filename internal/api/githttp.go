package api

import (
	"bytes"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (s *Server) gitHTTP(w http.ResponseWriter, r *http.Request) {
	repoPath := strings.TrimPrefix(r.URL.Path, "/git/")
	repoPath = strings.TrimPrefix(repoPath, "/")
	if repoPath == "" {
		writeError(w, http.StatusNotFound, "repository path is required")
		return
	}

	if !strings.Contains(repoPath, ".git") {
		writeError(w, http.StatusNotFound, "repository not found")
		return
	}

	backend, err := gitHTTPBackendPath()
	if err != nil {
		log.Printf("git http-backend: %v", err)
		writeError(w, http.StatusInternalServerError, "git http-backend is not available")
		return
	}

	var stderr bytes.Buffer
	handler := &cgi.Handler{
		Path:   backend,
		Root:   "/git",
		Dir:    s.config.Root,
		Env:    gitHTTPEnv(s.config.Root),
		Stderr: &stderr,
	}
	handler.ServeHTTP(w, r)

	if stderr.Len() > 0 {
		log.Printf("git http-backend: %s", stderr.String())
	}
}

func gitHTTPEnv(projectRoot string) []string {
	return []string{
		"GIT_PROJECT_ROOT=" + projectRoot,
		"GIT_HTTP_EXPORT_ALL=1",
	}
}

func gitHTTPBackendPath() (string, error) {
	if path, err := exec.LookPath("git-http-backend"); err == nil {
		return path, nil
	}

	out, err := exec.Command("git", "--exec-path").Output()
	if err != nil {
		return "", err
	}

	execPath := strings.TrimSpace(string(out))
	if execPath == "" {
		return "", exec.ErrNotFound
	}

	candidate := filepath.Join(execPath, "git-http-backend")
	if _, err := os.Stat(candidate); err != nil {
		return "", err
	}

	return candidate, nil
}
