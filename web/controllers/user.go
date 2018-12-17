package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"receipt-trade/web/models"
	"strconv"
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

type OrderSellAndNonce struct {
	OrderSell models.OrderSell
	Nonce     int
}

//获取nonce,价格，卖方账户地址
func (this *UserController) GetUserNonce() {
	var orderSell models.OrderSell
	var data OrderSellAndNonce
	user := new(models.User)
	userAddr := this.GetString("userAddr")
	user.PublicKey = userAddr
	listId := this.GetString("listId")

	err := o.Read(user, "PublicKey") //根据公钥查询user
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}

	err = o.QueryTable("order_sell").Filter("id", listId).One(&orderSell)
	if err == orm.ErrNoRows {
		fmt.Printf("没有找到挂单号为%d的挂单", listId)
	}

	data.OrderSell = orderSell
	data.Nonce = user.Nonce

	fmt.Println("orderSell: ", orderSell)
	fmt.Println("data: ", data)
	this.Data["json"] = &data
	this.ServeJSON()
}

func (this *UserController) GetFunds() {
	userId := this.GetString("userId")

	userFunds := new(models.UserFunds)
	userFunds.UserId, _ = strconv.Atoi(userId)
	o.Read(userFunds) //根据用户id查询资金信息

	this.Data["json"] = &userFunds
	this.ServeJSON()
}

func (this *UserController) GetReceipt() {
	userId := this.GetString("userId")
	var userReceipts []models.UserReceipt
	_, err := o.QueryTable("user_receipt").Filter("user_id", userId).All(&userReceipts)
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
	var orderSell models.OrderSell // OrderSell: Id, UserId, ReceiptId, Price, QtySell, NonceSell, SigSell, AddrSell
	var list models.List           // List:	ListId, UserId, ReceiptId, Price, QtySell, QtyDeal, QtyRemain

	fmt.Printf("<---------------------------->\n")
	fmt.Printf("user.go ListTrade() 挂牌执行输出如下：\n")
	//接收前端传来的订单数据
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
	list.UserId = orderSell.UserId
	list.ReceiptId = orderSell.ReceiptId
	list.Price = orderSell.Price
	list.QtySell = orderSell.QtySell
	list.QtyRemain = orderSell.QtySell

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
	userReceipt.QtyAvailable -= orderSell.QtySell
	userReceipt.QtyFrozen += orderSell.QtySell
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
	listId := orderSell.Insert()
	list.ListId = int(listId)
	list.Insert()

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

//撤单
func (this *UserController) Cancellation() {
	var webReply models.WebReply
	var list models.List

	userId := this.GetString("userId")
	listId := this.GetString("listId")

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

	//查询挂单的具体信息,添加排他锁
	err := oTX.Raw("select * from list where list_id = ? and user_id = ? for update", listId, userId).QueryRow(&list)
	if err == orm.ErrNoRows {
		fmt.Printf("Cancellation() 没有找到用户%d的%d号挂单\n", userId, listId)
		webReply.Reply = "撤单失败：没有找到该挂单"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	if err != nil {
		fmt.Println("阻塞中，正在排队", err)
		oTX.Rollback()
		webReply.Reply = "撤单失败：系统太过繁忙，请稍后重试"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}

	//解冻冻结的仓单
	_, err = oTX.QueryTable("user_receipt").Filter("user_id", list.UserId).Filter("receipt_id", list.ReceiptId).Update(orm.Params{
		"qty_available": orm.ColValue(orm.ColAdd, list.QtyRemain),
		"qty_frozen":    orm.ColValue(orm.ColMinus, list.QtyRemain),
	})
	if err != nil {
		fmt.Printf("更新user_receipt表 qty_available,qty_frozen字段失败 %v\n", err)
		oTX.Rollback()
		return
	}

	//提交事务
	oTX.Commit()

	//删除市场中的挂单
	_, err = o.QueryTable("list").Filter("list_id", list.ListId).Delete()
	if err != nil {
		fmt.Printf("删除list表失败 err: %v \n listId = %d \n", err, list.ListId)
		oTX.Rollback()
		webReply.Reply = "撤单失败"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}

	webReply.Reply = "撤单成功"
	this.Data["json"] = &webReply
	this.ServeJSON()
	return
}

func (this *UserController) GetUserList() {
	userId := this.GetString("userId")

	fmt.Println(userId)
	var lists []models.List
	_, err := o.QueryTable("list").Filter("user_id", userId).All(&lists)
	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
	}
	fmt.Println(lists)

	this.Data["json"] = &lists
	this.ServeJSON()
}
