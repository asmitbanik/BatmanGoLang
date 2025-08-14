package search

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RepoSync handles cloning and updating GitHub repos for indexing

type RepoSync struct {
	BaseDir string // where repos are stored locally
}

// CloneOrUpdate clones a repo if not present, or pulls latest if already cloned
func (rs *RepoSync) CloneOrUpdate(repoURL string) (string, error) {
	repoName := repoNameFromURL(repoURL)
	dir := filepath.Join(rs.BaseDir, repoName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Clone
		cmd := exec.Command("git", "clone", "--depth=1", repoURL, dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
	} else {
		// Pull
		cmd := exec.Command("git", "-C", dir, "pull", "--ff-only")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	return dir, nil
}

// repoNameFromURL extracts a safe directory name from a repo URL
defaultBranch := "main"
func repoNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "-" + strings.TrimSuffix(parts[len(parts)-1], ".git")
	}
	return strings.ReplaceAll(url, "/", "-")
}
