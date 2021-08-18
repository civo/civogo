package civogo

import (
	"bytes"
	"encoding/json"
	"time"
)

// Team is a named group of users (has many members)
type Team struct {
	ID             string    `json:"id"`
	Name           string    `json:"name,omitempty"`
	OrganisationID string    `json:"organisation_id,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// TeamMember is a link record between User and Team.
type TeamMember struct {
	ID          string    `json:"id"`
	TeamID      string    `json:"team_id,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Permissions string    `json:"permissions,omitempty"`
	Roles       string    `json:"roles,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (c *Client) ListTeams() ([]Team, error) {
	resp, err := c.SendGetRequest("/v2/teams")
	if err != nil {
		return nil, decodeERROR(err)
	}

	teams := make([]Team, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&teams); err != nil {
		return nil, err
	}

	return teams, nil
}

func (c *Client) CreateTeam(name, organisationID, accountID string) (*Team, error) {
	data := map[string]string{"name": name, "organisation_id": organisationID, "account_id": accountID}
	resp, err := c.SendPostRequest("/v2/teams", data)
	if err != nil {
		return nil, decodeERROR(err)
	}

	team := &Team{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (c *Client) RenameTeam(teamID, name string) (*Team, error) {
	data := map[string]string{"name": name}
	resp, err := c.SendPutRequest("/v2/teams/"+teamID, data)
	if err != nil {
		return nil, decodeERROR(err)
	}

	team := &Team{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (c *Client) DeleteTeam(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest("/v2/teams/" + id)
	if err != nil {
		return nil, decodeERROR(err)
	}

	return c.DecodeSimpleResponse(resp)
}

func (c *Client) ListTeamMembers(teamID string) ([]TeamMember, error) {
	resp, err := c.SendGetRequest("/v2/teams/" + teamID + "/members")
	if err != nil {
		return nil, decodeERROR(err)
	}

	teamMembers := make([]TeamMember, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&teamMembers); err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (c *Client) AddTeamMember(teamID, userID, permissions, roles string) ([]TeamMember, error) {
	data := map[string]string{"user_id": userID, "permissions": permissions, "roles": roles}
	_, err := c.SendPostRequest("/v2/teams/"+teamID+"/members", data)
	if err != nil {
		return nil, decodeERROR(err)
	}

	return c.ListTeamMembers(teamID)
}

func (c *Client) UpdateTeamMember(teamID, teamMemberID, permissions, roles string) (*TeamMember, error) {
	data := map[string]string{"permissions": permissions, "roles": roles}
	resp, err := c.SendPostRequest("/v2/teams/"+teamID+"/members/"+teamMemberID, data)
	if err != nil {
		return nil, decodeERROR(err)
	}

	teamMember := &TeamMember{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(teamMember); err != nil {
		return nil, err
	}

	return teamMember, nil
}

func (c *Client) RemoveTeamMember(teamID, teamMemberID string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest("/v2/teams/" + teamID + "/members/" + teamMemberID)
	if err != nil {
		return nil, decodeERROR(err)
	}

	return c.DecodeSimpleResponse(resp)
}
