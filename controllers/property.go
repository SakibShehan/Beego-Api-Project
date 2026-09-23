package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	beego.Controller
}

func (c *PropertyController) Get() {
	c.Ctx.WriteString("Hello Worlds")
}
