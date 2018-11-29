package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"receipt-trade/web/models"
)

//连接数据库
var o orm.Ormer //包的全局变量，定义Ormer接口变量
func init() {
	orm.RegisterDataBase("default", "mysql", "root:root@/receipt_trade?charset=utf8", 30)
	o = orm.NewOrm()
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "register.html"
}

func (c *MainController) Register() {
	c.TplName = "register.html"

	var user models.User

	body := c.Ctx.Input.RequestBody //获取json的二进制数据

	fmt.Println(string(body))
	//反序列化，并存入user
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(user)
	user.Insert() //写入数据库

}
