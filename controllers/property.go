package controllers

import (
	"Beego-Api-Project/models"
	"Beego-Api-Project/services"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	beego.Controller
}

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

	return params, nil
}
