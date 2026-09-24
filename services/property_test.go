package services

import (
	"reflect"
	"testing"

	"Beego-Api-Project/models"
)

/////////////////////// Transfrom Test //////////////////////////////////////////////////

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

/////////////////////////// AND filtering test ////////////////////////////

// define data from SourceProperty
func sampleStore() []models.SourceProperty {
	return []models.SourceProperty{
		{
			ID: "P1", Feed: 11, Published: true, USDPrice: 100,
			PropertyTypeCategory: "Hotel", StarRating: 4, ReviewScoreGeneral: 8.0,
			NumberOfReview: 20, BedroomCount: 2,
			AmenityCategories: []string{"Internet", "Pool"},
		},
		{
			ID: "P2", Feed: 12, Published: false, USDPrice: 200,
			PropertyTypeCategory: "Villa", StarRating: 5, ReviewScoreGeneral: 9.0,
			NumberOfReview: 5, BedroomCount: 4,
			AmenityCategories: []string{"Parking"},
		},
		{
			ID: "P3", Feed: 11, Published: false, USDPrice: 50,
			PropertyTypeCategory: "Hostel", StarRating: 2, ReviewScoreGeneral: 6.0,
			NumberOfReview: 2, BedroomCount: 1,
			AmenityCategories: []string{"Internet"},
		},
	}
}

// helpers to get a pointer
func ptrFloat(v float64) *float64 { return &v }
func ptrInt(v int) *int           { return &v }
func ptrBool(v bool) *bool        { return &v }

func TestMatches_ANDFilters(t *testing.T) {
	data := sampleStore()

	tests := []struct {
		name    string
		params  models.FilterParams
		wantIDs []string
	}{
		{
			name:    "feed AND published",
			params:  models.FilterParams{Feed: ptrInt(11), Published: ptrBool(false)},
			wantIDs: []string{"P3"},
		},
		{
			name:    "price range",
			params:  models.FilterParams{MinPrice: ptrFloat(60), MaxPrice: ptrFloat(150)},
			wantIDs: []string{"P1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, s := range data {
				if matches(s, tt.params) {
					got = append(got, s.ID)
				}
			}
			if !reflect.DeepEqual(got, tt.wantIDs) {
				t.Errorf("matches() AND filter mismatch\ngot:  %v\nwant: %v", got, tt.wantIDs)
			}
		})
	}
}

/////////////////// Amenities OR check /////////////////////////////////////

func TestMatches_AmenitiesOR(t *testing.T) {
	data := sampleStore()
	params := models.FilterParams{Amenities: []string{"Pool", "Parking"}}

	var got []string
	for _, s := range data {
		if matches(s, params) {
			got = append(got, s.ID)
		}
	}

	want := []string{"P1", "P2"} // p1 and p2 will mathc
	if !reflect.DeepEqual(got, want) {
		t.Errorf("matches() amenities OR mismatch\ngot:  %v\nwant: %v", got, want)
	}
}

///////////////// AND OR combined Test ///////////////////////////////

func TestCombinedAndOr(t *testing.T) {
	data := sampleStore()
	// feed=11 AND (has Internet OR has Pool)
	params := models.FilterParams{
		Feed:      ptrInt(11),
		Amenities: []string{"Internet", "Pool"},
	}

	var got []string
	for _, s := range data {
		if matches(s, params) {
			got = append(got, s.ID)
		}
	}

	want := []string{"P1", "P3"} // both are feed 11, both have Internet
	if !reflect.DeepEqual(got, want) {
		t.Errorf("matches() combined AND+OR mismatch\ngot:  %v\nwant: %v", got, want)
	}
}
