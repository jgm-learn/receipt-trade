package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"receipt-trade/web/models"
)

type DataController struct {
	beego.Controller
}

func (this *DataController) Get() {
	this.TplName = "allData.html"
}

func (this *DataController) GetUser() {
	var users []models.User

	_, err := o.QueryTable("user").All(&users)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &users
	this.ServeJSON()
}

func (this *DataController) GetFunds() {
	var userFunds []models.UserFunds

	_, err := o.QueryTable("user_funds").All(&userFunds)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &userFunds
	this.ServeJSON()
}

func (this *DataController) GetReceipt() {
	var userReceipts []models.UserReceipt

	_, err := o.QueryTable("user_receipt").All(&userReceipts)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &userReceipts
	this.ServeJSON()
}

func (this *DataController) GetOrderSell() {
	var orderSells []models.OrderSell

	_, err := o.QueryTable("order_sell").All(&orderSells)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &orderSells
	this.ServeJSON()
}
