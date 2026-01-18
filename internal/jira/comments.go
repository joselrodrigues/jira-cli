package jira

import (
	"encoding/json"
	"fmt"
)

type Comment struct {
	ID      string `json:"id"`
	Body    string `json:"body"`
	Author  User   `json:"author"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
}

type AddCommentRequest struct {
	Body string `json:"body"`
}

func (c *Client) GetComments(issueKey string) (*CommentsResponse, error) {
	data, err := c.Get(fmt.Sprintf("/issue/%s/comment", issueKey))
	if err != nil {
		return nil, err
	}

	var resp CommentsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse comments: %w", err)
	}

	return &resp, nil
}

func (c *Client) AddComment(issueKey, body string) (*Comment, error) {
	req := AddCommentRequest{Body: body}
	data, err := c.Post(fmt.Sprintf("/issue/%s/comment", issueKey), req)
	if err != nil {
		return nil, err
	}

	var comment Comment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, fmt.Errorf("failed to parse comment: %w", err)
	}

	return &comment, nil
}