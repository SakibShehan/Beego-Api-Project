package models

import "encoding/json"

type SourceProperty struct {
	ID                   string   `json:"id"`
	Feed                 int      `json:"feed"`
	Country              string   `json:"country"`
	CountryCode          string   `json:"country_code"`
	State                string   `json:"state"`
	StateAbbr            string   `json:"state_abbr"`
	City                 string   `json:"city"`
	Display              string   `json:"display"`
	LocationID           string   `json:"location_id"`
	PropertyName         string   `json:"property_name"`
	PropertySlug         string   `json:"property_slug"`
	PropertyTypeCategory string   `json:"property_type_category"`
	USDPrice             float64  `json:"usd_price"`
	Occupancy            int      `json:"occupancy"`
	BedroomCount         int      `json:"bedroom_count"`
	BathroomCount        int      `json:"bathroom_count"`
	NumberOfReview       int      `json:"number_of_review"`
	ReviewScoreGeneral   float64  `json:"review_score_general"`
	StarRating           int      `json:"star_rating"`
	AmenityCategories    []string `json:"amenity_categories"`
	LonLat               LonLat   `json:"lonlat"`
	Categories           string   `json:"categories"`
	Published            bool     `json:"published"`
	Images               []string `json:"images"`
}

type LonLat struct {
	Coordinates []float64 `json:"coordinates"`
}

// convert JSON Categories string into a actual string.
func (s SourceProperty) Breadcrumbs() ([]string, error) {
	if s.Categories == "" {
		return []string{}, nil
	}
	var out []string
	if err := json.Unmarshal([]byte(s.Categories), &out); err != nil {
		return []string{}, err
	}
	return out, nil
}
