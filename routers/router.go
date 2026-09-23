package routers

import (
	"Beego-Api-Project/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSRouter("/", &controllers.PropertyController{}),
		),
	)

	beego.AddNamespace(ns)
}
