package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"receipt-trade/web/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {
	this.TplName = "user.html"
}
func (this *UserController) GetFunds() {
	//orm.RegisterDataBase("default", "mysql", "root:root@/receipt_trade?charset=utf8", 30)
	//o := orm.NewOrm()

	user := new(models.User)
	pK := this.GetString("pucKey")
	user.PublicKey = pK

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	fmt.Println(user)

	userFunds := new(models.UserFunds)
	userFunds.UserId = user.UserId
	o.Read(userFunds) //根据用户id查询资金信息
	fmt.Printf("%d\n", userFunds.TotalFunds)

	this.Data["json"] = &userFunds
	this.ServeJSON()
}

func (this *UserController) GetReceipt() {
	//orm.RegisterDataBase("default", "mysql", "root:root@/receipt_trade?charset=utf8", 30) //连接数据库
	//o := orm.NewOrm()

	pK := this.GetString("pucKey")
	user := new(models.User)
	user.PublicKey = pK

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	fmt.Printf("%s receipt:\n", user.UserName)

	var userReceipts []models.UserReceipt
	_, err = o.QueryTable("UserReceipt").Filter("UserId", user.UserId).All(&userReceipts)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &userReceipts
	this.ServeJSON()
}
