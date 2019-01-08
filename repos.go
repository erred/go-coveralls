package coveralls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// https://coveralls.io/api/docs

// POST /api/repos
func (c *Client) AddRepository(repo Repository) (*Repository, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(repo); err != nil {
		return nil, fmt.Errorf("encode repository: %v", err)
	}

	req, err := c.NewRequest(http.MethodPost, URL("/api/repos", ""), &buf)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	var r Repository
	if err = c.Do(req, &r); err != nil {
		return nil, fmt.Errorf("exec request: %v", err)
	}
	return &r, nil
}

type Repository struct {
	Service         GitProvider `json:"service"`
	Name            string      `json:"name"`
	PRComment       *bool       `json:"comment_on_pull_request,omitempty"`
	SendStatus      *bool       `json:"send_build_status,omitempty"`
	FailThreshold   *float32    `json:"commit_status_fail_threshold,omitempty"`
	ChangeThreshold *float32    `json:"commit_status_change_threshold,omitempty"`

	// Only in responses
	// time.RFC3339
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
