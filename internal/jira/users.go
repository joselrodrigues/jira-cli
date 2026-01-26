package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UserSearchResult struct {
	AccountID    string `json:"accountId,omitempty"`
	Name         string `json:"name,omitempty"`
	Key          string `json:"key,omitempty"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress,omitempty"`
	Active       bool   `json:"active"`
	AccountType  string `json:"accountType,omitempty"`
}

func (u *UserSearchResult) GetIdentifier(isCloud bool) string {
	if isCloud {
		return u.AccountID
	}
	return u.Name
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

func (c *Client) AssignIssue(issueKey, userIdentifier string) error {
	var body map[string]interface{}
	if userIdentifier == "" {
		body = map[string]interface{}{"accountId": nil, "name": nil}
	} else if c.isCloud {
		body = map[string]interface{}{"accountId": userIdentifier}
	} else {
		body = map[string]interface{}{"name": userIdentifier}
	}

	_, err := c.Put(fmt.Sprintf("/issue/%s/assignee", issueKey), body)
	return err
}
