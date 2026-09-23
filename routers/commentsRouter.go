package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["Beego-Api-Project/controllers:PropertyController"] = append(beego.GlobalControllerRouter["Beego-Api-Project/controllers:PropertyController"],
        beego.ControllerComments{
            Method: "GetList",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["Beego-Api-Project/controllers:PropertyController"] = append(beego.GlobalControllerRouter["Beego-Api-Project/controllers:PropertyController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
