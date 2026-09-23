package routers

import (
	"Beego-Api-Project/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/properties",
			beego.NSRouter("/", &controllers.PropertyController{}, "get:GetList"),
			beego.NSRouter("/:id", &controllers.PropertyController{}, "get:GetOne"),
		),
	)

	beego.AddNamespace(ns)
}
