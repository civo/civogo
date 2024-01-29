package civogo

import (
	"fmt"
	"time"
)

// ListTeams implemented in a fake way for automated tests
func (c *FakeClient) ListTeams() ([]Team, error) {
	return c.OrganisationTeams, nil
}

// CreateTeam implemented in a fake way for automated tests
func (c *FakeClient) CreateTeam(name string) (*Team, error) {
	team := Team{
		ID:        c.generateID(),
		Name:      name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	c.OrganisationTeams = append(c.OrganisationTeams, team)
	return &team, nil
}

// RenameTeam implemented in a fake way for automated tests
func (c *FakeClient) RenameTeam(teamID, name string) (*Team, error) {
	for _, team := range c.OrganisationTeams {
		if team.ID == teamID {
			team.Name = name
			return &team, nil
		}
	}

	return nil, fmt.Errorf("unable to find that role")
}

// DeleteTeam implemented in a fake way for automated tests
func (c *FakeClient) DeleteTeam(id string) (*SimpleResponse, error) {
	for i, team := range c.OrganisationTeams {
		if team.ID == id {
			c.OrganisationTeams[len(c.OrganisationTeams)-1], c.OrganisationTeams[i] = c.OrganisationTeams[i], c.OrganisationTeams[len(c.OrganisationTeams)-1]
			c.OrganisationTeams = c.OrganisationTeams[:len(c.OrganisationTeams)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failure"}, fmt.Errorf("unable to find that team")
}

// ListTeamMembers implemented in a fake way for automated tests
func (c *FakeClient) ListTeamMembers(teamID string) ([]TeamMember, error) {
	return c.OrganisationTeamMembers[teamID], nil
}

// AddTeamMember implemented in a fake way for automated tests
func (c *FakeClient) AddTeamMember(teamID, userID, permissions, roles string) ([]TeamMember, error) {
	c.OrganisationTeamMembers[teamID] = append(c.OrganisationTeamMembers[teamID], TeamMember{
		ID:          c.generateID(),
		TeamID:      teamID,
		UserID:      userID,
		Permissions: permissions,
		Roles:       roles,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	return c.ListTeamMembers(teamID)
}

// UpdateTeamMember implemented in a fake way for automated tests
func (c *FakeClient) UpdateTeamMember(teamID, teamMemberID, permissions, roles string) (*TeamMember, error) {
	for _, teamMember := range c.OrganisationTeamMembers[teamID] {
		if teamMember.ID == teamMemberID {
			teamMember.Permissions = permissions
			teamMember.Roles = roles
			return &teamMember, nil
		}
	}

	return nil, fmt.Errorf("unable to find that role")
}

// RemoveTeamMember implemented in a fake way for automated tests
func (c *FakeClient) RemoveTeamMember(teamID, teamMemberID string) (*SimpleResponse, error) {
	for i, teamMember := range c.OrganisationTeamMembers[teamID] {
		if teamMember.ID == teamMemberID {
			c.OrganisationTeamMembers[teamID][len(c.OrganisationTeamMembers[teamID])-1], c.OrganisationTeamMembers[teamID][i] = c.OrganisationTeamMembers[teamID][i], c.OrganisationTeamMembers[teamID][len(c.OrganisationTeamMembers[teamID])-1]
			c.OrganisationTeamMembers[teamID] = c.OrganisationTeamMembers[teamID][:len(c.OrganisationTeamMembers[teamID])-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failure"}, fmt.Errorf("unable to find that team member")
}
