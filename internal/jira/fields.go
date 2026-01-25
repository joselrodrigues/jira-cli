package jira

import (
	"encoding/json"
	"fmt"
)

type Field struct {
	ID          string      `json:"id"`
	Key         string      `json:"key,omitempty"`
	Name        string      `json:"name"`
	Custom      bool        `json:"custom"`
	Schema      FieldSchema `json:"schema,omitempty"`
	Description string      `json:"description,omitempty"`
}

type FieldSchema struct {
	Type     string `json:"type,omitempty"`
	Custom   string `json:"custom,omitempty"`
	CustomID int    `json:"customId,omitempty"`
}

func (c *Client) GetFields() ([]Field, error) {
	data, err := c.Get("/field")
	if err != nil {
		return nil, err
	}

	var fields []Field
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, fmt.Errorf("failed to parse fields: %w", err)
	}

	return fields, nil
}
