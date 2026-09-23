package controllers

import (
	"Beego-Api-Project/models"
	"Beego-Api-Project/services"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// property validation
var validPropertyTypes = map[string]bool{
	"Hotel":     true,
	"House":     true,
	"Apartment": true,
	"Villa":     true,
	"Resort":    true,
	"Hostel":    true,
}

// feeds validation
var validFeeds = map[int]bool{
	11: true,
	12: true,
	22: true,
	24: true,
}

type PropertyController struct {
	beego.Controller
}

// GetList godoc
// @Title Get property list
// @Description List properties with optional AND filters, amenities OR filter, and limit
// @Param	min_price		query	float64	false	"Minimum price (usd_price >= min_price)"
// @Param	max_price		query	float64	false	"Maximum price (usd_price <= max_price)"
// @Param	min_star_rating		query	int	false	"Minimum star rating"
// @Param	min_review_score	query	float64	false	"Minimum review score"
// @Param	min_reviews		query	int	false	"Minimum number of reviews"
// @Param	published		query	bool	false	"Published status (true or false)"
// @Param	property_type		query	string	false	"Exact match: Hotel, House, Apartment, Villa, Resort, Hostel"
// @Param	feed			query	int	false	"Exact match: 11, 12, 22, or 24"
// @Param	min_bedroom		query	int	false	"Minimum bedroom count"
// @Param	amenities		query	string	false	"Comma-separated list, OR match, e.g. Internet,Parking"
// @Param	limit			query	int	false	"Cap number of results returned"
// @Success 200 {object} models.ListResponse
// @Failure 400 {object} models.ErrorResponse "invalid query parameter"
// @router / [get]
// responses in call of  GET /v1/properties
func (c *PropertyController) GetList() {
	params, errResp := c.parseFilterParams()
	if errResp != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = *errResp
		c.ServeJSON()
		return
	}

	items := services.GetFiltered(params)

	c.Data["json"] = models.ListResponse{
		Result: models.ListResult{
			Count: len(items),
			Items: items,
		},
	}
	c.ServeJSON()
}

// GetOne godoc
// @Title Get property by ID
// @Description Return a single transformed property object
// @Param	id		path	string	true	"Property ID, e.g. BC-1000001"
// @Success 200 {object} models.PropertyResponse
// @Failure 404 {object} models.ErrorResponse "property not found"
// @router /:id [get]
// returns by id
func (c *PropertyController) GetOne() {
	id := c.Ctx.Input.Param(":id")

	prop, ok := services.GetByID(id)
	if !ok {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = models.ErrorResponse{Error: "Property not found"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = prop
	c.ServeJSON()
}

//  reads and validates query params

func (c *PropertyController) parseFilterParams() (models.FilterParams, *models.ErrorResponse) {
	var params models.FilterParams

	if raw := c.GetString("min_price"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid min_price: must be a number"}
		}
		params.MinPrice = &v
	}

	if raw := c.GetString("max_price"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid max_price: must be a number"}
		}
		params.MaxPrice = &v
	}

	if raw := c.GetString("min_star_rating"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid min_star_rating: must be an integer"}
		}
		params.MinStarRating = &v
	}

	if raw := c.GetString("min_review_score"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid min_review_score: must be a number"}
		}
		params.MinReviewScore = &v
	}

	if raw := c.GetString("min_reviews"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid min_reviews: must be an integer"}
		}
		params.MinReviews = &v
	}

	if raw := c.GetString("published"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid published: must be true or false"}
		}
		params.Published = &v
	}

	if raw := c.GetString("property_type"); raw != "" {
		if !validPropertyTypes[raw] {
			return params, &models.ErrorResponse{Error: "invalid property_type: must be one of Hotel, House, Apartment, Villa, Resort, Hostel"}
		}
		params.PropertyType = &raw
	}

	if raw := c.GetString("feed"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid feed: must be an integer"}
		}
		if !validFeeds[v] {
			return params, &models.ErrorResponse{Error: "invalid feed: must be one of 11, 12, 22, 24"}
		}
		params.Feed = &v
	}

	if raw := c.GetString("min_bedroom"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid min_bedroom: must be an integer"}
		}
		params.MinBedroom = &v
	}

	if raw := c.GetString("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return params, &models.ErrorResponse{Error: "invalid limit: must be an integer"}
		}
		if v < 0 {
			return params, &models.ErrorResponse{Error: "invalid limit: must be zero or a positive integer"}
		}
		params.Limit = &v
	}

	//////// OR filtering ///////////////////

	if raw := c.GetString("amenities"); raw != "" {
		parts := strings.Split(raw, ",")
		amenities := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				return params, &models.ErrorResponse{Error: "invalid amenities: contains an empty value"}
			}
			amenities = append(amenities, p)
		}
		params.Amenities = amenities
	}

	return params, nil
}
