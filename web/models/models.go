package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId    int `orm:"PK;auto"`
	UserName  string
	PassWord  string
	PublicKey string
	Nonce     int
}

//仓单
type Receipt struct {
	ReceiptId      int    `orm:"PK;auto"` //仓单编号
	Class          string //品种
	ProductionDate string //产期
	Level          string //等级
	Warehouse      string //仓库
	Provenance     string //产地
}

//用户所拥有的仓单
type UserReceipt struct {
	Id           int //id 主键
	UserId       int //用户id
	ReceiptId    int //仓单编号
	QtyTotal     int //仓单总量
	QtyAvailable int //剩余数量
	QtyFrozen    int //冻结数量
}

//资金表
type UserFunds struct {
	UserId         int `orm:"PK"` //用户id
	TotalFunds     int //资金总量
	AvailableFunds int //可用资金
	FrozenFunds    int //冻结资金
}

//卖方订单表
type OrderSell struct {
	Id        int
	UserId    int
	ReceiptId int
	Price     int
	QtySell   int
	NonceSell int
	SigSell   string
	AddrSell  string
}

//市场表
type List struct {
	ListId    int `orm:"PK"`
	UserId    int
	ReceiptId int
	Price     int
	QtySell   int
	QtyDeal   int
	QtyRemain int
}

type OrderBuy struct {
	Id        int
	UserId    int
	ListId    int
	ReceiptId int
	Price     int
	QtyBuy    int
	NonceBuy  int
	SigBuy    string
	AddrBuy   string
}

type WebReply struct {
	Reply string
}

type TradeReq struct {
	ReceiptId int
	Price     int
	QtySell   int
	NonceSell int
	QtyBuy    int
	NonceBuy  int
	AddrSell  string
	AddrBuy   string
	SigSell   string
	SigBuy    string
}

var o orm.Ormer //定义Ormer接口变量

func init() {
	//连接数据库
	orm.RegisterDataBase("default", "mysql", "root:root@/receipt_trade?charset=utf8", 30)
	//注册声明的model
	orm.RegisterModel(new(User), new(Receipt), new(UserReceipt), new(UserFunds), new(OrderSell), new(List))
	//创建表
	orm.RunSyncdb("default", false, true)

	o = orm.NewOrm() //创建orm结构体实例
}

func (user User) Insert() {
	id, err := o.Insert(&user)
	if err == nil {
		fmt.Printf("插入数据库 id = %d\n", id)
	} else {
		fmt.Printf("插入数据库失败err: %v\n", err)
	}
}

func (receipt Receipt) Insert() {
	id, err := o.Insert(&receipt)
	if err == nil {
		fmt.Printf("插入数据库 id = %d\n", id)
	} else {
		fmt.Printf("插入数据库失败err: %v\n", err)
	}
}

func (userRct UserReceipt) TableUnique() [][]string {
	return [][]string{
		[]string{"user_id", "receipt_id"},
	}
}

func (userRct UserReceipt) Insert() {
	id, err := o.Insert(&userRct)
	if err == nil {
		fmt.Printf("UserReceipt.Insert() 为用户添加仓单成功 id = %d\n", id)
	} else {
		fmt.Printf("UserReceipt.Insert() 为用户添加仓单失败 err: %v\n", err)
	}
}

func (userFus UserFunds) Insert() {
	id, err := o.Insert(&userFus)
	if err == nil {
		fmt.Printf("UserFunds 插入数据库 id = %d\n", id)
	} else {
		fmt.Printf("UserFunds 插入数据库失败err: %v\n", err)
	}
}

func (orderSell OrderSell) Insert() int64 {
	id, err := o.Insert(&orderSell)
	if err == nil {
		fmt.Printf("OrderSell 插入数据库 id = %d\n", id)
	} else {
		fmt.Printf("OrderSell 插入数据库失败err: %v\n", err)
	}
	return id
}

func (list List) Insert() {
	id, err := o.Insert(&list)
	if err == nil {
		fmt.Printf("Market 插入数据库 id = %d\n", id)
	} else {
		fmt.Printf("Market 插入数据库失败err: %v\n", err)
	}
}
