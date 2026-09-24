package services

import (
	"reflect"
	"testing"

	"Beego-Api-Project/models"
)

// each entry are mapped from source property
func TestTransform(t *testing.T) {
	tests := []struct {
		name  string
		input models.SourceProperty
		want  models.PropertyResponse
	}{
		{
			name: "full mapping with categories, lonlat, amenities, images",
			input: models.SourceProperty{
				ID:                   "BC-1000001",
				Feed:                 11,
				Country:              "Japan",
				CountryCode:          "JP",
				State:                "Tokyo",
				StateAbbr:            "TK",
				City:                 "Shibuya",
				Display:              "Shibuya, Tokyo, Japan",
				LocationID:           "89",
				PropertyName:         "Sakura Inn",
				PropertySlug:         "sakura-inn",
				PropertyTypeCategory: "Hotel",
				USDPrice:             120.5,
				Occupancy:            4,
				BedroomCount:         2,
				BathroomCount:        1,
				NumberOfReview:       15,
				ReviewScoreGeneral:   8.5,
				StarRating:           4,
				AmenityCategories:    []string{"Internet", "Parking"},
				LonLat: models.LonLat{
					// spec: coordinates = [lon, lat]
					Coordinates: []float64{139.7, 35.6},
				},
				Categories: `[{"LocationID":"1","Name":"Japan","Type":"country","Slug":"japan","Display":["japan"]},{"LocationID":"2","Name":"Tokyo","Type":"state","Slug":"tokyo","Display":["tokyo"]}]`,
				Published:  true,
				Images:     []string{"img1.jpg", "img2.jpg"},
			},
			want: models.PropertyResponse{
				ID:   "BC-1000001",
				Feed: 11,
				GeoInfo: models.GeoInfo{
					Breadcrumbs: []string{"Japan", "Tokyo"},
					City:        "Shibuya",
					Country:     "Japan",
					CountryCode: "JP",
					Name:        "Shibuya, Tokyo, Japan",
					LocationID:  "89",
					Lat:         35.6,
					Lon:         139.7,
					State:       "Tokyo",
					StateAbbr:   "TK",
				},
				Property: models.PropertyInfo{
					Amenities:    []string{"Internet", "Parking"},
					Name:         "Sakura Inn",
					Slug:         "sakura-inn",
					PropertyType: "Hotel",
					Price:        120.5,
					ReviewScore:  8.5,
					StarRating:   4,
					Counts: models.Counts{
						Bathroom:  1,
						Bedroom:   2,
						Reviews:   15,
						Occupancy: 4,
					},
					Image: models.Image{
						Count:  2,
						Images: []string{"img1.jpg", "img2.jpg"},
					},
				},
				Published: true,
			},
		},
		{
			name: "empty categories, nil amenities and images produce empty slices not nil",
			input: models.SourceProperty{
				ID:         "BC-1000002",
				Feed:       12,
				Categories: "",
				LonLat:     models.LonLat{Coordinates: []float64{}},
				Published:  false,
				// AmenityCategories and Images left nil deliberately
			},
			want: models.PropertyResponse{
				ID:   "BC-1000002",
				Feed: 12,
				GeoInfo: models.GeoInfo{
					Breadcrumbs: []string{},
					Lat:         0,
					Lon:         0,
				},
				Property: models.PropertyInfo{
					Amenities: []string{},
					Counts:    models.Counts{},
					Image: models.Image{
						Count:  0,
						Images: []string{},
					},
				},
				Published: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Transform(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() mismatch\ngot:  %+v\nwant: %+v", got, tt.want)
			}
		})
	}
}
