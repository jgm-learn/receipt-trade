package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "receipt-trade/web/controllers/grpcPB"
	"receipt-trade/web/models"
	"time"
)

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {
	this.TplName = "admin.html"
}

//添加仓单信息
func (this *AdminController) AddReceipt() {
	this.TplName = "admin.html"
	var receipt models.Receipt //仓单结构体

	body := this.Ctx.Input.RequestBody //获取http数据
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &receipt); err == nil {
		fmt.Println(receipt) //打印仓单
		receipt.Insert()     //插入数据库
	} else {
		fmt.Println(err)
	}
}

//为用户添加仓单数量
func (this *AdminController) AddUserReceipt() {
	fmt.Printf("<---------------------------->\n")
	fmt.Printf("admin.go 第38行代码 为用户添加仓单执行输出如下：\n")
	this.TplName = "admin.html"
	var userReceipt models.UserReceipt //用户仓单结构体

	body := this.Ctx.Input.RequestBody //获取http数据
	fmt.Printf("前端传来数据如下：\n")
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &userReceipt); err == nil {
		userReceipt.Insert() //写入数据库
	} else {
		fmt.Println("json Unmarshal 出错：")
		fmt.Println(err)
	}

	fmt.Printf("\n")
	//查询用户的以太坊地址
	var userAddr = getUserAddr(userReceipt.UserId)
	if len(userAddr) == 0 {
		fmt.Printf("查询数据库失败，用户%d不存在\n", userReceipt.UserId)
		return
	}

	fmt.Printf("\n")
	//调用智能合约-建立连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("admin.go 第62行执行失败，grpc客户端连接失败！ err: %v\n", err)
	}

	defer conn.Close()

	client := pb.NewRPCServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var request pb.IstRctRequest
	request.UserAddr = userAddr
	request.ReceiptId = int64(userReceipt.ReceiptId)
	request.TotalQty = int64(userReceipt.TotalQuantity)
	rst, err := client.InsertUserReceipt(ctx, &request) //调用智能合约，为用户添加仓单

	if err != nil {
		fmt.Printf("admin.go 第78行执行失败，客户端rpc执行失败！err: %v\n", err)
		return
	}

	fmt.Printf("rpc客户端调用成功。返回结果：%s\n", rst.RstDetails)
	fmt.Printf("<---------------------------->\n")
}

//为用户增加资金
func (this *AdminController) AddUserFunds() {
	this.TplName = "admin.html"
	var userFunds models.UserFunds //用户资金结构体

	body := this.Ctx.Input.RequestBody //获取http数据
	fmt.Println(string(body))          //打印接收到的数据
	if err := json.Unmarshal(body, &userFunds); err == nil {
		fmt.Println(userFunds) //打印反序列化后的结果
		userFunds.Insert()     //插入数据库
	} else {
		fmt.Println("json Unmarshal 出错：")
		fmt.Println(err)
	}

	//查询用户的以太坊地址
	var userAddr = getUserAddr(userFunds.UserId)
	if len(userAddr) == 0 {
		fmt.Printf("查询数据库失败，用户%d不存在\n", userFunds.UserId)
		return
	}
	fmt.Printf("\n")
	//调用智能合约-建立连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("admin.go 第114行执行失败，grpc客户端连接失败！ err: %v\n", err)
	}
	defer conn.Close()
	client := pb.NewRPCServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var request pb.IstFundsRequest
	request.UserAddr = userAddr
	request.TotalFunds = int64(userFunds.TotalFunds)
	rst, err := client.InsertUserFunds(ctx, &request) //调用智能合约，为用户添加仓单
	if err != nil {
		fmt.Printf("admin.go 第125行执行失败，客户端rpc执行失败！err: %v\n", err)
		return
	}
	fmt.Printf("rpc客户端调用成功。返回结果：%s\n", rst.RstDetails)
	fmt.Printf("<---------------------------->\n")
}

//查询用户的以太坊地址
func getUserAddr(userId int) string {
	var userInfo models.User
	o.QueryTable("user").Filter("user_id", userId).One(&userInfo) //查询用户信息
	return userInfo.PublicKey
}
