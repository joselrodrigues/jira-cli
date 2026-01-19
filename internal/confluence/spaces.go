package confluence

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Space struct {
	ID     int        `json:"id"`
	Key    string     `json:"key"`
	Name   string     `json:"name"`
	Status string     `json:"status"`
	Type   string     `json:"type"`
	Links  SpaceLinks `json:"_links"`
}

type SpaceLinks struct {
	WebUI string `json:"webui"`
	Self  string `json:"self"`
}

type SpacesResponse struct {
	Results []Space `json:"results"`
	Start   int     `json:"start"`
	Limit   int     `json:"limit"`
	Size    int     `json:"size"`
	Links   struct {
		Next string `json:"next"`
	} `json:"_links"`
}

type SpaceContentResponse struct {
	Page     PageResults `json:"page"`
	Blogpost PageResults `json:"blogpost"`
}

type PageResults struct {
	Results []Page `json:"results"`
	Start   int    `json:"start"`
	Limit   int    `json:"limit"`
	Size    int    `json:"size"`
	Links   struct {
		Next string `json:"next"`
	} `json:"_links"`
}

func (c *Client) GetSpace(spaceKey string) (*Space, error) {
	endpoint := fmt.Sprintf("/space/%s", spaceKey)
	data, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var space Space
	if err := json.Unmarshal(data, &space); err != nil {
		return nil, fmt.Errorf("failed to parse space: %w", err)
	}

	return &space, nil
}

func (c *Client) ListSpaces(limit int) (*SpacesResponse, error) {
	params := url.Values{}
	params.Set("limit", fmt.Sprintf("%d", limit))

	endpoint := "/space?" + params.Encode()
	data, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var spaces SpacesResponse
	if err := json.Unmarshal(data, &spaces); err != nil {
		return nil, fmt.Errorf("failed to parse spaces: %w", err)
	}

	return &spaces, nil
}

func (c *Client) GetSpaceContent(spaceKey string, contentType string, limit int) (*PageResults, error) {
	params := url.Values{}
	params.Set("limit", fmt.Sprintf("%d", limit))

	var endpoint string
	if contentType != "" {
		endpoint = fmt.Sprintf("/space/%s/content/%s?%s", spaceKey, contentType, params.Encode())
	} else {
		endpoint = fmt.Sprintf("/space/%s/content?%s", spaceKey, params.Encode())
	}

	data, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		var pages PageResults
		if err := json.Unmarshal(data, &pages); err != nil {
			return nil, fmt.Errorf("failed to parse pages: %w", err)
		}
		return &pages, nil
	}

	var content SpaceContentResponse
	if err := json.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}
	return &content.Page, nil
}
