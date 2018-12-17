package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "receipt-trade/web/controllers/grpcPB"
	"receipt-trade/web/models"
	"time"
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

	//查询挂单的具体信息,添加排他锁
	err := oTX.Raw("select * from list where list_id = ? for update", orderBuy.ListId).QueryRow(&list)
	if err != nil {
		fmt.Println("阻塞中，正在排队", err)
		oTX.Rollback()
		webReply.Reply = "摘牌失败：系统太过繁忙，请稍后重试"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}

	//判断是否买卖双方为同一用户
	if list.UserId == orderBuy.UserId {
		fmt.Println("摘牌失败：买卖双方为同一用户")
		oTX.Rollback()
		webReply.Reply = "摘牌失败：不能摘自己的挂单"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	//判断可售仓单数量是否足够
	if list.QtyRemain < orderBuy.QtyBuy {
		oTX.Rollback()
		webReply.Reply = "摘牌失败：可售仓单数量不足"
		this.Data["json"] = &webReply
		this.ServeJSON()
		return
	}
	//判断买方是否有足够资金
	payment := orderBuy.QtyBuy * list.Price //应付款
	//ToDo fee := float64(payment) * feeRate       //手续费
	oTX.QueryTable("user_funds").Filter("user_id", orderBuy.UserId).One(&userFunds)
	if userFunds.AvailableFunds < payment {
		oTX.Rollback()
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
		oTX.Rollback()
		return
	}

	//更新卖方user_receipt表
	_, err = oTX.QueryTable("user_receipt").Filter("user_id", list.UserId).Filter("receipt_id", list.ReceiptId).Update(orm.Params{
		"qty_total":  orm.ColValue(orm.ColMinus, qtyBuy),
		"qty_frozen": orm.ColValue(orm.ColMinus, qtyBuy),
	})
	if err != nil {
		fmt.Printf("更新user_receipt表 qty_total,qty_frozen字段失败 %v\n", err)
		oTX.Rollback()
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
		oTX.Rollback()
		return
	}

	//更新买方user_funds表
	_, err = oTX.QueryTable("user_funds").Filter("user_id", orderBuy.UserId).Update(orm.Params{
		"total_funds":     orm.ColValue(orm.ColMinus, payment),
		"available_funds": orm.ColValue(orm.ColMinus, payment),
	})
	if err != nil {
		fmt.Printf("更新user_funds表 total_funds,available_funds字段失败 %v\n", err)
		oTX.Rollback()
		return
	}

	//提交事务
	oTX.Commit()

	//调用智能合约
	Trade(orderBuy, list.ListId)

	webReply.Reply = "摘牌成功"
	this.Data["json"] = &webReply
	this.ServeJSON()
	return
}

//调用智能合约处理摘牌交易
func Trade(orderBuy models.OrderBuy, listId int) {
	var orderSell models.OrderSell

	o.QueryTable("order_sell").Filter("id", listId).One(&orderSell)
	if orderSell.Id == 0 {
		fmt.Printf("查询数据库失败，%d 订单不存在\n", 1)
		return
	}

	//构建消息对象
	var req pb.TradeReq
	req.ReceiptId = int64(orderSell.ReceiptId)
	req.Price = int64(orderSell.Price)
	req.QtySell = int64(orderSell.QtySell)
	req.NonceSell = int64(orderSell.NonceSell)
	req.QtyBuy = int64(orderBuy.QtyBuy)
	req.NonceBuy = int64(orderBuy.NonceBuy)
	req.AddrSell = orderSell.AddrSell
	req.AddrBuy = orderBuy.AddrBuy
	req.SigSell = orderSell.SigSell
	req.SigBuy = orderBuy.SigBuy

	//调用智能合约-建立连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("admin.go grpc客户端连接失败！ err: %v\n", err)
	}

	defer conn.Close()

	client := pb.NewRPCServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rst, err := client.Trade(ctx, &req) //调用智能合约，为用户添加仓单

	if err != nil {
		fmt.Printf("admin.go 客户端调用rpc执行失败！err: %v\n", err)
		return
	}

	fmt.Printf("rpc客户端调用成功。返回结果：%s\n", rst.RstDetails)
	fmt.Printf("<---------------------------->\n")
}
