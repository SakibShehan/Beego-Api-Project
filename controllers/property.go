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
