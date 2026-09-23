package controllers

import (
	"Beego-Api-Project/models"
	"Beego-Api-Project/services"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	beego.Controller
}

// responses in call of  GET /v1/properties
func (c *PropertyController) GetList() {
	items := services.GetAll()

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
