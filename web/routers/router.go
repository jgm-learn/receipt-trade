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
	beego.Router("/admin/trade", &controllers.AdminController{}, "get:Trade")

	beego.Router("/user", &controllers.UserController{})
	beego.Router("/user/getUserId", &controllers.UserController{}, "get:GetUserId")
	beego.Router("/user/getUserNonce", &controllers.UserController{}, "get:GetUserNonce")
	beego.Router("/user/getFunds", &controllers.UserController{}, "get:GetFunds")
	beego.Router("/user/getReceipt", &controllers.UserController{}, "get:GetReceipt")
	beego.Router("/user/listTrade", &controllers.UserController{}, "post:ListTrade")
	beego.Router("/user/cancellation", &controllers.UserController{}, "post:Cancellation")
	beego.Router("/user/getUserList", &controllers.UserController{}, "get:GetUserList")

	beego.Router("/market", &controllers.MarketController{})
	beego.Router("/market/getReceiptList", &controllers.MarketController{}, "get:GetReceiptList")
	beego.Router("/market/delist", &controllers.MarketController{}, "post:Delist")

	beego.Router("/allData", &controllers.DataController{})
	beego.Router("/allData/getUser", &controllers.DataController{}, "get:GetUser")
	beego.Router("/allData/getFunds", &controllers.DataController{}, "get:GetFunds")
	beego.Router("/allData/getReceipt", &controllers.DataController{}, "get:GetReceipt")
	beego.Router("/allData/getOrderSell", &controllers.DataController{}, "get:GetOrderSell")
}
