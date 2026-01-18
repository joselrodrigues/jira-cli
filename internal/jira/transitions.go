package jira

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   Status `json:"to"`
}

type TransitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

type DoTransitionRequest struct {
	Transition TransitionID `json:"transition"`
}

type TransitionID struct {
	ID string `json:"id"`
}

func (c *Client) GetTransitions(issueKey string) (*TransitionsResponse, error) {
	data, err := c.Get(fmt.Sprintf("/issue/%s/transitions", issueKey))
	if err != nil {
		return nil, err
	}

	var resp TransitionsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse transitions: %w", err)
	}

	return &resp, nil
}

func (c *Client) DoTransition(issueKey, transitionNameOrID string) error {
	transitions, err := c.GetTransitions(issueKey)
	if err != nil {
		return err
	}

	var transitionID string
	for _, t := range transitions.Transitions {
		if t.ID == transitionNameOrID || strings.EqualFold(t.Name, transitionNameOrID) {
			transitionID = t.ID
			break
		}
	}

	if transitionID == "" {
		available := make([]string, len(transitions.Transitions))
		for i, t := range transitions.Transitions {
			available[i] = fmt.Sprintf("%s (%s)", t.Name, t.ID)
		}
		return fmt.Errorf("transition '%s' not found. Available: %s", transitionNameOrID, strings.Join(available, ", "))
	}

	req := DoTransitionRequest{Transition: TransitionID{ID: transitionID}}
	_, err = c.Post(fmt.Sprintf("/issue/%s/transitions", issueKey), req)
	return err
}
