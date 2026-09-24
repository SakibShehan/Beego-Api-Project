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

type CategoryEntry struct {
	LocationID string   `json:"LocationID"`
	Name       string   `json:"Name"`
	Type       string   `json:"Type"`
	Slug       string   `json:"Slug"`
	Display    []string `json:"Display"`
}

func (s SourceProperty) Breadcrumbs() ([]string, error) {
	if s.Categories == "" {
		return []string{}, nil
	}

	var entries []CategoryEntry
	if err := json.Unmarshal([]byte(s.Categories), &entries); err != nil {
		return []string{}, err
	}

	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, e.Name)
	}
	return out, nil
}
