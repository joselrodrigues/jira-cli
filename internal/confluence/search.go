package confluence

import (
	"encoding/json"
	"fmt"
)

type SearchResponse struct {
	Results []Page `json:"results"`
	Start   int    `json:"start"`
	Limit   int    `json:"limit"`
	Size    int    `json:"size"`
	Links   struct {
		Next string `json:"next"`
	} `json:"_links"`
}

func (c *Client) SearchContent(cql string, expand []string, limit int) (*SearchResponse, error) {
	data, err := c.Search(cql, expand, limit)
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &result, nil
}
