package coveralls

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// https://docs.coveralls.io/api-introduction
//
// Latest build
// /:provider/:owner/:repo.json
// 10 builds
// /:provider/:owner/:repo.json?page=:page
// Job by id
// /jobs/:id.json
// File by id
// /files/:id.json
// build by id
// /builds/:id.json
// build report for paths
// /builds/:id.json?paths=:comma_separated_paths_with_globbing
// source line coverage
// /builds/:id/source.json?filename=:filename

type BuildPage struct {
	Page   int     `json:"page"`
	Pages  int     `json:"pages"`
	Builds []Build `json:"builds"`
}

type Build struct {
	CreatedAt      time.Time `json:"created_at"`
	URL            string    `json:"url"`
	CommitMsg      string    `json:"commit_message"`
	Branch         string    `json:"branch"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	CommitSHA      string    `json:"commit_sha"`
	RepoName       string    `json:"repo_name"`
	BadgeURL       string    `json:"badge_url"`
	CoverageChange float32   `json:"coverage_change"`
	CoveredPercent float32   `json:"covered_percent"`

	// only for ?paths=...
	Paths string `json:"paths,omitempty"`
	Count int    `json:"selected_source_files_count,omitempty"`
}

// https://coveralls.io/:provider/:owner/:repo.json
func (c *Client) LatestBuild(prov GitProvider, owner, repo string) (Build, error) {
	u := URL(fmt.Sprintf("/%v/%v/%v.json", prov, owner, repo), "")
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Build{}, fmt.Errorf("create request: %v", err)
	}

	var build Build
	if err := c.Do(req, &build); err != nil {
		return Build{}, fmt.Errorf("exec request: %v", err)
	}
	return build, nil
}

// https://coveralls.io/:provider/:owner/:repo.json?page=:page
// page >= 1
func (c *Client) Builds(prov GitProvider, owner, repo string, page int) (BuildPage, error) {
	val := url.Values{}
	val.Set("page", strconv.Itoa(page))
	u := URL(fmt.Sprintf("/%v/%v/%v.json", prov, owner, repo), val.Encode())
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return BuildPage{}, fmt.Errorf("create request: %v", err)
	}

	var b BuildPage
	if err := c.Do(req, &b); err != nil {
		return BuildPage{}, fmt.Errorf("exec request: %v", err)
	}
	return b, nil
}

// specify 1
// https://coveralls.io/builds/:buildid.json
func (c *Client) Build(buildID int64, shasum string, paths ...string) (Build, error) {
	id := shasum
	if id == "" {
		id = strconv.FormatInt(buildID, 10)
	}
	q := url.QueryEscape(strings.Join(paths, ","))
	if q != "" {
		q = "pages=" + q
	}

	u := URL(fmt.Sprintf("/builds/%v.json", id), q)
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Build{}, fmt.Errorf("create request: %v", err)
	}

	var b Build
	if err := c.Do(req, &b); err != nil {
		return Build{}, fmt.Errorf("exec request: %v", err)
	}
	return b, nil
}

type Job struct {
	RepoName       string    `json:"repo_name"`
	FullNumber     float32   `json:"full_number"`
	Timestamp      time.Time `json:"timestamp"`
	CoveredPercent float32   `json:"covered_percent"`
}

func (c *Client) GetJob(jobID int64) (Job, error) {
	// https://coveralls.io/jobs/:jobid.json
	u := URL(fmt.Sprintf("/jobs/%v.json", jobID), "")
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Job{}, fmt.Errorf("create request: %v", err)
	}

	var j Job
	if err := c.Do(req, &j); err != nil {
		return Job{}, fmt.Errorf("exec request: %v", err)
	}
	return j, nil
}

type Coverage []*int

// null -> -1
func (c Coverage) Deref() []int {
	res := make([]int, len(c))
	for i, v := range c {
		if v == nil {
			res[i] = -1
		} else {
			res[i] = *v
		}
	}
	return res
}

func (c Coverage) String() string {
	b := make([]string, len(c))
	for i, v := range c {
		if v != nil {
			b[i] = strconv.Itoa(*v)
		} else {
			b[i] = "null"
		}
	}
	return "[" + strings.Join(b, " ") + "]"
}

func (c *Client) GetBuildFile(buildID int64, shasum string, filename string) (Coverage, error) {
	id := shasum
	if id == "" {
		id = strconv.FormatInt(buildID, 10)
	}
	q := "filename=" + url.QueryEscape(filename)
	u := URL(fmt.Sprintf("/builds/%v.json", id), q)
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Coverage{}, fmt.Errorf("create request: %v", err)
	}

	var cov Coverage
	if err := c.Do(req, &cov); err != nil {
		return Coverage{}, fmt.Errorf("exec request: %v", err)
	}
	return cov, nil
}

// DEPRECIATED??
// https://coveralls.io/files/:fileid.json
// redirects to
// jobs/:jobid/source_files/:fileid.json
func (c *Client) GetJobFile(jobID, fileID int64) (Coverage, error) {
	u := URL(fmt.Sprintf("/jobs/%v/source_files/%v.json", jobID, fileID), "")
	req, err := c.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return Coverage{}, fmt.Errorf("create request: %v", err)
	}

	var cov Coverage
	if err := c.Do(req, &cov); err != nil {
		return Coverage{}, fmt.Errorf("exec request: %v", err)
	}
	return cov, nil
}
