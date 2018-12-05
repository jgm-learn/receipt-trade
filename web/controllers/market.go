package controllers

import (
	"encoding/json"
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
	var lists []models.List
	_, err := o.QueryTable("list").All(&lists)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	this.Data["json"] = &lists
	this.ServeJSON()

}

func (this *MarketController) Delist() {
	var webReply models.WebReply
	var orderBuy models.OrderBuy
	//var user models.User
	var list models.List
	var userFunds models.UserFunds
	var userReceipt models.UserReceipt
	var qtyBuy int

	fmt.Printf("<---------------------------->\n")
	fmt.Printf("market.go Delist() 摘牌执行输出如下：\n")
	body := this.Ctx.Input.RequestBody // body: UsrId,ListId,QtyBuy,NonceBuy,SigBuy,AddrBuy
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &orderBuy); err != nil {
		fmt.Println("json Unmarshal 出错：")
		fmt.Println(err)
		webReply.Reply = "摘牌失败：json Unmarshal 出错 "
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	qtyBuy = orderBuy.QtyBuy

	//事务处理
	oTX := orm.NewOrm()
	errs := oTX.Begin()
	if errs != nil {
		fmt.Println("o.Begin ", errs)
		webReply.Reply = "摘牌失败：o.Begin()出错 "
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}

	err := oTX.Raw("select * from list where list_id = ? for update", orderBuy.ListId).QueryRow(&list)
	if err != nil {
		fmt.Println("阻塞中，正在排队", err)
		webReply.Reply = "摘牌失败：系统太过繁忙，请稍后重试"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	/*
		//查询挂单的具体信息
		err := o.QueryTable("list").Filter("list_id", orderBuy.ListId).One(&list)
		if err == orm.ErrNoRows {
			fmt.Println("查询挂单编号为%d的挂单", orderBuy.ListId)
			webReply.Reply = "摘牌失败：该挂单不存在"
			this.Data["json"] = &webReply
			this.ServeJSON()
			return
		}
	*/

	//判断可售仓单数量是否足够
	if list.QtyRemain < orderBuy.QtyBuy {
		webReply.Reply = "摘牌失败：可售仓单数量不足"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	//判断买方是否有足够资金
	payment := orderBuy.QtyBuy * list.Price //应付款
	oTX.QueryTable("user_funds").Filter("user_id", orderBuy.UserId).One(&userFunds)
	if userFunds.AvailableFunds < payment {
		webReply.Reply = "摘牌失败：资金不足"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}

	//更新list表
	list.QtyDeal += qtyBuy
	list.QtyRemain -= qtyBuy
	_, err = oTX.QueryTable("list").Filter("list_id", orderBuy.ListId).Update(orm.Params{
		"qty_deal":   list.QtyDeal,
		"qty_remain": list.QtyRemain,
	})
	if err != nil {
		fmt.Printf("更新user表 qty_deal,qty_remain字段失败 %v\n", err)
		return
	}

	//更新卖方user_receipt表
	/*
		err = o.QueryTable("user_receipt").Filter("user_id", list.UserId).Filter("receipt_id", list.ReceiptId).One(&userReceipt)
		userReceipt.QtyTotal -= qtyBuy
		userReceipt.QtyFrozen -= qtyBuy
	*/

	_, err = oTX.QueryTable("user_receipt").Filter("user_id", list.UserId).Filter("receipt_id", list.ReceiptId).Update(orm.Params{
		//"qty_total":  userReceipt.QtyTotal,
		//"qty_frozen": userReceipt.QtyFrozen,
		"qty_total":  orm.ColValue(orm.ColMinus, qtyBuy),
		"qty_frozen": orm.ColValue(orm.ColMinus, qtyBuy),
	})
	if err != nil {
		fmt.Printf("更新user_receipt表 qty_total,qty_frozen字段失败 %v\n", err)
		return
	}

	//更新买方user_receipt表
	err = oTX.QueryTable("user_receipt").Filter("user_id", orderBuy.UserId).Filter("receipt_id", orderBuy.ReceiptId).One(&userReceipt)
	if err == orm.ErrNoRows {
		fmt.Println("该用户没有该仓单")
		userReceipt.UserId = orderBuy.UserId
		userReceipt.ReceiptId = orderBuy.ReceiptId
		userReceipt.QtyTotal = qtyBuy
		userReceipt.QtyAvailable = qtyBuy
		userReceipt.Insert()
	} else {
		_, err = oTX.QueryTable("user_receipt").Filter("user_id", orderBuy.UserId).Filter("receipt_id", orderBuy.ReceiptId).Update(orm.Params{
			"qty_total":     orm.ColValue(orm.ColAdd, qtyBuy),
			"qty_available": orm.ColValue(orm.ColAdd, qtyBuy),
		})
	}

	//更新卖方user_funds表
	_, err = oTX.QueryTable("user_funds").Filter("user_id", list.UserId).Update(orm.Params{
		"total_funds":     orm.ColValue(orm.ColAdd, payment),
		"available_funds": orm.ColValue(orm.ColAdd, payment),
	})
	if err != nil {
		fmt.Printf("更新user_funds表 total_funds,available_funds字段失败 %v\n", err)
		return
	}

	//更新买方user_funds表
	_, err = oTX.QueryTable("user_funds").Filter("user_id", orderBuy.UserId).Update(orm.Params{
		"total_funds":     orm.ColValue(orm.ColMinus, payment),
		"available_funds": orm.ColValue(orm.ColMinus, payment),
	})
	if err != nil {
		fmt.Printf("更新user_funds表 total_funds,available_funds字段失败 %v\n", err)
		return
	}

	//提交事务
	oTX.Commit()

	webReply.Reply = "摘牌成功"
	this.Data["json"] = &webReply
	this.ServeJSON()
	return
}
