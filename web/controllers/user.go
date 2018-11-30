package controllers

import (
	"encoding/json"
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
func (this *UserController) GetUserId() {
	user := new(models.User)
	userAddr := this.GetString("userAddr")
	user.PublicKey = userAddr

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	this.Data["json"] = &user
	this.ServeJSON()
}

func (this *UserController) GetUserNonce() {
	user := new(models.User)
	userAddr := this.GetString("userAddr")
	user.PublicKey = userAddr

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	this.Data["json"] = &user
	this.ServeJSON()
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

func (this *UserController) ListTrade() {
	var webReply models.WebReply

	fmt.Printf("<---------------------------->\n")
	fmt.Printf("user.go ListTrade() 挂牌执行输出如下：\n")
	//接收前端传来的订单数据
	var orderSell models.OrderSell     // OrderSell: Id, UserId, ReceiptId, Price, QtySell, NonceSell, SigSell, AddrSell
	var market models.Market           // Market:	MarketId, ReceiptId, Price, QtySell, QtyDeal, QtyRemain
	body := this.Ctx.Input.RequestBody //body: UserId, ReceiptId, Price, QtySell, NonceSell, SigSell, AddrSell
	fmt.Printf("前端传来数据如下：\n")
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &orderSell); err != nil {
		fmt.Println("json Unmarshal 出错：")
		fmt.Println(err)
		webReply.Reply = "挂牌失败：json Unmarshal 出错 "
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	market.ReceiptId = orderSell.ReceiptId
	market.Price = orderSell.Price
	market.QtySell = orderSell.QtySell
	market.QtyRemain = orderSell.QtySell

	//查询用户的可用仓单数量，判断用户仓单数量是否足够
	var userReceipt models.UserReceipt
	err := o.QueryTable("UserReceipt").Filter("UserId", orderSell.UserId).Filter("ReceiptId", orderSell.ReceiptId).One(&userReceipt)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
		webReply.Reply = "挂牌失败：该用户没有该仓单"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
		return
	}

	fmt.Println("userReceipt:")
	fmt.Println(userReceipt)
	if orderSell.QtySell > userReceipt.QtyAvailable {
		fmt.Println("用户的可用仓单数量不足，无法挂牌")
		webReply.Reply = "挂牌失败：仓单数量不足"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}

	//冻结仓单，更新user_receipt表
	userReceipt.QtyAvailable = userReceipt.QtyAvailable - orderSell.QtySell
	userReceipt.QtyFrozen = orderSell.QtySell
	num, err := o.QueryTable("user_receipt").Filter("id", userReceipt.Id).Update(orm.Params{
		"qty_available": userReceipt.QtyAvailable,
		"qty_frozen":    userReceipt.QtyFrozen,
	})
	if err == nil {
		fmt.Printf("更新user_receipt表 qty_available,qty_frozen 字段成功 id = %d\n", num)
	} else {
		fmt.Printf("更新user_receipt表 qty_available,qty_frozen 字字段失败 %v", err)
	}

	//将订单数据写入order_sell market表
	orderSell.Insert()
	market.Insert()

	//更新user表的nonce字段
	num, err = o.QueryTable("user").Filter("user_id", orderSell.UserId).Update(orm.Params{
		"nonce": orderSell.NonceSell})
	if err == nil {
		fmt.Printf("更新user表 nonce字段成功 id = %d\n", num)
	} else {
		fmt.Printf("更新user表 nonce字段失败 %v", err)
	}

	webReply.Reply = "挂牌成功"
	this.Data["json"] = &webReply
	this.ServeJSON()
	fmt.Printf("<---------------------------->\n")
}
