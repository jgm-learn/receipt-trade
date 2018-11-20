package routers

import (
	"github.com/astaxie/beego"
	"receipt-trade/web/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.MainController{}, "post:Register")

	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/addReceipt", &controllers.AdminController{}, "post:AddReceipt")
	beego.Router("/admin/addUserReceipt", &controllers.AdminController{}, "post:AddUserReceipt")
	beego.Router("/admin/addUserFunds", &controllers.AdminController{}, "post:AddUserFunds")

	beego.Router("/user", &controllers.UserController{})
	beego.Router("/user/getFunds", &controllers.UserController{}, "get:GetFunds")
	beego.Router("/user/getReceipt", &controllers.UserController{}, "get:GetReceipt")
}
