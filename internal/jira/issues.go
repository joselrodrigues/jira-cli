package jira

import (
	"encoding/json"
	"fmt"
)

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description,omitempty"`
	Status      Status      `json:"status,omitempty"`
	Priority    Priority    `json:"priority,omitempty"`
	Assignee    *User       `json:"assignee,omitempty"`
	Reporter    *User       `json:"reporter,omitempty"`
	Project     Project     `json:"project,omitempty"`
	IssueType   IssueType   `json:"issuetype,omitempty"`
	StoryPoints float64     `json:"customfield_10106,omitempty"`
	Sprint      interface{} `json:"customfield_10104,omitempty"`
}

type Status struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
}

type User struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress,omitempty"`
}

type Project struct {
	Key  string `json:"key"`
	Name string `json:"name,omitempty"`
}

type IssueType struct {
	Name string `json:"name"`
}

type SearchResult struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}

type CreateIssueRequest struct {
	Fields CreateIssueFields `json:"fields"`
}

type CreateIssueFields struct {
	Project     Project   `json:"project"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	IssueType   IssueType `json:"issuetype"`
}

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type UpdateIssueRequest struct {
	Fields map[string]interface{} `json:"fields"`
}

func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	data, err := c.Get(fmt.Sprintf("/issue/%s", issueKey))
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse issue: %w", err)
	}

	return &issue, nil
}

func (c *Client) CreateIssue(project, issueType, summary, description string) (*CreateIssueResponse, error) {
	req := CreateIssueRequest{
		Fields: CreateIssueFields{
			Project:     Project{Key: project},
			Summary:     summary,
			Description: description,
			IssueType:   IssueType{Name: issueType},
		},
	}

	data, err := c.Post("/issue", req)
	if err != nil {
		return nil, err
	}

	var resp CreateIssueResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

func (c *Client) UpdateIssue(issueKey string, fields map[string]interface{}) error {
	req := UpdateIssueRequest{Fields: fields}
	_, err := c.Put(fmt.Sprintf("/issue/%s", issueKey), req)
	return err
}

func (c *Client) SearchIssues(jql string, maxResults int) (*SearchResult, error) {
	fields := []string{"key", "summary", "status", "priority", "assignee", "customfield_10106"}
	data, err := c.Search(jql, fields, maxResults)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &result, nil
}

func (c *Client) GetMyIssues() (*SearchResult, error) {
	jql := "assignee = currentUser() AND status NOT IN (Done, Closed, Listo, CERRADO)"
	return c.SearchIssues(jql, 50)
}

func (c *Client) GetSprintIssues(project string) (*SearchResult, error) {
	jql := fmt.Sprintf("project = %s AND sprint in openSprints()", project)
	return c.SearchIssues(jql, 100)
}