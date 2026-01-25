package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UserSearchResult struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress,omitempty"`
	Active       bool   `json:"active"`
	AccountType  string `json:"accountType,omitempty"`
}

func (c *Client) SearchUsers(query string) ([]UserSearchResult, error) {
	params := url.Values{}
	params.Set("username", query)
	params.Set("maxResults", "50")

	data, err := c.Get("/user/search?" + params.Encode())
	if err != nil {
		return nil, err
	}

	var users []UserSearchResult
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	return users, nil
}

func (c *Client) AssignIssue(issueKey, accountID string) error {
	body := map[string]interface{}{
		"accountId": accountID,
	}
	if accountID == "" {
		body["accountId"] = nil
	}

	_, err := c.Put(fmt.Sprintf("/issue/%s/assignee", issueKey), body)
	return err
}
