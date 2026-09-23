package services

import (
	"encoding/json"
	"os"

	"Beego-Api-Project/models"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// in memory storing, stays in package level
var store []models.SourceProperty

// read data from config file
func LoadData() error {
	path, err := beego.AppConfig.String("datafile")
	if err != nil || path == "" {
		path = "data/rental_properties.json"
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		logs.Error("failed to read data file %s: %v", path, err)
		return err
	}

	//keeps data here after unmarshal
	var records []models.SourceProperty
	if err := json.Unmarshal(raw, &records); err != nil {
		logs.Error("failed to parse data file: %v", err)
		return err
	}

	store = records
	logs.Info("loaded %d properties from %s", len(store), path)
	return nil
}

// API response shaping
func Transform(s models.SourceProperty) models.PropertyResponse {
	breadcrumbs, err := s.Breadcrumbs()
	if err != nil {
		logs.Error("failed to parse categories for %s: %v", s.ID, err)
		breadcrumbs = []string{}
	}

	var lon, lat float64
	if len(s.LonLat.Coordinates) == 2 {
		lon = s.LonLat.Coordinates[0]
		lat = s.LonLat.Coordinates[1]
	}

	amenities := s.AmenityCategories
	if amenities == nil {
		amenities = []string{}
	}
	images := s.Images
	if images == nil {
		images = []string{}
	}

	return models.PropertyResponse{
		ID:   s.ID,
		Feed: s.Feed,
		GeoInfo: models.GeoInfo{
			Breadcrumbs: breadcrumbs,
			City:        s.City,
			Country:     s.Country,
			CountryCode: s.CountryCode,
			Name:        s.Display,
			LocationID:  s.LocationID,
			Lat:         lat,
			Lon:         lon,
			State:       s.State,
			StateAbbr:   s.StateAbbr,
		},
		Property: models.PropertyInfo{
			Amenities:    amenities,
			Name:         s.PropertyName,
			Slug:         s.PropertySlug,
			PropertyType: s.PropertyTypeCategory,
			Price:        s.USDPrice,
			ReviewScore:  s.ReviewScoreGeneral,
			StarRating:   s.StarRating,
			Counts: models.Counts{
				Bathroom:  s.BathroomCount,
				Bedroom:   s.BedroomCount,
				Reviews:   s.NumberOfReview,
				Occupancy: s.Occupancy,
			},
			Image: models.Image{
				Count:  len(images),
				Images: images,
			},
		},
		Published: s.Published,
	}
}

// transforms every record currently in the store.
// Always returns a non-nil slice.
func GetAll() []models.PropertyResponse {
	result := []models.PropertyResponse{}
	for _, s := range store {
		result = append(result, Transform(s))
	}
	return result
}

// ok will false if property doesnt exists
func GetByID(id string) (resp models.PropertyResponse, ok bool) {
	for _, s := range store {
		if s.ID == id {
			return Transform(s), true
		}
	}
	return models.PropertyResponse{}, false
}
