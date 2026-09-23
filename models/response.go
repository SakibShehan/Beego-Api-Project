package models

// main structure
type PropertyResponse struct {
	ID        string       `json:"ID"`
	Feed      int          `json:"Feed"`
	GeoInfo   GeoInfo      `json:"GeoInfo"`   //its a strcuture
	Property  PropertyInfo `json:"Property"`  //its a strcuture
	Published bool         `json:"Published"` //its a strcuture
}

//geo info structure containing its information

type GeoInfo struct {
	Breadcrumbs []string `json:"Breadcrumbs"`
	City        string   `json:"City"`
	Country     string   `json:"Country"`
	CountryCode string   `json:"CountryCode"`
	Name        string   `json:"Name"`
	LocationID  string   `json:"LocationID"`
	Lat         float64  `json:"Lat"`
	Lon         float64  `json:"Lon"`
	State       string   `json:"State"`
	StateAbbr   string   `json:"StateAbbr"`
}

//property info structure containing its information
type PropertyInfo struct {
	Amenities    []string `json:"Amenities"`
	Name         string   `json:"Name"`
	Slug         string   `json:"Slug"`
	PropertyType string   `json:"PropertyType"`
	Price        float64  `json:"Price"`
	ReviewScore  float64  `json:"ReviewScore"`
	StarRating   int      `json:"StarRating"`
	Counts       Counts   `json:"Counts"`
	Image        Image    `json:"Image"`
}

// this structure will show different types of counts
type Counts struct {
	Bathroom  int `json:"Bathroom"`
	Bedroom   int `json:"Bedroom"`
	Reviews   int `json:"Reviews"`
	Occupancy int `json:"Occupancy"`
}

type Image struct {
	Count  int      `json:"Count"`
	Images []string `json:"Images"`
}

type ListResult struct {
	Count int                `json:"Count"`
	Items []PropertyResponse `json:"Items"`
}

type ListResponse struct {
	Result ListResult `json:"Result"`
}

// Sam shared error shape for all endpoints.
type ErrorResponse struct {
	Error string `json:"Error"`
}
