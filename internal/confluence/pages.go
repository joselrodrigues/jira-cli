package confluence

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Page struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Status  string    `json:"status"`
	Title   string    `json:"title"`
	Space   *Space    `json:"space,omitempty"`
	Version *Version  `json:"version,omitempty"`
	Body    *PageBody `json:"body,omitempty"`
	Links   PageLinks `json:"_links"`
}

type Version struct {
	Number    int    `json:"number"`
	By        *User  `json:"by,omitempty"`
	When      string `json:"when"`
	Message   string `json:"message"`
	MinorEdit bool   `json:"minorEdit"`
}

type User struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type PageBody struct {
	Storage *BodyContent `json:"storage,omitempty"`
	View    *BodyContent `json:"view,omitempty"`
}

type BodyContent struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type PageLinks struct {
	WebUI string `json:"webui"`
	Edit  string `json:"edit"`
	Self  string `json:"self"`
}

type CreatePageRequest struct {
	Type      string           `json:"type"`
	Title     string           `json:"title"`
	Space     CreatePageSpace  `json:"space"`
	Ancestors []CreateAncestor `json:"ancestors,omitempty"`
	Body      CreatePageBody   `json:"body"`
}

type CreatePageSpace struct {
	Key string `json:"key"`
}

type CreateAncestor struct {
	ID string `json:"id"`
}

type CreatePageBody struct {
	Storage BodyContent `json:"storage"`
}

type UpdatePageRequest struct {
	Type    string         `json:"type"`
	Title   string         `json:"title"`
	Body    CreatePageBody `json:"body"`
	Version UpdateVersion  `json:"version"`
}

type UpdateVersion struct {
	Number  int    `json:"number"`
	Message string `json:"message,omitempty"`
}

func (c *Client) GetPage(pageID string, expand []string) (*Page, error) {
	params := url.Values{}
	if len(expand) > 0 {
		params.Set("expand", strings.Join(expand, ","))
	}

	endpoint := fmt.Sprintf("/content/%s", pageID)
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	data, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("failed to parse page: %w", err)
	}

	return &page, nil
}

func (c *Client) CreatePage(spaceKey, title, body string, parentID string) (*Page, error) {
	req := CreatePageRequest{
		Type:  "page",
		Title: title,
		Space: CreatePageSpace{Key: spaceKey},
		Body: CreatePageBody{
			Storage: BodyContent{
				Value:          body,
				Representation: "storage",
			},
		},
	}

	if parentID != "" {
		req.Ancestors = []CreateAncestor{{ID: parentID}}
	}

	data, err := c.Post("/content", req)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("failed to parse created page: %w", err)
	}

	return &page, nil
}

func (c *Client) UpdatePage(pageID, title, body string, currentVersion int, message string) (*Page, error) {
	req := UpdatePageRequest{
		Type:  "page",
		Title: title,
		Body: CreatePageBody{
			Storage: BodyContent{
				Value:          body,
				Representation: "storage",
			},
		},
		Version: UpdateVersion{
			Number:  currentVersion + 1,
			Message: message,
		},
	}

	endpoint := fmt.Sprintf("/content/%s", pageID)
	data, err := c.Put(endpoint, req)
	if err != nil {
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("failed to parse updated page: %w", err)
	}

	return &page, nil
}

func (c *Client) DeletePage(pageID string) error {
	endpoint := fmt.Sprintf("/content/%s", pageID)
	_, err := c.Delete(endpoint)
	return err
}
