package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Region represents a geographical/DC region for Civo resources
type Region struct {
	Code          string  `json:"code"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	OutOfCapacity bool    `json:"out_of_capacity"`
	Country       string  `json:"country"`
	CountryName   string  `json:"country_name"`
	Features      Feature `json:"features"`
}

// Feature represent a all feature inside a region
type Feature struct {
	Iaas       bool `json:"iaas"`
	Kubernetes bool `json:"kubernetes"`
}

// ListRegions returns all load balancers owned by the calling API account
func (c *Client) ListRegions() ([]Region, error) {
	resp, err := c.SendGetRequest("/v2/regions")
	if err != nil {
		return nil, decodeERROR(err)
	}

	regions := make([]Region, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&regions); err != nil {
		return nil, err
	}

	return regions, nil
}

// FindRegion is a function to find a region
func (c *Client) FindRegion(search string) (*Region, error) {
	allregion, err := c.ListRegions()
	if err != nil {
		return nil, decodeERROR(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := Region{}

	for _, value := range allregion {
		if value.Name == search || value.Code == search || value.Country == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.Name, search) || strings.Contains(value.Code, search) || strings.Contains(value.Country, search) {
			if exactMatch == false {
				result = value
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}
