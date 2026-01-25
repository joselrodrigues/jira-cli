package jira

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Board struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Location BoardLocation `json:"location,omitempty"`
}

type BoardLocation struct {
	ProjectID   int    `json:"projectId,omitempty"`
	ProjectKey  string `json:"projectKey,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
}

type BoardsResponse struct {
	MaxResults int     `json:"maxResults"`
	StartAt    int     `json:"startAt"`
	Total      int     `json:"total"`
	IsLast     bool    `json:"isLast"`
	Values     []Board `json:"values"`
}

type Sprint struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	CompleteDate  string `json:"completeDate,omitempty"`
	OriginBoardID int    `json:"originBoardId,omitempty"`
	Goal          string `json:"goal,omitempty"`
}

type SprintsResponse struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

func (c *Client) doAgileRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := fmt.Sprintf("%s/rest/agile/1.0%s", c.baseURL, endpoint)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	}

	respBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) GetBoards(projectKey string) (*BoardsResponse, error) {
	params := url.Values{}
	params.Set("maxResults", "100")
	if projectKey != "" {
		params.Set("projectKeyOrId", projectKey)
	}

	data, err := c.doAgileRequest(http.MethodGet, "/board?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result BoardsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse boards: %w", err)
	}

	return &result, nil
}

func (c *Client) GetSprints(boardID int, state string) (*SprintsResponse, error) {
	params := url.Values{}
	params.Set("maxResults", "50")
	if state != "" {
		params.Set("state", state)
	}

	endpoint := fmt.Sprintf("/board/%d/sprint?%s", boardID, params.Encode())
	data, err := c.doAgileRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var result SprintsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse sprints: %w", err)
	}

	return &result, nil
}

func (c *Client) MoveToSprint(sprintID int, issueKeys []string) error {
	body := map[string]interface{}{
		"issues": issueKeys,
	}

	endpoint := fmt.Sprintf("/sprint/%d/issue", sprintID)
	_, err := c.doAgileRequest(http.MethodPost, endpoint, body)
	return err
}
