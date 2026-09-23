package main

import (
	_ "Beego-Api-Project/routers"
	"Beego-Api-Project/services"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if err := services.LoadData(); err != nil {
		logs.Error("startup: failed to load property data: %v", err)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
	}
	beego.Run()
}
