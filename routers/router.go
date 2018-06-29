// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"AddressServer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/AddressServer",
			beego.NSRouter("/UploadIp", &controllers.AddressServerController{}, "post:UploadIp"),
			beego.NSRouter("/SearchLastestLogin", &controllers.AddressServerController{}, "post:SearchLastestLogin"),
			beego.NSRouter("/SearchLastestRegister", &controllers.AddressServerController{}, "post:SearchLastestRegister"),
			beego.NSRouter("/UploadGps", &controllers.AddressServerController{}, "post:UploadGps"),
			beego.NSRouter("/SearchGps", &controllers.AddressServerController{}, "post:SearchGps"),
		),
	)
	beego.AddNamespace(ns)
}
