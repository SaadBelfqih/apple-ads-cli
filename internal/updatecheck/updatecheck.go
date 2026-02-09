package updatecheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

const defaultRepo = "SaadBelfqih/apple-ads-cli"

type cacheFile struct {
	LastChecked time.Time `json:"last_checked"`
	LatestTag   string    `json:"latest_tag"`
	LatestURL   string    `json:"latest_url"`
	Prerelease  bool      `json:"prerelease"`
}

func Enabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("AADS_UPDATE_CHECK")))
	if v == "" {
		return true
	}
	switch v {
	case "0", "false", "no", "off":
		return false
	default:
		return true
	}
}

func ttlFromEnv() time.Duration {
	v := strings.TrimSpace(os.Getenv("AADS_UPDATE_CHECK_TTL"))
	if v == "" {
		return 24 * time.Hour
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 24 * time.Hour
	}
	return d
}

func repoFromEnv() string {
	v := strings.TrimSpace(os.Getenv("AADS_UPDATE_CHECK_REPO"))
	if v == "" {
		return defaultRepo
	}
	return v
}

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".aads", "update_check.json"), nil
}

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func normalizeSemver(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if !strings.HasPrefix(s, "v") {
		s = "v" + s
	}
	return s
}

func compareVersions(latestTag, current string) (updateAvailable bool) {
	latest := normalizeSemver(latestTag)
	if !semver.IsValid(latest) {
		return false
	}

	cur := normalizeSemver(current)
	if !semver.IsValid(cur) {
		// Dev/commit builds: if there's a valid semver release, notify.
		return true
	}
	return semver.Compare(latest, cur) > 0
}

func readCache(path string) (*cacheFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c cacheFile
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func writeCache(path string, c *cacheFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

type githubRelease struct {
	TagName    string `json:"tag_name"`
	HTMLURL    string `json:"html_url"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
}

type apiError struct {
	StatusCode int
	Status     string
}

func (e *apiError) Error() string {
	if e == nil {
		return ""
	}
	if e.Status != "" {
		return "github api status " + e.Status
	}
	return fmt.Sprintf("github api status %d", e.StatusCode)
}

func fetchLatestRelease(ctx context.Context, repo string) (*githubRelease, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=5", repo), nil)
	if err != nil {
		return nil, err
	}
	// Unauthenticated is fine for low volume; this avoids requiring user tokens.
	req.Header.Set("Accept", "application/vnd.github+json")
	// If the repo is private, users can opt-in by providing a token.
	if tok := strings.TrimSpace(os.Getenv("AADS_GITHUB_TOKEN")); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	} else if tok := strings.TrimSpace(os.Getenv("GH_TOKEN")); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	} else if tok := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &apiError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}
	for _, r := range releases {
		if r.Draft {
			continue
		}
		if strings.TrimSpace(r.TagName) == "" {
			continue
		}
		return &r, nil
	}
	return nil, errors.New("no releases found")
}

// MaybeNotify checks GitHub for a newer release and prints a one-line notice to w.
// It never downloads or installs anything.
func MaybeNotify(ctx context.Context, w io.Writer, currentVersion string) {
	if !Enabled() {
		return
	}
	// Avoid corrupting piped output.
	if !isTTY(os.Stdout) {
		return
	}

	ttl := ttlFromEnv()
	repo := repoFromEnv()

	p, err := cachePath()
	if err == nil {
		if c, err := readCache(p); err == nil {
			if !c.LastChecked.IsZero() && time.Since(c.LastChecked) < ttl {
				if c.LatestTag != "" && compareVersions(c.LatestTag, currentVersion) {
					fmt.Fprintf(w, "Update available: %s (current: %s). See %s\n", c.LatestTag, currentVersion, c.LatestURL)
				}
				return
			}
		}
	}

	// Keep this fast; if the network is slow we just skip.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	r, err := fetchLatestRelease(ctx, repo)
	if err != nil {
		// Throttle repeated failed checks (e.g., private repo without token).
		if p != "" {
			_ = writeCache(p, &cacheFile{LastChecked: time.Now()})
		}
		return
	}

	if p != "" {
		_ = writeCache(p, &cacheFile{
			LastChecked: time.Now(),
			LatestTag:   r.TagName,
			LatestURL:   r.HTMLURL,
			Prerelease:  r.Prerelease,
		})
	}

	if compareVersions(r.TagName, currentVersion) {
		fmt.Fprintf(w, "Update available: %s (current: %s). See %s\n", r.TagName, currentVersion, r.HTMLURL)
	}
}

// CheckLatest prints the latest release info (intended for `aads update --check`).
func CheckLatest(ctx context.Context, w io.Writer, currentVersion string) error {
	repo := repoFromEnv()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	r, err := fetchLatestRelease(ctx, repo)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Current: %s\n", currentVersion)
	fmt.Fprintf(w, "Latest:  %s\n", r.TagName)
	fmt.Fprintf(w, "URL:     %s\n", r.HTMLURL)
	if compareVersions(r.TagName, currentVersion) {
		fmt.Fprintln(w, "Status:  update available")
	} else {
		fmt.Fprintln(w, "Status:  up-to-date")
	}
	return nil
}
