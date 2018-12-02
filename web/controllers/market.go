package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"receipt-trade/web/models"
)

type MarketController struct {
	beego.Controller
}

func (this *MarketController) Get() {
	this.TplName = "market.html"
}

func (this *MarketController) GetReceiptList() {
	var markets []models.Market
	_, err := o.QueryTable("market").All(&markets)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &markets
	this.ServeJSON()

}
